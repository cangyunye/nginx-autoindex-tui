package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"nginx-autoindex-tui/internal/fetcher"
	"nginx-autoindex-tui/internal/i18n"
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
	locale := i18n.DetectLocale()

	rootCmd := &cobra.Command{
		Use:   "nginx-autoindex-tui [URL]",
		Short: "nginx autoindex terminal file browser",
		Long:  `An interactive terminal file browser and downloader for nginx autoindex servers.` + "\n\n" + `Supports HTML and JSON autoindex listings, batch download via wget, and directory navigation.`,
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

			// 自动检测终端 locale：若为 GBK 编码环境且用户未显式指定边框风格，回退为 ascii
			// 确保 Unicode 制表符不会在 GBK 终端上显示为乱码。
			if !cmd.Flags().Changed("border-style") && locale.IsGBK {
				borderStyle = "ascii"
				msg := "GBK terminal detected, auto-switching to --border-style ascii"
				if locale.Lang == "zh" {
					msg = "检测到 GBK 终端编码，自动使用 --border-style ascii"
				}
				fmt.Fprintln(os.Stderr, msg)
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
				Locale:         locale,
			})

			if insecure {
				fetcher.SetInsecure(true)
			}
			if userAgent != "" {
				fetcher.SetUserAgent(userAgent)
			}

			output := i18n.WrapOutput(os.Stdout, locale)
			p := tea.NewProgram(m, tea.WithOutput(output))
			if _, err := p.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "bubble tea error:", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "force overwrite existing files, skip confirmation")
	rootCmd.Flags().IntVarP(&concurrent, "concurrent", "c", 3, "max concurrent wget download processes")
	rootCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "download directory (default current dir)")
	rootCmd.Flags().BoolVarP(&insecure, "insecure", "k", false, "skip HTTPS SSL certificate verification")
	rootCmd.Flags().StringVarP(&userAgent, "user-agent", "u", "", "HTTP User-Agent header (default Go-http-client)")
	rootCmd.Flags().StringVar(&borderStyle, "border-style", "normal", "table border style: normal / rounded / ascii")
	rootCmd.Flags().StringVarP(&theme, "theme", "t", "dark", "color scheme: dark / light / mono")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
