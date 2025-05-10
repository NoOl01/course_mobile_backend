package otp

import (
	"math/rand"
	"time"
)

var Store = make(map[string]int)

func StoreOTP(email string, otp int) {
	Store[email] = otp
	go func() {
		time.Sleep(2 * time.Minute)
		delete(Store, email)
	}()
}

func VerifyOTP(email string, otp int) bool {
	storedCode, exists := Store[email]
	return exists && storedCode == otp
}

func Generate() int {
	return rand.Intn(1000000)
}
