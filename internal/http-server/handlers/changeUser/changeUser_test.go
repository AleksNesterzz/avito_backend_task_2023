package changeUser_test

import (
	"avito_backend/internal/http-server/handlers/changeUser"
	"avito_backend/internal/http-server/handlers/changeUser/mocks"
	"avito_backend/internal/http-server/handlers/createSeg"
	"avito_backend/internal/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangeUserHandler(t *testing.T) {
	cases := []struct {
		name      string
		addseg    []string
		delseg    []string
		id        int
		respError string
		mockError error
	}{
		{
			name:   "change user",
			addseg: []string{"TESTING1", "TESTING2", "TESTING3"},
			delseg: []string{"TESTING3"},
			id:     1020,
		},
		{
			name:   "Add another one",
			addseg: []string{},
			delseg: []string{"TESTING2"},
			id:     1020,
		},
		{
			name:      "PKEY_constraint",
			addseg:    []string{"TESTING1"},
			delseg:    []string{},
			id:        1020,
			respError: "segment already exists",
			mockError: storage.ErrSegExists,
		},
		{
			name:      "Not found",
			addseg:    []string{},
			delseg:    []string{"TESTING2"},
			id:        1020,
			respError: "user not in segment",
			mockError: storage.ErrUserNotFound,
		},
		{
			name:   "change another one",
			addseg: []string{"TESTING123"},
			delseg: []string{},
			id:     1002,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			UserChangerMock := mocks.NewUserChanger(t)

			if tc.respError == "" || tc.mockError != nil {
				UserChangerMock.On("ChangeUser", tc.addseg, tc.delseg, tc.id).
					Return("", tc.mockError).
					Once()
			}

			handler := changeUser.New(UserChangerMock)

			input := `{"addseg":[`
			for i := 0; i < len(tc.addseg); i++ {
				input = input + "\"" + tc.addseg[i] + "\" ,"
			}

			input = strings.TrimSuffix(input, ",")

			input = input + "], \"delseg\":["

			for i := 0; i < len(tc.delseg); i++ {
				input = input + "\"" + tc.delseg[i] + "\","
			}

			input = strings.TrimSuffix(input, ",")

			input = input + "],"

			input_num := fmt.Sprintf(` "id":%d}`, tc.id)

			input = input + input_num

			req, err := http.NewRequest(http.MethodPut, "/user", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.name == "PKEY_constraint" || tc.name == "Not found" {
				require.Equal(t, rr.Code, http.StatusBadRequest)
			} else {
				require.Equal(t, rr.Code, http.StatusOK)
			}

			body := rr.Body.String()

			var resp createSeg.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
