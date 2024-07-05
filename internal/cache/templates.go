package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/container-labs/ada/internal/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const SoyTemplatesRepository = "ssh://github.com:7999/gst/ada-templates.git"

// const SoyTemplatesRepository = "https://github.com/massdriver-cloud/application-templates"

var logger = common.Logger()

// TemplateCacheDir is a reader function to access the local cache of templates.
// When developing templates, the cache source can be overwritten for reads by setting `MD_DEV_TEMPLATES_PATH`
func TemplateCacheDir() string {
	var templatesPath string
	localDevTemplatesPath := os.Getenv("ada_DEV_TEMPLATES_PATH")
	if localDevTemplatesPath == "" {
		dir, _ := templateCacheDir()
		templatesPath = dir
	} else {
		logger.Info(fmt.Sprintf("Reading templates for local development path: %s", localDevTemplatesPath))
		templatesPath = localDevTemplatesPath
	}

	return templatesPath
}

func Templates() ([]string, error) {
	templates := []string{}
	templateDirs, err := ioutil.ReadDir(TemplateCacheDir())
	if err != nil {
		return templates, err
	}

	for _, f := range templateDirs {
		// all directories that aren't .git
		// cheap way of listing templates
		if f.IsDir() && f.Name() != ".git" {
			templates = append(templates, f.Name())
		}
	}
	return templates, nil
}

func RefreshTemplates() error {
	if err := clearTemplateCache(); err != nil {
		return err
	}
	return downloadTemplates()
}

// templateCacheDir is the actual cache directory. This should be used internally when managing
// files so that development template directories aren't overwritten on accident.
func templateCacheDir() (string, error) {
	cacheDir, err := cacheDir()
	if err != nil {
		return "", err
	}

	templateDir := filepath.Join(cacheDir, "templates")
	if _, errDir := os.Stat(templateDir); !os.IsNotExist(errDir) {
		return templateDir, errDir
	}

	if errMkdir := os.Mkdir(templateDir, 0755); errMkdir != nil {
		return templateDir, errMkdir
	}
	return templateDir, nil
}

func clearTemplateCache() error {
	templateCacheDir, _ := templateCacheDir()
	if err := os.RemoveAll(templateCacheDir); err != nil {
		return err
	}
	return nil
}

func downloadTemplates() error {
	templateCacheDir, _ := templateCacheDir()
	// log.Debug().Msgf("Downloading templates to %s", templateCacheDir)

	// begin: remove for public repo
	usr, err := user.Current()
	if err != nil {
		return err
	}
	homeDir := usr.HomeDir

	var publicKeys *ssh.PublicKeys
	publicKeys, err = ssh.NewPublicKeysFromFile("git", filepath.Join(homeDir, ".ssh/id_ed25519"), "")
	if err != nil {
		return err
	}
	// end: remove for public repo

	_, cloneErr := git.PlainClone(templateCacheDir, false, &git.CloneOptions{
		URL:  SoyTemplatesRepository,
		Auth: publicKeys,
		// Progress: os.Stdout,
		Depth: 1,
	})
	if cloneErr != nil {
		return cloneErr
	}
	return nil
}
