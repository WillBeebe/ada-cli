package ada

import (
	"fmt"
	"io/ioutil"

	"github.com/container-labs/ada/internal/common"
	"gopkg.in/yaml.v2"
)

var logger = common.Logger()

func CheckExists() error {
	yamlFile := "paddle.yaml"

	_, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	return nil
}

func Load() AdaFile {
	yamlFile := "paddle.yaml"

	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		// logger.Errorf("yamlFile.Get err   #%v ", err)
	}

	var config AdaFile
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		// logger.Errorf("Unmarshal: %v", err)
	}

	if config.Type == "nodejs" {
		config.Port = 3000
	} else if config.Type == "python" {
		config.Port = 8000
	} else {
		config.Port = 8080
	}

	metadata := GetMetadata(config)
	config.Metadata = metadata

	logger.Debug(fmt.Sprintf("%+v\n", config))

	return config
}
