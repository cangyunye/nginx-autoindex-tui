package fetcher_test

import (
	"testing"

	"nginx-autoindex-tui/internal/fetcher"
)

func TestResolveURL(t *testing.T) {
	cases := []struct {
		base, href, want string
	}{
		{"http://example.com/files/", "docs/", "http://example.com/files/docs/"},
		{"http://example.com/files/", "../", "http://example.com/"},
		{"http://example.com/files/", "file.txt", "http://example.com/files/file.txt"},
		{"http://example.com/files/", "/other.txt", "http://example.com/other.txt"},
	}
	for _, c := range cases {
		got, err := fetcher.ResolveURL(c.base, c.href)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if got != c.want {
			t.Errorf("ResolveURL(%q, %q) = %q, want %q", c.base, c.href, got, c.want)
		}
	}
}
