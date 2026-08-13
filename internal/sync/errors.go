package sync

import "errors"

var (
	// ErrTabbedDocument is returned by export/preview when the Doc has multiple or nested tabs.
	ErrTabbedDocument = errors.New("document has multiple tabs; use track")
	// ErrTrackedPath is returned by import when the file is under a track root or marked track: true.
	ErrTrackedPath = errors.New("path is tracked; use track not import")
	// ErrDirtyTrack is returned when local tracked files changed since the last track.
	ErrDirtyTrack = errors.New("local tracked files have been modified; pass --force to overwrite")
	// ErrNotTrackRoot is returned when -out has content but no _track.toml.
	ErrNotTrackRoot = errors.New("directory is not a track root; use a new directory or --force to adopt")
	// ErrDocIDMismatch is returned when -doc does not match the sidecar (even with --force).
	ErrDocIDMismatch = errors.New("track root doc_id does not match -doc")
	// ErrOutNotDirectory is returned when -out exists as a file.
	ErrOutNotDirectory = errors.New("track -out must be a directory")
)
