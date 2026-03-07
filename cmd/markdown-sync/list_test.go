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
		t.Fatalf("listToWriter failed: %v, stderr=%s", err, errb.String())
	}
	got := out.String()
	if !strings.Contains(got, "doc-1") || !strings.Contains(got, "Sample Doc") {
		t.Fatalf("unexpected list output: %q", got)
	}
}
