package ada

import "os"

func BuildVersion() string {
	versionTag := "local-latest"
	if value, ok := os.LookupEnv("GIT_VERSION_TAG"); ok {
		versionTag = value
	}
	return versionTag
}
func GetGitVersionTag() string {
	versionTag := "local-latest"
	if value, ok := os.LookupEnv("GIT_VERSION_TAG"); ok {
		versionTag = value
	}
	return versionTag
}
func GetGitShortSHA() string {
	gitSHA := "local-latest"
	if value, ok := os.LookupEnv("GIT_SHORT_SHA"); ok {
		gitSHA = value
	}
	return gitSHA
}
