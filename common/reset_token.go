package common

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var tokens = make(map[string]string)

func StoreToken(email string) string {
	token := randomString(20)
	tokens[email] = token
	go func() {
		time.Sleep(30 * time.Minute)
		delete(tokens, email)
	}()
	return token
}

func VerifyToken(email string, token string) bool {
	storedToken, exists := tokens[email]
	return exists && storedToken == token
}

func randomString(length int) string {
	rand.Intn(1000000)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
