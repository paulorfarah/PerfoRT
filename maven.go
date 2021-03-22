package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func GetMavenDependenciesClasspath(path string) string {
	found := false
	classpath := ""
	logfile := "maven-classpath.log"

	cmd := exec.Command("mvn", "dependency:build-classpath") //, " > "+logfile)
	cmd.Dir = path
	// var out bytes.Buffer
	var stderr bytes.Buffer
	// cmd.Stdout = &out
	cmd.Stderr = &stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// open the out file for writing
	outfile, err := os.Create(path + string(os.PathSeparator) + logfile)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	writer := bufio.NewWriter(outfile)
	defer writer.Flush()
	io.Copy(writer, stdoutPipe)

	err = cmd.Run()
	if err != nil {
		fmt.Println("[>>ERROR]: Error getting maven dependencies classpath: ", err.Error())
		fmt.Println("Dir: " + path + " Command: " + "mvn dependency:build-classpath > " + logfile)
		fmt.Printf("%s\n", stderr.String())
	}

	f, err := os.Open(path + string(os.PathSeparator) + logfile)
	if err != nil {
		fmt.Println("[>>ERROR]: There has been an error getting maven dependencies classpath!: ", err.Error())
		fmt.Println("log file: " + path + string(os.PathSeparator) + logfile)
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
