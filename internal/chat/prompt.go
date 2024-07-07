package chat

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/container-labs/ada/internal/styles"
)

func Prompt(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(styles.PromptContentStyle.Render(prompt))
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(response, "\n"), nil
}
