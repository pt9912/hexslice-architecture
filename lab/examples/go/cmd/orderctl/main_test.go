package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"hexslice/example/internal/hexagon/domain/order"
)

// The end-to-end test drives the real thing: the composition root's own wiring,
// the real adapters, the real slices and the real CLI. Only the process boundary —
// main's exit code — is out of scope.
//
// It lives in package main, and that is forced rather than chosen: a test in the
// CLI adapter would have to import a driven adapter, which is a lateral adapter
// edge the architecture forbids. The composition root is the only place that sees
// both sides.

// lines returns the non-empty output lines, so the assertions below are about
// content and order without depending on trailing newlines.
func lines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}
	return out
}

func TestDemoRunsTheWholeLifecycle(t *testing.T) {
	var out bytes.Buffer
	app := newCLI(&out, &bytes.Buffer{})

	if err := app.Run(context.Background(), []string{"demo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		"== HexSlice order demo ==",
		"created order ORD-000001 | total 47.48 EUR | status CREATED",
		"notification: order ORD-000001 cancelled (reason: customer changed their mind)",
		"cancelled order ORD-000001 | status CANCELLED",
		"second cancel rejected by domain: order: order is already cancelled",
	}
	got := lines(out.String())
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d:\n%s", len(want), len(got), strings.Join(got, "\n"))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d:\n got: %s\nwant: %s", i+1, got[i], want[i])
		}
	}
}

func TestCreateParsesTheFlagSurfaceAndPersists(t *testing.T) {
	var out bytes.Buffer
	app := newCLI(&out, &bytes.Buffer{})

	err := app.Run(context.Background(), []string{
		"create",
		"-customer", "cust-1",
		"-line", "BOOK-1:2:1999:EUR",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"created order ORD-000001 | total 39.98 EUR | status CREATED"}
	got := lines(out.String())
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCancelSurfacesTheDomainError(t *testing.T) {
	app := newCLI(&bytes.Buffer{}, &bytes.Buffer{})

	err := app.Run(context.Background(), []string{"cancel", "-id", "does-not-exist"})
	if !errors.Is(err, order.ErrNotFound) {
		t.Fatalf("expected ErrNotFound from the domain, got %v", err)
	}
}

func TestMalformedFlagIsReportedOnStderr(t *testing.T) {
	var errOut bytes.Buffer
	app := newCLI(&bytes.Buffer{}, &errOut)

	err := app.Run(context.Background(), []string{"create", "-customer"})
	if err == nil {
		t.Fatal("expected an error for a missing flag argument")
	}
	// The usage text is the flag package's ("Usage of create: …"), so the
	// assertion is on what this example owns: the error reaches stderr and names
	// the flag.
	if errOut.Len() == 0 {
		t.Fatal("expected the flag error to be written to stderr")
	}
	if !strings.Contains(err.Error(), "-customer") {
		t.Fatalf("expected the error to name the flag, got %v", err)
	}
}
