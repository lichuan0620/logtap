package httputil

import (
	"encoding/json"
	"net/http"
)

// WriteGetResponse composes a response for a GET request.
func WriteGetResponse(w http.ResponseWriter, body interface{}, err error, headers ...HeaderField) {
	writeResponse(w, http.StatusOK, body, err, headers...)
}

// WritePostResponse composes a response for a POST request.
func WritePostResponse(w http.ResponseWriter, body interface{}, err error, headers ...HeaderField) {
	writeResponse(w, http.StatusCreated, body, err, headers...)
}

func writeResponse(w http.ResponseWriter, successCode int, body interface{}, err error, headers ...HeaderField) {
	for _, h := range headers {
		w.Header().Add(h.Key(), h.Value())
	}
	if err != nil {
		body = err.Error()
		if hErr, ok := err.(*httpError); ok {
			w.WriteHeader(hErr.Code)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(successCode)
	}
	if body != nil {
		data, err := json.MarshalIndent(body, "", "    ")
		if err == nil {
			w.Write(data)
		}
	}
}
