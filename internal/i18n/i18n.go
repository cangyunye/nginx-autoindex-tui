package i18n

import (
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Locale 表示用户的终端区域设置。
type Locale struct {
	IsGBK bool   // 终端是否为 GBK 编码环境
	Lang  string // "en" 或 "zh"
}

// DetectLocale 按 LC_ALL > LC_CTYPE > LANG 顺序检测区域设置。
// 默认返回英文 locale。
func DetectLocale() Locale {
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		val := os.Getenv(env)
		if val == "" || val == "C" || val == "POSIX" {
			continue
		}
		low := strings.ToLower(val)
		if strings.Contains(low, ".gbk") ||
			strings.Contains(low, ".gb2312") ||
			strings.Contains(low, ".gb18030") ||
			strings.Contains(low, ".cp936") ||
			strings.Contains(low, ".euccn") {
			return Locale{IsGBK: true, Lang: "zh"}
		}
		// 常见中文 locale 如 zh_CN.UTF-8 → 使用中文界面
		if strings.HasPrefix(low, "zh_") {
			return Locale{Lang: "zh"}
		}
	}
	return Locale{Lang: "en"}
}

// Strings 包含所有用户可见的 TUI 文字。
type Strings struct {
	TitleDefault string

	// 页面元素
	URLPrefix string

	// 加载状态
	Loading        string
	LoadingQuitMsg string

	// 错误状态
	ErrorPrefix    string
	ErrorRetryHelp string

	// 空白状态
	NoEntries  string
	RefreshMsg string

	// 表格头
	ColIdx      string
	ColType     string
	ColName     string
	ColDateTime string
	ColSize     string

	// 行类型
	RowFile string
	RowDir  string

	// wget 输出区
	WgetHeader string

	// 覆盖确认
	OverwriteConfirm string // Printf 格式，含 %s

	// 底部帮助栏
	HelpText     string
	ForceON      string
	QuitSuffix   string
	ConfirmQuit  string
	ConfirmRetry string

	// 状态栏消息
	ForceOverwriteON  string
	ForceOverwriteOFF string
	OutputSkipped     string // Printf 格式，含 %s
	ExecFailedPrefix  string

	// 安全 / 错误消息
	OutsideRootURL string // Printf 格式，含 %s
}

// GetStrings 返回指定语言的字符串表。
func GetStrings(lang string) Strings {
	if lang == "zh" {
		return zhStrings
	}
	return enStrings
}

// enStrings 英文（默认）字符串表。
var enStrings = Strings{
	TitleDefault: "Nginx Autoindex Browser",

	URLPrefix: "URL: ",

	Loading:        "  loading...",
	LoadingQuitMsg: "(ctrl+c to quit)",

	ErrorPrefix:    "ERROR: ",
	ErrorRetryHelp: "Press R to retry, ctrl+c to quit",

	NoEntries:  "  (no entries)",
	RefreshMsg: "R refresh, ctrl+c to quit",

	ColIdx:      "Idx",
	ColType:     "Type",
	ColName:     "Name",
	ColDateTime: "Date/Time",
	ColSize:     "Size",

	RowFile: "FILE",
	RowDir:  "DIR ",

	WgetHeader: "--- wget output ---",

	OverwriteConfirm: "\"%s\" exists. [Y] overwrite  [N] skip  [A] always",

	HelpText:    "↑/↓ or j/k move (wrap) · Enter open/dl · Esc/b parent · R refresh · F force-overwrite",
	ForceON:     " [FORCE ON]",
	QuitSuffix:  " · q/ctrl+c quit",
	ConfirmQuit: "q/ctrl+c",
	ConfirmRetry: "R",

	ForceOverwriteON:  "force overwrite ON",
	ForceOverwriteOFF: "force overwrite OFF",
	OutputSkipped:     "skipped: %s",
	ExecFailedPrefix:  "exec failed: ",
	OutsideRootURL:    "outside root URL: %s",
}

// zhStrings 中文（GBK 环境）字符串表。
var zhStrings = Strings{
	TitleDefault: "Nginx 自动索引浏览器",

	URLPrefix: "URL: ",

	Loading:        "  加载中...",
	LoadingQuitMsg: "(ctrl+c 退出)",

	ErrorPrefix:    "错误：",
	ErrorRetryHelp: "按 R 重试，ctrl+c 退出",

	NoEntries:  "  （无条目）",
	RefreshMsg: "R 刷新，ctrl+c 退出",

	ColIdx:      "序号",
	ColType:     "类型",
	ColName:     "名称",
	ColDateTime: "日期/时间",
	ColSize:     "大小",

	RowFile: "文件",
	RowDir:  "目录",

	WgetHeader: "--- wget 输出 ---",

	OverwriteConfirm: "“%s”已存在。[Y] 覆盖  [N] 跳过  [A] 始终覆盖",

	HelpText:    "↑/↓ 或 j/k 移动（循环）· Enter 进入/下载 · Esc/b 上级 · R 刷新 · F 强制覆盖",
	ForceON:     " [强制覆盖开]",
	QuitSuffix:  " · q/ctrl+c 退出",
	ConfirmQuit: "q/ctrl+c",
	ConfirmRetry: "R",

	ForceOverwriteON:  "强制覆盖已开启",
	ForceOverwriteOFF: "强制覆盖已关闭",
	OutputSkipped:     "已跳过：%s",
	ExecFailedPrefix:  "执行失败：",
	OutsideRootURL:    "超出根 URL 范围：%s",
}

// WrapOutput 在 GBK 环境下将 writer 包装为 UTF-8→GBK 编码器，
// 非 GBK 环境下原样返回。
func WrapOutput(w io.Writer, l Locale) io.Writer {
	if !l.IsGBK {
		return w
	}
	return transform.NewWriter(w, simplifiedchinese.GBK.NewEncoder())
}
