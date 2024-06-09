package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2/pkg/api"
)

var (
	border = lipgloss.Border{
		Top:         "._.:*:",
		Bottom:      "._.:*:",
		Left:        "|*",
		Right:       "|*",
		TopLeft:     "*",
		TopRight:    "*",
		BottomLeft:  "*",
		BottomRight: "*",
	}
)

type DateTime struct{ time.Time }

type errMsg error

type Contribution struct {
	isGreen bool
}

type model struct {
	spinner spinner.Model
	loading bool
	err     error
	green   bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner: s,
		loading: true,
		green:   false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Cmd(func() tea.Msg {
			return ExampleContributionGraphQL()
		}),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	case Contribution:
		m.loading = false
		m.green = msg.isGreen
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.loading {
		return fmt.Sprintf("\n\n %s Checking for Contributions!\n\n", m.spinner.View())
	}
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(border)

	var msg string
	if m.green {
		msg = "You are green for today!"
		msg = style.BorderForeground(lipgloss.Color("#04B575")).Render(msg)
	} else {
		msg = "You haven't made any contributions... yet!"
		msg = style.BorderForeground(lipgloss.Color("#E55353")).Render(msg)
	}
	return fmt.Sprintf("%s\n", msg)
}

func ExampleContributionGraphQL() Contribution {
	opts := api.ClientOptions{
		EnableCache: true,
		Timeout:     5 * time.Second,
	}
	client, err := api.NewGraphQLClient(opts)
	if err != nil {
		log.Fatal(err)
	}
	// Add attributes to get the contributions for commits, issues, and pull requests, and maybe more
	var Query struct {
		Viewer struct {
			ContributionsCollection struct {
				TotalCommitContributions int
				HasAnyContributions      bool
			} `graphql:"contributionsCollection(from: $from)"`
		} `graphql:"viewer"`
	}

	err = client.Query("contributionQuery", &Query, map[string]interface{}{"from": DateTime{time.Now()}})
	if err != nil {
		log.Fatal(err)
	}

	isGreen := Query.Viewer.ContributionsCollection.HasAnyContributions
	return Contribution{isGreen: isGreen}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
