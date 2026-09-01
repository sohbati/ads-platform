package model

type SendOtpResponse struct {
	Message            string `json:"message"`
	ResendAfterSeconds int    `json:"resend_after_seconds"`
}

type VerifyOtpRequest struct {
	Otp string `json:"otp" binding:"required,len=6,numeric"`
}

type VerifyOtpResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

type CacheOtpRequest struct {
	Otp string `json:"otp"`
}
