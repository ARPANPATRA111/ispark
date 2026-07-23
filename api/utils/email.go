package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// GenerateOTP generates a cryptographically secure 6-digit numeric string
func GenerateOTP() (string, error) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 || err != nil {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

// Delivery timeouts. Sending mail happens inside a user-facing request, so it
// must fail fast rather than hold the response open: hosting providers commonly
// block outbound SMTP, and a blocked connection otherwise stalls until the
// operating system's TCP timeout (observed at ~135s on Render, which made
// registration look like it had hung forever).
const (
	emailAPITimeout  = 15 * time.Second
	emailSMTPTimeout = 10 * time.Second
)

// brevoSendEndpoint is a variable so tests can point it at a local server.
var brevoSendEndpoint = "https://api.brevo.com/v3/smtp/email"

// brevoSendEndpointForTest redirects sends to url and returns a restore func.
func brevoSendEndpointForTest(url string) func() {
	previous := brevoSendEndpoint
	brevoSendEndpoint = url
	return func() { brevoSendEndpoint = previous }
}

// SendOTP delivers the OTP, preferring transports that work everywhere.
//
//  1. Brevo's HTTP API (BREVO_API_KEY) over HTTPS — port 443 is never blocked,
//     which SMTP ports frequently are on PaaS hosts.
//  2. SMTP (SMTP_HOST/USER/PASS), with an explicit dial timeout.
//  3. Console logging, so local development needs no mail credentials at all.
//
// A delivery failure is logged and returns nil: the OTP is already stored, and
// failing the whole registration because mail is down would be worse than
// letting the user request a resend.
func SendOTP(toEmail, otpCode, purpose string) error {
	subject := "Verification Code for iSpark"
	body := fmt.Sprintf("Your OTP code for %s is: %s\nThis code expires in 15 minutes.", purpose, otpCode)

	sender := strings.TrimSpace(os.Getenv("SMTP_SENDER"))

	if apiKey := strings.TrimSpace(os.Getenv("BREVO_API_KEY")); apiKey != "" {
		if err := sendViaBrevoAPI(apiKey, sender, toEmail, subject, body); err != nil {
			log.Printf("Brevo API send failed: %v", err)
			logEmailFallback(toEmail, subject, body)
			return nil
		}
		log.Printf("OTP email sent to %s via the Brevo API", toEmail)
		return nil
	}

	smtpHost := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	smtpUser := strings.TrimSpace(os.Getenv("SMTP_USER"))
	smtpPass := strings.TrimSpace(os.Getenv("SMTP_PASS"))
	if smtpHost == "" || smtpUser == "" || smtpPass == "" {
		logEmailFallback(toEmail, subject, body)
		return nil
	}

	smtpPort := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if smtpPort == "" {
		smtpPort = "587"
	}
	if err := sendViaSMTP(smtpHost, smtpPort, smtpUser, smtpPass, sender, toEmail, subject, body); err != nil {
		log.Printf("SMTP send failed (hosts often block outbound SMTP; set BREVO_API_KEY to send over HTTPS instead): %v", err)
		logEmailFallback(toEmail, subject, body)
		return nil
	}

	log.Printf("OTP email successfully sent to %s", toEmail)
	return nil
}

func logEmailFallback(toEmail, subject, body string) {
	log.Printf("\n--- [EMAIL NOT DELIVERED - LOGGING INSTEAD] ---\nTo: %s\nSubject: %s\nBody: %s\n----------------------------------------------\n",
		toEmail, subject, body)
}

type brevoContact struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type brevoPayload struct {
	Sender      brevoContact   `json:"sender"`
	To          []brevoContact `json:"to"`
	Subject     string         `json:"subject"`
	TextContent string         `json:"textContent"`
}

func sendViaBrevoAPI(apiKey, sender, toEmail, subject, body string) error {
	if sender == "" {
		return fmt.Errorf("SMTP_SENDER is not set; Brevo requires a verified sender address")
	}

	payload, err := json.Marshal(brevoPayload{
		Sender:      brevoContact{Email: sender, Name: "iSPARC"},
		To:          []brevoContact{{Email: toEmail}},
		Subject:     subject,
		TextContent: body,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, brevoSendEndpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: emailAPITimeout}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		detail, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("status %d: %s", res.StatusCode, strings.TrimSpace(string(detail)))
	}
	return nil
}

// sendViaSMTP mirrors smtp.SendMail but dials with a timeout, so a host that
// silently drops outbound SMTP fails in seconds instead of minutes.
func sendViaSMTP(host, port, user, pass, sender, toEmail, subject, body string) error {
	if sender == "" {
		sender = user
	}

	msg := []byte(
		"From: iSpark <" + sender + ">\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), emailSMTPTimeout)
	if err != nil {
		return fmt.Errorf("dial %s:%s: %w", host, port, err)
	}
	_ = conn.SetDeadline(time.Now().Add(emailSMTPTimeout))

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer func() { _ = client.Close() }()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
			return err
		}
	}
	if err := client.Auth(smtp.PlainAuth("", user, pass, host)); err != nil {
		return err
	}
	if err := client.Mail(sender); err != nil {
		return err
	}
	if err := client.Rcpt(toEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}
