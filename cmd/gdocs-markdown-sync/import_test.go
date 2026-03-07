package main

import (
	"bytes"
	"os"
	"testing"
)

func TestImportDiffOnly(t *testing.T) {
	tmp := t.TempDir()
	file := tmp + "/local.md"
	if err := os.WriteFile(file, []byte("# Test\n\nHello world"), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := importToWriter("oauth", file, "doc-1", true, &out, &errb); err != nil {
		t.Fatalf("import diff failed: %v, stderr=%s", err, errb.String())
	}
	if out.Len() == 0 {
		t.Fatalf("expected diff output")
	}
}
