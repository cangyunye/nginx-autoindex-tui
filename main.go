package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"nginx-autoindex-tui/internal/fetcher"
	"nginx-autoindex-tui/internal/tui"
)

var (
	forceOverwrite bool
	concurrent     int
	outputDir      string
	insecure       bool
	userAgent      string
	borderStyle    string
	theme          string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nginx-autoindex-tui [URL]",
		Short: "nignx autoindex 终端文件浏览器",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var initialURL string
			if len(args) >= 1 {
				initialURL = strings.TrimSpace(args[0])
			}
			if initialURL == "" {
				_ = cmd.Help()
				os.Exit(1)
			}

			m := tui.NewModel(tui.Config{
				InitialURL:     initialURL,
				ForceOverwrite: forceOverwrite,
				OutputDir:      outputDir,
				Concurrent:     concurrent,
				Insecure:       insecure,
				UserAgent:      userAgent,
				BorderStyle:    borderStyle,
				Theme:          theme,
			})

			if insecure {
				fetcher.SetInsecure(true)
			}
			if userAgent != "" {
				fetcher.SetUserAgent(userAgent)
			}

			p := tea.NewProgram(m)
			if _, err := p.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "bubble tea error:", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "强制覆盖已存在文件，跳过确认")
	rootCmd.Flags().IntVarP(&concurrent, "concurrent", "c", 3, "同时下载的 wget 进程最大数")
	rootCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "下载文件保存目录（默认当前目录）")
	rootCmd.Flags().BoolVarP(&insecure, "insecure", "k", false, "跳过 HTTPS 的 SSL 证书验证")
	rootCmd.Flags().StringVarP(&userAgent, "user-agent", "u", "", "HTTP 请求的 User-Agent 头（默认 Go-http-client）")
	rootCmd.Flags().StringVar(&borderStyle, "border-style", "normal", "表格边框风格：normal / rounded / ascii")
	rootCmd.Flags().StringVarP(&theme, "theme", "t", "dark", "颜色方案：dark / light / mono")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
