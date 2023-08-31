package deleteSeg_test

import (
	"avito_backend/internal/http-server/handlers/deleteSeg"
	"avito_backend/internal/http-server/handlers/deleteSeg/mocks"
	"avito_backend/internal/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateSegHandler(t *testing.T) {
	cases := []struct {
		name      string
		segment   string
		respError string
		mockError error
	}{
		{
			name:    "Delete segment",
			segment: "TESTING1",
		},
		{
			name:    "Delete another one",
			segment: "TESTING2",
		},
		{
			name:      "Empty delete",
			segment:   "TESTING1",
			respError: "segment not found",
			mockError: storage.ErrSegNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			SegDeleterMock := mocks.NewSegDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				SegDeleterMock.On("DeleteSeg", tc.segment).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := deleteSeg.New(SegDeleterMock)

			input := fmt.Sprintf(`{"name": "%s"}`, tc.segment)

			req, err := http.NewRequest(http.MethodDelete, "/seg/del", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if tc.name != "Empty delete" {
				require.Equal(t, rr.Code, http.StatusOK)
			} else {
				require.Equal(t, rr.Code, http.StatusBadRequest)
			}

			body := rr.Body.String()

			var resp deleteSeg.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
