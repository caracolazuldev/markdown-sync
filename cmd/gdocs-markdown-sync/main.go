package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
	gsync "github.com/caracolazuldev/gdocs-markdown-sync/internal/sync"
)

func usage() {
	fmt.Fprintf(os.Stderr, "gdocs-markdown-sync: simple CLI\n")
	fmt.Fprintf(os.Stderr, "Usage:\n  gdocs-markdown-sync <command> [flags]\nCommands: export, import, preview, list, track\n")
}

// activeFetcher is SampleFetcher in unit tests; main() sets APIFetcher for the binary.
var activeFetcher md.Fetcher = md.SampleFetcher{}

func fetchDoc(authMode, docID string) (*md.FetchedDocument, error) {
	return activeFetcher.Fetch(context.Background(), authMode, docID)
}

// previewToWriter renders markdown and writes at most maxLines to stdout.
// If maxLines <= 0, the entire document is written.
func previewToWriter(authMode, docID string, maxLines int, stdout io.Writer, stderr io.Writer) error {
	fetched, err := fetchDoc(authMode, docID)
	if err != nil {
		fmt.Fprintf(stderr, "failed to fetch doc: %v\n", err)
		return err
	}
	if fetched.IsTabbed() {
		err := gsync.TabbedHint(docID, fetched.TabCount())
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	docModel := fetched.Document
	markdown, err := md.DocumentToMarkdown(docModel)
	if err != nil {
		fmt.Fprintf(stderr, "conversion error: %v\n", err)
		return err
	}
	if maxLines <= 0 {
		_, err := io.WriteString(stdout, markdown)
		return err
	}
	// strip YAML frontmatter for preview counting if present
	body := markdown
	if strings.HasPrefix(markdown, "---\n") {
		if idx := strings.Index(markdown, "\n---\n"); idx != -1 {
			body = markdown[idx+len("\n---\n"):]
		}
	}
	// write up to maxLines lines from the body
	lines := strings.SplitN(body, "\n", maxLines+1)
	output := strings.Join(lines[:min(len(lines), maxLines)], "\n")
	if len(lines) > maxLines {
		output += "\n\n... (truncated)"
	}
	_, err = io.WriteString(stdout, output+"\n")
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	activeFetcher = md.APIFetcher{}
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	auth := fs.String("auth", "oauth", "auth flow: oauth|service")
	out := fs.String("out", "", "output directory or file")
	doc := fs.String("doc", "", "Google Doc ID")
	file := fs.String("file", "", "local markdown file")
	dry := fs.Bool("dry-run", false, "dry run")
	diff := fs.Bool("diff", false, "show diff between local file and remote doc")
	force := fs.Bool("force", false, "overwrite dirty tracked files or adopt an untracked directory")
	fs.Parse(os.Args[2:])

	switch cmd {
	case "export":
		if err := exportToWriter(*auth, *doc, *out, *dry, os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "export error: %v\n", err)
			os.Exit(1)
		}
	case "import":
		if *diff {
			if err := importToWriter(*auth, *file, *doc, true, os.Stdout, os.Stderr); err != nil {
				fmt.Fprintf(os.Stderr, "import error: %v\n", err)
				os.Exit(1)
			}
		} else if *dry {
			fmt.Printf("dry-run: would import file=%s doc=%s auth=%s\n", *file, *doc, *auth)
		} else {
			if err := importToWriter(*auth, *file, *doc, false, os.Stdout, os.Stderr); err != nil {
				fmt.Fprintf(os.Stderr, "import error: %v\n", err)
				os.Exit(1)
			}
		}
	case "preview":
		if err := previewToWriter(*auth, *doc, 20, os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "preview error: %v\n", err)
			os.Exit(1)
		}
	case "list":
		if err := listToWriter(*auth, os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "list error: %v\n", err)
			os.Exit(1)
		}
	case "track":
		if err := trackToWriter(*auth, *doc, *out, *dry, *force, os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "track error: %v\n", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

// exportToWriter performs the export flow and writes output or messages to the
// provided writers. It is separated from main for easier testing.
func exportToWriter(authMode, docID, out string, dry bool, stdout io.Writer, stderr io.Writer) error {
	if docID == "" {
		fmt.Fprintln(stderr, "export requires -doc <doc-id>")
		return fmt.Errorf("missing doc id")
	}
	fetched, err := fetchDoc(authMode, docID)
	if err != nil {
		fmt.Fprintf(stderr, "failed to fetch doc: %v\n", err)
		return err
	}
	if fetched.IsTabbed() {
		err := gsync.TabbedHint(docID, fetched.TabCount())
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	docModel := fetched.Document
	markdown, err := md.DocumentToMarkdown(docModel)
	if err != nil {
		fmt.Fprintf(stderr, "conversion error: %v\n", err)
		return err
	}
	if out == "" {
		_, err := io.WriteString(stdout, markdown)
		return err
	}
	if dry {
		_, err := fmt.Fprintf(stdout, "dry-run: would write %d bytes to %s\n", len(markdown), out)
		return err
	}
	if err := ioutil.WriteFile(out, []byte(markdown), 0644); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return err
	}
	_, err = fmt.Fprintf(stdout, "wrote %d bytes to %s\n", len(markdown), out)
	return err
}

// listToWriter writes a simple list of available documents to stdout.
func listToWriter(authMode string, stdout io.Writer, stderr io.Writer) error {
	list, err := md.ListDocuments(authMode)
	if err != nil {
		fmt.Fprintf(stderr, "failed to list documents: %v\n", err)
		return err
	}
	for _, ds := range list {
		if _, err := fmt.Fprintf(stdout, "%s\t%s\n", ds.ID, ds.Title); err != nil {
			return err
		}
	}
	return nil
}

// importToWriter compares a local markdown file with the remote document and
// writes a simple line-oriented diff to stdout when diffOnly is true.
func importToWriter(authMode, localFile, docID string, diffOnly bool, stdout io.Writer, stderr io.Writer) error {
	if docID == "" {
		fmt.Fprintln(stderr, "import requires -doc <doc-id>")
		return fmt.Errorf("missing doc id")
	}
	if localFile == "" {
		fmt.Fprintln(stderr, "import requires -file <path>")
		return fmt.Errorf("missing file path")
	}
	if err := gsync.CheckImportAllowed(localFile, docID); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	localBytes, err := ioutil.ReadFile(localFile)
	if err != nil {
		fmt.Fprintf(stderr, "failed to read local file: %v\n", err)
		return err
	}
	local := strings.Split(string(localBytes), "\n")

	fetched, err := fetchDoc(authMode, docID)
	if err != nil {
		fmt.Fprintf(stderr, "failed to fetch doc: %v\n", err)
		return err
	}
	if fetched.IsTabbed() {
		err := gsync.TabbedHint(docID, fetched.TabCount())
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	remoteMd, err := md.DocumentToMarkdown(fetched.Document)
	if err != nil {
		fmt.Fprintf(stderr, "conversion error: %v\n", err)
		return err
	}
	remote := strings.Split(remoteMd, "\n")

	if diffOnly {
		// Produce a minimal unified diff with a single chunk.
		fmt.Fprintln(stdout, "--- remote")
		fmt.Fprintln(stdout, "+++ local")
		// compute lengths
		rLen := len(remote)
		lLen := len(local)
		fmt.Fprintf(stdout, "@@ -1,%d +1,%d @@\n", rLen, lLen)
		max := rLen
		if lLen > max {
			max = lLen
		}
		for i := 0; i < max; i++ {
			var r, l string
			if i < rLen {
				r = remote[i]
			}
			if i < lLen {
				l = local[i]
			}
			if r == l {
				fmt.Fprintf(stdout, " %s\n", r)
			} else {
				if r != "" {
					fmt.Fprintf(stdout, "- %s\n", r)
				}
				if l != "" {
					fmt.Fprintf(stdout, "+ %s\n", l)
				}
			}
		}
		return nil
	}

	// Apply: parse local markdown into a Document and send to adapter.
	parsed, err := md.FromMarkdown(string(localBytes))
	if err != nil {
		fmt.Fprintf(stderr, "failed to parse local markdown: %v\n", err)
		return err
	}
	d, ok := parsed.(*md.Document)
	if !ok {
		return fmt.Errorf("parsed markdown returned unexpected type")
	}
	if err := md.ApplyDocument(authMode, docID, d); err != nil {
		fmt.Fprintf(stderr, "failed to apply document: %v\n", err)
		return err
	}
	_, err = fmt.Fprintf(stdout, "applied %d body elements to %s\n", len(d.Body), docID)
	return err
}

func trackToWriter(authMode, docID, out string, dry, force bool, stdout, stderr io.Writer) error {
	return gsync.Track(activeFetcher, gsync.Options{
		AuthMode: authMode,
		DocID:    docID,
		OutDir:   out,
		DryRun:   dry,
		Force:    force,
		Stdout:   stdout,
		Stderr:   stderr,
	})
}
