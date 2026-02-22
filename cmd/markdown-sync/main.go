package main

import (
    "flag"
    "fmt"
    "os"
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
        fmt.Printf("export: doc=%s out=%s auth=%s dry=%v\n", *doc, *out, *auth, *dry)
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
