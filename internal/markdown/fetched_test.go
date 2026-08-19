package markdown

import (
	"context"
	"strings"
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

func TestBodyToDocument_ListItems(t *testing.T) {
	lists := map[string]docs.List{
		"ul": {ListProperties: &docs.ListProperties{NestingLevels: []*docs.NestingLevel{
			{GlyphType: "DISC"},
			{GlyphType: "CIRCLE"},
		}}},
		"ol": {ListProperties: &docs.ListProperties{NestingLevels: []*docs.NestingLevel{
			{GlyphType: "DECIMAL"},
		}}},
	}
	body := &docs.Body{Content: []*docs.StructuralElement{
		{Paragraph: &docs.Paragraph{
			Bullet:         &docs.Bullet{ListId: "ul", NestingLevel: 0},
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "NORMAL_TEXT"},
			Elements:       []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "one\n"}}},
		}},
		{Paragraph: &docs.Paragraph{
			Bullet:         &docs.Bullet{ListId: "ul", NestingLevel: 1},
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "NORMAL_TEXT"},
			Elements:       []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "two\n"}}},
		}},
		{Paragraph: &docs.Paragraph{
			Bullet:         &docs.Bullet{ListId: "ol", NestingLevel: 0},
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "NORMAL_TEXT"},
			Elements: []*docs.ParagraphElement{
				{TextRun: &docs.TextRun{Content: "first "}},
				{TextRun: &docs.TextRun{Content: "bold", TextStyle: &docs.TextStyle{Bold: true}}},
			},
		}},
		{Paragraph: &docs.Paragraph{
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "NORMAL_TEXT"},
			Elements:       []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "plain\n"}}},
		}},
		{Paragraph: &docs.Paragraph{
			Bullet:         &docs.Bullet{ListId: "ul", NestingLevel: 0},
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "HEADING_1"},
			Elements:       []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Heading bullet\n"}}},
		}},
	}}
	doc := bodyToDocument("T", "doc-1", "tab-1", body, lists)
	if len(doc.Body) != 5 {
		t.Fatalf("body len=%d want 5: %+v", len(doc.Body), doc.Body)
	}
	li0, ok := doc.Body[0].(ListItem)
	if !ok || li0.Ordered || li0.Level != 0 || li0.Text != "one" {
		t.Fatalf("item0=%#v", doc.Body[0])
	}
	li1, ok := doc.Body[1].(ListItem)
	if !ok || li1.Ordered || li1.Level != 1 || li1.Text != "two" {
		t.Fatalf("item1=%#v", doc.Body[1])
	}
	li2, ok := doc.Body[2].(ListItem)
	if !ok || !li2.Ordered || li2.Level != 0 || li2.Text != "first **bold**" {
		t.Fatalf("item2=%#v", doc.Body[2])
	}
	p, ok := doc.Body[3].(Paragraph)
	if !ok || p.Text != "plain" {
		t.Fatalf("plain=%#v", doc.Body[3])
	}
	h, ok := doc.Body[4].(Heading)
	if !ok || h.Level != 1 || h.Text != "Heading bullet" {
		t.Fatalf("heading=%#v", doc.Body[4])
	}

	md, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(md, "- one\n") || !strings.Contains(md, "  - two\n") || !strings.Contains(md, "1. first **bold**\n") {
		t.Fatalf("markdown=%q", md)
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
	doc := bodyToDocument("T", "doc-1", "tab-1", body, nil)
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
