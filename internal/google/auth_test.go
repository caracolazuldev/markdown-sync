package google

import (
    "context"
    "os"
    "path/filepath"
    "testing"
)

func TestNewHTTPClient_Service_InvalidPath(t *testing.T) {
    ctx := context.Background()
    // non-existent file should produce an error
    _, err := NewHTTPClient(ctx, "service", filepath.Join(os.TempDir(), "nonexistent-sa.json"))
    if err == nil {
        t.Fatalf("expected error for nonexistent credentials path")
    }
}

func TestNewHTTPClient_Service_BadJSON(t *testing.T) {
    ctx := context.Background()
    // create a temp file with invalid JSON
    f := filepath.Join(os.TempDir(), "bad-sa.json")
    os.WriteFile(f, []byte("notjson"), 0o600)
    defer os.Remove(f)
    _, err := NewHTTPClient(ctx, "service", f)
    if err == nil {
        t.Fatalf("expected error for invalid JSON service account file")
    }
}
