package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nginx-autoindex-tui/internal/parser"
)

func TestParse(t *testing.T) {
	f, err := os.Open(filepath.Join("fixtures", "sample_autoindex.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	page, err := parser.Parse(f)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if !strings.Contains(page.Title, "Index of") {
		t.Errorf("title missing 'Index of', got %q", page.Title)
	}

	wantNames := []string{"../", "docs/", "file1.txt", "image.png"}
	if len(page.Entries) < len(wantNames) {
		t.Fatalf("want at least %d entries, got %d", len(wantNames), len(page.Entries))
	}

	for i, want := range wantNames {
		if page.Entries[i].Name != want {
			t.Errorf("entries[%d].Name = %q, want %q", i, page.Entries[i].Name, want)
		}
	}

	// directory entries should be marked IsDir (except ../)
	if !page.Entries[1].IsDir {
		t.Errorf("entries[1].IsDir = false, expected true (docs/)")
	}
	// files should not be marked as directory
	if page.Entries[2].IsDir {
		t.Errorf("entries[2].IsDir = true, expected false (file1.txt)")
	}
	if page.Entries[2].Size == "" {
		t.Errorf("entries[2].Size is empty, expected something like '1.2K'")
	}
	if page.Entries[2].DateTime == "" {
		t.Errorf("entries[2].DateTime is empty")
	}
}
