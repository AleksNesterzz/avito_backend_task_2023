package postgres

import (
	_ "avito_backend/internal/lib/csv_log"
	"avito_backend/internal/storage"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	_ "errors"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS segments(
		id SERIAL,
		name TEXT PRIMARY KEY);`)
	if err != nil {
		return nil, fmt.Errorf("preparing statement error: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("execution statement error: %s", err)
	}
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS clients(
		id integer,
		seg TEXT references segments(name),
		PRIMARY KEY(id, seg)
		);`)
	if err != nil {
		return nil, fmt.Errorf("preparing statement error: %s", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("execution statement error: %s", err)
	}

	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS logs(
		id serial PRIMARY KEY,
		client integer,
		operation TEXT,
		seg TEXT,
		dt timestamp
		);`)

	//ЛИБО не делать каскадное удаление, либо добавить логику при каскадном удалении
	if err != nil {
		return nil, fmt.Errorf("preparing statement error: %s", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("execution statement error: %s", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateSeg(name string) (int64, error) {
	//const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO segments(name) VALUES($1);")
	if err != nil {
		return 0, fmt.Errorf("error preparing insert %w", err)
	}
	_, err = stmt.Exec(name)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Constraint == "segments_pkey" {
			return 0, storage.ErrSegExists
		}
		return 0, fmt.Errorf("error executing inserting %w", err)
	}

	return 0, nil
}

//ДОБАВИТЬ ВСЕХ КЛИЕНТОВ у которых я удалил сегмент в логи

func (s *Storage) DeleteSeg(name string) (int64, error) {
	stmt_sel, err := s.db.Prepare("SELECT id FROM clients WHERE seg=$1;")
	if err != nil {
		fmt.Println("error preparing deletion")
		return 0, fmt.Errorf("error preparing deletion %w", err)
	}
	// stmt_insert, err := s.db.Prepare("INSERT INTO logs(client, operation, seg, dt) VALUES($1,$2,$3,$4)")
	// if err != nil {
	// 	fmt.Println("error inserting for cascade delete")
	// 	return 0, fmt.Errorf("error inserting for cascade delete %w", err)
	// }
	fields, err := stmt_sel.Query(name)
	if err != nil {
		fmt.Println("error selecting for cascade delete")
		return 0, fmt.Errorf("error selecting for cascade delete %w", err)
	}
	for fields.Next() {
		var temp int
		fields.Scan(&temp)
		_, err := s.ChangeUser([]string{}, []string{name}, temp)
		if err != nil {
			fmt.Println("error executing insert for cascade delete")
			return 0, fmt.Errorf("error executing insert for cascade delete %w", err)
		}
	}
	defer fields.Close()
	stmt_del, err := s.db.Prepare("DELETE FROM segments WHERE name=$1;")
	if err != nil {
		return 0, fmt.Errorf("error preparing deletion %w", err)
	}
	res, err := stmt_del.Exec(name)
	if err != nil {
		fmt.Println("error here")
		return 0, fmt.Errorf("error executing deletion %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error executing deletion %w", err)
	} else {
		if rows == 0 {
			return 0, storage.ErrSegNotFound
		}
		return rows, nil
	}
}

func (s *Storage) ChangeUser(addSeg []string, delSeg []string, id int) (string, error) {

	//СДЕЛАТЬ ЧЕРЕЗ UPDATE?

	stmt_add, err := s.db.Prepare("INSERT INTO clients VALUES($1, $2)")
	if err != nil {
		//TODO : будет ли функция ломаться при одном неверном значении в массиве?
		return "", fmt.Errorf("error preparing adding client to segment %w", err)
	}
	stmt_del, err := s.db.Prepare("DELETE FROM clients WHERE id=$1 AND seg=$2")
	if err != nil {
		return "", fmt.Errorf("error preparing delete client from segment %w", err)
	}
	stmt_log, err := s.db.Prepare(`INSERT INTO logs(client,operation,seg,dt) VALUES($1,$2,$3,$4)`)
	if err != nil {
		return "", fmt.Errorf("error preparing log field %w", err)
	}
	//
	//попробовать использовать горутины и waitGroup
	//
	for i := 0; i < len(addSeg); i++ {
		_, err := stmt_add.Exec(id, addSeg[i])
		if err != nil {
			pqErr := err.(*pq.Error)
			if pqErr.Constraint == "clients_pkey" {
				return "", fmt.Errorf("pkey constraint %w", storage.ErrSegExists)
			}
			if pqErr.Constraint == "clients_seg_fkey" {
				return "", fmt.Errorf("fkey constraint %w", storage.ErrSegNotExists)
			}
			return "", fmt.Errorf("error executing inserting %w", err)
		}
		_, err = stmt_log.Exec(id, "add", addSeg[i], time.Now())
		if err != nil {
			return "", fmt.Errorf("error adding log %w", err)
		}
	}
	for i := 0; i < len(delSeg); i++ {
		res, err := stmt_del.Exec(id, delSeg[i])
		if err != nil {
			return "", fmt.Errorf("error executing deleting %w", err)
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			return "", fmt.Errorf("user has no such segment %w", storage.ErrUserNotFound)
		}
		_, err = stmt_log.Exec(id, "del", delSeg[i], time.Now())
		if err != nil {
			return "", fmt.Errorf("error adding log %w", err)
		}
	}
	return "", nil
}

func (s *Storage) GetClientSeg(id int) ([]string, error) {
	var segments []string
	stmt, err := s.db.Prepare("SELECT seg from clients WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("error preparing GetClientSeg %w", err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("error executing GetClientSeg %w", err)
	} else {
		var temp string
		for rows.Next() {
			rows.Scan(&temp)
			segments = append(segments, temp)
		}

	}
	defer rows.Close()
	return segments, nil
}

func (s *Storage) GetLogs(id int, year_month string) ([][]string, error) {
	logs := make([][]string, 1)
	date := strings.Split(year_month, "-")
	year, _ := strconv.Atoi(date[0])
	month, _ := strconv.Atoi(date[1])
	stmt, err := s.db.Prepare("SELECT client,seg,operation,dt from logs WHERE client=$1 AND date_part('year',dt)=$2 AND date_part('month', dt)=$3")
	if err != nil {
		return nil, fmt.Errorf("error preparing GetLogs %w", err)
	}
	rows, err := stmt.Query(id, year, month)
	if err != nil {
		return nil, fmt.Errorf("error executing GetClientLogs %w", err)
	} else {
		i := 0
		for rows.Next() {
			logs[i] = make([]string, 4)
			rows.Scan(&logs[i][0], &logs[i][1], &logs[i][2], &logs[i][3])
			i++
			logs = append(logs, []string{})
		}
	}
	defer rows.Close()
	return logs, nil
}
