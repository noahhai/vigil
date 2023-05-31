package utils

import (
	"encoding/json"
	"net/http"
)

func WriteDataResponse(w http.ResponseWriter, data interface{}, err error) {
	resp := respObj{
		Data: data,
	}
	if err != nil  {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err.Error()
	}
	if marshalled, err := json.Marshal(resp); err == nil {
		w.Write(marshalled)
	} else {
		m, _ := json.Marshal(respObj{})
		w.Write(m)
	}
}
func WriteDataResponseWithStatus(w http.ResponseWriter, data interface{}, err error, statusCode int) {
	resp := respObj{
		Data: data,
	}
	if err != nil {
		w.WriteHeader(statusCode)
		resp.Error = err.Error()
	}

	if marshalled, err := json.Marshal(resp); err == nil {
		w.Write(marshalled)
	} else {
		m, _ := json.Marshal(respObj{})
		w.Write(m)
	}
}

type respObj struct {
	Data interface{}
	Error string
}
