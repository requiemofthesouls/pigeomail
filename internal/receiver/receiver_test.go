package receiver

import (
	"net/smtp"
	"testing"
)

func TestReceiveEmail(t *testing.T) {
	smtpHost := "127.0.0.1"
	smtpPort := "21025"
	authEmail := "aaa@gmail.com"

	// Устанавливаем соединение с SMTP-сервером
	addr := smtpHost + ":" + smtpPort
	conn, err := smtp.Dial(addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Указываем заголовки письма
	from := "Your Name <your_email@gmail.com>"
	to := []string{"keepo@pigeomail.ddns.net"}
	subject := "Test Email"
	body := "Hello,\n\nThis is a test email."

	// Формируем тело письма
	message := []byte("To: " + to[0] + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// Отправляем письмо через установленное SMTP-соединение
	if err := smtp.SendMail(addr, nil, authEmail, to, message); err != nil {
		panic(err)
	}
}
