package markdown

import (
    "fmt"
    "strings"
)

// Document represents a very small subset of a document model used for
// conversion unit tests and early integration. Real Google Docs models will
// be adapted to this structure by the google adapter.
type Document struct {
    Title string
    Body  []Element
}

// Element is a single block element in the document body.
type Element interface{
    element()
}

type Paragraph struct{ Text string }
func (Paragraph) element() {}

type Heading struct{ Level int; Text string }
func (Heading) element() {}

type CodeBlock struct{ Language string; Code string }
func (CodeBlock) element() {}

type Image struct{ URL string; Alt string }
func (Image) element() {}

// DocumentToMarkdown converts a Document into markdown text. This function is
// intentionally small and focuses on predictable, testable mappings used by
// the rest of the project.
func DocumentToMarkdown(doc *Document) (string, error) {
    if doc == nil {
        return "", fmt.Errorf("nil document")
    }
    var sb strings.Builder
    if doc.Title != "" {
        sb.WriteString("---\n")
        sb.WriteString("title: ")
        sb.WriteString(escapeString(doc.Title))
        sb.WriteString("\n---\n\n")
    }
    for _, e := range doc.Body {
        switch v := e.(type) {
        case Heading:
            if v.Level < 1 || v.Level > 6 {
                v.Level = 1
            }
            sb.WriteString(strings.Repeat("#", v.Level))
            sb.WriteString(" ")
            sb.WriteString(escapeString(v.Text))
            sb.WriteString("\n\n")
        case Paragraph:
            sb.WriteString(escapeString(v.Text))
            sb.WriteString("\n\n")
        case CodeBlock:
            sb.WriteString("```")
            if v.Language != "" {
                sb.WriteString(v.Language)
            }
            sb.WriteString("\n")
            sb.WriteString(v.Code)
            sb.WriteString("\n```\n\n")
        case Image:
            sb.WriteString("![")
            sb.WriteString(escapeString(v.Alt))
            sb.WriteString("](")
            sb.WriteString(v.URL)
            sb.WriteString(")\n\n")
        default:
            // unknown element, ignore
        }
    }
    return sb.String(), nil
}

func escapeString(s string) string {
    // minimal escaping for yaml/newlines in title; extend as needed
    return strings.ReplaceAll(s, "\n", " ")
}

// ToMarkdown is a compatibility wrapper accepting an interface. If the
// provided value is a *Document it will be converted; otherwise an error is
// returned. This keeps the API small while allowing tests to operate on the
// concrete Document type.
func ToMarkdown(doc interface{}) (string, error) {
    if d, ok := doc.(*Document); ok {
        return DocumentToMarkdown(d)
    }
    return "", fmt.Errorf("unsupported document type")
}

// FromMarkdown parses markdown and returns a docs-compatible model (placeholder).
// For now this is intentionally unimplemented and returns an error. Import
// workflows will add parsing as needed.
func FromMarkdown(md string) (interface{}, error) {
    return nil, fmt.Errorf("FromMarkdown not implemented")
}
