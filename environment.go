package main

import (
	"os"
	"path"
	"path/filepath"
)

var globalProfile = "default"

func SetProfile(profile string) {
	globalProfile = profile
}

func GetProfile() string {
	return globalProfile
}

func GetLastRunFileName() string {
	return path.Join(localAppDataDir(), "last-run")
}

func localAppDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	nonProfilePath := path.Join(homeDir, ".local/share/daily-runner/")

	if globalProfile == "" {
		return nonProfilePath
	} else {
		return path.Join(nonProfilePath, globalProfile)
	}
}
func LockFilePath() string {
	return filepath.Join(os.TempDir(), "daily-runner"+globalProfile+".lck")
}

func WrapperLogFilePath() string {
	return path.Join(localAppDataDir(), "log/daily-runner.log")
}
func CommandLogFilePath() string {
	return path.Join(localAppDataDir(), "log/command-output.log")
}
