package environment

import (
	"os"
	"path"
)

var globalProfile = "default"

func SetProfile(profile string) {
	globalProfile = profile
}

func GetLastRunFileName() string {
	return path.Join(GetLocalAppDataDir(), "last-run")
}

func GetLocalAppDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	nonProfilePath := path.Join(homeDir, ".local/share/daily-run-wrapper/")

	if globalProfile == "" {
		return nonProfilePath
	} else {
		return path.Join(nonProfilePath, globalProfile)
	}
}
