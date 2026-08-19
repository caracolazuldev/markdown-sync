package markdown

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// Document represents a very small subset of a document model used for
// conversion unit tests and early integration. Real Google Docs models will
// be adapted to this structure by the google adapter.
type Document struct {
	Title string
	DocID string
	TabID string
	Track bool
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
	Level   int
	Text    string
}

func (ListItem) element() {}

type Quote struct {
	Text string
}

func (Quote) element() {}

type HorizontalRule struct{}

func (HorizontalRule) element() {}

type Table struct {
	Header []string
	Rows   [][]string
}

func (Table) element() {}

// DocumentToMarkdown converts a Document into markdown text. This function is
// intentionally small and focuses on predictable, testable mappings used by
// the rest of the project.
func DocumentToMarkdown(doc *Document) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("nil document")
	}
	var sb strings.Builder
	if hasFrontMatter(doc) {
		sb.WriteString("---\n")
		if doc.Title != "" {
			sb.WriteString("title: ")
			sb.WriteString(escapeYAML(doc.Title))
			sb.WriteString("\n")
		}
		if doc.DocID != "" {
			sb.WriteString("doc_id: ")
			sb.WriteString(escapeYAML(doc.DocID))
			sb.WriteString("\n")
		}
		if doc.TabID != "" {
			sb.WriteString("tab_id: ")
			sb.WriteString(escapeYAML(doc.TabID))
			sb.WriteString("\n")
		}
		if doc.Track {
			sb.WriteString("track: true\n")
		}
		sb.WriteString("---\n\n")
	}
	for _, e := range doc.Body {
		switch v := e.(type) {
		case Heading:
			if v.Level < 1 || v.Level > 6 {
				v.Level = 1
			}
			sb.WriteString(strings.Repeat("#", v.Level))
			sb.WriteString(" ")
			sb.WriteString(v.Text)
			sb.WriteString("\n\n")
		case Paragraph:
			sb.WriteString(v.Text)
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
			sb.WriteString(escapeYAML(v.Alt))
			sb.WriteString("](")
			sb.WriteString(v.URL)
			sb.WriteString(")\n\n")
		case ListItem:
			if v.Level > 0 {
				sb.WriteString(strings.Repeat("  ", v.Level))
			}
			if v.Ordered {
				sb.WriteString("1. ")
			} else {
				sb.WriteString("- ")
			}
			sb.WriteString(v.Text)
			sb.WriteString("\n")
		case Quote:
			sb.WriteString("> ")
			sb.WriteString(v.Text)
			sb.WriteString("\n\n")
		case HorizontalRule:
			sb.WriteString("---\n\n")
		case Table:
			tableMD := tableToMarkdown(v)
			if tableMD != "" {
				sb.WriteString(tableMD)
				sb.WriteString("\n\n")
			}
		default:
			// unknown element, ignore
		}
	}
	return sb.String(), nil
}

func hasFrontMatter(doc *Document) bool {
	return doc.Title != "" || doc.DocID != "" || doc.TabID != "" || doc.Track
}

func escapeYAML(s string) string {
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
	mdp := goldmark.New(goldmark.WithExtensions(extension.Table))
	parser := mdp.Parser()
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
			appendListItems(out, v, []byte(clean), 0)
		case *ast.Blockquote:
			q := strings.TrimSpace(extractText(v, []byte(clean)))
			if q != "" {
				out.Body = append(out.Body, Quote{Text: q})
			}
		case *ast.ThematicBreak:
			out.Body = append(out.Body, HorizontalRule{})
		case *extast.Table:
			header, rows := extractTable(v, []byte(clean))
			if len(header) > 0 || len(rows) > 0 {
				out.Body = append(out.Body, Table{Header: header, Rows: rows})
			}
		}
	}

	return out, nil
}

func appendListItems(out *Document, list *ast.List, source []byte, level int) {
	ordered := list.IsOrdered()
	for li := list.FirstChild(); li != nil; li = li.NextSibling() {
		item, ok := li.(*ast.ListItem)
		if !ok {
			continue
		}
		itemText := extractListItemText(item, source)
		if itemText != "" {
			out.Body = append(out.Body, ListItem{Ordered: ordered, Level: level, Text: itemText})
		}
		for c := item.FirstChild(); c != nil; c = c.NextSibling() {
			if sub, ok := c.(*ast.List); ok {
				appendListItems(out, sub, source, level+1)
			}
		}
	}
}

func extractListItemText(item *ast.ListItem, source []byte) string {
	var parts []string
	for c := item.FirstChild(); c != nil; c = c.NextSibling() {
		if _, isList := c.(*ast.List); isList {
			continue
		}
		t := strings.TrimSpace(extractText(c, source))
		if t != "" {
			parts = append(parts, t)
		}
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

// extractText serializes block content to a string, preserving inline markdown
// markers so downstream Google Docs apply logic can detect bold, italic, code,
// and links via parseInline.
func extractText(n ast.Node, source []byte) string {
	var sb strings.Builder
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			switch v := c.(type) {
			case *ast.Text:
				sb.Write(v.Segment.Value(source))
				if v.HardLineBreak() {
					sb.WriteString("  \n")
				} else if v.SoftLineBreak() {
					sb.WriteByte(' ')
				}
			case *ast.CodeSpan:
				sb.WriteByte('`')
				sb.Write(v.Text(source))
				sb.WriteByte('`')
			case *ast.Emphasis:
				marker := strings.Repeat("*", v.Level)
				sb.WriteString(marker)
				walk(v)
				sb.WriteString(marker)
			case *ast.Link:
				sb.WriteByte('[')
				walk(v)
				sb.WriteByte(']')
				sb.WriteByte('(')
				sb.Write(v.Destination)
				sb.WriteByte(')')
			default:
				walk(v)
			}
		}
	}
	walk(n)
	return sb.String()
}

func stripFrontMatter(in string) (string, string) {
	body, fm := SplitFrontMatter(in)
	return body, fm.Title
}

// FrontMatter is the YAML header parsed from a Markdown file.
type FrontMatter struct {
	Title string
	DocID string
	TabID string
	Track bool
}

// SplitFrontMatter separates a leading YAML block from the Markdown body.
func SplitFrontMatter(in string) (string, FrontMatter) {
	var fm FrontMatter
	if !strings.HasPrefix(in, "---\n") {
		return in, fm
	}
	end := strings.Index(in[4:], "\n---\n")
	if end == -1 {
		return in, fm
	}
	raw := in[4 : 4+end]
	rest := in[4+end+5:]
	for _, line := range strings.Split(raw, "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		switch key {
		case "title":
			fm.Title = val
		case "doc_id":
			fm.DocID = val
		case "tab_id":
			fm.TabID = val
		case "track":
			fm.Track = val == "true" || val == "yes" || val == "1"
		}
	}
	return rest, fm
}

func extractTable(tbl *extast.Table, source []byte) ([]string, [][]string) {
	var header []string
	var rows [][]string
	for n := tbl.FirstChild(); n != nil; n = n.NextSibling() {
		switch v := n.(type) {
		case *extast.TableHeader:
			header = extractTableCells(v, source)
		case *extast.TableRow:
			rows = append(rows, extractTableRow(v, source))
		}
	}
	return header, rows
}

func extractTableCells(n ast.Node, source []byte) []string {
	var cells []string
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if cell, ok := c.(*extast.TableCell); ok {
			cells = append(cells, strings.TrimSpace(extractText(cell, source)))
		}
	}
	return cells
}

func extractTableRow(row *extast.TableRow, source []byte) []string {
	return extractTableCells(row, source)
}

func tableToMarkdown(t Table) string {
	if len(t.Header) == 0 && len(t.Rows) == 0 {
		return ""
	}
	header := t.Header
	if len(header) == 0 && len(t.Rows) > 0 {
		header = make([]string, len(t.Rows[0]))
	}
	var sb strings.Builder
	writeRow := func(cols []string) {
		sb.WriteString("| ")
		sb.WriteString(strings.Join(cols, " | "))
		sb.WriteString(" |\n")
	}
	writeRow(header)
	sep := make([]string, len(header))
	for i := range sep {
		sep[i] = "---"
	}
	writeRow(sep)
	for _, r := range t.Rows {
		writeRow(r)
	}
	return strings.TrimSuffix(sb.String(), "\n")
}
