//go:build integration
// +build integration

package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	gauth "github.com/caracolazuldev/gdocs-markdown-sync/internal/google"
	docs "google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

const integrationMarker = "INTEGRATION_TABLE_MARKER_V1"

// This integration test updates a real Google Doc. It only runs when the
// environment variables `GOOGLE_APPLICATION_CREDENTIALS` and
// `INTEGRATION_DOC_ID` are set and when `-tags=integration` is provided to
// `go test`.
func TestApplyDocumentIntegration(t *testing.T) {
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	docID := os.Getenv("INTEGRATION_DOC_ID")
	if creds == "" || docID == "" {
		t.Skip("integration test skipped; set GOOGLE_APPLICATION_CREDENTIALS and INTEGRATION_DOC_ID to enable")
	}

	// build a simple local file to apply
	tmp := t.TempDir()
	file := tmp + "/apply.md"
	if err := os.WriteFile(file, []byte("# Integration Test\n\nApplied by test."), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	var out bytes.Buffer
	var errb bytes.Buffer
	if err := importToWriter("service", file, docID, false, &out, &errb); err != nil {
		t.Fatalf("importToWriter apply failed: %v, stderr=%s", err, errb.String())
	}
}

func TestApplyDocumentIntegration_TableContent(t *testing.T) {
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	docID := os.Getenv("INTEGRATION_DOC_ID")
	if creds == "" || docID == "" {
		t.Skip("integration test skipped; set GOOGLE_APPLICATION_CREDENTIALS and INTEGRATION_DOC_ID to enable")
	}

	tmp := t.TempDir()
	file := tmp + "/apply-table.md"
	content := "# Integration Table Test\n\n" +
		integrationMarker + "\n\n" +
		"| Name | Value |\n" +
		"| --- | --- |\n" +
		"| A | **1** |\n" +
		"| B | [2](https://example.com) |\n"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	var out bytes.Buffer
	var errb bytes.Buffer
	if err := importToWriter("service", file, docID, false, &out, &errb); err != nil {
		t.Fatalf("importToWriter apply failed: %v, stderr=%s", err, errb.String())
	}
	if out.Len() == 0 {
		t.Fatalf("expected apply output message")
	}

	ctx := context.Background()
	httpClient, err := gauth.NewHTTPClient(ctx, "service", creds)
	if err != nil {
		t.Fatalf("create auth client: %v", err)
	}
	svc, err := docs.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("create docs service: %v", err)
	}
	remote, err := svc.Documents.Get(docID).Do()
	if err != nil {
		t.Fatalf("fetch updated document: %v", err)
	}

	text, hasBold, hasLink := collectDocStats(remote)
	for _, want := range []string{"Integration Table Test", integrationMarker, "Name", "Value", "A", "B", "1", "2"} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected updated doc text to contain %q, got %q", want, text)
		}
	}
	if strings.Count(text, integrationMarker) != 1 {
		t.Fatalf("expected marker %q exactly once, got count=%d text=%q", integrationMarker, strings.Count(text, integrationMarker), text)
	}
	if !hasBold {
		t.Fatalf("expected at least one bold text run in updated doc")
	}
	if !hasLink {
		t.Fatalf("expected at least one link text run in updated doc")
	}
}

func collectDocStats(doc *docs.Document) (allText string, hasBold bool, hasLink bool) {
	if doc == nil || doc.Body == nil {
		return "", false, false
	}
	var sb strings.Builder
	var walkElements func(elements []*docs.StructuralElement)
	walkElements = func(elements []*docs.StructuralElement) {
		for _, se := range elements {
			if se == nil {
				continue
			}
			if se.Paragraph != nil {
				for _, pe := range se.Paragraph.Elements {
					if pe == nil || pe.TextRun == nil {
						continue
					}
					sb.WriteString(pe.TextRun.Content)
					if pe.TextRun.TextStyle != nil {
						if pe.TextRun.TextStyle.Bold {
							hasBold = true
						}
						if pe.TextRun.TextStyle.Link != nil && pe.TextRun.TextStyle.Link.Url != "" {
							hasLink = true
						}
					}
				}
			}
			if se.Table != nil {
				for _, row := range se.Table.TableRows {
					if row == nil {
						continue
					}
					for _, cell := range row.TableCells {
						if cell == nil {
							continue
						}
						walkElements(cell.Content)
					}
				}
			}
		}
	}
	walkElements(doc.Body.Content)
	return sb.String(), hasBold, hasLink
}
