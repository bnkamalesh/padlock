// Package totp generates the Time based OTP, based on:
// RFC 6238 https://tools.ietf.org/html/rfc6238#section-4.1
// RFC 4226 https://tools.ietf.org/html/rfc4226
package totp

import (
	"strings"
	"time"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type TOTP struct {
	appCtx    *appcontext.AppContext
	Issuer    string
	Label     string
	CreatedAt *time.Time
}

// Secret generates a new secret for TOTP for the given identifiers
func (t *TOTP) Secret(parts ...string) string {
	return strings.Join(parts, "")
}

func (t *TOTP) Verify(secret string) bool {
	return false
}

func (t *TOTP) URI() string {
	return ""
}
