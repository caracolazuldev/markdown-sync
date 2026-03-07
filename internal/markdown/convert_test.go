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
	if !ok || !strings.Contains(p.Text, "Paragraph with") {
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
