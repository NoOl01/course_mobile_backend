package otp

import (
	"math/rand"
	"time"
)

var OtpStore = make(map[string]int)

func StoreOTP(email string, otp int) {
	OtpStore[email] = otp
	go func() {
		time.Sleep(2 * time.Minute)
		delete(OtpStore, email)
	}()
}

func VerifyOTP(email string, otp int) bool {
	storedCode, exists := OtpStore[email]
	return exists && storedCode == otp
}

func OtpGenerate() int {
	return rand.Intn(1000000)
}
