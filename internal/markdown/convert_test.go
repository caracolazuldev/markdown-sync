package markdown

import (
	"strings"
	"testing"
)

func TestDocumentToMarkdown_TitleAndParagraph(t *testing.T) {
	doc := &Document{
		Title: "My Title",
		Body:  []Element{Paragraph{Text: "Hello world"}},
	}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "title: My Title") {
		t.Fatalf("missing title in output: %q", got)
	}
	if !strings.Contains(got, "Hello world") {
		t.Fatalf("missing paragraph text: %q", got)
	}
}

func TestDocumentToMarkdown_HeadingAndLevels(t *testing.T) {
	doc := &Document{
		Body: []Element{
			Heading{Level: 2, Text: "Section"},
			Heading{Level: 0, Text: "Fallback"},
		},
	}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "## Section") {
		t.Fatalf("expected H2 heading, got: %q", got)
	}
	if !strings.Contains(got, "# Fallback") {
		t.Fatalf("expected fallback H1 for invalid level, got: %q", got)
	}
}

func TestDocumentToMarkdown_CodeBlockAndImage(t *testing.T) {
	doc := &Document{
		Body: []Element{
			CodeBlock{Language: "go", Code: "fmt.Println(\"hi\")"},
			Image{URL: "https://example.com/img.png", Alt: "An example"},
		},
	}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "```go") || !strings.Contains(got, "fmt.Println(\"hi\")") {
		t.Fatalf("code block not rendered correctly: %q", got)
	}
	if !strings.Contains(got, "![An example](https://example.com/img.png)") {
		t.Fatalf("image not rendered correctly: %q", got)
	}
}

func TestToMarkdown_WrapperErrors(t *testing.T) {
	if _, err := ToMarkdown(nil); err == nil {
		t.Fatalf("expected error for nil input")
	}
	if _, err := ToMarkdown(struct{}{}); err == nil {
		t.Fatalf("expected error for unsupported type")
	}
}

func TestFromMarkdown_GoldmarkBasicMapping(t *testing.T) {
	input := "---\n" +
		"title: Parsed Title\n" +
		"---\n\n" +
		"# Heading One\n\n" +
		"Paragraph with **bold** and `code`.\n\n" +
		"```go\nfmt.Println(\"hi\")\n```\n\n" +
		"![Alt text](https://example.com/a.png)\n"

	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc, ok := v.(*Document)
	if !ok {
		t.Fatalf("expected *Document, got %T", v)
	}
	if doc.Title != "Parsed Title" {
		t.Fatalf("unexpected title: %q", doc.Title)
	}
	if len(doc.Body) < 4 {
		t.Fatalf("expected at least 4 body elements, got %d", len(doc.Body))
	}

	h, ok := doc.Body[0].(Heading)
	if !ok || h.Level != 1 || h.Text != "Heading One" {
		t.Fatalf("unexpected heading: %#v", doc.Body[0])
	}

	p, ok := doc.Body[1].(Paragraph)
	if !ok || p.Text != "Paragraph with **bold** and `code`." {
		t.Fatalf("unexpected paragraph: %#v", doc.Body[1])
	}

	cb, ok := doc.Body[2].(CodeBlock)
	if !ok || cb.Language != "go" || !strings.Contains(cb.Code, "fmt.Println") {
		t.Fatalf("unexpected codeblock: %#v", doc.Body[2])
	}

	img, ok := doc.Body[3].(Image)
	if !ok || img.URL != "https://example.com/a.png" {
		t.Fatalf("unexpected image: %#v", doc.Body[3])
	}
}

func TestFromMarkdown_NoFrontMatter(t *testing.T) {
	v, err := FromMarkdown("Simple paragraph\n")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	if doc.Title != "" {
		t.Fatalf("expected empty title, got %q", doc.Title)
	}
	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 body element, got %d", len(doc.Body))
	}
}

func TestFromMarkdown_ListItems(t *testing.T) {
	input := "- one\n- two\n\n1. first\n2. second\n"
	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	if len(doc.Body) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(doc.Body))
	}

	li0, ok := doc.Body[0].(ListItem)
	if !ok || li0.Ordered || li0.Level != 0 || li0.Text != "one" {
		t.Fatalf("unexpected first list item: %#v", doc.Body[0])
	}
	li2, ok := doc.Body[2].(ListItem)
	if !ok || !li2.Ordered || li2.Level != 0 || li2.Text != "first" {
		t.Fatalf("unexpected third list item: %#v", doc.Body[2])
	}
}

func TestFromMarkdown_NestedListItems(t *testing.T) {
	input := "- parent\n  - child\n    - grandchild\n"
	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	if len(doc.Body) != 3 {
		t.Fatalf("expected 3 list items, got %d", len(doc.Body))
	}
	li0 := doc.Body[0].(ListItem)
	li1 := doc.Body[1].(ListItem)
	li2 := doc.Body[2].(ListItem)
	if li0.Level != 0 || li1.Level != 1 || li2.Level != 2 {
		t.Fatalf("unexpected levels: %#v %#v %#v", li0, li1, li2)
	}
}

func TestDocumentToMarkdown_ListItems(t *testing.T) {
	doc := &Document{Body: []Element{
		ListItem{Ordered: false, Text: "alpha"},
		ListItem{Ordered: true, Level: 1, Text: "beta"},
	}}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "- alpha") {
		t.Fatalf("missing bullet item: %q", got)
	}
	if !strings.Contains(got, "1. beta") {
		t.Fatalf("missing ordered item: %q", got)
	}
	if !strings.Contains(got, "  1. beta") {
		t.Fatalf("missing nested indentation for ordered item: %q", got)
	}
}

func TestFromMarkdown_QuoteAndRule(t *testing.T) {
	input := "> quoted line\n\n---\n"
	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	if len(doc.Body) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(doc.Body))
	}
	q, ok := doc.Body[0].(Quote)
	if !ok || q.Text != "quoted line" {
		t.Fatalf("unexpected quote: %#v", doc.Body[0])
	}
	if _, ok := doc.Body[1].(HorizontalRule); !ok {
		t.Fatalf("expected HorizontalRule, got %#v", doc.Body[1])
	}
}

func TestDocumentToMarkdown_QuoteAndRule(t *testing.T) {
	doc := &Document{Body: []Element{
		Quote{Text: "quoted line"},
		HorizontalRule{},
	}}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "> quoted line") {
		t.Fatalf("missing quote markdown: %q", got)
	}
	if !strings.Contains(got, "---") {
		t.Fatalf("missing horizontal rule markdown: %q", got)
	}
}

func TestFromMarkdown_Table(t *testing.T) {
	input := "| Name | Value |\n| --- | --- |\n| A | 1 |\n| B | 2 |\n"
	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 body element, got %d", len(doc.Body))
	}
	tbl, ok := doc.Body[0].(Table)
	if !ok {
		t.Fatalf("expected Table element, got %#v", doc.Body[0])
	}
	if len(tbl.Header) != 2 || tbl.Header[0] != "Name" || tbl.Header[1] != "Value" {
		t.Fatalf("unexpected table header: %#v", tbl.Header)
	}
	if len(tbl.Rows) != 2 || len(tbl.Rows[0]) != 2 || tbl.Rows[0][0] != "A" || tbl.Rows[1][1] != "2" {
		t.Fatalf("unexpected table rows: %#v", tbl.Rows)
	}
}

func TestDocumentToMarkdown_Table(t *testing.T) {
	doc := &Document{Body: []Element{
		Table{Header: []string{"Name", "Value"}, Rows: [][]string{{"A", "1"}}},
	}}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "| A | 1 |") || !strings.Contains(got, "| Name | Value |") {
		t.Fatalf("unexpected table markdown: %q", got)
	}
}

func TestFromMarkdown_PreservesInlineFormatting(t *testing.T) {
	input := "Hello **bold** *italic* `code` and [link](https://example.com)\n"
	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 body element, got %d", len(doc.Body))
	}
	p := doc.Body[0].(Paragraph)
	want := "Hello **bold** *italic* `code` and [link](https://example.com)"
	if p.Text != want {
		t.Fatalf("inline markers not preserved: got %q want %q", p.Text, want)
	}
}

func TestFromMarkdown_TablePreservesInlineFormatting(t *testing.T) {
	input := "| Col |\n| --- |\n| **bold** |\n"
	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	tbl := doc.Body[0].(Table)
	if tbl.Rows[0][0] != "**bold**" {
		t.Fatalf("table cell markers not preserved: got %q", tbl.Rows[0][0])
	}
}

func TestFromMarkdown_ListItemPreservesInlineFormatting(t *testing.T) {
	input := "- item with **bold**\n"
	v, err := FromMarkdown(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	doc := v.(*Document)
	li := doc.Body[0].(ListItem)
	if li.Text != "item with **bold**" {
		t.Fatalf("list item markers not preserved: got %q", li.Text)
	}
}

func TestFromMarkdown_PreservesHardLineBreak(t *testing.T) {
	v, err := FromMarkdown("a  \n**Home**\n")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	p := v.(*Document).Body[0].(Paragraph)
	if p.Text != "a  \n**Home**" {
		t.Fatalf("hard break not preserved: got %q", p.Text)
	}
}

func TestDocumentToMarkdown_PreservesHardLineBreak(t *testing.T) {
	doc := &Document{Body: []Element{Paragraph{Text: "a  \nHome"}}}
	got, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "a  \nHome") {
		t.Fatalf("hard break flattened: %q", got)
	}
	if strings.Contains(got, "a   Home") {
		t.Fatalf("hard break became spaces: %q", got)
	}
}
