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

	message := []byte(fmt.Sprintf("Subject: Ваш код подтверждения\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		`<!DOCTYPE html>
		<html>
		<head>
			<style>
		        body {
		            font-family: Arial, sans-serif;
		            background-color: #f5f5f5;
		            padding: 20px;
		        }
		        .container {
		            max-width: 500px;
		            margin: auto;
		            background-color: #F7F7F9;
		            padding: 30px;
		            border-radius: 8px;
		        }
		        .header {
		            font-size: 20px;
		            margin-bottom: 20px;
		        }
		        .code {
		            font-size: 32px;
		            font-weight: bold;
		            color: #D8D8D8;
		            background-color: #48B2E7;
		            padding: 10px 20px;
		            display: inline-block;
		            border-radius: 6px;
		            margin: 20px 0;
		        }
		        .footer {
		            font-size: 14px;
		            color: #707B81;
		        }
		    </style>
		</head>
		<body>
		  	<div class="container">
		    	<div class="header">Запрос на смену пароля</div>
		    	<p>Вы отправили запрос на смену пароля. Если это были не вы, просто проигнорируйте это письмо.</p>
		    	<p>Ваш код:</p>
		    	<div class="code">%d</div>
		    	<p class="footer">Если вы не запрашивали смену пароля, игнорируйте это письмо.</p>
		  	</div>
		</body>
		</html>`, otp))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, message)
}
