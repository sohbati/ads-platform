package model

import "encoding/json"

type PublishRequest struct {
	Data json.RawMessage `json:"data" binding:"required"`
}

type PublishResponse struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}
