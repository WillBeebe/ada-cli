package ada

type AdaFile struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Team        string `yaml:"team"`
	Metadata    AdaMetadata
	Lint        AdaLint        `yaml:"lint"`
	Install     AdaInstall     `yaml:"install"`
	InstallDeps AdaInstallDeps `yaml:"installDeps"`
	Tests       AdaTests       `yaml:"tests"`
	Port        int
}

type AdaMetadata struct {
	GarImageName             string
	ArtifactoryImageName     string
	ProjectId                string
	ProjectNumber            string
	KubernetesDevClusterName string
	GCP                      GCPMetadata
}

type AdaLint struct {
	Disabled bool `yaml:"disabled"`
}

type AdaInstall struct {
	Disabled bool `yaml:"disabled"`
}

type AdaInstallDeps struct {
	Disabled bool `yaml:"disabled"`
}

type AdaTests struct {
	Unit AdaUnitTests `yaml:"unit"`
	E2E  AdaE2ETests  `yaml:"e2e"`
}

type AdaUnitTests struct {
	Command  string `yaml:"command"`
	Disabled bool   `yaml:"disabled"`
}

type AdaE2ETests struct {
	Command  string `yaml:"command"`
	Disabled bool   `yaml:"disabled"`
}

type GCPMetadata struct {
	DevProjectID string
	QAProjectID  string
	StgProjectID string
	PrdProjectID string
}
