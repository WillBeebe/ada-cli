package create

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/api"
	"github.com/container-labs/ada/internal/theme"
)

var (
	adaptiveTextColor   = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	adaptiveAccentColor = lipgloss.AdaptiveColor{Light: "#1E90FF", Dark: "#87CEFA"}
	adaptiveErrorColor  = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF6347"}

	focusedStyle = lipgloss.NewStyle().Foreground(adaptiveAccentColor)
	blurredStyle = lipgloss.NewStyle().Foreground(adaptiveTextColor)
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type inputField int

const (
	nameField inputField = iota
	pathField
	providerField
	providerModelField
	submitButton
)

type item struct {
	name  string
	value string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.value }
func (i item) FilterValue() string { return i.name }

type model struct {
	inputs        []textinput.Model
	providerList  list.Model
	modelList     list.Model
	focused       inputField
	err           error
	selectedItems map[inputField]item
}

var providerOptions = []item{
	{name: "Anthropic", value: "anthropic"},
	{name: "Cohere", value: "cohere"},
	{name: "Google", value: "google"},
	{name: "Groq", value: "groq"},
	{name: "Ollama", value: "ollama"},
	{name: "OpenAI", value: "openai"},
	{name: "Perplexity AI", value: "perplexity"},
}

var providerModelOptions = map[string][]item{
	"anthropic": {
		{name: "Claude 3.5 Sonnet", value: "claude-3-5-sonnet-20240620"},
		{name: "Claude 3 Opus", value: "claude-3-opus-20240229"},
		{name: "Claude 3 Sonnet", value: "claude-3-sonnet-20240229"},
		{name: "Claude 3 Haiku", value: "claude-3-haiku-20240307"},
	},
	"cohere": {
		{name: "Command-R", value: "command-r"},
		{name: "Command-R Plus", value: "command-r-plus"},
	},
	"groq": {
		{name: "Llama 3 8b", value: "llama3-8b-8192"},
		{name: "Llama 3 70b", value: "llama3-70b-8192"},
		{name: "Llama 2 70b", value: "llama2-70b-4096"},
		{name: "Mixtral 8x7b", value: "mixtral-8x7b-32768"},
		{name: "Gemma 7b", value: "gemma-7b-it"},
	},
	"google": {
		{name: "Gemini 1.0 Pro", value: "gemini-1.0-pro-latest"},
		{name: "Gemini 1.0 Pro - Vision", value: "gemini-1.0-pro-vision-latest"},
		{name: "Gemini 1.5 Pro", value: "gemini-1.5-pro-latest"},
		{name: "Chat Bison", value: "chat-bison-001"},
		{name: "Text Bison", value: "text-bison-001"},
	},
	"ollama": {
		{name: "Codellama", value: "codellama:latest"},
		{name: "Llama3 8b", value: "llama3:8b"},
		{name: "Llama2 7b", value: "llama2:7b"},
		{name: "Mixtral", value: "mixtral:latest"},
	},
	"openai": {
		{name: "GPT 4o", value: "gpt-4o"},
		{name: "GPT 4 Turbo 2024-04-09", value: "gpt-4-turbo-2024-04-09"},
		{name: "GPT 3.5 Turbo", value: "gpt-3.5-turbo-0125"},
	},
	"perplexity": {
		{name: "Mistral 7b Instruct", value: "mistral-7b-instruct"},
		{name: "Sonar Medium - Chat", value: "sonar-medium-chat"},
		{name: "Sonar Medium - Online", value: "sonar-medium-online"},
	},
}

func initialModel() model {
	var inputs []textinput.Model
	inputs = make([]textinput.Model, 2)

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.Cursor.Style = theme.CurrentTheme.NormalText.Copy().Foreground(theme.CurrentTheme.AccentColor)
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Project Name"
			t.Focus()
			t.PromptStyle = theme.CurrentTheme.NormalText.Copy().Foreground(theme.CurrentTheme.AccentColor)
			t.TextStyle = theme.CurrentTheme.NormalText.Copy().Foreground(theme.CurrentTheme.AccentColor)
		case 1:
			t.Placeholder = "Project Path"
			t.CharLimit = 64
		}

		inputs[i] = t
	}

	providerItems := make([]list.Item, len(providerOptions))
	for i, option := range providerOptions {
		providerItems[i] = option
	}

	providerList := list.New(providerItems, list.NewDefaultDelegate(), 0, 0)
	providerList.Title = "Select Provider"
	providerList.SetShowStatusBar(false)
	providerList.SetFilteringEnabled(false)
	providerList.Styles.Title = theme.CurrentTheme.BoldText.Copy().Foreground(theme.CurrentTheme.AccentColor)

	modelList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	modelList.Title = "Select Provider Model"
	modelList.SetShowStatusBar(false)
	modelList.SetFilteringEnabled(false)
	modelList.Styles.Title = theme.CurrentTheme.BoldText.Copy().Foreground(theme.CurrentTheme.AccentColor)

	return model{
		inputs:        inputs,
		providerList:  providerList,
		modelList:     modelList,
		focused:       nameField,
		selectedItems: make(map[inputField]item),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == submitButton {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > submitButton {
				m.focused = nameField
			} else if m.focused < nameField {
				m.focused = submitButton
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if inputField(i) == m.focused {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}

			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.providerList.SetSize(msg.Width-h, msg.Height-v)
		m.modelList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.focused {
	case nameField, pathField:
		m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
		cmds = append(cmds, cmd)
	case providerField:
		m.providerList, cmd = m.providerList.Update(msg)
		cmds = append(cmds, cmd)
		if item, ok := m.providerList.SelectedItem().(item); ok {
			m.selectedItems[providerField] = item
			modelItems := make([]list.Item, len(providerModelOptions[item.value]))
			for i, option := range providerModelOptions[item.value] {
				modelItems[i] = option
			}
			m.modelList.SetItems(modelItems)
		}
	case providerModelField:
		m.modelList, cmd = m.modelList.Update(msg)
		cmds = append(cmds, cmd)
		if item, ok := m.modelList.SelectedItem().(item); ok {
			m.selectedItems[providerModelField] = item
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	if m.focused == providerField {
		b.WriteString(m.providerList.View())
	} else if m.focused == providerModelField {
		b.WriteString(m.modelList.View())
	}

	button := &blurredButton
	if m.focused == submitButton {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n", *button)

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(adaptiveErrorColor)
		b.WriteString(fmt.Sprintf("\n\nError: %s", errorStyle.Render(m.err.Error())))
	}

	return b.String()
}

func CreateProject() (*api.Project, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running program: %w", err)
	}

	finalModel := m.(model)

	if finalModel.err != nil {
		return nil, finalModel.err
	}

	name := finalModel.inputs[0].Value()
	path := finalModel.inputs[1].Value()
	provider := finalModel.selectedItems[providerField].value
	providerModel := finalModel.selectedItems[providerModelField].value

	return &api.Project{
		Name:          name,
		Path:          path,
		Provider:      provider,
		ProviderModel: providerModel,
	}, nil
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func loadColorPreferences() {
	cfg := ada.LoadConfig()
	if cfg != nil && cfg.ColorScheme != "" {
		// Apply color scheme preferences
		// This is just an example, adjust according to your config structure
		switch cfg.ColorScheme {
		case "light":
			adaptiveTextColor.Dark = "#000000"
			adaptiveAccentColor.Dark = "#1E90FF"
			adaptiveErrorColor.Dark = "#FF0000"
		case "dark":
			adaptiveTextColor.Light = "#FFFFFF"
			adaptiveAccentColor.Light = "#87CEFA"
			adaptiveErrorColor.Light = "#FF6347"
		}
	}
}

func init() {
	loadColorPreferences()
}
