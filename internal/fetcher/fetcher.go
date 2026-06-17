package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}

// FetchResult 包含 HTTP 响应体与 Content-Type。
type FetchResult struct {
	Body        string
	ContentType string
}

// Fetch 对 rawURL 发起 GET，返回响应体文本与 Content-Type。
// 当 rawURL 缺失 scheme 时，默认补 http://。
func Fetch(rawURL string) (*FetchResult, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url %q: %w", rawURL, err)
	}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", u.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %s from %s", resp.Status, u.String())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return &FetchResult{
		Body:        string(body),
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
}

// FetchHTML 是 Fetch 的便捷包装，只返回 body（兼容旧调用）。
func FetchHTML(rawURL string) (string, error) {
	res, err := Fetch(rawURL)
	if err != nil {
		return "", err
	}
	return res.Body, nil
}

// ResolveURL 基于 baseURL 把 href 解析为绝对 URL。
func ResolveURL(baseURL, href string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(rel).String(), nil
}
