package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type MvnTestResult struct {
	ClassName   string
	TestsRun    int
	Failures    int
	Errors      int
	Skipped     int
	TimeElapsed float64
}

func GetMavenDependenciesClasspath(path string) string {
	logfile := "maven-classpath.log"

	fmt.Println("mvn dependency:build-classpath")
	cmd := exec.Command("mvn", "dependency:build-classpath")
	cmd.Dir = path
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

	return getClasspath(path)
}

func getClasspath(path string) string {
	found := false
	classpath := ""
	logfile := "maven-classpath.log"
	f, err := os.Open(path + string(os.PathSeparator) + logfile)
	if err != nil {
		fmt.Println("[>>ERROR]: There has been an error getting maven dependencies classpath!: ", err.Error())
		fmt.Println("log file: " + path + string(os.PathSeparator) + logfile)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		row := scanner.Bytes()
		if len(row) > 5 {
			if bytes.Equal(row[:6], []byte("[INFO]")) {
				found = false
			}

			if found {
				classpath += strings.Trim(string(row), " ")
			}
			if bytes.Equal(row[7:], []byte("Dependencies classpath:")) {
				found = true
			}
		}
	}
	return classpath
}

func MvnCompile(path string) bool {
	logfile := "maven-compiler.log"

	fmt.Println("------------------------------------------------ mvn compile")
	cmd := exec.Command("mvn", "compile")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		fmt.Printf("cmd.Run() failed with %s\n", err)
		log.Printf("Compilation out:\n%s\n", string(output))
		fmt.Printf("Compilation out:\n%s\n", string(output))
		return false
	}
	log.Printf("Compilation out:\n%s\n", string(output))
	fmt.Printf("Compilation out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		log.Printf(err.Error())
		fmt.Printf(err.Error())
		return false
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
	return true
}

func MvnTest(path string) ([]MvnTestResult, bool) {
	ok := true
	logfile := "maven-test.log"

	log.Println("------------------------------------------------ mvn test")
	cmd := exec.Command("mvn", "-Drat.skip=true", "test")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("Mvn test out:\n%s\n", string(output))
	log.Printf("Mvn test out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		ok = false
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
	return readMavenTestResults(path), ok
}

func readMavenTestResults(path string) []MvnTestResult {
	logfile := "maven-test.log"
	f, err := os.Open(path + string(os.PathSeparator) + logfile)
	if err != nil {
		log.Println("[>>ERROR]: There has been an error running mvn test: ", err.Error())
		log.Println("log file: " + path + string(os.PathSeparator) + logfile)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var tests []MvnTestResult
	for scanner.Scan() {
		row := scanner.Bytes()
		elements := strings.Split(string(row), " ")
		if len(elements) > 9 {
			if bytes.Equal(row[:10], []byte("Tests run:")) ||
				bytes.Equal(row[:17], []byte("[INFO] Tests run:")) {
				cls := strings.Split(string(row), " ")
				cl := cls[len(cls)-1]
				re := regexp.MustCompile("[0-9]+(.[0-9]+)*")
				res := re.FindAllString(string(row), -1)
				tr, err := strconv.Atoi(res[0])
				if err != nil {
					tr = -1
				}
				f, err := strconv.Atoi(res[1])
				if err != nil {
					f = -1
				}
				e, err := strconv.Atoi(res[2])
				if err != nil {
					e = -1
				}
				s, err := strconv.Atoi(res[3])
				if err != nil {
					s = -1
				}
				te, err := strconv.ParseFloat(res[4], 64)
				if err != nil {
					te = float64(-1.0)
				}

				test := &MvnTestResult{ClassName: cl,
					TestsRun:    tr,
					Failures:    f,
					Errors:      e,
					Skipped:     s,
					TimeElapsed: te}
				tests = append(tests, *test)

			}
		}
	}
	return tests
}
