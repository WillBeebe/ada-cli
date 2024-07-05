package internal

import (
	"fmt"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/golang"
	"github.com/container-labs/ada/internal/nodejs"
	"github.com/container-labs/ada/internal/python"
	"github.com/container-labs/ada/internal/terraform"
)

// LanguageStrategy defines the interface for language-specific operations
type LanguageStrategy interface {
	Install() error
	AddDependency(dep string) error
	Start() error // New method for starting the application
}

// LanguageFactory is responsible for creating the appropriate LanguageStrategy
func LanguageFactory(adaFile *ada.AdaFile) (LanguageStrategy, error) {
	switch adaFile.Type {
	case "python":
		return &PythonStrategy{}, nil
	case "terraform":
		return &TerraformStrategy{}, nil
	case "nodejs":
		return &NodeJSStrategy{}, nil
	case "go":
		return &GoStrategy{}, nil
	default:
		return nil, fmt.Errorf("unsupported language type: %s", adaFile.Type)
	}
}

// Implement strategies for each language

type PythonStrategy struct{}

func (p *PythonStrategy) Install() error {
	return python.Install()
}

func (p *PythonStrategy) AddDependency(dep string) error {
	return python.AddDependency(dep)
}

func (p *PythonStrategy) Start() error {
	return python.Start() // Assuming this method exists in the python package
}

type TerraformStrategy struct{}

func (t *TerraformStrategy) Install() error {
	return terraform.Install()
}

func (t *TerraformStrategy) AddDependency(dep string) error {
	return terraform.AddDependency(dep)
}

func (t *TerraformStrategy) Start() error {
	return terraform.Start() // Assuming this method exists in the terraform package
}

type NodeJSStrategy struct{}

func (n *NodeJSStrategy) Install() error {
	return nodejs.Install()
}

func (n *NodeJSStrategy) AddDependency(dep string) error {
	return nodejs.AddDependency(dep)
}

func (n *NodeJSStrategy) Start() error {
	return nodejs.Start() // Assuming this method exists in the nodejs package
}

type GoStrategy struct{}

func (g *GoStrategy) Install() error {
	return golang.Install()
}

func (g *GoStrategy) AddDependency(dep string) error {
	return golang.AddDependency(dep)
}

func (g *GoStrategy) Start() error {
	return golang.Start() // Assuming this method exists in the golang package
}
