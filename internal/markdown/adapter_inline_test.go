package markdown

import (
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
