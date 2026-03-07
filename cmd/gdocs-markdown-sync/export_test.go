package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportToStdout(t *testing.T) {
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := exportToWriter("oauth", "doc-1", "", false, &out, &errb); err != nil {
		t.Fatalf("export failed: %v, stderr=%s", err, errb.String())
	}
	got := out.String()
	if !strings.Contains(got, "Sample Doc") && !strings.Contains(got, "Introduction") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestExportWriteFile(t *testing.T) {
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "out.md")
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := exportToWriter("oauth", "doc-1", outPath, false, &out, &errb); err != nil {
		t.Fatalf("export failed: %v, stderr=%s", err, errb.String())
	}
	// verify file written
	data, err := ioutil.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "Sample Doc") {
		t.Fatalf("unexpected file contents: %q", string(data))
	}
}

func TestExportDryRunDoesNotWrite(t *testing.T) {
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "out.md")
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := exportToWriter("oauth", "doc-1", outPath, true, &out, &errb); err != nil {
		t.Fatalf("export dry-run failed: %v, stderr=%s", err, errb.String())
	}
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Fatalf("expected no file written in dry-run, got stat err=%v", err)
	}
	if !strings.Contains(out.String(), "dry-run: would write") {
		t.Fatalf("expected dry-run message, got: %q", out.String())
	}
}
