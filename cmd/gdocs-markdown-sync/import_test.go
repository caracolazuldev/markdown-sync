package main

import (
	"bytes"
	"errors"
	"os"
	"testing"

	gsync "github.com/caracolazuldev/gdocs-markdown-sync/internal/sync"
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

func TestImportRefusesTrackedFile(t *testing.T) {
	tmp := t.TempDir()
	file := tmp + "/local.md"
	if err := os.WriteFile(file, []byte("---\ntitle: X\ntrack: true\n---\n\nHi\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	err := importToWriter("oauth", file, "doc-1", true, &out, &errb)
	if !errors.Is(err, gsync.ErrTrackedPath) {
		t.Fatalf("err=%v", err)
	}
}
