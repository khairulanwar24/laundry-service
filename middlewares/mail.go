package middleware

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// Mail mengirim email polos via SMTP memakai kredensial dari environment
// variable (SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM).
// Jika SMTP_HOST belum dikonfigurasi, email tidak dikirim dan isinya
// hanya dicatat ke log — supaya alur development/testing tidak terblokir.
func Mail(to, subject, body string) (string, error) {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		log.Printf("[mail] SMTP belum dikonfigurasi, lewati pengiriman. To=%s Subject=%s Body=%s", to, subject, body)
		return "SMTP belum dikonfigurasi, email tidak dikirim", nil
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = username
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n", from, to, subject, body)

	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		return "", fmt.Errorf("gagal mengirim email: %w", err)
	}

	return "Email berhasil dikirim", nil
}
