package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
)

func TestTrackToWriter(t *testing.T) {
	old := activeFetcher
	t.Cleanup(func() { activeFetcher = old })
	activeFetcher = md.FakeFetcher{Doc: &md.FetchedDocument{
		ID:    "doc-1",
		Title: "Minutes",
		Tabs: []*md.Tab{
			{ID: "w", Title: "Welcome", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "hi"}}}},
			{ID: "x", Title: "Extra", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "x"}}}},
		},
	}}
	dir := t.TempDir()
	var out, errb bytes.Buffer
	if err := trackToWriter("oauth", "doc-1", dir, false, false, &out, &errb); err != nil {
		t.Fatalf("track: %v stderr=%s", err, errb.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "welcome.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "_track.toml")); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "wrote") {
		t.Fatalf("stdout=%q", out.String())
	}
}
