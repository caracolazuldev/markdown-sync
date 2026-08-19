package markdown

import (
	"strings"
	"testing"
)

func TestParseInline_GoldmarkSpans(t *testing.T) {
	input := "hello **bold** *italic* `code` and [link](https://example.com)"
	clean, spans := parseInline(input)

	expected := "hello bold italic code and link"
	if clean != expected {
		t.Fatalf("unexpected cleaned text: got %q want %q", clean, expected)
	}

	if len(spans) == 0 {
		t.Fatalf("expected non-empty spans")
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
		t.Fatalf("missing styles: bold=%v italic=%v code=%v link=%v", foundBold, foundItalic, foundCode, foundLink)
	}
}

func TestParseInline_Strikethrough(t *testing.T) {
	clean, spans := parseInline("→ ~~strike~~")
	if clean != "→ strike" {
		t.Fatalf("clean=%q", clean)
	}
	sp, ok := spanOf(spans, "strike")
	if !ok {
		t.Fatal("missing strike span")
	}
	wantOff := int(utf16Len("→ "))
	if sp.Offset != wantOff || sp.Length != int(utf16Len("strike")) {
		t.Fatalf("strike span offset=%d length=%d want offset=%d length=%d", sp.Offset, sp.Length, wantOff, utf16Len("strike"))
	}
}

func TestParseInline_UTF16OffsetAfterArrow(t *testing.T) {
	clean, spans := parseInline("Committee → **Home**")
	if clean != "Committee → Home" {
		t.Fatalf("clean=%q", clean)
	}
	sp, ok := spanOf(spans, "bold")
	if !ok {
		t.Fatal("missing bold span")
	}
	wantOff := int(utf16Len("Committee → "))
	if sp.Offset != wantOff || sp.Length != int(utf16Len("Home")) {
		t.Fatalf("bold span offset=%d length=%d want offset=%d length=%d", sp.Offset, sp.Length, wantOff, utf16Len("Home"))
	}
}

func TestParseInline_UTF16OffsetAfterCurlyQuotes(t *testing.T) {
	clean, spans := parseInline("“” **after**")
	if clean != "“” after" {
		t.Fatalf("clean=%q", clean)
	}
	sp, ok := spanOf(spans, "bold")
	if !ok {
		t.Fatal("missing bold span")
	}
	wantOff := int(utf16Len("“” "))
	if sp.Offset != wantOff || sp.Length != int(utf16Len("after")) {
		t.Fatalf("bold span offset=%d length=%d want offset=%d length=%d", sp.Offset, sp.Length, wantOff, utf16Len("after"))
	}
}

func TestParseInline_UTF16OffsetAfterEmoji(t *testing.T) {
	clean, spans := parseInline("😀 **x**")
	if clean != "😀 x" {
		t.Fatalf("clean=%q", clean)
	}
	sp, ok := spanOf(spans, "bold")
	if !ok {
		t.Fatal("missing bold span")
	}
	wantOff := int(utf16Len("😀 "))
	if wantOff != 3 {
		t.Fatalf("emoji prefix utf16=%d want 3", wantOff)
	}
	if sp.Offset != wantOff || sp.Length != 1 {
		t.Fatalf("bold span offset=%d length=%d want offset=%d length=1", sp.Offset, sp.Length, wantOff)
	}
}

func TestParseInline_HardLineBreakIsVerticalTab(t *testing.T) {
	clean, spans := parseInline("a  \n**Home**")
	if clean != "a\u000bHome" {
		t.Fatalf("clean=%q want %q", clean, "a\u000bHome")
	}
	if strings.Contains(clean, "\n") {
		t.Fatal("hard break must not become a paragraph newline")
	}
	sp, ok := spanOf(spans, "bold")
	if !ok {
		t.Fatal("missing bold span")
	}
	wantOff := int(utf16Len("a\u000b"))
	if sp.Offset != wantOff || sp.Length != 4 {
		t.Fatalf("bold span offset=%d length=%d want offset=%d length=4", sp.Offset, sp.Length, wantOff)
	}
}

func spanOf(spans []inlineSpan, kind string) (inlineSpan, bool) {
	for _, sp := range spans {
		if sp.Kind == kind {
			return sp, true
		}
	}
	return inlineSpan{}, false
}
