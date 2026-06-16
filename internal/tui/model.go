package tui

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"nginx-autoindex-tui/internal/fetcher"
	"nginx-autoindex-tui/internal/parser"
)

// fetchMsg 用于在 Cmd 中异步返回解析结果或错误。
type fetchMsg struct {
	page *parser.Page
	url  string
	err  error
}

type execMsg struct {
	output string
	err    error
}

// Model 持有 TUI 的全部状态。
type Model struct {
	currentURL string
	rootURL    string
	page       *parser.Page
	cursor     int
	loading    bool
	err        string
	lastOutput string // wget 执行的摘要输出
}

func NewModel(initialURL string) *Model {
	return &Model{currentURL: initialURL, rootURL: initialURL, loading: true}
}

// ---- Elm hooks ----

func (m *Model) Init() tea.Cmd {
	return fetchCmd(m.currentURL)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.page != nil && m.cursor < len(m.page.Entries)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, m.openSelected()
		case "r":
			m.loading = true
			m.err = ""
			return m, fetchCmd(m.currentURL)
		}
	case fetchMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.currentURL = msg.url
		m.page = msg.page
		m.cursor = 0
		m.lastOutput = ""
	case execMsg:
		if msg.err != nil {
			m.lastOutput = "exec failed: " + msg.err.Error()
		} else {
			m.lastOutput = msg.output
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.pageOrFallbackTitle()) + "\n")
	b.WriteString(metaStyle.Render("URL: "+m.currentURL) + "\n\n")

	if m.loading {
		b.WriteString("  loading...\n")
		b.WriteString(helpStyle.Render("(ctrl+c 退出)") + "\n")
		return b.String()
	}
	if m.err != "" {
		b.WriteString(errStyle.Render("ERROR: "+m.err) + "\n")
		b.WriteString(helpStyle.Render("按 R 重试，ctrl+c 退出") + "\n")
		return b.String()
	}
	if m.page == nil || len(m.page.Entries) == 0 {
		b.WriteString("  (no entries)\n")
		b.WriteString(helpStyle.Render("R 刷新，ctrl+c 退出") + "\n")
		return b.String()
	}

	// 构造表格数据
	rows := make([][]string, 0, len(m.page.Entries))
	for i, e := range m.page.Entries {
		idx := fmt.Sprintf("%d", i+1)
		kind := "FILE"
		if e.IsDir {
			kind = "DIR "
		}
		rows = append(rows, []string{idx, kind, e.Name, e.DateTime, e.Size})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))).
		Headers("Idx", "Type", "Name", "Date/Time", "Size").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == -1 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FEF3C7")).Bold(true).Padding(0, 1)
			}
			base := lipgloss.NewStyle().Padding(0, 1)
			if row == m.cursor {
				base = base.Foreground(lipgloss.Color("#FFCC66")).Background(lipgloss.Color("#3C3C3C")).Bold(true)
			}
			if col == 1 && rows[row][1] == "DIR " {
				base = base.Foreground(lipgloss.Color("#66DDAA"))
			}
			return base
		})

	b.WriteString(t.String())
	b.WriteString("\n")

	if m.lastOutput != "" {
		b.WriteString(metaStyle.Render("--- wget output ---") + "\n")
		b.WriteString(m.lastOutput + "\n")
	}

	b.WriteString(helpStyle.Render("↑/↓ or j/k 移动 · Enter 进入目录/下载文件 · R 刷新 · q/ctrl+c 退出") + "\n")
	return b.String()
}

// ---- helpers ----

func (m *Model) pageOrFallbackTitle() string {
	if m.page != nil && m.page.Title != "" {
		return m.page.Title
	}
	return "Nginx Autoindex Browser"
}

// openSelected 决定当前选中条目是目录进入还是文件下载，并返回对应 Cmd。
func (m *Model) openSelected() tea.Cmd {
	if m.page == nil || m.cursor < 0 || m.cursor >= len(m.page.Entries) {
		return nil
	}
	e := m.page.Entries[m.cursor]
	absURL, err := fetcher.ResolveURL(m.currentURL, e.Href)
	if err != nil {
		return func() tea.Msg { return execMsg{err: err} }
	}
	if e.IsDir {
		m.loading = true
		m.err = ""
		return fetchCmd(absURL)
	}
	return execWgetCmd(absURL)
}

// fetchCmd 发起 HTTP 获取并解析页面。
func fetchCmd(url string) tea.Cmd {
	return func() tea.Msg {
		html, err := fetcher.FetchHTML(url)
		if err != nil {
			return fetchMsg{url: url, err: err}
		}
		page, err := parser.Parse(strings.NewReader(html))
		if err != nil {
			return fetchMsg{url: url, err: err}
		}
		return fetchMsg{url: url, page: page}
	}
}

// execWgetCmd 打印 wget 命令到 stdout，然后执行 wget 下载覆盖当前文件。
func execWgetCmd(url string) tea.Cmd {
	return func() tea.Msg {
		cmdText := fmt.Sprintf("wget -N -q %q", url)
		cmd := exec.Command("wget", "-N", "-q", url)
		out, err := cmd.CombinedOutput()
		return execMsg{
			output: "$ " + cmdText + "\n" + string(out),
			err:    err,
		}
	}
}
