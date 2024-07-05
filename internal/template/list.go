package template

import (
	"fmt"
	"strings"

	"github.com/container-labs/ada/internal/cache"
	"github.com/container-labs/ada/internal/common"
)

var logger = common.Logger()

func List() error {
	templates, err := cache.Templates()

	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Templates:\n  %s\n", strings.Join(templates, "\n  ")))

	return nil
}
