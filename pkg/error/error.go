package handler_error

import "encoding/json"

type errorHandler struct {
	Error string `json:"error"`
}

func ErrorHandler(error error) []byte {
	errStruct := &errorHandler{}
	errStruct.Error = error.Error()

	body, err := json.Marshal(&errStruct)

	if err != nil {
		return nil
	}

	return body
}
