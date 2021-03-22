package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func GetMavenDependenciesClasspath(path string) string {
	found := false
	classpath := ""
	logfile := "maven-classpath.log"

	fmt.Println("mvn dependency:build-classpath")
	cmd := exec.Command("mvn", "dependency:build-classpath") //, " > "+logfile)
	cmd.Dir = path
	// var out bytes.Buffer
	// var stderr bytes.Buffer
	// cmd.Stdout = &out
	// cmd.Stderr = &stderr

	// err := cmd.Run()
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	// fmt.Printf("combined out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		panic(err)
	}
	// if err != nil {
	// 	fmt.Println("[>>ERROR]: Error getting maven dependencies classpath: ", err.Error())
	// 	fmt.Println("Dir: " + path + " Command: " + "mvn dependency:build-classpath > " + logfile)
	// } else {
	// 	fmt.Println("executed successfully")
	// }
	// fmt.Println("------")
	// fmt.Printf("%s\n", out.String())
	// fmt.Println("^^^ out ^^^ - vvv error vvv")
	// fmt.Printf("%s\n", stderr.String())

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
