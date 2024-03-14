package Controller

import (
	"encoding/json"
	"net/http"

	m "Week6/Model"
)

func SendSuccessResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	var response m.ErrorResponse
	response.Status = code
	response.Message = message

	json.NewEncoder(w).Encode(response)
}

func SendErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	var response m.ErrorResponse
	response.Status = code
	response.Message = message

	json.NewEncoder(w).Encode(response)
}
