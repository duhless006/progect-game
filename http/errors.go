package http

import (
	"encoding/json"
	"time"
)

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewError(message error) ErrorDTO {
	return ErrorDTO{
		Message: message.Error(),
		Time:    time.Now(),
	}
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return `{"error": "failed to marshal"}`
	}
	return string(b)
}
