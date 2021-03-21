package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetMavenDependenciesClasspath(path string) string {
	found := false
	classpath := ""
	logfile := "maven-classpath.log"
	cmd := exec.Command("mvn", "dependency:build-classpath", ">", logfile)
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("[>>ERROR]: Error getting maven dependencies classpath: ", err.Error())
	}

	f, err := os.Open(path + string(os.PathSeparator) + logfile)
	if err != nil {
		fmt.Println("[>>ERROR]: There has been an error getting maven dependencies classpath!: ", err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 5 {
			if bytes.Equal(line[:6], []byte("[INFO]")) {
				found = false
			}

			if found {
				classpath += strings.Trim(string(line), " ")
			}
			if bytes.Equal(line[7:], []byte("Dependencies classpath:")) {
				found = true
			}
		}
	}
	return classpath
}
