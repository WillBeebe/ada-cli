package ada

type AdaFile struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`

	Metadata    AdaMetadata
	Lint        AdaLint        `yaml:"lint"`
	Install     AdaInstall     `yaml:"install"`
	InstallDeps AdaInstallDeps `yaml:"installDeps"`
	Tests       AdaTests       `yaml:"tests"`
	Port        int
}

type AdaMetadata struct {
	GarImageName         string
	ArtifactoryImageName string
	ProjectId            string
	ProjectNumber        string
	GCP                  GCPMetadata
	Team                 string `yaml:"team"`
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
}

type AdaUnitTests struct {
	Command  string `yaml:"command"`
	Disabled bool   `yaml:"disabled"`
}

type GCPMetadata struct {
	ProjectID string
}
