package authdto

type SendOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
	Code  string `json:"code" validate:"required,number"`
}
