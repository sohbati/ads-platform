package model

type SessionHttpRequestModel struct {
	Data string `json:"data" binding:"required"`
}
