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

func TestParseJSON(t *testing.T) {
	f, err := os.Open(filepath.Join("fixtures", "sample_autoindex.json"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	page, err := parser.ParseJSON(f)
	if err != nil {
		t.Fatalf("ParseJSON returned error: %v", err)
	}

	wantNames := []string{"../", "docs", "file1.txt", "image.png"}
	if len(page.Entries) < len(wantNames) {
		t.Fatalf("want at least %d entries, got %d", len(wantNames), len(page.Entries))
	}

	for i, want := range wantNames {
		if page.Entries[i].Name != want {
			t.Errorf("entries[%d].Name = %q, want %q", i, page.Entries[i].Name, want)
		}
	}

	// directory entries should be marked IsDir
	if !page.Entries[1].IsDir {
		t.Errorf("entries[1].IsDir = false, expected true (docs, type=directory)")
	}
	// directory href should end with /
	if page.Entries[1].Href != "docs/" {
		t.Errorf("entries[1].Href = %q, want %q", page.Entries[1].Href, "docs/")
	}
	// ../ also has / suffix
	if page.Entries[0].Href != "../" {
		t.Errorf("entries[0].Href = %q, want %q", page.Entries[0].Href, "../")
	}
	// file should not be directory
	if page.Entries[2].IsDir {
		t.Errorf("entries[2].IsDir = true, expected false (file1.txt)")
	}
	// file should have size
	if page.Entries[2].Size != "1229" {
		t.Errorf("entries[2].Size = %q, want %q", page.Entries[2].Size, "1229")
	}
	// directory should have mtime
	if page.Entries[1].DateTime == "" {
		t.Errorf("entries[1].DateTime is empty")
	}
	// href should match name
	if page.Entries[2].Href != "file1.txt" {
		t.Errorf("entries[2].Href = %q, want %q", page.Entries[2].Href, "file1.txt")
	}

	// file with numeric size should parse without error
	numericJSON := `[{"name":"f.bin","type":"file","mtime":"...","size":63969}]`
	r := strings.NewReader(numericJSON)
	page2, err := parser.ParseJSON(r)
	if err != nil {
		t.Fatalf("ParseJSON with numeric size returned error: %v", err)
	}
	if len(page2.Entries) != 1 || page2.Entries[0].Size != "63969" {
		t.Errorf("numeric size: got Size=%q, want %q", page2.Entries[0].Size, "63969")
	}
}

func TestParseAuto(t *testing.T) {
	t.Run("json content type", func(t *testing.T) {
		r := strings.NewReader(`[{"name":"f.txt","type":"file","mtime":"...","size":"99"}]`)
		page, err := parser.ParseAuto(r, "application/json")
		if err != nil {
			t.Fatalf("ParseAuto(json) returned error: %v", err)
		}
		if len(page.Entries) != 1 || page.Entries[0].Name != "f.txt" {
			t.Errorf("unexpected entries: %+v", page.Entries)
		}
	})

	t.Run("json numeric size", func(t *testing.T) {
		r := strings.NewReader(`[{"name":"f.bin","type":"file","mtime":"...","size":12345}]`)
		page, err := parser.ParseAuto(r, "application/json")
		if err != nil {
			t.Fatalf("ParseAuto(json numeric size) returned error: %v", err)
		}
		if page.Entries[0].Size != "12345" {
			t.Errorf("numeric size: got %q, want %q", page.Entries[0].Size, "12345")
		}
	})

	t.Run("html content type", func(t *testing.T) {
		html := `<html><body><h1>Index</h1><pre><a href="f.txt">f.txt</a> 01-Jan-2024 10:00 1K</pre></body></html>`
		r := strings.NewReader(html)
		page, err := parser.ParseAuto(r, "text/html")
		if err != nil {
			t.Fatalf("ParseAuto(html) returned error: %v", err)
		}
		if len(page.Entries) != 1 || page.Entries[0].Name != "f.txt" {
			t.Errorf("unexpected entries: %+v", page.Entries)
		}
	})

	t.Run("unclear falls back to json then html", func(t *testing.T) {
		// valid JSON should succeed
		r := strings.NewReader(`[{"name":"a.json","type":"file","mtime":"..."}]`)
		page, err := parser.ParseAuto(r, "")
		if err != nil {
			t.Fatalf("ParseAuto(unclear→json) returned error: %v", err)
		}
		if len(page.Entries) != 1 || page.Entries[0].Name != "a.json" {
			t.Errorf("unexpected entries: %+v", page.Entries)
		}
	})
}
