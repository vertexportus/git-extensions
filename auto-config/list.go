package auto_config

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

func ChooseFromList(items []list.Item, title string) (GpgKey, error) {
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	response, err := tea.NewProgram(m).Run()
	if err != nil {
		return GpgKey{}, err
	}
	return response.(model).choice.value, nil
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

type itemDelegate struct{}
type item struct {
	label string
	value GpgKey
}

func (i item) FilterValue() string                             { return "" }
func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
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
	choice   item
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
			i, ok := m.list.SelectedItem().(item)
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
	nilItem := item{}
	if m.choice != nilItem {
		return fmt.Sprintf("%s%s\n",
			quitTextStyle.Render("Configuring git with: "),
			selectedItemStyle.Render(m.choice.label))
	}
	if m.quitting {
		return quitTextStyle.Render("Canceling auto-config!")
	}
	return "\n" + m.list.View()
}
