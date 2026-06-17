package tui

import "github.com/charmbracelet/lipgloss"

// ThemeStyles 是一组按主题区分的样式。
type ThemeStyles struct {
	Title      lipgloss.Style
	Help       lipgloss.Style
	Error      lipgloss.Style
	Meta       lipgloss.Style
	Border     lipgloss.Border
	BorderFg   lipgloss.Color
	Header     lipgloss.Style
	CursorFg   lipgloss.Color
	CursorBg   lipgloss.Color
	DirFg      lipgloss.Color
	NormalFg   lipgloss.Color
}

// StylesForTheme 返回对应主题的样式集合。
func StylesForTheme(theme string) ThemeStyles {
	switch theme {
	case "light":
		return ThemeStyles{
			Title:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D00B0")).Padding(0, 1),
			Help:     lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Padding(0, 1),
			Error:    lipgloss.NewStyle().Foreground(lipgloss.Color("#CC0000")).Padding(0, 1),
			Meta:     lipgloss.NewStyle().Foreground(lipgloss.Color("#0066CC")).Padding(0, 1),
			Border:   lipgloss.NormalBorder(),
			BorderFg: lipgloss.Color("#AAAAAA"),
			Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
			CursorFg: lipgloss.Color("#FFFFFF"),
			CursorBg: lipgloss.Color("#3366CC"),
			DirFg:    lipgloss.Color("#008844"),
			NormalFg: lipgloss.Color(""),
		}
	case "mono":
		return ThemeStyles{
			Title:    lipgloss.NewStyle().Bold(true).Padding(0, 1),
			Help:     lipgloss.NewStyle().Padding(0, 1),
			Error:    lipgloss.NewStyle().Bold(true).Padding(0, 1),
			Meta:     lipgloss.NewStyle().Padding(0, 1),
			Border:   lipgloss.NormalBorder(),
			BorderFg: lipgloss.Color(""),
			Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
			CursorFg: lipgloss.Color(""),
			CursorBg: lipgloss.Color(""),
			DirFg:    lipgloss.Color(""),
			NormalFg: lipgloss.Color(""),
		}
	default: // dark
		return ThemeStyles{
			Title:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Padding(0, 1),
			Help:     lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Padding(0, 1),
			Error:    lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Padding(0, 1),
			Meta:     lipgloss.NewStyle().Foreground(lipgloss.Color("#66CCFF")).Padding(0, 1),
			Border:   lipgloss.NormalBorder(),
			BorderFg: lipgloss.Color("#555555"),
			Header:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FEF3C7")).Bold(true).Padding(0, 1),
			CursorFg: lipgloss.Color("#FFCC66"),
			CursorBg: lipgloss.Color("#3C3C3C"),
			DirFg:    lipgloss.Color("#66DDAA"),
			NormalFg: lipgloss.Color(""),
		}
	}
}

// BorderForStyle 返回对应边框风格的 lipgloss.Border。
func BorderForStyle(s string) lipgloss.Border {
	switch s {
	case "rounded":
		return lipgloss.RoundedBorder()
	case "ascii":
		return lipgloss.Border{
			Top:          "-",
			Bottom:       "-",
			Left:         "|",
			Right:        "|",
			TopLeft:      "+",
			TopRight:     "+",
			BottomLeft:   "+",
			BottomRight:  "+",
			MiddleLeft:   "+",
			MiddleRight:  "+",
			Middle:       "+",
			MiddleTop:    "+",
		}
	default: // normal
		return lipgloss.NormalBorder()
	}
}
