// Package notify is an outbound adapter implementing the cancel-order Notifier
// port by writing human-readable messages to an io.Writer.
package notify

import (
	"context"
	"fmt"
	"io"

	"hexslice/example/internal/hexagon/domain/order"
)

// Writer implements the Notifier port using an io.Writer sink (e.g. stdout).
type Writer struct {
	out io.Writer
}

// NewWriter creates a Writer notifier over the given sink.
func NewWriter(out io.Writer) *Writer { return &Writer{out: out} }

// OrderCancelled writes a cancellation notice.
func (w *Writer) OrderCancelled(_ context.Context, id order.ID, reason string) error {
	_, err := fmt.Fprintf(w.out, "notification: order %s cancelled (reason: %s)\n", id, reason)
	return err
}
