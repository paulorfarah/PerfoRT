package main

import (
	// "bytes"
	// "fmt"
	// "go-repo-downloader/models"
	// "log"
	// "os"
	// "os/exec"
	// "strings"
	// "github.com/wcharczuk/go-chart/v2"
	// "github.com/wcharczuk/go-chart/v2/drawing"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
)

func main() {
	// plotRandoopResults()

	// db := models.GetDB()
	// models.GetRandoopMetrics()
	// getReleaseList("junit-team", "junit4")
	MvnCompile("/mnt/sda4/go-work/src/github.com/paulorfarah/junit4")
}

func getReleaseList(owner, repo string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{Page: 2, PerPage: 10}
	releases, rsp, err := client.Repositories.ListReleases(ctx, owner, repo, opt)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n%+v\n", releases)
	// for _, r := range releases {
	// 	fmt.Println("--------------------------")
	// 	fmt.Println(*r.Name)
	// 	fmt.Println(r.PublishedAt.String())
	// 	fmt.Println(*r.URL)

	// }
	fmt.Printf("\n%+v\n", rsp)
}

func randoop() {
	fmt.Println("teste")
	// cmd := exec.Command("java", "-classpath", "/mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR", "randoop.main.Main", "gentests", "--testclass=testproject.Test")
	// script := CreateRandoopScript("testproject.Test")
	// cmd := exec.Command("bash " + script)
	c := "java -classpath /mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR randoop.main.Main gentests --testclass=testproject.Test > testproject.Test.txt"
	cmd := exec.Command("bash", "-c", c)
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
		// fmt.Println(ReadRandoopResults("testproject.Test.txt"))

	}
}

func CreateRandoopScript(class string) string {
	fn := strings.ReplaceAll(class, ".", "_") + ".sh"
	// Create new file
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = os.Chmod(fn, 0700)
	if err != nil {
		log.Fatal(err)
	}

	c := "java -classpath /mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR randoop.main.Main gentests --testclass=" + class + " > " + class + ".txt"
	_, err2 := f.WriteString(c)

	if err2 != nil {
		log.Fatal(err2)
	}

	return fn
}
