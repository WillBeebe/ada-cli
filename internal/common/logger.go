package common

import (
	"fmt"

	"github.com/container-labs/ada/internal/styles"
)

var logger *AdaLogger

func Logger() *AdaLogger {
	if logger == nil {
		logger = new(AdaLogger)
	}

	return logger
}

type AdaLogger struct {
	logLevel string
}

func (ada *AdaLogger) Debug(msg string) {
	if ada.logLevel == "debug" {
		fmt.Println(styles.ContentStyle.Width(styles.Width()).Render(msg))
	}
}

func (ada *AdaLogger) Error(msg string) {
	fmt.Println(styles.ContentErrorStyle.Width(styles.Width()).Render(msg))
}

func (ada *AdaLogger) Info(msg string) {
	fmt.Println(styles.ContentStyle.Width(styles.Width()).Render(msg))
}

func (ada *AdaLogger) Warn(msg string) {
	fmt.Println(styles.ContentInfoStyle.Width(styles.Width()).Render(msg))
}

func (ada *AdaLogger) SetLevel(lvl string) {
	ada.logLevel = lvl
}
