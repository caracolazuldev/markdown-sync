package markdown

import (
	"testing"

	docs "google.golang.org/api/docs/v1"
)

func TestBuildNativeTableFillRequests_OrderAndContent(t *testing.T) {
	model := []Table{{
		Header: []string{"H1", "H2"},
		Rows:   [][]string{{"A", "B"}},
	}}

	remote := &docs.Document{
		Body: &docs.Body{Content: []*docs.StructuralElement{
			{Table: &docs.Table{TableRows: []*docs.TableRow{
				{TableCells: []*docs.TableCell{{StartIndex: 10}, {StartIndex: 20}}},
				{TableCells: []*docs.TableCell{{StartIndex: 30}, {StartIndex: 40}}},
			}}},
		}},
	}

	reqs := buildNativeTableFillRequests(model, remote)
	if len(reqs) < 4 {
		t.Fatalf("expected at least 4 fill requests, got %d", len(reqs))
	}

	// descending index order expected: 41,31,21,11
	want := []int64{41, 31, 21, 11}
	var gotIdx []int64
	for _, r := range reqs {
		if r.InsertText != nil && r.InsertText.Location != nil {
			gotIdx = append(gotIdx, r.InsertText.Location.Index)
		}
	}
	if len(gotIdx) != 4 {
		t.Fatalf("expected 4 InsertText requests, got %d", len(gotIdx))
	}
	for i, w := range want {
		if gotIdx[i] != w {
			t.Fatalf("insert request %d index mismatch: got %d want %d", i, gotIdx[i], w)
		}
	}
}

func TestBuildNativeTableFillRequests_InlineStyles(t *testing.T) {
	model := []Table{{
		Header: []string{"**Bold**", "[Link](https://example.com)"},
	}}

	remote := &docs.Document{
		Body: &docs.Body{Content: []*docs.StructuralElement{
			{Table: &docs.Table{TableRows: []*docs.TableRow{
				{TableCells: []*docs.TableCell{{StartIndex: 10}, {StartIndex: 20}}},
			}}},
		}},
	}

	reqs := buildNativeTableFillRequests(model, remote)
	var hasBold, hasLink bool
	for _, r := range reqs {
		if r.UpdateTextStyle != nil && r.UpdateTextStyle.TextStyle != nil {
			if r.UpdateTextStyle.TextStyle.Bold {
				hasBold = true
			}
			if r.UpdateTextStyle.TextStyle.Link != nil && r.UpdateTextStyle.TextStyle.Link.Url == "https://example.com" {
				hasLink = true
			}
		}
	}
	if !hasBold || !hasLink {
		t.Fatalf("expected bold and link text style updates, got hasBold=%v hasLink=%v", hasBold, hasLink)
	}
}

func TestBuildNativeTableFillRequests_UTF16BoldRangeAfterArrow(t *testing.T) {
	model := []Table{{
		Header: []string{"→ **X**"},
	}}
	remote := &docs.Document{
		Body: &docs.Body{Content: []*docs.StructuralElement{
			{Table: &docs.Table{TableRows: []*docs.TableRow{
				{TableCells: []*docs.TableCell{{StartIndex: 10}}},
			}}},
		}},
	}
	reqs := buildNativeTableFillRequests(model, remote)
	start, end, ok := firstBoldRange(reqs)
	if !ok {
		t.Fatal("missing bold fill style")
	}
	wantStart := int64(11) + utf16Len("→ ")
	if start != wantStart || end != wantStart+1 {
		t.Fatalf("bold range [%d,%d) want [%d,%d)", start, end, wantStart, wantStart+1)
	}
}

func TestNativeTablesFromDocument(t *testing.T) {
	doc := &Document{Body: []Element{
		Paragraph{Text: "x"},
		Table{Header: []string{"A"}, Rows: [][]string{{"1"}}},
		Table{Header: []string{"A", "B"}, Rows: [][]string{{"1"}}}, // irregular, should skip
	}}
	got := nativeTablesFromDocument(doc)
	if len(got) != 1 {
		t.Fatalf("expected 1 native table, got %d", len(got))
	}
}
