package markdown

import (
	"context"
	"testing"

	docs "google.golang.org/api/docs/v1"
)

func TestFetchedDocument_IsTabbed(t *testing.T) {
	legacy := &FetchedDocument{Document: sampleDocument()}
	if legacy.IsTabbed() {
		t.Fatal("legacy body should not be tabbed")
	}
	one := &FetchedDocument{Tabs: []*Tab{{ID: "a", Title: "A", Document: sampleDocument()}}}
	if one.IsTabbed() {
		t.Fatal("single leaf tab should not be tabbed")
	}
	two := &FetchedDocument{Tabs: []*Tab{
		{ID: "a", Title: "A"},
		{ID: "b", Title: "B"},
	}}
	if !two.IsTabbed() {
		t.Fatal("two top-level tabs should be tabbed")
	}
	nested := &FetchedDocument{Tabs: []*Tab{
		{ID: "a", Title: "A", Children: []*Tab{{ID: "a1", Title: "A1"}}},
	}}
	if !nested.IsTabbed() {
		t.Fatal("nested childTabs should be tabbed")
	}
	if nested.TabCount() != 2 {
		t.Fatalf("TabCount=%d want 2", nested.TabCount())
	}
}

func TestBodyToDocument_ParagraphAndHeading(t *testing.T) {
	body := &docs.Body{Content: []*docs.StructuralElement{
		{Paragraph: &docs.Paragraph{
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "HEADING_1"},
			Elements:       []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Hello\n"}}},
		}},
		{Paragraph: &docs.Paragraph{
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "NORMAL_TEXT"},
			Elements:       []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "World"}}},
		}},
	}}
	doc := bodyToDocument("T", "doc-1", "tab-1", body)
	if len(doc.Body) != 2 {
		t.Fatalf("body len=%d want 2: %+v", len(doc.Body), doc.Body)
	}
	h, ok := doc.Body[0].(Heading)
	if !ok || h.Level != 1 || h.Text != "Hello" {
		t.Fatalf("heading=%v ok=%v", doc.Body[0], ok)
	}
	p, ok := doc.Body[1].(Paragraph)
	if !ok || p.Text != "World" {
		t.Fatalf("paragraph=%v ok=%v", doc.Body[1], ok)
	}
}

func TestFetchedFromAPI_LegacyAndTabs(t *testing.T) {
	legacy := fetchedFromAPI("id", &docs.Document{Title: "L", Body: &docs.Body{}})
	if legacy.IsTabbed() || legacy.Title != "L" {
		t.Fatalf("legacy: %+v", legacy)
	}
	api := fetchedFromAPI("id", &docs.Document{
		Title: "Multi",
		Tabs: []*docs.Tab{
			{TabProperties: &docs.TabProperties{TabId: "t1", Title: "One"}, DocumentTab: &docs.DocumentTab{Body: &docs.Body{}}},
			{TabProperties: &docs.TabProperties{TabId: "t2", Title: "Two"}, DocumentTab: &docs.DocumentTab{Body: &docs.Body{}}},
		},
	})
	if !api.IsTabbed() || api.TabCount() != 2 {
		t.Fatalf("tabbed: count=%d isTabbed=%v", api.TabCount(), api.IsTabbed())
	}
	if api.Tabs[0].ID != "t1" || api.Tabs[1].Title != "Two" {
		t.Fatalf("tabs=%+v", api.Tabs)
	}
}

func TestSampleFetcher_SingleTab(t *testing.T) {
	fd, err := SampleFetcher{}.Fetch(context.Background(), "oauth", "doc-1")
	if err != nil {
		t.Fatal(err)
	}
	if fd.IsTabbed() {
		t.Fatal("sample should not be tabbed")
	}
	if fd.Document == nil || fd.Document.Title != "Sample Doc" {
		t.Fatalf("sample doc: %+v", fd.Document)
	}
}
