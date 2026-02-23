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
