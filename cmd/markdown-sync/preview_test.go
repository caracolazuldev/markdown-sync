package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPreviewDefaultLines(t *testing.T) {
	var out bytes.Buffer
	var errb bytes.Buffer
	// default maxLines in main is 20 but we call helper directly
	if err := previewToWriter("oauth", "doc-1", 2, &out, &errb); err != nil {
		t.Fatalf("preview failed: %v, stderr=%s", err, errb.String())
	}
	got := out.String()
	if !strings.Contains(got, "Introduction") {
		t.Fatalf("expected preview to include heading, got: %q", got)
	}
	if !strings.Contains(got, "(truncated)") {
		t.Fatalf("expected truncated marker, got: %q", got)
	}
}

func TestPreviewFull(t *testing.T) {
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := previewToWriter("oauth", "doc-1", 0, &out, &errb); err != nil {
		t.Fatalf("preview failed: %v, stderr=%s", err, errb.String())
	}
	if !strings.Contains(out.String(), "Sample Doc") {
		t.Fatalf("expected full preview to contain title, got: %q", out.String())
	}
}
