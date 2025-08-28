package service

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

type Validator struct {
	otpLength  int
	phoneRegex string
}

func newValidator(otpLength int, phoneRegex string) Validator {
	return Validator{otpLength: otpLength, phoneRegex: phoneRegex}
}

func (v Validator) validateSendOtp(req SendOtpRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Phone, validation.Required, validation.Match(regexp.MustCompile(v.phoneRegex))),
	)
}

func (v Validator) validateVerifyOtp(req VerifyOtpRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Otp, validation.Required, validation.Length(6, 6)),
		validation.Field(&req.Phone, validation.Required, validation.Match(regexp.MustCompile(v.phoneRegex))),
	)
}

func (v Validator) validateRefreshToken(req RefreshRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.RefreshToken, validation.Required),
	)
}
