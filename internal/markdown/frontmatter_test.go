package markdown

import (
	"strings"
	"testing"
)

func TestDocumentToMarkdown_TrackFrontMatter(t *testing.T) {
	doc := &Document{
		Title: "Minutes",
		DocID: "doc-1",
		TabID: "tab-9",
		Track: true,
		Body:  []Element{Paragraph{Text: "Hello"}},
	}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"title: Minutes", "doc_id: doc-1", "tab_id: tab-9", "track: true", "Hello"} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in %q", want, got)
		}
	}
}

func TestSplitFrontMatter_Track(t *testing.T) {
	in := "---\ntitle: A\ndoc_id: d\ntab_id: t\ntrack: true\n---\n\nBody\n"
	body, fm := SplitFrontMatter(in)
	if strings.TrimSpace(body) != "Body" {
		t.Fatalf("body=%q", body)
	}
	if fm.Title != "A" || fm.DocID != "d" || fm.TabID != "t" || !fm.Track {
		t.Fatalf("fm=%+v", fm)
	}
}
