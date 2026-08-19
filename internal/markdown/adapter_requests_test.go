package markdown

import (
	"testing"

	docs "google.golang.org/api/docs/v1"
)

func TestBuildDocsRequests_ListItemsCreateBullets(t *testing.T) {
	doc := &Document{Body: []Element{
		ListItem{Ordered: false, Text: "one"},
		ListItem{Ordered: false, Level: 1, Text: "two"},
		Paragraph{Text: "break"},
		ListItem{Ordered: true, Text: "first"},
		ListItem{Ordered: true, Level: 2, Text: "second"},
	}}

	reqs, err := buildDocsRequests("", doc)
	if err != nil {
		t.Fatalf("buildDocsRequests error: %v", err)
	}

	var bulletCount int
	var hasUnordered, hasOrdered bool
	for _, r := range reqs {
		if r.CreateParagraphBullets != nil {
			bulletCount++
			if r.CreateParagraphBullets.BulletPreset == "BULLET_DISC_CIRCLE_SQUARE" {
				hasUnordered = true
			}
			if r.CreateParagraphBullets.BulletPreset == "NUMBERED_DECIMAL_ALPHA_ROMAN" {
				hasOrdered = true
			}
		}
	}

	if bulletCount < 2 || !hasUnordered || !hasOrdered {
		t.Fatalf("expected unordered and ordered bullet requests, got count=%d unordered=%v ordered=%v reqs=%v", bulletCount, hasUnordered, hasOrdered, summarizeReqKinds(reqs))
	}

	// ensure tab-prefixed inserts are present for nested list levels
	var hasOneTab, hasTwoTabs bool
	for _, r := range reqs {
		if r.InsertText != nil {
			if len(r.InsertText.Text) > 0 && r.InsertText.Text[0] == '\t' {
				hasOneTab = true
			}
			if len(r.InsertText.Text) > 1 && r.InsertText.Text[0] == '\t' && r.InsertText.Text[1] == '\t' {
				hasTwoTabs = true
			}
		}
	}
	if !hasOneTab || !hasTwoTabs {
		t.Fatalf("expected nested list inserts with tab prefixes, got oneTab=%v twoTabs=%v", hasOneTab, hasTwoTabs)
	}
}

func TestBuildDocsRequests_QuoteAndRule(t *testing.T) {
	doc := &Document{Body: []Element{
		Quote{Text: "quoted"},
		HorizontalRule{},
	}}

	reqs, err := buildDocsRequests("", doc)
	if err != nil {
		t.Fatalf("buildDocsRequests error: %v", err)
	}

	var hasQuoteItalic, hasRuleInsert bool
	for _, r := range reqs {
		if r.UpdateTextStyle != nil && r.UpdateTextStyle.TextStyle != nil && r.UpdateTextStyle.TextStyle.Italic {
			hasQuoteItalic = true
		}
		if r.InsertText != nil && r.InsertText.Text == "──────────\n" {
			hasRuleInsert = true
		}
	}

	if !hasQuoteItalic || !hasRuleInsert {
		t.Fatalf("expected quote italic and rule insert, got italic=%v rule=%v reqs=%v", hasQuoteItalic, hasRuleInsert, summarizeReqKinds(reqs))
	}
}

func TestBuildDocsRequests_ParagraphInlineStyles(t *testing.T) {
	v, err := FromMarkdown("Hello **bold** *italic* `code` [link](https://example.com)\n")
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	doc := v.(*Document)

	reqs, err := buildDocsRequests("", doc)
	if err != nil {
		t.Fatalf("buildDocsRequests error: %v", err)
	}

	var hasBold, hasItalic, hasCode, hasLink bool
	for _, r := range reqs {
		if r.UpdateTextStyle == nil || r.UpdateTextStyle.TextStyle == nil {
			continue
		}
		ts := r.UpdateTextStyle.TextStyle
		if ts.Bold {
			hasBold = true
		}
		if ts.Italic {
			hasItalic = true
		}
		if ts.WeightedFontFamily != nil && ts.WeightedFontFamily.FontFamily == "Courier New" {
			hasCode = true
		}
		if ts.Link != nil && ts.Link.Url == "https://example.com" {
			hasLink = true
		}
	}
	if !hasBold || !hasItalic || !hasCode || !hasLink {
		t.Fatalf("missing inline styles: bold=%v italic=%v code=%v link=%v", hasBold, hasItalic, hasCode, hasLink)
	}
}

func TestBuildDocsRequests_Strikethrough(t *testing.T) {
	v, err := FromMarkdown("Hello ~~strike~~\n")
	if err != nil {
		t.Fatalf("FromMarkdown: %v", err)
	}
	reqs, err := buildDocsRequests("", v.(*Document))
	if err != nil {
		t.Fatalf("buildDocsRequests: %v", err)
	}
	var hasStrike bool
	for _, r := range reqs {
		if r.UpdateTextStyle != nil && r.UpdateTextStyle.TextStyle != nil && r.UpdateTextStyle.TextStyle.Strikethrough {
			hasStrike = true
			start := r.UpdateTextStyle.Range.StartIndex
			end := r.UpdateTextStyle.Range.EndIndex
			wantStart := int64(1) + utf16Len("Hello ")
			if start != wantStart || end != wantStart+utf16Len("strike") {
				t.Fatalf("strike range [%d,%d) want [%d,%d)", start, end, wantStart, wantStart+utf16Len("strike"))
			}
		}
	}
	if !hasStrike {
		t.Fatal("missing strikethrough UpdateTextStyle")
	}
}

func TestBuildDocsRequests_HeadingInlineStyles(t *testing.T) {
	v, err := FromMarkdown("# Title with **bold**\n")
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	doc := v.(*Document)

	reqs, err := buildDocsRequests("", doc)
	if err != nil {
		t.Fatalf("buildDocsRequests error: %v", err)
	}

	var hasBold, hasHeadingStyle bool
	for _, r := range reqs {
		if r.UpdateTextStyle != nil && r.UpdateTextStyle.TextStyle != nil && r.UpdateTextStyle.TextStyle.Bold {
			hasBold = true
		}
		if r.UpdateParagraphStyle != nil && r.UpdateParagraphStyle.ParagraphStyle != nil &&
			r.UpdateParagraphStyle.ParagraphStyle.NamedStyleType == "HEADING_1" {
			hasHeadingStyle = true
		}
	}
	if !hasBold || !hasHeadingStyle {
		t.Fatalf("expected heading bold and heading style, got bold=%v heading=%v", hasBold, hasHeadingStyle)
	}
}

func TestBuildDocsRequests_UTF16BoldRangeAfterArrow(t *testing.T) {
	v, err := FromMarkdown("Committee → **Home**\n")
	if err != nil {
		t.Fatalf("FromMarkdown: %v", err)
	}
	reqs, err := buildDocsRequests("", v.(*Document))
	if err != nil {
		t.Fatalf("buildDocsRequests: %v", err)
	}
	start, end, ok := firstBoldRange(reqs)
	if !ok {
		t.Fatal("missing bold UpdateTextStyle")
	}
	wantStart := int64(1) + utf16Len("Committee → ")
	wantEnd := wantStart + utf16Len("Home")
	if start != wantStart || end != wantEnd {
		t.Fatalf("bold range [%d,%d) want [%d,%d)", start, end, wantStart, wantEnd)
	}
}

func TestBuildDocsRequests_UTF16BoldRangeAfterEmoji(t *testing.T) {
	v, err := FromMarkdown("😀 **x**\n")
	if err != nil {
		t.Fatalf("FromMarkdown: %v", err)
	}
	reqs, err := buildDocsRequests("", v.(*Document))
	if err != nil {
		t.Fatalf("buildDocsRequests: %v", err)
	}
	start, end, ok := firstBoldRange(reqs)
	if !ok {
		t.Fatal("missing bold UpdateTextStyle")
	}
	wantStart := int64(1) + utf16Len("😀 ")
	if start != wantStart || end != wantStart+1 {
		t.Fatalf("bold range [%d,%d) want [%d,%d)", start, end, wantStart, wantStart+1)
	}
}

func TestBuildDocsRequests_HardLineBreakVerticalTab(t *testing.T) {
	v, err := FromMarkdown("a  \n**Home**\n")
	if err != nil {
		t.Fatalf("FromMarkdown: %v", err)
	}
	doc := v.(*Document)
	p, ok := doc.Body[0].(Paragraph)
	if !ok || p.Text != "a  \n**Home**" {
		t.Fatalf("paragraph text=%#v", doc.Body[0])
	}
	reqs, err := buildDocsRequests("", doc)
	if err != nil {
		t.Fatalf("buildDocsRequests: %v", err)
	}
	var inserted string
	for _, r := range reqs {
		if r.InsertText != nil {
			inserted += r.InsertText.Text
		}
	}
	if inserted != "a\u000bHome\n" {
		t.Fatalf("insert=%q want %q", inserted, "a\u000bHome\n")
	}
	start, end, ok := firstBoldRange(reqs)
	if !ok {
		t.Fatal("missing bold UpdateTextStyle")
	}
	wantStart := int64(1) + utf16Len("a\u000b")
	if start != wantStart || end != wantStart+4 {
		t.Fatalf("bold range [%d,%d) want [%d,%d)", start, end, wantStart, wantStart+4)
	}
}

func TestBuildDocsRequests_HardBreakAfterArrowThenBold(t *testing.T) {
	v, err := FromMarkdown("→  \n**Home**\n")
	if err != nil {
		t.Fatalf("FromMarkdown: %v", err)
	}
	reqs, err := buildDocsRequests("", v.(*Document))
	if err != nil {
		t.Fatalf("buildDocsRequests: %v", err)
	}
	start, end, ok := firstBoldRange(reqs)
	if !ok {
		t.Fatal("missing bold UpdateTextStyle")
	}
	wantStart := int64(1) + utf16Len("→\u000b")
	if start != wantStart || end != wantStart+4 {
		t.Fatalf("bold range [%d,%d) want [%d,%d)", start, end, wantStart, wantStart+4)
	}
}

func firstBoldRange(reqs []*docs.Request) (start, end int64, ok bool) {
	for _, r := range reqs {
		if r.UpdateTextStyle != nil && r.UpdateTextStyle.TextStyle != nil && r.UpdateTextStyle.TextStyle.Bold && r.UpdateTextStyle.Range != nil {
			return r.UpdateTextStyle.Range.StartIndex, r.UpdateTextStyle.Range.EndIndex, true
		}
	}
	return 0, 0, false
}

func TestBuildDocsRequests_TableNativeInsert(t *testing.T) {
	doc := &Document{Body: []Element{
		Table{Header: []string{"Name", "Value"}, Rows: [][]string{{"A", "1"}}},
	}}

	reqs, err := buildDocsRequests("", doc)
	if err != nil {
		t.Fatalf("buildDocsRequests error: %v", err)
	}

	var found bool
	for _, r := range reqs {
		if r.InsertTable != nil {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected native InsertTable request, req kinds=%v", summarizeReqKinds(reqs))
	}
}

func TestBuildDocsRequests_TableFallbackInsertForIrregularRows(t *testing.T) {
	doc := &Document{Body: []Element{
		Table{Header: []string{"Name", "Value"}, Rows: [][]string{{"A"}}},
	}}

	reqs, err := buildDocsRequests("", doc)
	if err != nil {
		t.Fatalf("buildDocsRequests error: %v", err)
	}

	var hasFallbackInsert bool
	for _, r := range reqs {
		if r.InsertText != nil && r.InsertText.Text != "" && r.InsertText.Text[0] == '|' {
			hasFallbackInsert = true
			break
		}
	}
	if !hasFallbackInsert {
		t.Fatalf("expected markdown fallback insert for irregular table, req kinds=%v", summarizeReqKinds(reqs))
	}
}

func summarizeReqKinds(reqs []*docs.Request) []string {
	var out []string
	for _, r := range reqs {
		switch {
		case r.InsertText != nil:
			out = append(out, "InsertText")
		case r.UpdateTextStyle != nil:
			out = append(out, "UpdateTextStyle")
		case r.UpdateParagraphStyle != nil:
			out = append(out, "UpdateParagraphStyle")
		case r.CreateParagraphBullets != nil:
			out = append(out, "CreateParagraphBullets")
		case r.InsertInlineImage != nil:
			out = append(out, "InsertInlineImage")
		case r.InsertTable != nil:
			out = append(out, "InsertTable")
		default:
			out = append(out, "Other")
		}
	}
	return out
}
