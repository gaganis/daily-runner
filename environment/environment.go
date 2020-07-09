package environment

import (
	"os"
	"path"
)

func GetLastRunFileName() string {
	return path.Join(GetLocalAppDataDir(), "last-run")
}

func GetLocalAppDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homeDir, ".local/share/daily-run-wrapper/")
}
