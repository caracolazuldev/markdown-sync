package sync

import (
	"path"
	"strconv"
	"strings"
	"unicode"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
)

// MappedTab is one tab's destination path inside a track root.
type MappedTab struct {
	TabID   string
	Title   string
	RelPath string
	Doc     *md.Document
}

// MapTabTree maps a Docs tab tree to nested relative paths.
// A leaf is <slug>.md. A parent with children is <slug>/<slug>.md plus children.
func MapTabTree(docID string, tabs []*md.Tab) []MappedTab {
	var out []MappedTab
	mapLevel("", "", docID, tabs, &out)
	return out
}

func mapLevel(dir, reservedSlug, docID string, tabs []*md.Tab, out *[]MappedTab) {
	used := map[string]bool{}
	if reservedSlug != "" {
		used[reservedSlug] = true
	}
	for _, t := range tabs {
		if t == nil {
			continue
		}
		slug := uniqueSlug(slugify(t.Title), t.ID, used)
		used[slug] = true
		doc := decorateTab(t, docID)
		if len(t.Children) > 0 {
			rel := path.Join(dir, slug, slug+".md")
			*out = append(*out, MappedTab{TabID: t.ID, Title: t.Title, RelPath: rel, Doc: doc})
			mapLevel(path.Join(dir, slug), slug, docID, t.Children, out)
			continue
		}
		rel := path.Join(dir, slug+".md")
		*out = append(*out, MappedTab{TabID: t.ID, Title: t.Title, RelPath: rel, Doc: doc})
	}
}

func decorateTab(t *md.Tab, docID string) *md.Document {
	d := &md.Document{Title: t.Title, DocID: docID, TabID: t.ID, Track: true}
	if t.Document != nil {
		d.Body = t.Document.Body
		if d.Title == "" {
			d.Title = t.Document.Title
		}
	}
	return d
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevHyphen := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			prevHyphen = false
		case !prevHyphen:
			b.WriteByte('-')
			prevHyphen = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func uniqueSlug(slug, tabID string, used map[string]bool) string {
	if slug == "" {
		slug = "tab"
	}
	if !used[slug] {
		return slug
	}
	suf := shortTabID(tabID)
	cand := slug + "-" + suf
	if !used[cand] {
		return cand
	}
	n := 2
	for {
		c := cand + "-" + strconv.Itoa(n)
		if !used[c] {
			return c
		}
		n++
	}
}

func shortTabID(id string) string {
	var b strings.Builder
	for _, r := range id {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	s := b.String()
	if s == "" {
		return "id"
	}
	if len(s) > 8 {
		s = s[len(s)-8:]
	}
	return strings.ToLower(s)
}
