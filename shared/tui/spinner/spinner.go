package spinner

import (
	"fmt"
	"git_extensions/shared/errors"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ################################################################################################## Public/Configs ###

type Config struct {
	Label   string
	Spinner spinner.Spinner
	Color   lipgloss.TerminalColor
	Cmd     func(channel chan<- tea.Cmd)
}

func Show(config *Config) {
	p := tea.NewProgram(model{
		spinner: createModel(config),
		config:  config,
		channel: make(chan tea.Cmd)})
	if _, err := p.Run(); err != nil {
		errors.HandleError(err)
	}
}

func createModel(config *Config) spinner.Model {
	handleConfigDefaults(config)
	s := spinner.New()
	s.Spinner = config.Spinner
	s.Style = lipgloss.NewStyle().Foreground(config.Color)
	return s
}

func handleConfigDefaults(config *Config) {
	if config.Color == nil {
		config.Color = defaultColor
	}
	if len(config.Spinner.Frames) == 0 {
		config.Spinner = defaultSpinner
	}
}

// ####################################################################################################### Constants ###
var (
	defaultSpinner = spinner.Dot
	defaultColor   = lipgloss.Color("205")
)

// ######################################################################################################### Spinner ###

type model struct {
	config   *Config
	spinner  spinner.Model
	quitting bool
	channel  chan tea.Cmd
}

func (m model) Init() tea.Cmd {
	if m.config.Cmd != nil {
		go m.config.Cmd(m.channel)
	}
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// run through messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case tea.QuitMsg:
		m.quitting = true
		return m, tea.Quit

	default:
		// check command channel
		select {
		case cmd, ok := <-m.channel:
			if !ok {
				return m, tea.Quit
			}
			return m, cmd
		default:
		}
		// default to updating spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("   %s %s", m.spinner.View(), m.config.Label)
}
