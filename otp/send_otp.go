package otp

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOtp(email string, otp int) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte(fmt.Sprintf(
		"Subject: Ваш код подтверждения\r\n"+
			"\r\n"+
			"Вы отправляли запрос на смену пароля\r\n"+
			"Если вы этого не делали просто проигнорируйте это сообщение\r\n"+
			"Код: %d", otp))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, message)
}
