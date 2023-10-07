package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

type Config struct {
	Title            string
	Items            []list.Item
	Width            int
	Height           int
	ShowStatusBar    bool
	FilteringEnabled bool

	SelectedItemText string
	CancelText       string

	SuppressQuitText bool
}

func Choose(config *Config) (ItemValue, error) {
	handleConfigDefaults(config)
	l := list.New(config.Items, itemDelegate{}, config.Width, config.Height)
	l.Title = config.Title
	l.SetShowStatusBar(config.ShowStatusBar)
	l.SetFilteringEnabled(config.FilteringEnabled)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, config: config}

	response, err := tea.NewProgram(m).Run()
	if err != nil {
		return itemDelegate{}, err
	}
	return response.(model).choice.value, nil
}

func handleConfigDefaults(config *Config) {
	if config.Width == 0 {
		config.Width = defaultWidth
	}
	if config.Height == 0 {
		config.Height = listHeight
	}

	if config.SelectedItemText == "" {
		config.SelectedItemText = "Selected item: "
	}

	if config.CancelText == "" {
		config.CancelText = "Cancelled"
	}
}

// ####################################################################################################### Constants ###

const defaultWidth = 50
const listHeight = 14

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// ############################################################################################################ List ###

type ItemValue interface{}
type itemDelegate struct{}
type Item[V ItemValue] struct {
	label string
	value V
}

func NewListSimpleItem(value string) Item[ItemValue] {
	return Item[ItemValue]{label: value, value: value}
}
func NewListItem(label string, value ItemValue) Item[ItemValue] {
	return Item[ItemValue]{label: label, value: value}
}

func (i Item[ListItemValue]) FilterValue() string              { return "" }
func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item[ItemValue])
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.label)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, err := fmt.Fprint(w, fn(str))
	if err != nil {
		panic(err)
	}
}

type model struct {
	list     list.Model
	config   *Config
	choice   Item[ItemValue]
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter", " ":
			i, ok := m.list.SelectedItem().(Item[ItemValue])
			if ok {
				m.choice = i
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	nilItem := Item[ItemValue]{}
	if m.choice != nilItem {
		if m.config.SuppressQuitText {
			return ""
		}
		return fmt.Sprintf("%s%s\n",
			quitTextStyle.Render(m.config.SelectedItemText),
			selectedItemStyle.Render(m.choice.label))
	}
	if m.quitting {
		if m.config.SuppressQuitText {
			return ""
		}
		return quitTextStyle.Render(m.config.CancelText)
	}
	return "\n" + m.list.View()
}
