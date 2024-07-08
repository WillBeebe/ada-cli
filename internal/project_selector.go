package internal

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/container-labs/ada/internal/api"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	currentProjectStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#1E90FF")).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000"))

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#FF69B4")).
				Bold(true)

	normalDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	selectedDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#FF69B4"))
)

type item struct {
	id        int
	name      string
	isCurrent bool
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return fmt.Sprintf("Project ID: %d", i.id) }
func (i item) FilterValue() string { return i.name }

type model struct {
	list     list.Model
	choice   *item
	quitting bool
}

// Custom delegate that embeds the default delegate
type customDelegate struct {
	list.DefaultDelegate
}

// Override the Render method
func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	var title, desc string

	if i.isCurrent {
		title = currentProjectStyle.Render(i.Title() + " (current)")
	} else if index == m.Index() {
		title = selectedItemStyle.Render(i.Title())
	} else {
		title = normalItemStyle.Render(i.Title())
	}

	if index == m.Index() {
		desc = selectedDescStyle.Render(i.Description())
	} else {
		desc = normalDescStyle.Render(i.Description())
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

func initialModel(projects []api.Project, currentProjectID int) model {
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = item{id: p.ID, name: p.Name, isCurrent: p.ID == currentProjectID}
	}

	delegate := customDelegate{list.NewDefaultDelegate()}

	l := list.New(items, delegate, 0, 0)
	l.Title = "Select a project"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = l.Styles.Title.
		Foreground(lipgloss.Color("#FF69B4")).
		Bold(true)

	return model{
		list: l,
	}
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
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = &i
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != nil {
		return fmt.Sprintf("You chose %s (%d)\n", m.choice.name, m.choice.id)
	}
	if m.quitting {
		return "Bye!\n"
	}
	return docStyle.Render(m.list.View())
}

func SelectProject(projects []api.Project, currentProjectID int) (*api.Project, error) {
	p := tea.NewProgram(initialModel(projects, currentProjectID), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running program: %w", err)
	}

	if m.(model).choice != nil {
		selectedItem := m.(model).choice
		for _, p := range projects {
			if p.ID == selectedItem.id {
				return &p, nil
			}
		}
	}

	return nil, nil // No selection made
}
