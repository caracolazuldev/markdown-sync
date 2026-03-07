package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPreviewToWriter(t *testing.T) {
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := previewToWriter("oauth", "doc-1", 5, &out, &errb); err != nil {
		t.Fatalf("preview failed: %v, stderr=%s", err, errb.String())
	}
	if !strings.Contains(out.String(), "Introduction") {
		t.Fatalf("unexpected preview output: %q", out.String())
	}
}
