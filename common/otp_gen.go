package common

import "math/rand"

func OtpGenerate() int {
	return rand.Intn(1000000)
}
