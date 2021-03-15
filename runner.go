package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("teste")
	cmd := exec.Command("java", "-classpath", "/mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR", "randoop.main.Main", "gentests", "--testclass=testproject.Test")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("\n[>>ERROR]: Cannot run randoop gentests: ", err)
		fmt.Println(out)
	} else {
		fmt.Println("\n [>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>SUCCESS]: Randoop executed successully!")
		fmt.Println(out)

	}
}
