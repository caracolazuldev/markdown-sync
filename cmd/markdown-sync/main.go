package main

import (
    "flag"
    "fmt"
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
        if *doc == "" {
            fmt.Fprintln(os.Stderr, "export requires -doc <doc-id>")
            os.Exit(2)
        }
        docModel, err := md.FetchDocument(*auth, *doc)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failed to fetch doc: %v\n", err)
            os.Exit(1)
        }
        markdown, err := md.DocumentToMarkdown(docModel)
        if err != nil {
            fmt.Fprintf(os.Stderr, "conversion error: %v\n", err)
            os.Exit(1)
        }
        if *out == "" {
            fmt.Print(markdown)
            return
        }
        if *dry {
            fmt.Printf("dry-run: would write %d bytes to %s\n", len(markdown), *out)
            return
        }
        if err := ioutil.WriteFile(*out, []byte(markdown), 0644); err != nil {
            fmt.Fprintf(os.Stderr, "failed to write output: %v\n", err)
            os.Exit(1)
        }
        fmt.Fprintf(os.Stdout, "wrote %d bytes to %s\n", len(markdown), *out)
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
