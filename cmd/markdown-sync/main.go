package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	md "github.com/caracolazuldev/markdown-sync/internal/markdown"
)

func usage() {
	fmt.Fprintf(os.Stderr, "markdown-sync: simple CLI\n")
	fmt.Fprintf(os.Stderr, "Usage:\n  markdown-sync <command> [flags]\nCommands: export, import, preview, list\n")
}

func main() {
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
	fs.Parse(os.Args[2:])

	switch cmd {
	case "export":
		if err := exportToWriter(*auth, *doc, *out, *dry, os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "export error: %v\n", err)
			os.Exit(1)
		}
	case "import":
		fmt.Printf("import: file=%s doc=%s auth=%s dry=%v\n", *file, *doc, *auth, *dry)
	case "preview":
		fmt.Printf("preview: doc=%s auth=%s\n", *doc, *auth)
	case "list":
		fmt.Printf("list: auth=%s\n", *auth)
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
	docModel, err := md.FetchDocument(authMode, docID)
	if err != nil {
		fmt.Fprintf(stderr, "failed to fetch doc: %v\n", err)
		return err
	}
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
