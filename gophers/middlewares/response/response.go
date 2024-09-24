package response

import (
	"encoding/json"
	"net/http"
)

// ResponseMiddleware wraps a handler to standardize the JSON response
func ResponseMiddleware(process func(w http.ResponseWriter, r *http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to capture the response
		rw := NewResponseWriter(w)
		resp, err := process(w, r)

		// Err Response
		if err != nil {
			code := getErrCode(err)
			res := struct {
				Status  string `json:"status"`
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{
				Status:  "fail",
				Code:    code,
				Message: err.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		// Succesful Response
		res := struct {
			Status string      `json:"status"`
			Data   interface{} `json:"data,omitempty"`
			Code   int         `json:"code"`
		}{
			Status: "success",
			Code:   200,
			Data:   resp,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(rw.StatusCode)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}

// ResponseWriter is a custom ResponseWriter that captures the response body
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       interface{}
}

// NewResponseWriter creates a new instance of ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		Body:           nil,
	}
}

func getErrCode(err error) int {
	return 500
}
