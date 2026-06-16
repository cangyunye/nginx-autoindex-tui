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

// FetchHTML 对 rawURL 发起 GET，返回响应体文本。
// 当 rawURL 缺失 scheme 时，默认补 http://。
func FetchHTML(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid url %q: %w", rawURL, err)
	}
	resp, err := client.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("get %s: %w", u.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("status %s from %s", resp.Status, u.String())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	return string(body), nil
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
