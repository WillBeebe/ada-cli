package template

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/osteele/liquid"
)

type RenderData struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	OutputDir string
}

func Render(templateDir string, data *RenderData) error {
	// templateDir := fmt.Sprintf("%s/%s", cache.TemplateCacheDir(), "helm-chart")
	fmt.Println(data)

	files, readDirErr := os.ReadDir(templateDir)
	if readDirErr != nil {
		return readDirErr
	}

	if _, checkDirExistsErr := os.Stat(data.OutputDir); errors.Is(checkDirExistsErr, os.ErrNotExist) {
		mkdirErr := os.MkdirAll(data.OutputDir, os.ModePerm)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	for _, file := range files {
		srcPath := templateDir + "/" + file.Name()
		destPath := data.OutputDir + "/" + file.Name()

		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("%s exists. Overwrite? (y|N): ", destPath)
			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" && response != "yes" {
				continue
			}
		}

		if file.IsDir() {
			dataCopy := RenderData{}
			// copy(dataCopy, data)
			// TODO: better copy
			dataCopy.Name = data.Name
			dataCopy.OutputDir = fmt.Sprintf("%s/%s", data.OutputDir, file.Name())

			errRender := Render(fmt.Sprintf("%s/%s", templateDir, file.Name()), &dataCopy)
			if errRender != nil {
				return errRender
			}
		} else {

			contents, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			err = WriteToFile(destPath, contents, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func WriteToFile(filePath string, template []byte, data *RenderData) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	engine := liquid.NewEngine()

	var bindings map[string]interface{}
	inrec, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		return errMarshal
	}
	json.Unmarshal(inrec, &bindings)

	out, renderErr := engine.ParseAndRender(template, bindings)
	if renderErr != nil {
		return renderErr
	}

	// hack hack hack
	// find-replace { { with {{ and } } with }}
	newSlice := bytes.Replace(out, []byte("{ {"), []byte("{{"), -1)
	lastSlice := bytes.Replace(newSlice, []byte("} }"), []byte("}}"), -1)

	return os.WriteFile(filePath, lastSlice, 0600)
}
