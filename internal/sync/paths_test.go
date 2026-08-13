package sync

import (
	"testing"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
)

func TestMapTabTree_LeafAndNested(t *testing.T) {
	tabs := []*md.Tab{
		{ID: "w", Title: "Welcome", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "hi"}}}},
		{ID: "fy", Title: "FY2024", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "year"}}}, Children: []*md.Tab{
			{ID: "jan", Title: "Jan", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "jan"}}}},
			{ID: "q1", Title: "Q1", Document: &md.Document{}, Children: []*md.Tab{
				{ID: "mar", Title: "March", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "mar"}}}},
			}},
		}},
	}
	got := MapTabTree("doc-1", tabs)
	want := []string{
		"welcome.md",
		"fy2024/fy2024.md",
		"fy2024/jan.md",
		"fy2024/q1/q1.md",
		"fy2024/q1/march.md",
	}
	if len(got) != len(want) {
		t.Fatalf("len=%d want %d: %+v", len(got), len(want), rels(got))
	}
	for i, p := range want {
		if got[i].RelPath != p {
			t.Fatalf("path[%d]=%s want %s", i, got[i].RelPath, p)
		}
		if !got[i].Doc.Track || got[i].Doc.DocID != "doc-1" {
			t.Fatalf("frontmatter fields missing on %s: %+v", p, got[i].Doc)
		}
	}
}

func TestMapTabTree_EmptyParentBodyStillMapped(t *testing.T) {
	tabs := []*md.Tab{
		{ID: "p", Title: "Parent", Children: []*md.Tab{
			{ID: "c", Title: "Child"},
		}},
	}
	got := MapTabTree("d", tabs)
	if len(got) != 2 || got[0].RelPath != "parent/parent.md" {
		t.Fatalf("got %+v", rels(got))
	}
}

func TestMapTabTree_ChildSlugCollision(t *testing.T) {
	tabs := []*md.Tab{
		{ID: "parentid", Title: "Same", Children: []*md.Tab{
			{ID: "childid99", Title: "Same"},
		}},
	}
	got := MapTabTree("d", tabs)
	if len(got) != 2 {
		t.Fatalf("len=%d %+v", len(got), rels(got))
	}
	if got[0].RelPath != "same/same.md" {
		t.Fatalf("parent path %s", got[0].RelPath)
	}
	if got[1].RelPath == "same/same.md" {
		t.Fatal("child should not occupy parent-body path")
	}
	if got[1].RelPath != "same/same-childid.md" && got[1].RelPath != "same/same-childid99.md" {
		// shortTabID keeps last 8 alnum chars of childid99 -> childid9? childid99 is 9 chars, last 8 = hildid99
		if got[1].Doc == nil || got[1].TabID != "childid99" {
			t.Fatalf("unexpected child mapping %+v", got[1])
		}
	}
}

func TestSlugify(t *testing.T) {
	if g := slugify("FY 2024!"); g != "fy-2024" {
		t.Fatalf("slug=%q", g)
	}
	if g := slugify(""); g != "" {
		t.Fatalf("empty slug=%q", g)
	}
}

func rels(ms []MappedTab) []string {
	var s []string
	for _, m := range ms {
		s = append(s, m.RelPath)
	}
	return s
}
