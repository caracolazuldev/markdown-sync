package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Document represents a very small subset of a document model used for
// conversion unit tests and early integration. Real Google Docs models will
// be adapted to this structure by the google adapter.
type Document struct {
	Title string
	Body  []Element
}

// Element is a single block element in the document body.
type Element interface {
	element()
}

type Paragraph struct{ Text string }

func (Paragraph) element() {}

type Heading struct {
	Level int
	Text  string
}

func (Heading) element() {}

type CodeBlock struct {
	Language string
	Code     string
}

func (CodeBlock) element() {}

type Image struct {
	URL string
	Alt string
}

func (Image) element() {}

type ListItem struct {
	Ordered bool
	Text    string
}

func (ListItem) element() {}

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
		case ListItem:
			if v.Ordered {
				sb.WriteString("1. ")
			} else {
				sb.WriteString("- ")
			}
			sb.WriteString(escapeString(v.Text))
			sb.WriteString("\n")
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
	clean, title := stripFrontMatter(md)
	parser := goldmark.DefaultParser()
	docNode := parser.Parse(text.NewReader([]byte(clean)))

	out := &Document{Title: title}
	for n := docNode.FirstChild(); n != nil; n = n.NextSibling() {
		switch v := n.(type) {
		case *ast.Heading:
			text := strings.TrimSpace(extractText(v, []byte(clean)))
			if text != "" {
				out.Body = append(out.Body, Heading{Level: v.Level, Text: text})
			}
		case *ast.Paragraph:
			// If paragraph contains only image node, map to Image element.
			if v.FirstChild() != nil && v.FirstChild() == v.LastChild() {
				if img, ok := v.FirstChild().(*ast.Image); ok {
					alt := strings.TrimSpace(extractText(img, []byte(clean)))
					out.Body = append(out.Body, Image{URL: string(img.Destination), Alt: alt})
					continue
				}
			}
			text := strings.TrimSpace(extractText(v, []byte(clean)))
			if text != "" {
				out.Body = append(out.Body, Paragraph{Text: text})
			}
		case *ast.FencedCodeBlock:
			lang := strings.TrimSpace(string(v.Language([]byte(clean))))
			var lines []string
			for i := 0; i < v.Lines().Len(); i++ {
				line := v.Lines().At(i)
				lines = append(lines, string(line.Value([]byte(clean))))
			}
			code := strings.TrimSuffix(strings.Join(lines, ""), "\n")
			out.Body = append(out.Body, CodeBlock{Language: lang, Code: code})
		case *ast.CodeBlock:
			var lines []string
			for i := 0; i < v.Lines().Len(); i++ {
				line := v.Lines().At(i)
				lines = append(lines, string(line.Value([]byte(clean))))
			}
			code := strings.TrimSuffix(strings.Join(lines, ""), "\n")
			out.Body = append(out.Body, CodeBlock{Language: "", Code: code})
		case *ast.Image:
			alt := strings.TrimSpace(extractText(v, []byte(clean)))
			out.Body = append(out.Body, Image{URL: string(v.Destination), Alt: alt})
		case *ast.List:
			ordered := v.IsOrdered()
			for li := v.FirstChild(); li != nil; li = li.NextSibling() {
				item, ok := li.(*ast.ListItem)
				if !ok {
					continue
				}
				text := strings.TrimSpace(extractText(item, []byte(clean)))
				if text == "" {
					continue
				}
				out.Body = append(out.Body, ListItem{Ordered: ordered, Text: text})
			}
		}
	}

	return out, nil
}

func extractText(n ast.Node, source []byte) string {
	var sb strings.Builder
	ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch t := node.(type) {
		case *ast.Text:
			sb.Write(t.Segment.Value(source))
		case *ast.CodeSpan:
			sb.Write(t.Text(source))
		}
		return ast.WalkContinue, nil
	})
	return sb.String()
}

func stripFrontMatter(in string) (string, string) {
	if !strings.HasPrefix(in, "---\n") {
		return in, ""
	}
	end := strings.Index(in[4:], "\n---\n")
	if end == -1 {
		return in, ""
	}
	fm := in[4 : 4+end]
	rest := in[4+end+5:]
	re := regexp.MustCompile(`(?m)^title:\s*(.+)\s*$`)
	m := re.FindStringSubmatch(fm)
	title := ""
	if len(m) > 1 {
		title = strings.TrimSpace(m[1])
		title = strings.Trim(title, `"'`)
	}
	return rest, title
}
