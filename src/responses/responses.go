// Agent responses
package responses

import (
	"encoding/json"
	"net/http"
)

type ResponseComm struct {
	ID       uint64 `json:"id"`
	Response string `json:"response"`
}

func Response(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}
