package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// APIMiddleware wraps a handler to standardize the JSON response
func APIMiddleware(process func(w http.ResponseWriter, r *http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create a custom Writer to capture the response
		rw := NewResponseWriter(w)
		resp, err := process(w, r.WithContext(r.Context()))

		// Err Response
		if err != nil {
			// TODO: Translate errors into http Errors

			code := getErrCode(err)
			res := Response{
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

		res := Response{
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

// Writer is a custom Writer that captures the response body
type Writer struct {
	http.ResponseWriter
	StatusCode int
	Body       interface{}
}

// NewResponseWriter creates a new instance of Writer
func NewResponseWriter(w http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		Body:           nil,
	}
}

func getErrCode(_ error) int {
	return 500
}
