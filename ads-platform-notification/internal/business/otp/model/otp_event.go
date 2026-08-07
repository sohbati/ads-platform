package model

const OtpSubject = "notifications.otp.send"

type OtpEvent struct {
	Mobile string `json:"mobile"`
	Otp    string `json:"otp"`
}
