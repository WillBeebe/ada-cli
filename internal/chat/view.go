package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#E0E0E0")).
			Padding(0, 1).
			MarginTop(1).
			Align(lipgloss.Right)

	aiStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#4A4A4A")).
		Padding(0, 1).
		MarginTop(1).
		Align(lipgloss.Left)

	timestampStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginLeft(1).
			MarginRight(1).
			Italic(true)

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF69B4")).
				Bold(true)

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF69B4")).
			Padding(0, 1)

	waitingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Italic(true)
)

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.inputView(),
	)
}

func (m model) inputView() string {
	if m.isWaiting {
		return waitingStyle.Render("AI is thinking...")
	}

	cursor := " "
	if m.cursorY == len(m.lines)-1 && m.cursorX == len(m.lines[m.cursorY]) {
		cursor = "█"
	}

	input := strings.Join(m.lines, "\n")
	if input == "" {
		input = "Type your message here..."
	}

	return fmt.Sprintf(
		"%s\n%s",
		inputPromptStyle.Render("Enter your message:"),
		inputStyle.Render(input+cursor),
	)
}

func (m model) formatMessages() string {
	var sb strings.Builder

	for _, msg := range m.messages {
		var style lipgloss.Style
		var prefix string

		if msg.isUser {
			style = userStyle
			prefix = "You"
		} else {
			style = aiStyle
			prefix = "AI"
		}

		content := fmt.Sprintf("%s: %s", prefix, msg.content)
		wrappedContent := lipgloss.NewStyle().Width(m.viewport.Width - 4).Render(content)
		styledContent := style.Render(wrappedContent)

		timestamp := timestampStyle.Render(msg.time.Format("15:04"))

		if msg.isUser {
			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, timestamp, styledContent))
		} else {
			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, styledContent, timestamp))
		}
		sb.WriteString("\n\n")
	}

	return sb.String()
}

func highlightMarkdown(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	codeBlock := false
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			codeBlock = !codeBlock
			result = append(result, codeStyle.Render(line))
			continue
		}

		if codeBlock {
			result = append(result, codeStyle.Render(line))
			continue
		}

		// Headings
		if strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				result = append(result, headingStyle.Render(parts[0]+" "+parts[1]))
				continue
			}
		}

		// Inline code
		line = highlightInlineCode(line)

		// Links
		line = highlightLinks(line)

		// List items
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			result = append(result, listItemStyle.Render(line))
			continue
		}

		// Bold and italic
		line = strings.ReplaceAll(line, "**", boldStyle.Render("**"))
		line = strings.ReplaceAll(line, "*", italicStyle.Render("*"))

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func highlightInlineCode(line string) string {
	parts := strings.Split(line, "`")
	for i := 1; i < len(parts); i += 2 {
		if i < len(parts) {
			parts[i] = codeStyle.Render(parts[i])
		}
	}
	return strings.Join(parts, "")
}

func highlightLinks(line string) string {
	var result []string
	for {
		openBracket := strings.Index(line, "[")
		closeBracket := strings.Index(line, "]")
		openParen := strings.Index(line, "(")
		closeParen := strings.Index(line, ")")

		if openBracket == -1 || closeBracket == -1 || openParen == -1 || closeParen == -1 ||
			closeBracket < openBracket || openParen < closeBracket || closeParen < openParen {
			result = append(result, line)
			break
		}

		result = append(result, line[:openBracket])
		linkText := line[openBracket+1 : closeBracket]
		linkURL := line[openParen+1 : closeParen]
		result = append(result, linkStyle.Render(fmt.Sprintf("[%s](%s)", linkText, linkURL)))
		line = line[closeParen+1:]
	}
	return strings.Join(result, "")
}
