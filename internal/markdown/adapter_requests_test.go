package markdown

import (
	"testing"

	docs "google.golang.org/api/docs/v1"
)

func TestBuildDocsRequests_ListItemsCreateBullets(t *testing.T) {
	doc := &Document{Body: []Element{
		ListItem{Ordered: false, Text: "one"},
		ListItem{Ordered: false, Text: "two"},
		Paragraph{Text: "break"},
		ListItem{Ordered: true, Text: "first"},
		ListItem{Ordered: true, Text: "second"},
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
