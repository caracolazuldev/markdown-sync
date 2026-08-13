package sync

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
)

const trackedFileMode fs.FileMode = 0444
const writableFileMode fs.FileMode = 0644

// Options configures a track run.
type Options struct {
	AuthMode string
	DocID    string
	OutDir   string
	Force    bool
	DryRun   bool
	Stdout   io.Writer
	Stderr   io.Writer
}

// Track pulls a tabbed (or single-tab) Doc into OutDir.
func Track(f md.Fetcher, opts Options) error {
	if opts.Stdout == nil {
		opts.Stdout = io.Discard
	}
	if opts.Stderr == nil {
		opts.Stderr = io.Discard
	}
	if opts.OutDir == "" {
		return fmt.Errorf("track requires -out <directory>")
	}
	if f == nil {
		return fmt.Errorf("missing fetcher")
	}

	root, err := filepath.Abs(opts.OutDir)
	if err != nil {
		return err
	}

	kind, existing, err := classifyOut(root)
	if err != nil {
		return err
	}

	docID := opts.DocID
	switch kind {
	case outIsFile:
		return ErrOutNotDirectory
	case outUntracked:
		if !opts.Force {
			return ErrNotTrackRoot
		}
	case outTrackRoot:
		if !isTrackSidecar(existing) {
			return ErrNotTrackRoot
		}
		if docID == "" {
			docID = existing.DocID
		} else if docID != existing.DocID {
			return ErrDocIDMismatch
		}
	case outMissing, outEmpty:
		if docID == "" {
			return fmt.Errorf("track requires -doc <doc-id> on first run")
		}
	}

	fd, err := f.Fetch(context.Background(), opts.AuthMode, docID)
	if err != nil {
		return fmt.Errorf("fetch document: %w", err)
	}
	if fd == nil {
		return fmt.Errorf("fetch document: empty result")
	}

	tabs := fd.Tabs
	if len(tabs) == 0 && fd.Document != nil {
		tabs = []*md.Tab{{ID: "legacy", Title: fd.Title, Document: fd.Document}}
	}
	mapped := MapTabTree(docID, tabs)
	title := fd.Title
	if title == "" {
		title = docID
	}

	if kind == outTrackRoot && !opts.Force {
		if dirty := dirtyFiles(root, existing); len(dirty) > 0 {
			fmt.Fprintf(opts.Stderr, "dirty files: %s\n", strings.Join(dirty, ", "))
			return ErrDirtyTrack
		}
	}

	if opts.DryRun {
		for _, m := range mapped {
			fmt.Fprintf(opts.Stdout, "dry-run: would write %s\n", m.RelPath)
		}
		fmt.Fprintf(opts.Stdout, "dry-run: would write %s\n", sidecarName)
		return nil
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	byID := map[string]SidecarTab{}
	for _, t := range existing.Tabs {
		byID[t.TabID] = t
	}

	newSidecar := Sidecar{Mode: "track", DocID: docID, Title: title}
	written := map[string]bool{}
	for _, m := range mapped {
		markdown, err := md.DocumentToMarkdown(m.Doc)
		if err != nil {
			return fmt.Errorf("convert %s: %w", m.TabID, err)
		}
		abs := filepath.Join(root, filepath.FromSlash(m.RelPath))
		if old, ok := byID[m.TabID]; ok && old.Path != "" && old.Path != m.RelPath {
			oldAbs := filepath.Join(root, filepath.FromSlash(old.Path))
			if err := removeTrackedPath(oldAbs); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove old path %s: %w", old.Path, err)
			}
			cleanupEmptyParents(root, oldAbs)
		}
		if err := writeReadOnlyFile(abs, []byte(markdown)); err != nil {
			return err
		}
		sum := sha256.Sum256([]byte(markdown))
		newSidecar.Tabs = append(newSidecar.Tabs, SidecarTab{
			TabID:  m.TabID,
			Path:   m.RelPath,
			Title:  m.Title,
			SHA256: fmt.Sprintf("%x", sum),
		})
		written[m.TabID] = true
		fmt.Fprintf(opts.Stdout, "wrote %s\n", m.RelPath)
	}

	for _, old := range existing.Tabs {
		if written[old.TabID] {
			continue
		}
		fmt.Fprintf(opts.Stderr, "warning: local file %s is no longer a remote tab; leaving in place\n", old.Path)
		newSidecar.Tabs = append(newSidecar.Tabs, old)
	}

	if err := writeSidecar(root, newSidecar); err != nil {
		return err
	}
	fmt.Fprintf(opts.Stdout, "wrote %s\n", sidecarName)
	return nil
}

type outKind int

const (
	outMissing outKind = iota
	outEmpty
	outTrackRoot
	outUntracked
	outIsFile
)

func classifyOut(root string) (outKind, Sidecar, error) {
	st, err := os.Stat(root)
	if os.IsNotExist(err) {
		return outMissing, Sidecar{}, nil
	}
	if err != nil {
		return 0, Sidecar{}, err
	}
	if !st.IsDir() {
		return outIsFile, Sidecar{}, nil
	}
	sc, err := readSidecar(root)
	if err == nil && isTrackSidecar(sc) {
		return outTrackRoot, sc, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return 0, Sidecar{}, err
	}
	empty, err := dirEffectivelyEmpty(root)
	if err != nil {
		return 0, Sidecar{}, err
	}
	if empty {
		return outEmpty, Sidecar{}, nil
	}
	return outUntracked, Sidecar{}, nil
}

func dirEffectivelyEmpty(root string) (bool, error) {
	var foundFile bool
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == root {
			return nil
		}
		if !d.IsDir() {
			foundFile = true
			return fs.SkipAll
		}
		return nil
	})
	return !foundFile, err
}

func dirtyFiles(root string, sc Sidecar) []string {
	var dirty []string
	for _, t := range sc.Tabs {
		if t.Path == "" || t.SHA256 == "" {
			continue
		}
		p := filepath.Join(root, filepath.FromSlash(t.Path))
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		sum := fmt.Sprintf("%x", sha256.Sum256(data))
		if sum != t.SHA256 {
			dirty = append(dirty, t.Path)
		}
	}
	return dirty
}

func writeReadOnlyFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	_ = os.Chmod(path, writableFileMode)
	if err := os.WriteFile(path, data, writableFileMode); err != nil {
		return err
	}
	return os.Chmod(path, trackedFileMode)
}

func removeTrackedPath(path string) error {
	_ = os.Chmod(path, writableFileMode)
	return os.Remove(path)
}

func cleanupEmptyParents(root, removedFile string) {
	dir := filepath.Dir(removedFile)
	for dir != root && strings.HasPrefix(dir, root) {
		ents, err := os.ReadDir(dir)
		if err != nil || len(ents) > 0 {
			return
		}
		if err := os.Remove(dir); err != nil {
			return
		}
		dir = filepath.Dir(dir)
	}
}
