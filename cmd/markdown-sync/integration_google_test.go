//go:build integration
// +build integration

package main

import (
	"bytes"
	"os"
	"testing"
)

// This integration test updates a real Google Doc. It only runs when the
// environment variables `GOOGLE_APPLICATION_CREDENTIALS` and
// `INTEGRATION_DOC_ID` are set and when `-tags=integration` is provided to
// `go test`.
func TestApplyDocumentIntegration(t *testing.T) {
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	docID := os.Getenv("INTEGRATION_DOC_ID")
	if creds == "" || docID == "" {
		t.Skip("integration test skipped; set GOOGLE_APPLICATION_CREDENTIALS and INTEGRATION_DOC_ID to enable")
	}

	// build a simple local file to apply
	tmp := t.TempDir()
	file := tmp + "/apply.md"
	if err := os.WriteFile(file, []byte("# Integration Test\n\nApplied by test."), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	var out bytes.Buffer
	var errb bytes.Buffer
	if err := importToWriter("service", file, docID, false, &out, &errb); err != nil {
		t.Fatalf("importToWriter apply failed: %v, stderr=%s", err, errb.String())
	}
}
