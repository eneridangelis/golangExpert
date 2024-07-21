package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetTemp(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/temperatura/", HandleGetTemp)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	tests := []struct {
		cep          string
		expectedCode int
	}{
		{"50710465", http.StatusOK},
		{"99999999", http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.cep, func(t *testing.T) {
			resp, err := http.Get(ts.URL + "/temperatura/" + tt.cep)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read response body: %v", err)
				}

				if !bytes.Contains(body, []byte("temp_C")) {
					t.Errorf("Response body does not contain expected key 'temp_C'")
				}
			}
		})
	}
}
