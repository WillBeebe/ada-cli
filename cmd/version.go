package cmd

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v45/github"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

const currentVersion = "1.0.0"
const repoOwner = "your-github-username-or-org"
const repoName = "your-repo-name"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of Ada CLI",
	Run: func(cmd *cobra.Command, args []string) {
		checkVersion()
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version of Ada CLI",
	Run: func(cmd *cobra.Command, args []string) {
		upgradeVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upgradeCmd)
}

func checkVersion() {
	fmt.Printf("Current version: %s\n", currentVersion)
	latestVersion, err := getLatestVersion()
	if err != nil {
		fmt.Println("Error fetching the latest version:", err)
		return
	}

	currentSemVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		fmt.Println("Error parsing current version:", err)
		return
	}

	latestSemVer, err := semver.NewVersion(latestVersion)
	if err != nil {
		fmt.Println("Error parsing latest version:", err)
		return
	}

	if currentSemVer.LessThan(latestSemVer) {
		fmt.Printf("A new version is available: %s\n", latestVersion)
	} else {
		fmt.Println("You are using the latest version.")
	}
}

func upgradeVersion() {
	latestVersion, err := getLatestVersion()
	if err != nil {
		fmt.Println("Error fetching the latest version:", err)
		return
	}

	currentSemVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		fmt.Println("Error parsing current version:", err)
		return
	}

	latestSemVer, err := semver.NewVersion(latestVersion)
	if err != nil {
		fmt.Println("Error parsing latest version:", err)
		return
	}

	if currentSemVer.LessThan(latestSemVer) {
		fmt.Printf("Upgrading to version %s\n", latestVersion)
		if err := downloadAndUpdate(latestVersion); err != nil {
			fmt.Println("Error upgrading:", err)
			return
		}
		fmt.Println("Upgrade successful!")
	} else {
		fmt.Println("You are already using the latest version.")
	}
}

func getLatestVersion() (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, repoOwner, repoName)
	if err != nil {
		return "", err
	}

	return release.GetTagName(), nil
}

func downloadAndUpdate(version string) error {
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/ada-%s-%s", repoOwner, repoName, version, version, runtime.GOOS, runtime.GOARCH)
	fmt.Println("Downloading from:", url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling here, e.g. rollback
		if rerr := update.RollbackError(err); rerr != nil {
			return fmt.Errorf("rollback failed: %v", rerr)
		}
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}
