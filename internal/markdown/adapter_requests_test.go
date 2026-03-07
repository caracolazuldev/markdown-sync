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
			if r.CreateParagraphBullets.BulletPreset == "BULLET_DISC_CIRCLE" {
				hasUnordered = true
			}
			if r.CreateParagraphBullets.BulletPreset == "NUMBERED_DECIMAL" {
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
		default:
			out = append(out, "Other")
		}
	}
	return out
}
