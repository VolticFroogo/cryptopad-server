package helper

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the type used for error JSON responses.
type ErrorResponse struct {
	Error string
}

// JSONResponse sends a client a JSON response.
func JSONResponse(data interface{}, status int, w http.ResponseWriter) (err error) {
	// Set the status header of the response.
	w.WriteHeader(status)

	dataJSON, err := json.Marshal(data) // Encode response into JSON.
	if err != nil {
		return
	}

	w.Write(dataJSON) // Write JSON data to response writer.
	return
}

// ThrowErr is used for throwing errors via a JSON response.
func ThrowErr(err error, status int, w http.ResponseWriter) {
	// Send the error as a JSON response.
	JSONResponse(ErrorResponse{
		Error: err.Error(),
	}, status, w)
}
