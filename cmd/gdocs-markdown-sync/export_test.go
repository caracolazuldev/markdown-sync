package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
	gsync "github.com/caracolazuldev/gdocs-markdown-sync/internal/sync"
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

func TestExportTabbedFailsAndDoesNotWrite(t *testing.T) {
	old := activeFetcher
	t.Cleanup(func() { activeFetcher = old })
	activeFetcher = md.FakeFetcher{Doc: &md.FetchedDocument{
		ID:    "doc-tabbed",
		Title: "Tabbed",
		Tabs: []*md.Tab{
			{ID: "a", Title: "A", Document: &md.Document{Title: "A", Body: []md.Element{md.Paragraph{Text: "one"}}}},
			{ID: "b", Title: "B", Document: &md.Document{Title: "B", Body: []md.Element{md.Paragraph{Text: "two"}}}},
		},
	}}
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "out.md")
	var out bytes.Buffer
	var errb bytes.Buffer
	err := exportToWriter("oauth", "doc-tabbed", outPath, false, &out, &errb)
	if !errors.Is(err, gsync.ErrTabbedDocument) {
		t.Fatalf("err=%v", err)
	}
	if !strings.Contains(err.Error(), "track") || !strings.Contains(errb.String(), "track") {
		t.Fatalf("hint missing: err=%v stderr=%s", err, errb.String())
	}
	if _, statErr := os.Stat(outPath); !os.IsNotExist(statErr) {
		t.Fatal("export must not write first-tab-only markdown")
	}
}
