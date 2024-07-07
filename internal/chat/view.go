package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.inputMode {
		return m.inputView()
	}
	return m.chatView()
}

func (m model) inputView() string {
	s := strings.Builder{}
	s.WriteString(m.viewport.View())
	s.WriteString("\n\n")
	s.WriteString(promptStyle.Render("Enter multiple lines of text (press Esc to send, write 'bye' then press Esc to exit):"))
	s.WriteString("\n")
	for i, line := range m.lines {
		if i == m.cursorY {
			beforeCursor := m.highlightSyntax(line[:m.cursorX])
			afterCursor := m.highlightSyntax(line[m.cursorX:])
			cursorLine := lipgloss.JoinHorizontal(lipgloss.Left,
				beforeCursor,
				cursorStyle.Render(cursorText),
				afterCursor,
			)
			s.WriteString(cursorLine)
		} else {
			s.WriteString(m.highlightSyntax(line))
		}
		s.WriteString("\n")
	}
	return s.String()
}

func (m model) chatView() string {
	s := strings.Builder{}
	s.WriteString(m.viewport.View())
	s.WriteString("\n\n")
	if m.isWaiting {
		s.WriteString(m.spinner.View() + " Waiting for response...")
	} else {
		s.WriteString("Press Enter to start typing")
	}
	return s.String()
}

func (m model) formatMessages() string {
	var sb strings.Builder
	for _, msg := range m.messages {
		content := highlightMarkdown(msg.content)
		if msg.isUser {
			sb.WriteString(userMsgStyle.Render("User: " + content))
		} else {
			sb.WriteString(aiMsgStyle.Render("AI: " + content))
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

func (m model) highlightSyntax(line string) string {
	words := strings.Fields(line)
	for i, word := range words {
		if strings.HasPrefix(word, "/") {
			words[i] = commandStyle.Render(word)
		} else if strings.HasPrefix(word, "*") && strings.HasSuffix(word, "*") {
			words[i] = boldStyle.Render(word)
		} else if strings.HasPrefix(word, "_") && strings.HasSuffix(word, "_") {
			words[i] = italicStyle.Render(word)
		}
	}
	return textStyle.Render(strings.Join(words, " "))
}
