package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("teste")
	// cmd := exec.Command("java", "-classpath", "/mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR", "randoop.main.Main", "gentests", "--testclass=testproject.Test")
	script := CreateRandoopScript("testproject.Test")

	cmd := exec.Command(script)
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

func CreateRandoopScript(class string) string {
	f := strings.ReplaceAll(class, ".", "_") + ".sh"
	// Create new file
	new, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}
	defer new.Close()

	stats, err := os.Stat(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Permission File Before: %s\n", stats.Mode())
	err = os.Chmod(f, 0700)
	if err != nil {
		log.Fatal(err)
	}

	stats, err = os.Stat(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Permission File After: %s\n", stats.Mode())
	return f
}
