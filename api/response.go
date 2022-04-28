package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"shopingList/internal"
)

var statusCtxKey = NewContextKey("Status")

// JSON map alias
type JSON map[string]interface{}

type response struct {
	Data  interface{} `json:"data"`
	Error JSON        `json:"error"`
}

// Status sets http status code in request
func Status(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), statusCtxKey, status))
}

// SendErrorJSON writes error data response into ResponseWriter
func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string, errCode int) {
	resp := response{
		Error: JSON{"code": errCode, "message": internal.WrappedError{Mess: details, Err: err}.Error()},
	}
	Status(r, httpStatusCode)
	SendJSON(w, r, resp)
}

// SendDataJSON writes data response into ResponseWriter
func SendDataJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, data interface{}) {
	resp := response{Data: data}
	Status(r, httpStatusCode)
	SendJSON(w, r, resp)
}

// SendJSON encodes input data to json object and write it into ResponseWriter
func SendJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(buf.Bytes()) // nolint: errcheck, gosec - not critic here
}
