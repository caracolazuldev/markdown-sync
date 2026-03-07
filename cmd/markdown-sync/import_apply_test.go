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

	// Require service credentials for apply; skip otherwise to avoid
	// attempting network calls in unit test environments.
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("skipping apply test; set GOOGLE_APPLICATION_CREDENTIALS to enable")
	}
	var errb bytes.Buffer
	var out bytes.Buffer
	if err := importToWriter("service", outPath, "doc-1", false, &out, &errb); err != nil {
		t.Fatalf("importToWriter apply failed: %v, stderr=%s", err, errb.String())
	}
}
