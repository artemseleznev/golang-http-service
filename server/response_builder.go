package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, httpResponseCode int, msg string) {
	errorCodeMap := map[int]string{
		http.StatusBadRequest:          "request_data_error",
		http.StatusMethodNotAllowed:    "method_not_allowed",
		http.StatusInternalServerError: "internal_server_error",
		http.StatusBadGateway:          "bad_gateway",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpResponseCode)
	errorResp := ErrorResp{Code: errorCodeMap[httpResponseCode], Message: msg}
	log.Printf("error: %d: %+v", httpResponseCode, errorResp)
	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("error: could not encode response body: %v", err)
	}
}

func okResponse(w http.ResponseWriter, data map[string]string) {
	resp, err := json.Marshal(data)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		log.Printf("error: could not write to response: %v", err)
	}
}
