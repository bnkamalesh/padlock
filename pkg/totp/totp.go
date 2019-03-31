// Package totp generates the Time based OTP, based on:
// RFC 6238 https://tools.ietf.org/html/rfc6238#section-4.1
// RFC 4226 https://tools.ietf.org/html/rfc4226
package totp

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/bnkamalesh/padlock/pkg/appcontext"
)

type algo string

const (
	StatusActive   = "active"
	StatusInactive = "inactive"

	AlgoSHA1   = algo("SHA1")
	AlgoSHA256 = algo("SHA256")
	AlgoSHA512 = algo("SHA512")
)

type TOTP struct {
	appCtx *appcontext.AppContext `json:"-"`
	// Issuer is the URI encoded name of the service for which TOTP is generated
	Issuer string `json:"issuer,omitempty"`
	// Algorithm is the hashing algorithm used to generate the secret
	Algorithm algo `json:"algorithm,omitempty"`
	// Digits is the number of digits to be generated for the TOTP
	Digits int `json:"digits,omitempty"`
	// Period is the number of seconds TOTP is valid for
	Period int `json:"period,omitempty"`
	// Drift is the possible deviation in time between the client and server
	Drift time.Duration `json:"drift,omitempty"`
	// DriftChecks is the maximum number of time time counter will be adjusted for a given client
	DriftChecks int `json:"driftChecks,omitempty"`
}

// checkDrift will back track the time counter twice, by 30 seconds.
// The OTP will be verified for each backtracking. This is done taking into account the possible
// clock drift between server and client.
func (t *TOTP) checkDrift() {

}

// Secret generates a new secret for TOTP for the given identifiers
func Secret(parts ...string) string {
	return strings.Join(parts, "")
}

func (t *TOTP) label(userID string) string {
	return fmt.Sprintf("%s:%s", url.QueryEscape(t.Issuer), userID)
}

// check checks if the given otp is valid
func (t *TOTP) check(otp string) bool {
	return false
}

// Check checks if the given otp is valid or not. It will also perform a drift check.
// The drift would only be allowed 'DriftChecks' times, post which the client has fix its clock
func (t *TOTP) Check(secret, otp string) bool {
	return false
}

// URI generates the URI representing all the required details to be consumed by authenticator
// apps, while scanning the QR code
func (t *TOTP) URI(userID string) string {
	params := url.Values{}
	params.Add("issuer", t.Issuer)
	params.Add("secret", "none")
	params.Add("algorithm", string(t.Algorithm))
	params.Add("digits", fmt.Sprintf("%d", t.Digits))
	params.Add("period", fmt.Sprintf("%d", t.Period))
	return fmt.Sprintf("otpauth://totp/%s?%s", t.label(userID), params.Encode())
}

// QR returns the QR code with the required URI payload
func (t *TOTP) QR() {

}

func New(issuer string, digits, period int, alg algo) *TOTP {

	t := &TOTP{
		Issuer:    issuer,
		Algorithm: alg,
		Digits:    digits,
		Period:    period,
	}

	return t
}
