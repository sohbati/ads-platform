package model

type SessionUser struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	NationalId string `json:"national_id"`
}

type VerifyOtpRequest struct {
	Otp string `json:"otp" binding:"required,len=6,numeric"`
}

type LoginResponse struct {
	Authenticated bool        `json:"authenticated"`
	User          SessionUser `json:"user"`
}

type MeResponse struct {
	Authenticated bool         `json:"authenticated"`
	User          *SessionUser `json:"user,omitempty"`
}
