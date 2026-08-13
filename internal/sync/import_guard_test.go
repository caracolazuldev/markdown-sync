package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckImportAllowed_UntrackedOK(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "a.md")
	if err := os.WriteFile(p, []byte("# Hi\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := CheckImportAllowed(p, "doc-1"); err != nil {
		t.Fatal(err)
	}
}

func TestCheckImportAllowed_TrackFrontMatter(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "a.md")
	if err := os.WriteFile(p, []byte("---\ntitle: A\ntrack: true\n---\n\nHi\n"), 0644); err != nil {
		t.Fatal(err)
	}
	err := CheckImportAllowed(p, "doc-1")
	if !errors.Is(err, ErrTrackedPath) {
		t.Fatalf("err=%v", err)
	}
}

func TestCheckImportAllowed_TrackRoot(t *testing.T) {
	dir := t.TempDir()
	if err := writeSidecar(dir, Sidecar{Mode: "track", DocID: "doc-1", Title: "T"}); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "nested", "a.md")
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("Hi\n"), 0644); err != nil {
		t.Fatal(err)
	}
	err := CheckImportAllowed(p, "doc-1")
	if !errors.Is(err, ErrTrackedPath) {
		t.Fatalf("err=%v", err)
	}
}
