package main

import (
	"fmt"
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
	Commits int
	Issues  int
	PRs     int
	Reviews int
	Repos   int
}

type model struct {
	spinner      spinner.Model
	loading      bool
	err          error
	contribution Contribution
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	return model{
		spinner:      s,
		loading:      true,
		contribution: Contribution{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Cmd(func() tea.Msg {
			c, err := getContributions()
			if err != nil {
				return err
			}
			return c
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
		return m, tea.Quit
	case Contribution:
		m.loading = false
		m.contribution = msg
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %s\n", m.err)
	}
	if m.loading {
		return fmt.Sprintf("\n\n %s Checking for Contributions!\n\n", m.spinner.View())
	}
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(border)

	var msg string
	if m.contribution.isGreen {
		msg = "You are green for today!"
		msg += fmt.Sprintf(
			"\n\ntotals:\n- %d commits\n- %d issues\n- %d PRs\n- %d PR reviews\n- %d repositories",
			m.contribution.Commits, m.contribution.Issues, m.contribution.PRs, m.contribution.Reviews, m.contribution.Repos,
		)
		msg = style.BorderForeground(lipgloss.Color("#04B575")).Render(msg)
	} else {
		msg = "You haven't made any contributions... yet!"
		msg = style.BorderForeground(lipgloss.Color("#E55353")).Render(msg)
	}
	return fmt.Sprintf("%s\n", msg)
}

func getContributions() (Contribution, errMsg) {
	opts := api.ClientOptions{
		EnableCache: true,
		Timeout:     5 * time.Second,
	}
	client, err := api.NewGraphQLClient(opts)
	if err != nil {
		return Contribution{}, errMsg(err)
	}
	var Query struct {
		Viewer struct {
			ContributionsCollection struct {
				TotalCommitContributions            int
				TotalIssueContributions             int
				TotalPullRequestContributions       int
				TotalPullRequestReviewContributions int
				TotalRepositoryContributions        int
				HasAnyContributions                 bool
			} `graphql:"contributionsCollection(from: $from)"`
		} `graphql:"viewer"`
	}
	t := time.Now()
	err = client.Query("contributionQuery", &Query, map[string]interface{}{"from": DateTime{t}})
	if err != nil {
		return Contribution{}, errMsg(err)
	}

	c := Contribution{
		isGreen: Query.Viewer.ContributionsCollection.HasAnyContributions,
		Commits: Query.Viewer.ContributionsCollection.TotalCommitContributions,
		Issues:  Query.Viewer.ContributionsCollection.TotalIssueContributions,
		PRs:     Query.Viewer.ContributionsCollection.TotalPullRequestContributions,
		Reviews: Query.Viewer.ContributionsCollection.TotalPullRequestReviewContributions,
		Repos:   Query.Viewer.ContributionsCollection.TotalRepositoryContributions,
	}

	return c, nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
