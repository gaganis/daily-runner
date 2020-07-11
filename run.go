package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func WrappedCommand(configuration Configuration) bool {
	logFile, err := openCommandLogFile()
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := logFile.Close(); e != nil {
			panic(e)
		}
	}()

	log.Printf("Running process")
	splitCommand := strings.Split(configuration.Command, " ")
	command := exec.Command(splitCommand[0], splitCommand[1:]...)

	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderrPipe, err := command.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err = command.Start(); err != nil {
		panic(err)
	}

	_, err = io.Copy(logFile, stdoutPipe)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(logFile, stderrPipe)
	if err != nil {
		panic(err)
	}

	err = command.Wait()
	log.Print(err)

	return err == nil
}

func openCommandLogFile() (*os.File, error) {
	logFilePath := CommandLogFilePath()
	if err := os.MkdirAll(path.Dir(logFilePath), 0755); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		panic(err)
	}

	return logFile, err
}
