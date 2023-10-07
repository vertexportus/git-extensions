package spinner

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

// ################################################################################################## Public/Configs ###

type Config struct {
	Label   string
	Spinner spinner.Spinner
	Color   lipgloss.TerminalColor
	quit    bool
}

type Spinner struct {
	teaProgram *tea.Program
}

func Show(config *Config) Spinner {
	handleConfigDefaults(config)
	s := spinner.New()
	s.Spinner = config.Spinner
	s.Style = lipgloss.NewStyle().Foreground(config.Color)

	p := tea.NewProgram(model{spinner: s, config: config})
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return Spinner{teaProgram: p}
}

func handleConfigDefaults(config *Config) {
	if config.Color == nil {
		config.Color = defaultColor
	}
	if len(config.Spinner.Frames) == 0 {
		config.Spinner = defaultSpinner
	}
}

func (c *Config) Quit() {
	c.quit = true
	//s.teaProgram.Send(tea.QuitMsg{})
}

// ####################################################################################################### Constants ###
var (
	defaultSpinner = spinner.Dot
	defaultColor   = lipgloss.Color("205")
)

// ######################################################################################################### Spinner ###

type errMsg error

type model struct {
	config   *Config
	spinner  spinner.Model
	quitting bool
	err      error
}

func (m model) Tick() tea.Msg {
	if m.config.quit {
		return tea.QuitMsg{}
	}
	return m.spinner.Tick
}

func (m model) Init() tea.Cmd {
	return m.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case tea.QuitMsg:
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
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}
