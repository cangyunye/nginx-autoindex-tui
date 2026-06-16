package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"nginx-autoindex-tui/internal/tui"
)

func main() {
	var initialURL string
	if len(os.Args) >= 2 {
		initialURL = strings.TrimSpace(os.Args[1])
	} else {
		fmt.Print("Enter nginx autoindex URL: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "read input:", err)
			os.Exit(1)
		}
		initialURL = strings.TrimSpace(line)
	}
	if initialURL == "" {
		fmt.Fprintln(os.Stderr, "no url provided")
		os.Exit(1)
	}

	m := tui.NewModel(initialURL)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "bubble tea error:", err)
		os.Exit(1)
	}
}
