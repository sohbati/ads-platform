package model

type SessionCustomer struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	NationalId string `json:"national_id"`
}

type VerifyOtpRequest struct {
	Otp string `json:"otp" binding:"required,len=6,numeric"`
}

type LoginResponse struct {
	Authenticated bool             `json:"authenticated"`
	Customer      SessionCustomer  `json:"customer"`
}

type MeResponse struct {
	Authenticated bool              `json:"authenticated"`
	Customer      *SessionCustomer  `json:"customer,omitempty"`
}
