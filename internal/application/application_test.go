package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/mushtaev-a/rpn/internal/application"
)

func TestHandleCalculator(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		body           map[string]string
		expectedStatus int
		expectedResult float64
		expectError    bool
	}{
		{
			name:           "valid simple expression",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "2+2"},
			expectedStatus: http.StatusOK,
			expectedResult: 4,
			expectError:    false,
		},
		{
			name:           "valid complex expression",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "(2+2)*3"},
			expectedStatus: http.StatusOK,
			expectedResult: 12,
			expectError:    false,
		},
		{
			name:           "empty expression",
			method:         http.MethodPost,
			body:           map[string]string{"expression": ""},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
		{
			name:           "invalid expression",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "2++2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name:           "invalid parentheses",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "(2+2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name:           "division by zero",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "1/0"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				bodyBytes, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, "/api/v1/calculate", bytes.NewBuffer(bodyBytes))
			} else {
				req = httptest.NewRequest(tc.method, "/api/v1/calculate", nil)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.HandleCalculation)

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.expectedStatus)
			}

			if !tc.expectError && tc.expectedStatus == http.StatusOK {
				var response struct {
					Result float64 `json:"result"`
				}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Could not decode response: %v", err)
				}
				if response.Result != tc.expectedResult {
					t.Errorf("handler returned wrong result: got %v want %v", response.Result, tc.expectedResult)
				}
			}
		})
	}
}
