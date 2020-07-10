package configuration

import (
	"cloud.google.com/go/civil"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

type Configuration struct {
	Profile         string
	Command         string
	Interval        time.Duration
	PreferedRunTime civil.Time
}

func ParseConfigFromFlags() Configuration {
	configuration := Configuration{}
	flag.StringVar(&configuration.Profile, "profile", "default", "Profile to use. If emtpy defaults to 'default'")
	flag.StringVar(&configuration.Command, "command",
		"echo 'daily-run-wrapper has run echo printing this text'",
		"The command that runner will execute")
	flag.DurationVar(&configuration.Interval, "interval",
		4*time.Minute,
		"The interval that daily-run-wrapper will use to check if it needs to run. "+
			"Can accept values acceptable to golang time.ParseDuration function")

	preferedRunTimePtr := flag.String("preferedTime", "",
		"Set a preferred time for the runner to run command. This time overrides the daily logic and the command will always "+
			"run if the system is up at that time.")

	flag.Parse()
	isValid, message := validateProfile(configuration.Profile)
	if !isValid {
		fmt.Print(message)
		os.Exit(2)
	}

	parsedCivilTime, err2 := civil.ParseTime(*preferedRunTimePtr)
	if err2 != nil {
		fmt.Printf("Parsing PreferedRunTime flag failed with error: %v", err2)
		os.Exit(2)
	}
	configuration.PreferedRunTime = parsedCivilTime
	return configuration
}

func validateProfile(profile string) (bool, string) {
	if profile == "default" {
		return false, "Wrong profile name argument provided 'default'. 'default' is reserved and cannot be used. " +
			"Please provide a different name. Exiting"
	}

	matched, _ := regexp.MatchString(`^[0-9a-zA-Z.\-_]*$`, profile)
	if !matched {
		return false, fmt.Sprintf("Wrong profile name argument provided: '%v'. Please provide a name containing only latin "+
			"letters, numbers and the characters '.-_'. Exiting", profile)
	}

	return true, ""
}
