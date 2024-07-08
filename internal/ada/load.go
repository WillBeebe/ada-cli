package ada

import (
	"fmt"
	"io/ioutil"

	"github.com/container-labs/ada/internal/common"
	"gopkg.in/yaml.v2"
)

var logger = common.Logger()

const (
	ADA_FILE = "ada.yaml"
)

func CheckExists() error {
	yamlFile := ADA_FILE

	_, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	return nil
}

func Load() *AdaFile {
	globalConfig := LoadConfig()
	logger.Debug(fmt.Sprintf("%+v\n", globalConfig))

	yamlFile := ADA_FILE

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

	return &config
}
