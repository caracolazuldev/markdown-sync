package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
	gsync "github.com/caracolazuldev/gdocs-markdown-sync/internal/sync"
)

func TestPreviewToWriter(t *testing.T) {
	var out bytes.Buffer
	var errb bytes.Buffer
	if err := previewToWriter("oauth", "doc-1", 5, &out, &errb); err != nil {
		t.Fatalf("preview failed: %v, stderr=%s", err, errb.String())
	}
	if !strings.Contains(out.String(), "Introduction") {
		t.Fatalf("unexpected preview output: %q", out.String())
	}
}

func TestPreviewTabbedFails(t *testing.T) {
	old := activeFetcher
	t.Cleanup(func() { activeFetcher = old })
	activeFetcher = md.FakeFetcher{Doc: &md.FetchedDocument{
		ID: "doc-tabbed",
		Tabs: []*md.Tab{
			{ID: "a", Title: "A"},
			{ID: "b", Title: "B"},
		},
	}}
	var out, errb bytes.Buffer
	err := previewToWriter("oauth", "doc-tabbed", 5, &out, &errb)
	if !errors.Is(err, gsync.ErrTabbedDocument) {
		t.Fatalf("err=%v", err)
	}
	if out.Len() != 0 {
		t.Fatalf("preview wrote %q", out.String())
	}
}
