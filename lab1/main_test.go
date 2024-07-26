package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetTemp(t *testing.T) {
	tests := []struct {
		name             string
		cep              string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "valid cep",
			cep:            "50710465",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid cep format",
			cep:            "0100",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "non-existent cep",
			cep:            "00000000",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/temperatura/", nil)
			req.SetPathValue("cep", tt.cep)
			w := httptest.NewRecorder()

			HandleGetTemp(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}
