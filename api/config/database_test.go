package config

import (
	"strings"
	"testing"
)

// Values pasted into a hosting dashboard often carry a trailing newline or
// space. In the libpq key=value DSN format whitespace terminates a value, so an
// untrimmed password used to be truncated mid-DSN and surfaced as
// "password authentication failed" even when the credential was correct.
func TestEnvTrimsPastedWhitespace(t *testing.T) {
	cases := map[string]string{
		"trailing newline": "s3cret\n",
		"trailing CRLF":    "s3cret\r\n",
		"trailing space":   "s3cret ",
		"leading space":    " s3cret",
		"surrounded":       "\t s3cret \n",
	}

	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TEST_DB_PASSWORD", raw)
			if got := env("TEST_DB_PASSWORD"); got != "s3cret" {
				t.Fatalf("env() = %q, want %q", got, "s3cret")
			}
		})
	}
}

func TestQuoteDSNProtectsDSNParsing(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  string
	}{
		{"plain value is left bare", "postgres", "postgres"},
		{"empty value is quoted", "", "''"},
		{"space is quoted", "pass word", "'pass word'"},
		{"single quote is escaped", "pa'ss", `'pa\'ss'`},
		{"backslash is escaped", `pa\ss`, `'pa\\ss'`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := quoteDSN(tc.value); got != tc.want {
				t.Fatalf("quoteDSN(%q) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}

// A password containing a space must not leak into the following DSN key.
func TestDSNKeepsFieldsSeparate(t *testing.T) {
	dsn := "host=" + quoteDSN("db.example.com") +
		" user=" + quoteDSN("postgres.ref") +
		" password=" + quoteDSN("pass word") +
		" dbname=" + quoteDSN("postgres")

	if !strings.Contains(dsn, "password='pass word' dbname=postgres") {
		t.Fatalf("password was not isolated from the next key: %s", dsn)
	}
}
