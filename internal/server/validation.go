package server

import (
	"strings"

	"github.com/plainq/plainq/internal/shared/pqerr"
	"github.com/plainq/servekit/idkit"
)

// validateQueueIDFromRequest performs validation of the queue identifier.
func validateQueueIDFromRequest(r interface{ GetQueueId() string }) error {
	if r == nil {
		return pqerr.ErrInvalidID
	}

	return validateQueueID(r.GetQueueId())
}

// validateQueueID validates given queue identifier.
func validateQueueID(queueID string) error {
	if queueID == "" {
		return pqerr.ErrInvalidID
	}

	if err := idkit.ValidateXID(strings.ToLower(queueID)); err != nil {
		return pqerr.ErrInvalidID
	}

	return nil
}
