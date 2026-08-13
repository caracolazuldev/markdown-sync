package sync

import (
	"fmt"
	"os"
	"path/filepath"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
)

// CheckImportAllowed refuses import of tracked files or trees.
func CheckImportAllowed(file, docID string) error {
	if file == "" {
		return nil
	}
	abs, err := filepath.Abs(file)
	if err != nil {
		abs = file
	}
	data, err := os.ReadFile(abs)
	if err == nil {
		_, fm := md.SplitFrontMatter(string(data))
		if fm.Track {
			return ErrTrackedPath
		}
	}
	dir := abs
	if st, err := os.Stat(abs); err == nil && !st.IsDir() {
		dir = filepath.Dir(abs)
	}
	for {
		sc, err := readSidecar(dir)
		if err == nil && isTrackSidecar(sc) {
			return ErrTrackedPath
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil
}

// TabbedHint wraps ErrTabbedDocument with a track command suggestion.
func TabbedHint(docID string, tabCount int) error {
	return fmt.Errorf("%w: document %s has %d tabs; use: gdocs-markdown-sync track -doc %s -out <dir>", ErrTabbedDocument, docID, tabCount, docID)
}
