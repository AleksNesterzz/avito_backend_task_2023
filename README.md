# Сборка и запуск приложения происходит с помощью команды   
`docker build -t go-app . && docker-compose up --build go-app`
# Для запуска тестов тестов можно воспользоваться командой   
 `docker-compose run go-app go test ./internal/http-server/handlers/createSeg ./internal/http-server/handlers/deleteSeg ./internal/http-server/handlers/changeUser ./internal/http-server/handlers/getClientSeg`
