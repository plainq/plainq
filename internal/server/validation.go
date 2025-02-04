package server

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/plainq/servekit/errkit"
	"github.com/plainq/servekit/idkit"
)

const (
	minNameLen, maxNameLen = 2, 50
	minPassLen, maxPassLen = 6, 50
)

// validateQueueIDFromRequest performs validation of the queue identifier.
func validateQueueIDFromRequest(r interface{ GetQueueId() string }) error {
	if r == nil {
		return errkit.ErrInvalidID
	}

	return validateQueueID(r.GetQueueId())
}

// validateQueueID validates given queue identifier.
func validateQueueID(queueID string) error {
	if queueID == "" {
		return errkit.ErrInvalidID
	}

	if err := idkit.ValidateXID(strings.ToLower(queueID)); err != nil {
		return errkit.ErrInvalidID
	}

	return nil
}

// validateUserName validates given name.
func validateUserName(name string) error {
	if err := charLenMinMax(name, minNameLen, maxNameLen); err != nil {
		return errkit.ErrValidation
	}

	if strings.ContainsAny(name, "/\"'?,:;<>!#$%^&*()={}[]|") {
		return errkit.ErrValidation
	}

	return nil
}

// validatePassword validates given password.
func validatePassword(password string) error {
	if err := charLenMinMax(password, minPassLen, maxPassLen); err != nil {
		return errkit.ErrValidation
	}

	return nil
}

// validateEmail validates given email.
func validateEmail(email string) error {
	// If contains special characters then invalid.
	if strings.ContainsAny(email, "/\"'?,:;<>!#$%^&*()={}[]|") || !strings.Contains(email, "@") {
		return errkit.ErrValidation
	}

	if tail := email[len(email)-1]; strings.Contains(string(tail), ".") {
		return errkit.ErrValidation
	}

	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return errkit.ErrValidation
	}

	domparts := strings.Split(parts[1], ".")
	if len(domparts) < 2 || domparts[0] == "" || domparts[1] == "" {
		return errkit.ErrValidation
	}

	return nil
}

// charLenMinMax returns non nil error if given password is too short or to long.
func charLenMinMax(value string, min, max int) error {
	valueLen := utf8.RuneCountInString(value)
	if valueLen < min {
		return fmt.Errorf("%w: value is to short", errkit.ErrValidation)
	}

	if valueLen > max {
		return fmt.Errorf("%w: value is to long", errkit.ErrValidation)
	}

	return nil
}
