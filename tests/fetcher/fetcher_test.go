package fetcher_test

import (
	"testing"

	"nginx-autoindex-tui/internal/fetcher"
)

func TestResolveURL(t *testing.T) {
	cases := []struct {
		base, href, want string
	}{
		// 标准导航
		{"http://example.com/files/", "docs/", "http://example.com/files/docs/"},
		{"http://example.com/files/", "../", "http://example.com/"},
		{"http://example.com/files/", "file.txt", "http://example.com/files/file.txt"},
		{"http://example.com/files/", "/other.txt", "http://example.com/other.txt"},

		// 多级导航：带尾部 / 时逐级深入
		{"http://example.com/files/", "a/", "http://example.com/files/a/"},
		{"http://example.com/files/a/", "b/", "http://example.com/files/a/b/"},
		{"http://example.com/files/a/b/", "../", "http://example.com/files/a/"},
		{"http://example.com/files/a/b/", "../../", "http://example.com/files/"},

		// JSON 目录名称无尾部 /，需要外部补上
		{"http://example.com/json/", "subdir/", "http://example.com/json/subdir/"},
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

// TestResolveURLNoTrailingSlash 验证 base URL 缺少尾部 / 时解析错误——这就是
// currentURL 必须归一化的原因。测试确保未来不会有人删掉归一化逻辑。
func TestResolveURLNoTrailingSlash(t *testing.T) {
	// 不带尾部 / → 最后一段被当文件，子目录解析到错误路径
	got, _ := fetcher.ResolveURL("http://example.com/files", "docs/")
	want := "http://example.com/docs/" // 丢失了 files 段！
	if got != want {
		t.Errorf("ResolveURL without trailing / = %q, want %q (proves normalization matters)", got, want)
	}
}
