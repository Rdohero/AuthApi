package controllers

import (
	"errors"
	"sync"
)

type OtpString struct {
	TokenString string
}

var (
	otpStore map[string]OtpString
	m        sync.Mutex
)

func init() {
	otpStore = make(map[string]OtpString)
}

func SimpanOtp(token, tokenString string) {
	m.Lock()
	defer m.Unlock()

	otpStore[token] = OtpString{
		TokenString: tokenString,
	}
}

func DapatkanOtpString(token string) (string, error) {
	m.Lock()
	defer m.Unlock()

	otp, ok := otpStore[token]
	if !ok {
		return "", errors.New("Otp salah")
	}

	return otp.TokenString, nil
}

func HapusOtp(token string) error {
	m.Lock()
	defer m.Unlock()

	_, ok := otpStore[token]
	if !ok {
		return errors.New("Otp salah")
	}

	delete(otpStore, token)

	return nil
}
