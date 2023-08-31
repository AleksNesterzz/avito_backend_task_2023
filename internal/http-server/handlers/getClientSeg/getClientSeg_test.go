package getClientSeg_test

import (
	"avito_backend/internal/http-server/handlers/getClientSeg"
	"avito_backend/internal/http-server/handlers/getClientSeg/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetClientSegHandler(t *testing.T) {
	cases := []struct {
		name      string
		id        int
		respError string
		mockError error
	}{
		{
			name: "Many segments",
			id:   1001,
		},
		{
			name: "1 segment",
			id:   123,
		},
		{
			name: "No segments",
			id:   1004,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			SegGetterMock := mocks.NewSegGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				SegGetterMock.On("GetClientSeg", tc.id).
					Return([]string{}, tc.mockError).
					Once()
			}

			handler := getClientSeg.New(SegGetterMock)

			input := fmt.Sprintf(`{"id": %d}`, tc.id)

			req, err := http.NewRequest(http.MethodGet, "/user", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp getClientSeg.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
