package utils

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGenerateOTPIsSixDigits(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		otp, err := GenerateOTP()
		if err != nil {
			t.Fatalf("GenerateOTP: %v", err)
		}
		if len(otp) != 6 {
			t.Fatalf("otp %q is not 6 characters", otp)
		}
		if strings.Trim(otp, "0123456789") != "" {
			t.Fatalf("otp %q contains non-digits", otp)
		}
		seen[otp] = true
	}
	if len(seen) < 2 {
		t.Fatal("GenerateOTP returned the same value every time")
	}
}

// The Brevo HTTP transport must be preferred when a key is present, because
// SMTP ports are frequently blocked on hosting platforms.
func TestSendOTPUsesBrevoAPIWhenKeyPresent(t *testing.T) {
	var gotPath, gotKey string
	var payload brevoPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotKey = r.Header.Get("api-key")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &payload)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"messageId":"ok"}`))
	}))
	defer srv.Close()

	restore := brevoSendEndpointForTest(srv.URL + "/v3/smtp/email")
	defer restore()

	t.Setenv("BREVO_API_KEY", "test-key")
	t.Setenv("SMTP_SENDER", "sender@example.com")

	if err := SendOTP("student@example.com", "123456", "registration"); err != nil {
		t.Fatalf("SendOTP: %v", err)
	}

	if gotPath != "/v3/smtp/email" {
		t.Fatalf("called %q", gotPath)
	}
	if gotKey != "test-key" {
		t.Fatalf("api-key header = %q", gotKey)
	}
	if payload.Sender.Email != "sender@example.com" {
		t.Fatalf("sender = %q", payload.Sender.Email)
	}
	if len(payload.To) != 1 || payload.To[0].Email != "student@example.com" {
		t.Fatalf("recipient = %+v", payload.To)
	}
	if !strings.Contains(payload.TextContent, "123456") {
		t.Fatalf("body missing the OTP: %q", payload.TextContent)
	}
}

// A failing provider must not fail the caller: the OTP is already stored and
// the user can request a resend.
func TestSendOTPSucceedsWhenProviderRejects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"bad key"}`))
	}))
	defer srv.Close()

	restore := brevoSendEndpointForTest(srv.URL)
	defer restore()

	t.Setenv("BREVO_API_KEY", "bad")
	t.Setenv("SMTP_SENDER", "sender@example.com")

	if err := SendOTP("student@example.com", "123456", "registration"); err != nil {
		t.Fatalf("SendOTP should swallow provider errors, got %v", err)
	}
}

// The regression this guards: an unreachable SMTP host used to stall the
// request until the OS TCP timeout (~135s in production). It must now give up
// quickly and fall back to logging.
func TestSMTPFailsFastWhenHostIsUnreachable(t *testing.T) {
	// A listener that accepts connections but never speaks SMTP, so the client
	// blocks on the greeting until the deadline fires.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() { _ = ln.Close() }()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			defer func() { _ = conn.Close() }()
		}
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())

	t.Setenv("BREVO_API_KEY", "")
	t.Setenv("SMTP_HOST", host)
	t.Setenv("SMTP_PORT", port)
	t.Setenv("SMTP_USER", "user")
	t.Setenv("SMTP_PASS", "pass")
	t.Setenv("SMTP_SENDER", "sender@example.com")

	start := time.Now()
	if err := SendOTP("student@example.com", "123456", "registration"); err != nil {
		t.Fatalf("SendOTP: %v", err)
	}
	elapsed := time.Since(start)

	if elapsed > emailSMTPTimeout+5*time.Second {
		t.Fatalf("SMTP send took %s; it must fail fast rather than stall the request", elapsed)
	}
}

func TestSendOTPFallsBackToLoggingWithoutCredentials(t *testing.T) {
	t.Setenv("BREVO_API_KEY", "")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_USER", "")
	t.Setenv("SMTP_PASS", "")

	if err := SendOTP("student@example.com", "123456", "registration"); err != nil {
		t.Fatalf("SendOTP with no credentials should log and succeed, got %v", err)
	}
}
