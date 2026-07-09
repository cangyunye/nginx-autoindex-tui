package tui

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"nginx-autoindex-tui/internal/fetcher"
	"nginx-autoindex-tui/internal/i18n"
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

// Config 是创建 Model 的配置参数。
type Config struct {
	InitialURL     string
	ForceOverwrite bool
	OutputDir      string
	Concurrent     int
	Insecure       bool
	UserAgent      string
	BorderStyle    string
	Theme          string
	Locale         i18n.Locale
}

// Model 持有 TUI 的全部状态。
type Model struct {
	currentURL     string
	rootURL        string
	page           *parser.Page
	cursor         int
	loading        bool
	err            string
	lastOutput     string // wget 执行的摘要输出
	forceOverwrite bool   // 是否跳过覆盖确认
	outputDir      string // 下载保存目录
	confirmPending bool   // 等待用户覆盖确认
	pendingURL     string // 待确认的文件 URL
	pendingName    string // 待确认的文件显示名
	concurrent     int    // 同时下载的 wget 进程最大数
	insecure       bool   // 跳过 SSL 证书验证
	userAgent      string // 自定义 User-Agent
	borderStyle    string // 边框风格
	theme          string // 颜色方案
	locale         i18n.Locale
}

func NewModel(cfg Config) *Model {
	// 确保 rootURL 以 / 结尾，使 url.ResolveReference 行为一致
	root := cfg.InitialURL
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	return &Model{
		currentURL:     root,  // 和 rootURL 使用同样的归一化值
		rootURL:        root,
		loading:        true,
		forceOverwrite: cfg.ForceOverwrite,
		outputDir:      cfg.OutputDir,
		concurrent:     cfg.Concurrent,
		insecure:       cfg.Insecure,
		userAgent:      cfg.UserAgent,
		borderStyle:    cfg.BorderStyle,
		theme:          cfg.Theme,
		locale:         cfg.Locale,
	}
}

// ---- Elm hooks ----

func (m *Model) Init() tea.Cmd {
	return fetchCmd(m.currentURL)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 覆盖确认模式：拦截 Y/N/A
		if m.confirmPending {
			switch msg.String() {
			case "y", "Y":
				m.confirmPending = false
				return m, m.execWgetCmd(m.pendingURL)
			case "n", "N":
				m.confirmPending = false
				skipStr := i18n.GetStrings(m.locale.Lang)
				m.lastOutput = fmt.Sprintf(skipStr.OutputSkipped, m.pendingName)
			case "a", "A":
				m.forceOverwrite = true
				m.confirmPending = false
				return m, m.execWgetCmd(m.pendingURL)
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.page == nil || len(m.page.Entries) == 0 {
				break
			}
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.page.Entries) - 1 // 到顶后跳到最后
			}
		case "down", "j":
			if m.page == nil || len(m.page.Entries) == 0 {
				break
			}
			if m.cursor < len(m.page.Entries)-1 {
				m.cursor++
			} else {
				m.cursor = 0 // 到底后跳到开头
			}
		case "esc", "backspace", "b":
			// 返回上级目录：查找 ../ 条目
			cmd := m.goToParent()
			if cmd != nil {
				return m, cmd
			}
		case "enter", " ":
			return m, m.openSelected()
		case "f", "F":
			m.forceOverwrite = !m.forceOverwrite
			str := i18n.GetStrings(m.locale.Lang)
			if m.forceOverwrite {
				m.lastOutput = str.ForceOverwriteON
			} else {
				m.lastOutput = str.ForceOverwriteOFF
			}
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
			m.lastOutput = i18n.GetStrings(m.locale.Lang).ExecFailedPrefix + msg.err.Error()
		} else {
			m.lastOutput = msg.output
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder

	s := StylesForTheme(m.theme)
	border := BorderForStyle(m.borderStyle)
	str := i18n.GetStrings(m.locale.Lang)

	b.WriteString(s.Title.Render(m.pageOrFallbackTitle()) + "\n")
	b.WriteString(s.Meta.Render(str.URLPrefix+m.currentURL) + "\n\n")

	if m.loading {
		b.WriteString(s.Meta.Render(str.Loading) + "\n")
		b.WriteString(s.Help.Render(str.LoadingQuitMsg) + "\n")
		return b.String()
	}
	if m.err != "" {
		b.WriteString(s.Error.Render(str.ErrorPrefix+m.err) + "\n")
		b.WriteString(s.Help.Render(str.ErrorRetryHelp) + "\n")
		return b.String()
	}
	if m.page == nil || len(m.page.Entries) == 0 {
		b.WriteString(s.Meta.Render(str.NoEntries) + "\n")
		b.WriteString(s.Help.Render(str.RefreshMsg) + "\n")
		return b.String()
	}

	// 构造表格数据
	rows := make([][]string, 0, len(m.page.Entries))
	for i, e := range m.page.Entries {
		idx := fmt.Sprintf("%d", i+1)
		kind := str.RowFile
		if e.IsDir {
			kind = str.RowDir
		}
		rows = append(rows, []string{idx, kind, e.Name, e.DateTime, e.Size})
	}

	t := table.New().
		Border(border).
		BorderStyle(lipgloss.NewStyle().Foreground(s.BorderFg)).
		Headers(str.ColIdx, str.ColType, str.ColName, str.ColDateTime, str.ColSize).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == -1 {
				return s.Header.Padding(0, 1)
			}
			base := lipgloss.NewStyle().Padding(0, 1).Foreground(s.NormalFg)
			if row == m.cursor {
				base = base.Foreground(s.CursorFg).Background(s.CursorBg).Bold(true)
			}
			if col == 1 && rows[row][1] == str.RowDir {
				base = base.Foreground(s.DirFg)
			}
			return base
		})

	b.WriteString(t.String())
	b.WriteString("\n")

	if m.lastOutput != "" {
		b.WriteString(s.Meta.Render(str.WgetHeader) + "\n")
		b.WriteString(m.lastOutput + "\n")
	}

	if m.confirmPending {
		b.WriteString(s.Error.Render(fmt.Sprintf(str.OverwriteConfirm, m.pendingName)) + "\n")
	} else {
		helpText := str.HelpText
		if m.forceOverwrite {
			helpText += str.ForceON
		}
		helpText += str.QuitSuffix
		b.WriteString(s.Help.Render(helpText) + "\n")
	}
	return b.String()
}

// ---- helpers ----

func (m *Model) pageOrFallbackTitle() string {
	if m.page != nil && m.page.Title != "" {
		return m.page.Title
	}
	return i18n.GetStrings(m.locale.Lang).TitleDefault
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
	// 根路径约束：不允许脱离 rootURL 范围
	if !strings.HasPrefix(absURL, m.rootURL) {
		m.err = fmt.Sprintf(i18n.GetStrings(m.locale.Lang).OutsideRootURL, absURL)
		return nil
	}
	if e.IsDir {
		m.loading = true
		m.err = ""
		return fetchCmd(absURL)
	}

	// 文件下载：检查是否已存在（仅在非强制覆盖时）
	if !m.forceOverwrite {
		filename := filenameFromURL(absURL)
		dlPath := filename
		if m.outputDir != "" {
			dlPath = filepath.Join(m.outputDir, filename)
		}
		if _, err := os.Stat(dlPath); err == nil {
			// 文件已存在，进入确认模式
			m.confirmPending = true
			m.pendingURL = absURL
			m.pendingName = filename
			return nil
		}
	}
	return m.execWgetCmd(absURL)
}

// fetchCmd 发起 HTTP 获取并根据 Content-Type 自动选择解析方式。
func fetchCmd(url string) tea.Cmd {
	return func() tea.Msg {
		res, err := fetcher.Fetch(url)
		if err != nil {
			return fetchMsg{url: url, err: err}
		}
		page, err := parser.ParseAuto(strings.NewReader(res.Body), res.ContentType)
		if err != nil {
			return fetchMsg{url: url, err: err}
		}
		return fetchMsg{url: url, page: page}
	}
}

// execWgetCmd 打印 wget 命令到 stdout，然后执行 wget 下载覆盖当前文件。
func (m *Model) execWgetCmd(url string) tea.Cmd {
	return func() tea.Msg {
		args := []string{"-N", "-q"}
		if m.outputDir != "" {
			args = append(args, "-P", m.outputDir)
		}
		args = append(args, url)
		cmdText := "wget " + strings.Join(args, " ")
		cmd := exec.Command("wget", args...)
		out, err := cmd.CombinedOutput()
		return execMsg{
			output: "$ " + cmdText + "\n" + string(out),
			err:    err,
		}
	}
}

// filenameFromURL 从下载 URL 中提取文件名。
func filenameFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return path.Base(u.Path)
}

// goToParent 查找当前页面的 ../ 条目并导航返回上级目录。
// 如果没有 ../（已在根目录），返回 nil。
func (m *Model) goToParent() tea.Cmd {
	if m.page == nil {
		return nil
	}
	for i, e := range m.page.Entries {
		if e.Name == "../" && e.Href == "../" {
			m.cursor = i
			return m.openSelected()
		}
	}
	return nil
}
