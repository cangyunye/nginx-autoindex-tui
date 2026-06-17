package parser

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Page 表示一次 autoindex 页面解析结果。
type Page struct {
	Title   string
	Entries []Entry
}

// Entry 表示一行条目。
type Entry struct {
	Href     string
	Name     string
	DateTime string
	Size     string
	IsDir    bool
}

// Parse 从 reader 读取 HTML，返回 Page。
func Parse(r io.Reader) (*Page, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	page := &Page{
		Title: strings.TrimSpace(doc.Find("h1").First().Text()),
	}

	// 找到 h1 之后的第一个 <pre>
	pre := doc.Find("h1").First().NextAllFiltered("pre").First()
	if pre.Length() == 0 {
		// 回退：直接找文档中第一个 <pre>
		pre = doc.Find("pre").First()
	}
	if pre.Length() == 0 {
		return page, nil
	}

	// 遍历 pre 内部所有 <a>
	pre.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		name := strings.TrimSpace(s.Text())

		// 从 <a> 之后的相邻文本节点里拆出 date time 与 size
		tail := extractTailAfter(s)
		fields := strings.Fields(tail)
		var dateTime, size string
		if len(fields) >= 2 {
			size = fields[len(fields)-1]
			dateTime = fields[len(fields)-2]
			// 若倒数第二列像 "HH:MM"，且倒数第三列形如日期，则合并
			if len(fields) >= 3 && looksLikeTime(fields[len(fields)-2]) {
				dateTime = fields[len(fields)-3] + " " + fields[len(fields)-2]
			}
		} else if len(fields) == 1 {
			size = fields[0]
		}

		page.Entries = append(page.Entries, Entry{
			Href:     href,
			Name:     name,
			DateTime: dateTime,
			Size:     size,
			IsDir:    strings.HasSuffix(href, "/"),
		})
	})

	return page, nil
}

// extractTailAfter 提取选择器 s 所在 <a> 节点之后、同一父节点内后续的文本内容。
func extractTailAfter(s *goquery.Selection) string {
	if s.Length() == 0 {
		return ""
	}
	var sb strings.Builder
	node := s.Get(0)
	for n := node.NextSibling; n != nil; n = n.NextSibling {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		} else if n.Type == html.ElementNode {
			// 遇到非文本元素就停止（比如后面还有 <br> 或别的标签）
			break
		}
	}
	return sb.String()
}

// jsonEntry 对应 nginx autoindex JSON 格式的一个元素。
type jsonEntry struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Mtime string `json:"mtime"`
	Size  string `json:"size,omitempty"`
}

// ParseJSON 从 reader 读取 JSON 格式的 autoindex 列表，返回 Page。
func ParseJSON(r io.Reader) (*Page, error) {
	var entries []jsonEntry
	if err := json.NewDecoder(r).Decode(&entries); err != nil {
		return nil, err
	}
	page := &Page{
		Title: "Index",
	}
	for _, e := range entries {
		isDir := e.Type == "directory" || strings.HasSuffix(e.Name, "/")
		href := e.Name
		// 目录 href 必须有尾部 /，使后续 url.ResolveReference 行为一致
		if isDir && !strings.HasSuffix(href, "/") {
			href += "/"
		}
		page.Entries = append(page.Entries, Entry{
			Href:     href,
			Name:     e.Name,
			DateTime: e.Mtime,
			Size:     e.Size,
			IsDir:    isDir,
		})
	}
	return page, nil
}

// ParseAuto 根据 Content-Type 自动选择解析方式，不明确时先试 JSON 再回退 HTML。
func ParseAuto(r io.Reader, contentType string) (*Page, error) {
	if strings.Contains(contentType, "application/json") {
		return ParseJSON(r)
	}
	if strings.Contains(contentType, "text/html") {
		return Parse(r)
	}
	// Content-Type 不明确：先试 JSON，失败后回退 HTML
	page, err := ParseJSON(r)
	if err == nil {
		return page, nil
	}
	// 重置 reader 重试 HTML 解析
	if rs, ok := r.(io.Seeker); ok {
		_, _ = rs.Seek(0, io.SeekStart)
	}
	return Parse(r)
}

// looksLikeTime 粗略判断是否形如 "HH:MM"。
func looksLikeTime(s string) bool {
	if len(s) != 5 || s[2] != ':' {
		return false
	}
	for i, r := range s {
		if i == 2 {
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
