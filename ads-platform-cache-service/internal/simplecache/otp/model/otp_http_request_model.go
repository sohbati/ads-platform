package model

type OtpHttpRequestModel struct {
	Otp string `json:"otp" binding:"required"`
}
