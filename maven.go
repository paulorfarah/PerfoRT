package main

import (
	"bufio"
	"bytes"

	// "context"
	"PerfoRT/models"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	// "runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gorm.io/gorm"
)

type MvnTestResult struct {
	ClassName   string
	TestsRun    int
	Failures    int
	Errors      int
	Skipped     int
	TimeElapsed float64
}

func GetMavenDependenciesClasspath(path, javaVer string) string {
	logfile := "maven-classpath.log"

	fmt.Println("JAVA_HOME=" + javaVer + " && mvn dependency:build-classpath")
	cmd := exec.Command("mvn", "dependency:build-classpath")
	cmd.Dir = path
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JAVA_HOME="+javaVer)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed in GetMavenDependenciesClasspath with %s\n", err)
		fmt.Println(string(output))
	}
	// fmt.Printf("combined out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		panic(err)
	}

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
			if len(row) > 6 {
				if bytes.Equal(row[7:], []byte("Dependencies classpath:")) {
					found = true
				}
			}
		}
	}
	return classpath
}

func MvnCompile(path, javaVer string) bool {
	logfile := "maven-compiler.log"

	log.Println("mvn -fn -Drat.skip=true clean compile")
	cmd := exec.Command("JAVA_HOME="+javaVer, "&&", "mvn", "-fn", "-Drat.skip=true", "clean", "compile")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("mvn -Drat.skip=true clean compile failed with %s\n", err)
		fmt.Printf("cmd.Run() failed with %v\n", err)
		log.Printf("Compilation out:\n%s\n", string(output))
		// fmt.Printf("Compilation out:\n%s\n", string(output))
		return false
	}
	// log.Printf("Compilation out:\n%s\n", string(output))
	// fmt.Printf("Compilation out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		log.Printf(err.Error())
		fmt.Printf(err.Error())
		return false
	}
	if strings.Contains(string(output), "BUILD SUCCESS") {
		return true
	} else {
		return false
	}
}

func MvnTest(db *gorm.DB, path, javaVer string, measurementID, commitID uint) bool {

	ok := true
	logfile := "maven-test.log"
	var output []byte
	var err error

	// ATTENTION ABOUT JACOCO:
	//
	// 1) configure export MAVEN_OPTS="-javaagent:/home/usuario/go-work/src/github.com/paulorfarah/PerfoRT/jacoco-0.8.6/jacocoagent.jar"
	//    It will generate jacoco.exec in the application folder
	// deprecated:
	// jacoco_exec := localpath + "/coverage/jacoco-" + strconv.Itoa(int(commitID)) + ".exec"
	// testStr := "- mvn -fn -Drat.skip=true -Djacoco.destFile=" + jacoco_exec + " clean org.jacoco:jacoco-maven-plugin:0.8.7:prepare-agent test"
	// log.Println(testStr)
	// fmt.Println(testStr)
	// cmd := exec.Command("mvn", "-fn", "-Drat.skip=true", "-Djacoco.destFile="+jacoco_exec, "clean", "org.jacoco:jacoco-maven-plugin:0.8.7:prepare-agent", "test")

	// coverage needs module, so can't collect coverage in this func

	cmd := exec.Command("mvn", "-fn", "-Drat.skip=true", "clean", "test")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JAVA_HOME="+javaVer)
	// cmd := exec.Command("mvn", "-fn", "-Drat.skip=true", "test")
	// log.Println("- mvn -fn -Drat.skip=true test")
	// fmt.Println("- mvn -fn -Drat.skip=true test")

	// fmt.Println("path: ", path)
	cmd.Dir = path

	// output, err = cmd.CombinedOutput()
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("mvn -Drat.skip=true test failed with %s\n", err.Error())
		fmt.Printf("mvn -Drat.skip=true test failed with %s\n", err.Error())
	}

	fmt.Printf("Mvn test out:\n%s\n", string(output))
	// log.Printf("Mvn test out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		ok = false
		// panic(err)
		log.Println("ERROR writing results of maven test to file: ", err)
		fmt.Println("ERROR writing results of maven test to file: ", err)
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
	// return readMavenTestResults(path), ok

	return ok
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
				if len(res) >= 4 {
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
	}
	return tests
}

func MvnInstall(path, javaVer string) bool {
	fmt.Println("MvnInstall: " + javaVer)
	cmd := exec.Command("mvn", "-fn", "-Drat.skip=true", "clean", "install")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JAVA_HOME="+javaVer)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("mvn -fn -Drat.skip=true clean install failed with %s\n", err)
		fmt.Printf("mvn -fn -Drat.skip=true clean install failed with %s\n", err)
		log.Printf("Compilation out:\n%s\n", string(output))
		// fmt.Printf("Compilation out:\n%s\n", string(output))
		return false
	}
	log.Println("mvn -fn -Drat.skip=true clean install failed")
	log.Printf("Compilation out:\n%s\n", string(output))
	// fmt.Printf("Compilation out:\n%s\n", string(output))
	if strings.Contains(string(output), "BUILD SUCCESS") {
		return true
	} else {
		return false
	}
}

// type MavenTestSuites struct {
// 	// XMLName    xml.Name         `xml:"name"`
// 	TestSuites []MavenTestSuite `xml:"testsuites"`
// }

type MavenTestSuite struct {
	XMLName xml.Name `xml:"testsuite"`
	// Properties []MavenTestProperty `xml:"properties"`
	TestCases []MavenTestCase `xml:"testcase"`
}

// type MavenTestProperty struct {
// 	// XMLName xml.Name `xml:"property"`
// 	Name  string `xml:"name,attr"`
// 	Value string `xml:"value,attr"`
// }

type MavenTestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Name      string   `xml:"name,attr"`
	ClassName string   `xml:"classname,attr"`
	Time      string   `xml:"time,attr"`
}

func ParseMavenTestResults(f string) MavenTestSuite {
	// fmt.Println("ParseMavenTestResults: ", f)
	xmlFile, err := os.Open(f)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	rawData, _ := ioutil.ReadAll(xmlFile)

	var testSuite MavenTestSuite
	xml.Unmarshal(rawData, &testSuite)
	return testSuite
}

func RunMavenTestCase(db *gorm.DB, path, module string, tc *models.TestCase, measurementID uint, commit models.Commit) {
	// # Executes a single specified test in SomeTestClass
	// Test only a certain testcase inside test class with
	// “mvn  -Dtest=TestSurefire#testcaseFirst test“ (-pl module)
	// This command will execute only single test case method i.e. testcaseFirst().

	// logfile := "maven-test.log"
	resultsPath := path

	className := tc.ClassName[strings.LastIndex(tc.ClassName, ".")+1:]
	testName := tc.Name[strings.LastIndex(tc.Name, ".")+1:]

	//set environment variable to activate profiler during testcases execution
	localpath, errPath := os.Getwd()
	if errPath != nil {
		log.Println(errPath)
		fmt.Println("error getting current path: ", errPath.Error())
	}
	// profiler_output := localpath + string(os.PathSeparator) + "profiler" + string(os.PathSeparator) + module + "_" + className + "_" + testName + ".txt"
	// os.Setenv("MAVEN_OPTS", "-agentpath:async-profiler-2.5.1-linux-x64/build/libasyncProfiler.so=start,event=wall,file="+profiler_output)
	// fmt.Println("export MAVEN_OPTS=-agentpath:" + localpath + string(os.PathSeparator) + "async-profiler-2.5.1-linux-x64/build/libasyncProfiler.so=start,event=wall,file=" + profiler_output)
	// cmdEnv := exec.Command("bash", "-c", "export", "MAVEN_OPTS=-agentpath:"+localpath+string(os.PathSeparator)+"async-profiler-2.5.1-linux-x64/build/libasyncProfiler.so=start,event=wall,file="+profiler_output)
	// _, errEnv := cmdEnv.Output()

	// if errEnv != nil {
	// 	fmt.Println("Error exporting environment variable MAVEN_OPTS: ", errEnv.Error())
	// 	return
	// }

	// log.Println(">>>------------------------------------------------ maven testcase", path, className, testName)
	// fmt.Println(">>>------------------------------------------------ maven testcase", path, className, testName)

	var cmd *exec.Cmd
	var cmdStr string
	// fmt.Println(path)
	param := "-Dtest=" + className + "#" + testName
	if module != "" {
		cmdStr = "mvn -Drat.skip=true test  -pl " + module + " " + param
		// fmt.Println(cmdStr)
		// log.Println(cmdStr)

		cmd = exec.Command("mvn", "-Drat.skip=true", "test", "-pl", module, param)
		resultsPath += "/" + module
	} else {
		cmdStr = "mvn -Drat.skip=true test " + param
		fmt.Println(cmdStr)
		log.Println(cmdStr)
		cmd = exec.Command("mvn", "test", param)
	}

	resultsPath += "/target/surefire-reports/"
	// resultsPath += discoverTestFilename(resultsPath, tc.ClassName)
	resultsPath += "TEST-" + tc.ClassName + ".xml"

	log.Println("resultsPath: ", resultsPath)
	cmd.Dir = path
	// fmt.Println("path: ", path)

	// var output []byte
	var err error

	mr := &models.Run{
		MeasurementID: measurementID,
		Type:          "maven",
		TestCaseID:    tc.ID,
	}
	models.CreateRun(db, mr)

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command: ", err.Error())
		log.Fatal(err)
	}
	pid := cmd.Process.Pid

	monitoringTime := 1.0
	monitoringTimeStr, ok := os.LookupEnv("monitoring_time")
	if ok {
		monitoringTime, _ = strconv.ParseFloat(monitoringTimeStr, 32)
	}
	stop := make(chan bool)
	go func() {
		defer close(stop)
		resources := []models.Resource{}
		for {
			select {
			case <-stop:
				db.CreateInBatches(resources, 1000)
				return
			default:
				resource, err := MonitorProcess(pid, mr.ID)
				if err == nil {
					resources = append(resources, resource)

				}
				time.Sleep(time.Duration(monitoringTime) * time.Second)
			}
		}
	}()

	err = cmd.Wait()

	stop <- true

	if err != nil {
		pid = cmd.Process.Pid
		process, err := os.FindProcess(int(pid))
		if err != nil {
			fmt.Printf("Failed to find process: %s\n", err)
		} else {
			errPid := process.Signal(syscall.Signal(0))
			fmt.Printf("process.Signal on pid %d returned: %v\n", pid, errPid)
			resPid := fmt.Sprintf("%v", errPid)
			if resPid != "os: process already finished" {
				fmt.Printf("maven test failed with %s\n", err.Error())
				log.Printf("Command finished with error: %s", err.Error())
			}
		}
	}

	profiler_output := localpath + string(os.PathSeparator) + "profiler" + string(os.PathSeparator) + fmt.Sprint(commit.CommitHash) + "_"

	if module != "" {
		profiler_output += module + "_"
	}

	profiler_output += className + "_" + testName //+ ".jfr"

	// jcmd 56822 JFR.start filename=teste.jfr

	// strPid := fmt.Sprintf("%s", pid)
	// cmdJCMD := exec.Command("jcmd", strPid, "JFR.start", "filename=/home/farah/go-work/src/github.com/paulorfarah/go-repo-downloader/profiler/test.jfr") //+profiler_output)

	// if err := cmdJCMD.Run(); err != nil {
	// 	log.Printf("Failed to start cmd: %v", err)
	// 	return
	// }

	// stop2 := make(chan bool)
	// go func() {

	// 	for {
	// 		select {
	// 		case <-stop2:
	// 			strPid := fmt.Sprintf("%s", pid)
	// 			cmd := exec.Command("jcmd", fmt.Sprintf("%s", strPid), "JFR.stop")

	// 			if err := cmd.Run(); err != nil {
	// 				log.Printf("Failed to stop cmd: %v", err)
	// 				return
	// 			}
	// 			return
	// 		default:
	// 			// jcmd 7060 JFR.start name=MyRecording settings=profile delay=20s duration=2m filename=C:\TEMP\myrecording.jfr
	// 			strPid := fmt.Sprintf("%s", pid)
	// 			cmd := exec.Command("jcmd", strPid, "JFR.start", "name=teste", "settings=profile", "delay=1s", "filename=/home/farah/go-work/src/github.com/paulorfarah/go-repo-downloader/profiler/test.jfr") //+profiler_output)

	// 			if err := cmd.Run(); err != nil {
	// 				log.Printf("Failed to start cmd: %v", err)
	// 				return
	// 			}
	// 		}
	// 	}
	// }()

	// err2 := cmd.Wait()

	// stop2 <- true

	// if err2 != nil {
	// 	pid = cmd.Process.Pid
	// 	// fmt.Println(pid)
	// 	process, err := os.FindProcess(int(pid))
	// 	if err != nil {
	// 		fmt.Printf("Failed to find process: %s\n", err)
	// 	} else {
	// 		errPid := process.Signal(syscall.Signal(0))
	// 		fmt.Printf("process.Signal on pid %d returned: %v\n", pid, errPid)
	// 		resPid := fmt.Sprintf("%v", errPid)
	// 		if resPid != "os: process already finished" {
	// 			fmt.Printf("maven test failed with %s\n", err.Error())
	// 			log.Printf("Command finished with error: %s", err.Error())
	// 		}
	// 	}
	// }

	// ParseProfilingClock(db, commit, *tc, profiler_output)
	// ParseProfilingAlloc(db, commit, *tc, profiler_output)

	// fmt.Printf("Mvn test out:\n%s\n", string(output))
	// log.Printf("gradle test out:\n%s\n", string(output))
	// err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	// if err != nil {
	// 	// ok = false
	// 	fmt.Println("ERROR writing logfile: ", err.Error())
	// 	log.Println("ERROR writing logfile: ", err.Error())
	// 	// panic(err)
	// }
	suite := ParseMavenTestResults(resultsPath)
	for _, test := range suite.TestCases {
		if test.Name == tc.Name {
			// fmt.Printf("  %s\n", test.Name)
			// fmt.Printf("time:  %s\n", test.Time)

			dur, errD := time.ParseDuration(strings.Replace(test.Time, ",", "", -1) + "s")
			if errD != nil {
				fmt.Println("ERROR parsing test time to duration: ", errD.Error())
			}
			// fmt.Printf("duration: %s\n", dur)
			mr.TestCaseTime = dur
			err := models.SaveRun(db, mr)
			if err != nil {
				fmt.Println("ERROR saving run: ", err.Error())
			}
		}
	}
}

func discoverTestFilename(path, className string) string {
	// TEST-" + tc.ClassName + ".xml"
	prefixes := [1]string{"TEST"}
	suffixes := [1]string{"xml"}

	for _, pref := range prefixes {
		for _, suf := range suffixes {
			testFile := pref + className + "." + suf
			if _, err := os.Stat(path + testFile); err == nil {
				log.Println("Test filename found: " + path + testFile)
				return testFile
			}
			log.Println("Test filename not found: " + path + testFile)
		}
	}

	return ""
}
