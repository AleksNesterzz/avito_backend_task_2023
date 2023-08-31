package tests

import (
	"avito_backend/internal/http-server/handlers/getLogs"
	"avito_backend/internal/http-server/handlers/getLogs/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLogHandler(t *testing.T) {
	cases := []struct {
		name      string
		id        int
		time      string
		respError string
		mockError error
	}{
		{
			name: "Log",
			id:   1001,
			time: "2023-08",
		},
		{
			name: "Other log",
			id:   123,
			time: "2023-08",
		},
		{
			name: "Empty log",
			id:   1004,
			time: "2024-08",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			LogGetterMock := mocks.NewLogGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				LogGetterMock.On("GetLogs", tc.id, tc.time).
					Return([][]string{}, tc.mockError).
					Once()
			}

			handler := getLogs.New(LogGetterMock)

			input := fmt.Sprintf(`{"id":%d, "time":"%s"}`, tc.id, tc.time)

			req, err := http.NewRequest(http.MethodGet, "/log", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp getLogs.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
