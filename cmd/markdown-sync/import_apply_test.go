package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestImportApplySucceeds(t *testing.T) {
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "local.md")
	// create a local file that will be parsed and applied
	localContent := "# Imported\n\nSome content to import.\n"
	if err := os.WriteFile(outPath, []byte(localContent), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	var errb bytes.Buffer
	var out bytes.Buffer
	if err := importToWriter("oauth", outPath, "doc-1", false, &out, &errb); err != nil {
		t.Fatalf("importToWriter apply failed: %v, stderr=%s", err, errb.String())
	}
}
