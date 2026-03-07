package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestImportApplyLocalFile(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "apply.md")
	if err := ioutil.WriteFile(file, []byte("# Apply\n\nContent"), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}
	var out bytes.Buffer
	var errb bytes.Buffer
	// Use diff-only mode in unit tests to avoid requiring Google credentials.
	if err := importToWriter("oauth", file, "doc-1", true, &out, &errb); err != nil {
		t.Fatalf("import apply (diff-only) failed: %v, stderr=%s", err, errb.String())
	}
}
