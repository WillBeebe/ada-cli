package internal

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/container-labs/ada/internal/api"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	id   int
	name string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return fmt.Sprintf("Project ID: %d", i.id) }
func (i item) FilterValue() string { return i.name }

type model struct {
	list     list.Model
	choice   *item
	quitting bool
}

func initialModel(projects []api.Project) model {
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = item{id: p.ID, name: p.Name}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a project"

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

func SelectProject(projects []api.Project) (*api.Project, error) {
	p := tea.NewProgram(initialModel(projects), tea.WithAltScreen())
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
