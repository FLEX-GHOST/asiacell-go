package asiacell

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidPasscode = errors.New("invalid sms passcode")
	ErrUnauthorized   = errors.New("unauthorized: missing or invalid session token")
	ErrMissingPID     = errors.New("missing PID in response")
	ErrRequestFailed  = errors.New("request failed")
	ErrCaptchaRequired = errors.New("captcha required by server")
)

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("asiacell api error (status %d): %s", e.StatusCode, e.Message)
}
