package service

type SendOtpRequest struct {
	Phone string `json:"phone"`
}
type SendOtpResponse struct {
	Message string `json:"message"`
}

type VerifyOtpRequest struct {
	Phone    string `json:"phone"`
	Otp      string `json:"otp"`
	DeviceID string `json:"device_id,omitempty"`
}
type VerifyOtpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

type MeRequest struct {
	Claims any `json:"claims"`
}
type MeResponse struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
}

type LogoutRequest struct {
	Claims any `json:"claims"`
}
type LogoutResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}
