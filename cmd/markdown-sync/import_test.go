package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportDiffShowsDifferences(t *testing.T) {
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "local.md")
	// create a local file that differs from the adapter stub
	localContent := "Sample Doc\n\nThis is a local extra line.\n"
	if err := os.WriteFile(outPath, []byte(localContent), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	var out bytes.Buffer
	var errb bytes.Buffer
	if err := importToWriter("oauth", outPath, "doc-1", true, &out, &errb); err != nil {
		t.Fatalf("importToWriter failed: %v, stderr=%s", err, errb.String())
	}
	got := out.String()
	if !strings.Contains(got, "+++ local") {
		t.Fatalf("expected diff header, got: %q", got)
	}
	if !strings.Contains(got, "+ This is a local extra line.") {
		t.Fatalf("expected added local line in diff, got: %q", got)
	}
}
