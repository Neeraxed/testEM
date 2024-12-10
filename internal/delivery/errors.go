package delivery

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	detail string
}

func ReturnHttpError(w http.ResponseWriter, e error) {
	he := &HttpError{
		detail: e.Error(),
	}
	resp, err := json.Marshal(he)
	if err != nil {
		return
	}
	w.Write(resp)
}
