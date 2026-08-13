package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const sidecarName = "_track.toml"

// Sidecar is the track-root metadata file.
type Sidecar struct {
	Mode  string
	DocID string
	Title string
	Tabs  []SidecarTab
}

// SidecarTab records one written tab file.
type SidecarTab struct {
	TabID  string
	Path   string
	Title  string
	SHA256 string
}

func sidecarPath(root string) string {
	return filepath.Join(root, sidecarName)
}

func writeSidecar(root string, sc Sidecar) error {
	sc.Mode = "track"
	var b strings.Builder
	b.WriteString(fmt.Sprintf("mode = %q\n", sc.Mode))
	b.WriteString(fmt.Sprintf("doc_id = %q\n", sc.DocID))
	b.WriteString(fmt.Sprintf("title = %q\n", sc.Title))
	for _, t := range sc.Tabs {
		b.WriteString("\n[[tabs]]\n")
		b.WriteString(fmt.Sprintf("tab_id = %q\n", t.TabID))
		b.WriteString(fmt.Sprintf("path = %q\n", t.Path))
		b.WriteString(fmt.Sprintf("title = %q\n", t.Title))
		b.WriteString(fmt.Sprintf("sha256 = %q\n", t.SHA256))
	}
	return os.WriteFile(sidecarPath(root), []byte(b.String()), 0644)
}

func readSidecar(root string) (Sidecar, error) {
	data, err := os.ReadFile(sidecarPath(root))
	if err != nil {
		return Sidecar{}, err
	}
	return parseSidecar(string(data))
}

func parseSidecar(s string) (Sidecar, error) {
	var sc Sidecar
	var cur *SidecarTab
	flush := func() {
		if cur != nil {
			sc.Tabs = append(sc.Tabs, *cur)
			cur = nil
		}
	}
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "[[tabs]]" {
			flush()
			cur = &SidecarTab{}
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = unquoteTOML(strings.TrimSpace(val))
		if cur == nil {
			switch key {
			case "mode":
				sc.Mode = val
			case "doc_id":
				sc.DocID = val
			case "title":
				sc.Title = val
			}
			continue
		}
		switch key {
		case "tab_id":
			cur.TabID = val
		case "path":
			cur.Path = val
		case "title":
			cur.Title = val
		case "sha256":
			cur.SHA256 = val
		}
	}
	flush()
	return sc, nil
}

func unquoteTOML(v string) string {
	if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
		inner := v[1 : len(v)-1]
		inner = strings.ReplaceAll(inner, `\"`, `"`)
		return inner
	}
	return v
}

func isTrackSidecar(sc Sidecar) bool {
	return sc.Mode == "track" && sc.DocID != ""
}
