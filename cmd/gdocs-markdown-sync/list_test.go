package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestListToWriter(t *testing.T) {
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := listToWriter("oauth", &out, &errb); err != nil {
		t.Fatalf("list failed: %v, stderr=%s", err, errb.String())
	}
	if !strings.Contains(out.String(), "doc-1") && !strings.Contains(out.String(), "Sample Doc") {
		t.Fatalf("unexpected list output: %q", out.String())
	}
}
