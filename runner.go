package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("teste")
	// cmd := exec.Command("java", "-classpath", "/mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR", "randoop.main.Main", "gentests", "--testclass=testproject.Test")
	cmd := exec.Command("randoop.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("\n[>>ERROR]: Cannot run randoop gentests (" + fmt.Sprint(err) + "): " + stderr.String())
		fmt.Println(out)
	} else {
		fmt.Println("\n [>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>SUCCESS]: Randoop executed successully!")
		fmt.Println(out.String())

	}
}
