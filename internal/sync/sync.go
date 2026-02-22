package sync

import "context"

// Orchestrator for export/import workflows. This package wires together auth,
// google adapter and markdown conversion.

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Export(ctx context.Context, docID string, outDir string) error {
    // TODO: fetch doc, convert to markdown, write files
    return nil
}

func (s *Service) Import(ctx context.Context, file string, docID string) error {
    // TODO: read file, convert to docs model, push to Google Docs
    return nil
}
