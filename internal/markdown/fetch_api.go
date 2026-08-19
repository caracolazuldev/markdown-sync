package markdown

import (
	"context"
	"fmt"
	"os"
	"strings"

	gauth "github.com/caracolazuldev/gdocs-markdown-sync/internal/google"
	docs "google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// APIFetcher loads a document from the Google Docs API with tab content.
type APIFetcher struct{}

func (APIFetcher) Fetch(ctx context.Context, authMode, docID string) (*FetchedDocument, error) {
	if docID == "" {
		return nil, fmt.Errorf("missing doc id")
	}
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	client, err := gauth.NewHTTPClient(ctx, authMode, creds)
	if err != nil {
		return nil, fmt.Errorf("creating google http client: %w", err)
	}
	svc, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("creating docs service: %w", err)
	}
	remote, err := svc.Documents.Get(docID).IncludeTabsContent(true).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("fetching remote document: %w", err)
	}
	return fetchedFromAPI(docID, remote), nil
}

func fetchedFromAPI(docID string, remote *docs.Document) *FetchedDocument {
	out := &FetchedDocument{ID: docID, Title: remote.Title}
	if len(remote.Tabs) == 0 {
		out.Document = bodyToDocument(remote.Title, docID, "", remote.Body, remote.Lists)
		return out
	}
	out.Tabs = convertAPITabs(docID, remote.Tabs)
	out.Document = primaryDocument(out)
	return out
}

func convertAPITabs(docID string, tabs []*docs.Tab) []*Tab {
	var out []*Tab
	for _, t := range tabs {
		if t == nil {
			continue
		}
		tab := &Tab{}
		if t.TabProperties != nil {
			tab.ID = t.TabProperties.TabId
			tab.Title = t.TabProperties.Title
		}
		if t.DocumentTab != nil {
			tab.Document = bodyToDocument(tab.Title, docID, tab.ID, t.DocumentTab.Body, t.DocumentTab.Lists)
		} else {
			tab.Document = bodyToDocument(tab.Title, docID, tab.ID, nil, nil)
		}
		if len(t.ChildTabs) > 0 {
			tab.Children = convertAPITabs(docID, t.ChildTabs)
		}
		out = append(out, tab)
	}
	return out
}

func primaryDocument(f *FetchedDocument) *Document {
	if f.Document != nil {
		return f.Document
	}
	if d := firstLeafDocument(f.Tabs); d != nil {
		return d
	}
	return &Document{Title: f.Title, DocID: f.ID}
}

func firstLeafDocument(tabs []*Tab) *Document {
	for _, t := range tabs {
		if t == nil {
			continue
		}
		if len(t.Children) == 0 && t.Document != nil {
			return t.Document
		}
		if d := firstLeafDocument(t.Children); d != nil {
			return d
		}
		if t.Document != nil {
			return t.Document
		}
	}
	return nil
}

func bodyToDocument(title, docID, tabID string, body *docs.Body, lists map[string]docs.List) *Document {
	doc := &Document{Title: title, DocID: docID, TabID: tabID}
	if body == nil {
		return doc
	}
	for _, se := range body.Content {
		if se == nil {
			continue
		}
		if se.Paragraph != nil {
			text := strings.TrimLeft(paragraphText(se.Paragraph), "\t")
			if text == "" {
				continue
			}
			style := ""
			if se.Paragraph.ParagraphStyle != nil {
				style = se.Paragraph.ParagraphStyle.NamedStyleType
			}
			if level, ok := headingLevel(style); ok {
				doc.Body = append(doc.Body, Heading{Level: level, Text: text})
				continue
			}
			if se.Paragraph.Bullet != nil {
				level := int(se.Paragraph.Bullet.NestingLevel)
				doc.Body = append(doc.Body, ListItem{
					Ordered: listItemOrdered(lists, se.Paragraph.Bullet),
					Level:   level,
					Text:    text,
				})
				continue
			}
			doc.Body = append(doc.Body, Paragraph{Text: text})
			continue
		}
		if se.Table != nil {
			if tbl, ok := tableFromAPI(se.Table); ok {
				doc.Body = append(doc.Body, tbl)
			}
		}
	}
	return doc
}

func listItemOrdered(lists map[string]docs.List, bullet *docs.Bullet) bool {
	if bullet == nil || lists == nil {
		return false
	}
	lst, ok := lists[bullet.ListId]
	if !ok || lst.ListProperties == nil {
		return false
	}
	levels := lst.ListProperties.NestingLevels
	i := int(bullet.NestingLevel)
	if i < 0 || i >= len(levels) || levels[i] == nil {
		return false
	}
	return glyphTypeOrdered(levels[i].GlyphType)
}

func glyphTypeOrdered(g string) bool {
	switch g {
	case "DECIMAL", "ZERO_DECIMAL", "ALPHA", "UPPER_ALPHA", "LOWER_ALPHA", "ROMAN", "UPPER_ROMAN", "LOWER_ROMAN":
		return true
	default:
		return false
	}
}

func headingLevel(named string) (int, bool) {
	switch named {
	case "HEADING_1", "TITLE":
		return 1, true
	case "HEADING_2", "SUBTITLE":
		return 2, true
	case "HEADING_3":
		return 3, true
	case "HEADING_4":
		return 4, true
	case "HEADING_5":
		return 5, true
	case "HEADING_6":
		return 6, true
	default:
		return 0, false
	}
}

type runStyle struct {
	bold   bool
	italic bool
	code   bool
	link   string
}

type inlineRun struct {
	text  string
	style runStyle
}

func paragraphText(p *docs.Paragraph) string {
	if p == nil {
		return ""
	}
	var runs []inlineRun
	for _, el := range p.Elements {
		if el == nil || el.TextRun == nil {
			continue
		}
		content := strings.ReplaceAll(el.TextRun.Content, "\n", "")
		if content == "" {
			continue
		}
		st := styleFromTextStyle(el.TextRun.TextStyle)
		if n := len(runs); n > 0 && runs[n-1].style == st {
			runs[n-1].text += content
			continue
		}
		runs = append(runs, inlineRun{text: content, style: st})
	}
	var b strings.Builder
	for _, r := range runs {
		b.WriteString(styledRunToMarkdown(r.text, r.style))
	}
	return strings.TrimSpace(b.String())
}

func styleFromTextStyle(ts *docs.TextStyle) runStyle {
	if ts == nil {
		return runStyle{}
	}
	st := runStyle{bold: ts.Bold, italic: ts.Italic}
	if ts.WeightedFontFamily != nil && ts.WeightedFontFamily.FontFamily == "Courier New" {
		st.code = true
	}
	if ts.Link != nil {
		st.link = ts.Link.Url
	}
	return st
}

// styledRunToMarkdown encodes a Docs text run as inline markdown so re-import
// via parseInline can recover bold, italic, code, and links.
func styledRunToMarkdown(content string, st runStyle) string {
	if content == "" {
		return ""
	}
	parts := strings.Split(content, "\u000b")
	for i, part := range parts {
		parts[i] = wrapStyledSegment(part, st)
	}
	return strings.Join(parts, "  \n")
}

func wrapStyledSegment(content string, st runStyle) string {
	if content == "" {
		return ""
	}
	inner := content
	if st.code {
		inner = "`" + content + "`"
	} else {
		inner = escapeMarkdownPunctuation(content)
	}
	if st.bold && st.italic {
		inner = "***" + inner + "***"
	} else if st.bold {
		inner = "**" + inner + "**"
	} else if st.italic {
		inner = "*" + inner + "*"
	}
	if st.link != "" {
		inner = "[" + inner + "](" + st.link + ")"
	}
	return inner
}

func escapeMarkdownPunctuation(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\', '*', '_', '`', '[', ']':
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func tableFromAPI(t *docs.Table) (Table, bool) {
	if t == nil || len(t.TableRows) == 0 {
		return Table{}, false
	}
	var rows [][]string
	for _, row := range t.TableRows {
		if row == nil {
			continue
		}
		var cells []string
		for _, cell := range row.TableCells {
			cells = append(cells, cellText(cell))
		}
		rows = append(rows, cells)
	}
	if len(rows) == 0 {
		return Table{}, false
	}
	return Table{Header: rows[0], Rows: rows[1:]}, true
}

func cellText(cell *docs.TableCell) string {
	if cell == nil {
		return ""
	}
	var parts []string
	for _, se := range cell.Content {
		if se != nil && se.Paragraph != nil {
			if s := paragraphText(se.Paragraph); s != "" {
				parts = append(parts, s)
			}
		}
	}
	return strings.Join(parts, " ")
}
