package container

import (
	"fmt"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/cmd"
)

// .PHONY: docker-push-gar
// docker-push-gar: ## push docker image to Google Artifact Registry
//
//  docker tag $(APP_NAME):$(VERSION) $(GAR_IMAGE_NAME):$(VERSION)
//  docker push $(GAR_IMAGE_NAME):$(VERSION)
func Push(AdaFile *ada.AdaFile) error {
	versionTag := ada.GetGitVersionTag()
	gitSHA := ada.GetGitShortSHA()
	command := fmt.Sprintf("docker tag %s:%s %s:%s", AdaFile.Name, gitSHA, AdaFile.Metadata.GarImageName, versionTag)
	output, err := cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		logger.Error(string(output))
		return err
	}
	logger.Info(string(output))
	command = fmt.Sprintf("docker tag %s:%s %s:%s", AdaFile.Name, gitSHA, AdaFile.Metadata.GarImageName, gitSHA)
	output, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		logger.Error(string(output))
		return err
	}
	logger.Info(string(output))
	command = fmt.Sprintf("docker tag %s:%s %s:local-latest", AdaFile.Name, gitSHA, AdaFile.Metadata.GarImageName)
	output, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		logger.Error(string(output))
		return err
	}
	logger.Info(string(output))
	command = fmt.Sprintf("docker push %s:%s", AdaFile.Metadata.GarImageName, versionTag)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		logger.Error(string(output))
		return err
	}
	logger.Info(string(output))
	command = fmt.Sprintf("docker push %s:%s", AdaFile.Metadata.GarImageName, gitSHA)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		logger.Error(string(output))
		return err
	}
	logger.Info(string(output))
	command = fmt.Sprintf("docker push %s:local-latest", AdaFile.Metadata.GarImageName)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		logger.Error(string(output))
		return err
	}
	logger.Info(string(output))
	return nil
}

// .PHONY: docker-push-artifactory
// docker-push-artifactory: ## push docker image to artifactory
//
//    jf docker tag $(APP_NAME):$(VERSION) $(ARTIFACTORY_IMAGE_NAME):$(VERSION)
//   jf rt dp $(ARTIFACTORY_IMAGE_NAME):$(VERSION) \
//      docker-canoo-local-dev \
//      --build-name=$(APP_NAME)-docker-build \
//      --build-number=$(VERSION)
//
// jf rt bp $(APP_NAME)-docker-build $(VERSION)
func PushArtifactory(AdaFile *ada.AdaFile) error {
	versionTag := ada.GetGitVersionTag()
	gitSHA := ada.GetGitShortSHA()
	command := fmt.Sprintf("jf docker tag %s:%s %s:%s", AdaFile.Name, gitSHA, AdaFile.Metadata.ArtifactoryImageName, versionTag)
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	command = fmt.Sprintf("jf docker tag %s:%s %s:%s", AdaFile.Name, gitSHA, AdaFile.Metadata.ArtifactoryImageName, gitSHA)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	command = fmt.Sprintf("jf docker tag %s:%s %s:local-latest", AdaFile.Name, gitSHA, AdaFile.Metadata.ArtifactoryImageName)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	command = fmt.Sprintf("jf rt dp %s:%s docker-canoo-local-dev --build-name=%s-docker-dev --build-number=%s", AdaFile.Metadata.ArtifactoryImageName, gitSHA, AdaFile.Name, gitSHA)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	command = fmt.Sprintf("jf rt dp %s:%s docker-canoo-local-dev --build-name=%s-docker-dev --build-number=%s", AdaFile.Metadata.ArtifactoryImageName, versionTag, AdaFile.Name, gitSHA)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	command = fmt.Sprintf("jf rt dp %s:local-latest docker-canoo-local-dev --build-name=%s-docker-dev --build-number=%s", AdaFile.Metadata.ArtifactoryImageName, AdaFile.Name, gitSHA)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	command = fmt.Sprintf("jf rt bag %s-docker-dev %s", AdaFile.Name, gitSHA)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	command = fmt.Sprintf("jf rt bp %s-docker-dev %s", AdaFile.Name, gitSHA)
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	if err != nil {
		return err
	}
	return err
}
