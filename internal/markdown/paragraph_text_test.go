package markdown

import (
	"strings"
	"testing"

	docs "google.golang.org/api/docs/v1"
)

func TestParagraphText_InlineStyles(t *testing.T) {
	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "Hello "}},
		{TextRun: &docs.TextRun{Content: "bold", TextStyle: &docs.TextStyle{Bold: true}}},
		{TextRun: &docs.TextRun{Content: " "}},
		{TextRun: &docs.TextRun{Content: "italic", TextStyle: &docs.TextStyle{Italic: true}}},
		{TextRun: &docs.TextRun{Content: " "}},
		{TextRun: &docs.TextRun{Content: "code", TextStyle: &docs.TextStyle{
			WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: "Courier New"},
		}}},
		{TextRun: &docs.TextRun{Content: " "}},
		{TextRun: &docs.TextRun{Content: "x", TextStyle: &docs.TextStyle{
			Link: &docs.Link{Url: "https://example.com"},
		}}},
		{TextRun: &docs.TextRun{Content: "\n"}},
	}}
	got := paragraphText(p)
	want := "Hello **bold** *italic* `code` [x](https://example.com)"
	if got != want {
		t.Fatalf("paragraphText=%q want %q", got, want)
	}
}

func TestParagraphText_BoldItalicCombined(t *testing.T) {
	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "both", TextStyle: &docs.TextStyle{Bold: true, Italic: true}}},
	}}
	got := paragraphText(p)
	if got != "***both***" {
		t.Fatalf("paragraphText=%q want %q", got, "***both***")
	}
}

func TestParagraphText_MergesAdjacentSameStyle(t *testing.T) {
	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "hel", TextStyle: &docs.TextStyle{Bold: true}}},
		{TextRun: &docs.TextRun{Content: "lo", TextStyle: &docs.TextStyle{Bold: true}}},
	}}
	got := paragraphText(p)
	if got != "**hello**" {
		t.Fatalf("paragraphText=%q want %q", got, "**hello**")
	}
}

func TestParagraphText_EscapesUnstyledPunctuation(t *testing.T) {
	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "*not-bold*"}},
	}}
	got := paragraphText(p)
	if got != `\*not-bold\*` {
		t.Fatalf("paragraphText=%q want %q", got, `\*not-bold\*`)
	}
}

func TestParagraphText_DoesNotEscapeInsideCode(t *testing.T) {
	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "a*b", TextStyle: &docs.TextStyle{
			WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: "Courier New"},
		}}},
	}}
	got := paragraphText(p)
	if got != "`a*b`" {
		t.Fatalf("paragraphText=%q want %q", got, "`a*b`")
	}
}

func TestBodyToDocument_HeadingAndTableInlineStyles(t *testing.T) {
	body := &docs.Body{Content: []*docs.StructuralElement{
		{Paragraph: &docs.Paragraph{
			ParagraphStyle: &docs.ParagraphStyle{NamedStyleType: "HEADING_1"},
			Elements: []*docs.ParagraphElement{
				{TextRun: &docs.TextRun{Content: "Title with "}},
				{TextRun: &docs.TextRun{Content: "bold", TextStyle: &docs.TextStyle{Bold: true}}},
				{TextRun: &docs.TextRun{Content: "\n"}},
			},
		}},
		{Table: &docs.Table{TableRows: []*docs.TableRow{
			{TableCells: []*docs.TableCell{{Content: []*docs.StructuralElement{
				{Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{
					{TextRun: &docs.TextRun{Content: "cell", TextStyle: &docs.TextStyle{Bold: true}}},
				}}},
			}}}},
		}}},
	}}
	doc := bodyToDocument("T", "doc-1", "tab-1", body, nil)
	if len(doc.Body) != 2 {
		t.Fatalf("body len=%d want 2: %+v", len(doc.Body), doc.Body)
	}
	h, ok := doc.Body[0].(Heading)
	if !ok || h.Text != "Title with **bold**" {
		t.Fatalf("heading=%#v", doc.Body[0])
	}
	tbl, ok := doc.Body[1].(Table)
	if !ok || len(tbl.Header) != 1 || tbl.Header[0] != "**cell**" {
		t.Fatalf("table=%#v", doc.Body[1])
	}
}

func TestParagraphText_RoundTripInlineStyles(t *testing.T) {
	input := "Hello **bold** *italic* `code` [x](https://example.com)"
	v, err := FromMarkdown(input + "\n")
	if err != nil {
		t.Fatalf("FromMarkdown: %v", err)
	}
	src := v.(*Document)
	if _, err := buildDocsRequests("", src); err != nil {
		t.Fatalf("buildDocsRequests: %v", err)
	}

	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "Hello "}},
		{TextRun: &docs.TextRun{Content: "bold", TextStyle: &docs.TextStyle{Bold: true}}},
		{TextRun: &docs.TextRun{Content: " "}},
		{TextRun: &docs.TextRun{Content: "italic", TextStyle: &docs.TextStyle{Italic: true}}},
		{TextRun: &docs.TextRun{Content: " "}},
		{TextRun: &docs.TextRun{Content: "code", TextStyle: &docs.TextStyle{
			WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: "Courier New"},
		}}},
		{TextRun: &docs.TextRun{Content: " "}},
		{TextRun: &docs.TextRun{Content: "x", TextStyle: &docs.TextStyle{
			Link: &docs.Link{Url: "https://example.com"},
		}}},
	}}
	encoded := paragraphText(p)
	if encoded != input {
		t.Fatalf("paragraphText=%q want %q", encoded, input)
	}

	parsed, err := FromMarkdown(encoded + "\n")
	if err != nil {
		t.Fatalf("FromMarkdown encoded: %v", err)
	}
	doc := parsed.(*Document)
	if len(doc.Body) != 1 {
		t.Fatalf("body len=%d want 1", len(doc.Body))
	}
	para, ok := doc.Body[0].(Paragraph)
	if !ok || para.Text != input {
		t.Fatalf("reparsed paragraph=%#v", doc.Body[0])
	}

	out, err := DocumentToMarkdown(doc)
	if err != nil {
		t.Fatalf("DocumentToMarkdown: %v", err)
	}
	if !strings.Contains(out, input) {
		t.Fatalf("exported markdown missing inline markers: %q", out)
	}

	clean, spans := parseInline(para.Text)
	if clean != "Hello bold italic code x" {
		t.Fatalf("cleaned text=%q", clean)
	}
	var foundBold, foundItalic, foundCode, foundLink bool
	for _, sp := range spans {
		switch sp.Kind {
		case "bold":
			foundBold = true
		case "italic":
			foundItalic = true
		case "code":
			foundCode = true
		case "link":
			if sp.Data == "https://example.com" {
				foundLink = true
			}
		}
	}
	if !foundBold || !foundItalic || !foundCode || !foundLink {
		t.Fatalf("missing styles after round-trip: bold=%v italic=%v code=%v link=%v", foundBold, foundItalic, foundCode, foundLink)
	}
}

func TestParagraphText_HardLineBreak(t *testing.T) {
	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "a\u000bHome\n"}},
	}}
	got := paragraphText(p)
	if got != "a  \nHome" {
		t.Fatalf("paragraphText=%q want %q", got, "a  \nHome")
	}
}

func TestParagraphText_HardLineBreakSplitsStyledWrap(t *testing.T) {
	p := &docs.Paragraph{Elements: []*docs.ParagraphElement{
		{TextRun: &docs.TextRun{Content: "one\u000btwo", TextStyle: &docs.TextStyle{Bold: true}}},
	}}
	got := paragraphText(p)
	if got != "**one**  \n**two**" {
		t.Fatalf("paragraphText=%q want %q", got, "**one**  \n**two**")
	}
}
