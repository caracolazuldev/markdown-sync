package sync

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	md "github.com/caracolazuldev/gdocs-markdown-sync/internal/markdown"
)

func nestedFetched(docID string) *md.FetchedDocument {
	return &md.FetchedDocument{
		ID:    docID,
		Title: "Minutes",
		Tabs: []*md.Tab{
			{ID: "w", Title: "Welcome", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "hello"}}}},
			{ID: "fy", Title: "FY2024", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "year"}}}, Children: []*md.Tab{
				{ID: "jan", Title: "Jan", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "january"}}}},
			}},
		},
	}
}

func TestTrack_WritesNestedReadOnlyTree(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "standing")
	var stdout, stderr bytes.Buffer
	err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{
		DocID:  "doc-1",
		OutDir: out,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		t.Fatalf("track: %v stderr=%s", err, stderr.String())
	}
	welcome := filepath.Join(out, "welcome.md")
	parent := filepath.Join(out, "fy2024", "fy2024.md")
	jan := filepath.Join(out, "fy2024", "jan.md")
	for _, p := range []string{welcome, parent, jan, sidecarPath(out)} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing %s: %v", p, err)
		}
	}
	st, err := os.Stat(welcome)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0444 {
		t.Fatalf("welcome mode=%o want 0444", st.Mode().Perm())
	}
	dst, err := os.Stat(filepath.Join(out, "fy2024"))
	if err != nil {
		t.Fatal(err)
	}
	if dst.Mode().Perm()&0200 == 0 {
		t.Fatalf("dir should be writable, mode=%o", dst.Mode().Perm())
	}
	sc, err := readSidecar(out)
	if err != nil {
		t.Fatal(err)
	}
	if sc.Mode != "track" || sc.DocID != "doc-1" || len(sc.Tabs) != 3 {
		t.Fatalf("sidecar %+v", sc)
	}
	data, _ := os.ReadFile(welcome)
	if !strings.Contains(string(data), "track: true") || !strings.Contains(string(data), "tab_id: w") {
		t.Fatalf("welcome content %s", data)
	}
}

func TestTrack_DryRunNoWrite(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "standing")
	var stdout bytes.Buffer
	err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{
		DocID:  "doc-1",
		OutDir: out,
		DryRun: true,
		Stdout: &stdout,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("dry-run should not create out: %v", err)
	}
	if !strings.Contains(stdout.String(), "dry-run: would write") {
		t.Fatalf("stdout=%q", stdout.String())
	}
}

func TestTrack_UntrackedDirRefused(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "notes.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{DocID: "doc-1", OutDir: dir})
	if !errors.Is(err, ErrNotTrackRoot) {
		t.Fatalf("err=%v", err)
	}
}

func TestTrack_ForceAdoptsUntracked(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "notes.md"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{DocID: "doc-1", OutDir: dir, Force: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sidecarPath(dir)); err != nil {
		t.Fatal(err)
	}
}

func TestTrack_DocIDMismatchEvenWithForce(t *testing.T) {
	dir := t.TempDir()
	if err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{DocID: "doc-1", OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	err := Track(md.FakeFetcher{Doc: nestedFetched("doc-2")}, Options{DocID: "doc-2", OutDir: dir, Force: true})
	if !errors.Is(err, ErrDocIDMismatch) {
		t.Fatalf("err=%v", err)
	}
}

func TestTrack_OutFileRefused(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "file.md")
	if err := os.WriteFile(p, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{DocID: "doc-1", OutDir: p})
	if !errors.Is(err, ErrOutNotDirectory) {
		t.Fatalf("err=%v", err)
	}
}

func TestTrack_RefreshChmodWriteChmod(t *testing.T) {
	dir := t.TempDir()
	f := md.FakeFetcher{Doc: nestedFetched("doc-1")}
	if err := Track(f, Options{DocID: "doc-1", OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	updated := nestedFetched("doc-1")
	updated.Tabs[0].Document.Body = []md.Element{md.Paragraph{Text: "updated hello"}}
	f.Doc = updated
	if err := Track(f, Options{OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "welcome.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "updated hello") {
		t.Fatalf("content %s", data)
	}
	st, _ := os.Stat(filepath.Join(dir, "welcome.md"))
	if st.Mode().Perm() != 0444 {
		t.Fatalf("mode=%o", st.Mode().Perm())
	}
}

func TestTrack_DirtyRequiresForce(t *testing.T) {
	dir := t.TempDir()
	if err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{DocID: "doc-1", OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "welcome.md")
	_ = os.Chmod(p, 0644)
	if err := os.WriteFile(p, []byte("local edit\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_ = os.Chmod(p, 0444)
	err := Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{OutDir: dir})
	if !errors.Is(err, ErrDirtyTrack) {
		t.Fatalf("err=%v", err)
	}
	err = Track(md.FakeFetcher{Doc: nestedFetched("doc-1")}, Options{OutDir: dir, Force: true})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(p)
	if strings.Contains(string(data), "local edit") {
		t.Fatal("force should overwrite dirty file")
	}
	st, _ := os.Stat(p)
	if st.Mode().Perm() != 0444 {
		t.Fatalf("mode=%o", st.Mode().Perm())
	}
}

func TestTrack_PromoteAndFlatten(t *testing.T) {
	dir := t.TempDir()
	leaf := &md.FetchedDocument{
		ID:    "doc-1",
		Title: "D",
		Tabs:  []*md.Tab{{ID: "p", Title: "Parent", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "p"}}}}},
	}
	if err := Track(md.FakeFetcher{Doc: leaf}, Options{DocID: "doc-1", OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "parent.md")); err != nil {
		t.Fatal(err)
	}
	nested := &md.FetchedDocument{
		ID:    "doc-1",
		Title: "D",
		Tabs: []*md.Tab{{ID: "p", Title: "Parent", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "p"}}}, Children: []*md.Tab{
			{ID: "c", Title: "Child", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "c"}}}},
		}}},
	}
	if err := Track(md.FakeFetcher{Doc: nested}, Options{OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "parent", "parent.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "parent.md")); !os.IsNotExist(err) {
		t.Fatal("expected leaf parent.md to be removed after promote")
	}
	if err := Track(md.FakeFetcher{Doc: leaf}, Options{OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "parent.md")); err != nil {
		t.Fatal(err)
	}
}

func TestTrack_RemovedTabLeftInPlace(t *testing.T) {
	dir := t.TempDir()
	full := nestedFetched("doc-1")
	if err := Track(md.FakeFetcher{Doc: full}, Options{DocID: "doc-1", OutDir: dir}); err != nil {
		t.Fatal(err)
	}
	onlyWelcome := &md.FetchedDocument{
		ID:    "doc-1",
		Title: "Minutes",
		Tabs:  []*md.Tab{{ID: "w", Title: "Welcome", Document: &md.Document{Body: []md.Element{md.Paragraph{Text: "hello"}}}}},
	}
	var stderr bytes.Buffer
	if err := Track(md.FakeFetcher{Doc: onlyWelcome}, Options{OutDir: dir, Stderr: &stderr}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "fy2024", "jan.md")); err != nil {
		t.Fatal("removed remote tab should keep local file")
	}
	if !strings.Contains(stderr.String(), "no longer a remote tab") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestSidecarRoundTrip(t *testing.T) {
	dir := t.TempDir()
	sc := Sidecar{Mode: "track", DocID: "d", Title: "T", Tabs: []SidecarTab{{TabID: "a", Path: "a.md", Title: "A", SHA256: "abc"}}}
	if err := writeSidecar(dir, sc); err != nil {
		t.Fatal(err)
	}
	got, err := readSidecar(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.DocID != "d" || len(got.Tabs) != 1 || got.Tabs[0].SHA256 != "abc" {
		t.Fatalf("%+v", got)
	}
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte("x")))
	if len(sum) != 64 {
		t.Fatalf("sha len %d", len(sum))
	}
}
