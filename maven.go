package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"go-repo-downloader/models"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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
	cmd := exec.Command("mvn", "-Drat.skip=true", "clean", "compile")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("mvn -Drat.skip=true clean compile failed with %s\n", err)
		fmt.Printf("cmd.Run() failed with %v\n", err)
		log.Printf("Compilation out:\n%s\n", string(output))
		// fmt.Printf("Compilation out:\n%s\n", string(output))
		return false
	}
	log.Printf("Compilation out:\n%s\n", string(output))
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

func MvnTest(db *gorm.DB, path string, measurementID uint) bool {
	ok := true
	logfile := "maven-test.log"

	log.Println("------------------------------------------------ mvn test")
	fmt.Println("------------------------------------------------ mvn test")
	// cmd := exec.Command("mvn", "-fn", "-Drat.skip=true", "-Djacoco.destFile=jacoco.exec", "clean", "org.jacoco:jacoco-maven-plugin:0.7.8:prepare-agent", "test")
	cmd := exec.Command("mvn", "-fn", "-Drat.skip=true", "clean", "test")
	cmd.Dir = path

	var output []byte
	var err error

	// output, err = cmd.CombinedOutput()
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	// pid := cmd.Process.Pid

	// stop := make(chan bool)
	// go func() {
	// 	perfMetrics := []PerfMetrics{}
	// 	for {
	// 		select {
	// 		case <-stop:
	// 			//save
	// 			for _, perfMetric := range perfMetrics {
	// 				mr := &models.Run{
	// 					MeasurementID: measurementID,
	// 					Type:          "maven",
	// 					Resources: models.Resources{
	// 						CpuPercent:        perfMetric.CpuPercent,
	// 						MemPercent:        perfMetric.MemoryPercent,
	// 						MemoryInfoStat:    *perfMetric.MemoryInfo,
	// 						IOCountersStat:    *perfMetric.IOCounters,
	// 						PageFaultsStat:    *perfMetric.PageFaults,
	// 						AvgStat:           *perfMetric.Load,
	// 						VirtualMemoryStat: *perfMetric.VirtualMemory,
	// 						SwapMemory:        *perfMetric.SwapMemory,
	// 						// CPUTime:           perfMetric.CPUTime,
	// 						// DiskIOCounters:    perfMetric.DiskIOCounters,
	// 						// NetIOCounters:     perfMetric.NetIOCounters,
	// 					},
	// 				}
	// 				models.CreateRun(db, mr)
	// 				for _, cpuTime := range perfMetric.CPUTimes {
	// 					models.CreateCPUTimes(db, &models.CPUTimes{
	// 						RunID:     mr.ID,
	// 						CPU:       cpuTime.CPU,
	// 						User:      cpuTime.User,
	// 						System:    cpuTime.System,
	// 						Idle:      cpuTime.Idle,
	// 						Nice:      cpuTime.Nice,
	// 						Iowait:    cpuTime.Iowait,
	// 						Irq:       cpuTime.Irq,
	// 						Softirq:   cpuTime.Softirq,
	// 						Steal:     cpuTime.Steal,
	// 						Guest:     cpuTime.Guest,
	// 						GuestNice: cpuTime.GuestNice,
	// 					})
	// 				}

	// 				for i, diskIOCounter := range perfMetric.DiskIOCounters {
	// 					models.CreateDiskIOCounters(db, &models.DiskIOCounters{
	// 						RunID:            mr.ID,
	// 						Device:           i,
	// 						ReadCount:        diskIOCounter.ReadCount,
	// 						MergedReadCount:  diskIOCounter.MergedReadCount,
	// 						WriteCount:       diskIOCounter.WriteCount,
	// 						MergedWriteCount: diskIOCounter.MergedWriteCount,
	// 						ReadBytes:        diskIOCounter.ReadBytes,
	// 						WriteBytes:       diskIOCounter.WriteBytes,
	// 						ReadTime:         diskIOCounter.ReadTime,
	// 						WriteTime:        diskIOCounter.WriteTime,
	// 						IopsInProgress:   diskIOCounter.IopsInProgress,
	// 						IoTime:           diskIOCounter.IoTime,
	// 						WeightedIO:       diskIOCounter.WeightedIO,
	// 						Name:             diskIOCounter.Name,
	// 						SerialNumber:     diskIOCounter.SerialNumber,
	// 						Label:            diskIOCounter.Label,
	// 					})
	// 				}

	// 				for i, netIOCounter := range perfMetric.NetIOCounters {
	// 					models.CreateNetIOCounters(db, &models.NetIOCounters{
	// 						RunID:       mr.ID,
	// 						NICID:       uint(i),
	// 						Name:        netIOCounter.Name,
	// 						BytesSent:   netIOCounter.BytesSent,
	// 						BytesRecv:   netIOCounter.BytesRecv,
	// 						PacketsSent: netIOCounter.PacketsSent,
	// 						PacketsRecv: netIOCounter.PacketsRecv,
	// 						Errin:       netIOCounter.Errin,
	// 						Errout:      netIOCounter.Errout,
	// 						Dropin:      netIOCounter.Dropin,
	// 						Dropout:     netIOCounter.Dropout,
	// 						Fifoin:      netIOCounter.Fifoin,
	// 						Fifoout:     netIOCounter.Fifoout,
	// 					})
	// 				}
	// 			}
	// 			return
	// 		default:
	// 			perfMetric, err := MonitorProcess(pid)
	// 			if err == nil {
	// 				perfMetrics = append(perfMetrics, perfMetric)
	// 			}

	// 		}
	// 	}
	// }()

	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v\n", err)
		fmt.Printf("Command finished with error: %v\n", err)
	}
	// stop <- true

	if err != nil {
		fmt.Printf("mvn -Drat.skip=true test failed with %s\n", err.Error())
	}

	// fmt.Printf("Mvn test out:\n%s\n", string(output))
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

func MvnInstall(path string) bool {
	logfile := "maven-install.log"

	fmt.Println("------------------------------------------------ mvn install")
	cmd := exec.Command("mvn", "-Drat.skip=true", "clean", "install")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("mvn -Drat.skip=true clean install failed with %s\n", err)
		fmt.Printf("mvn -Drat.skip=true clean install failed with %s\n", err)
		log.Printf("Compilation out:\n%s\n", string(output))
		// fmt.Printf("Compilation out:\n%s\n", string(output))
		return false
	}
	log.Printf("Compilation out:\n%s\n", string(output))
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

func RunMavenTestCase(db *gorm.DB, path, module string, tc *models.TestCase, measurementID uint) {
	// # Executes a single specified test in SomeTestClass
	// Test only a certain testcase inside test class with
	// “mvn  -Dtest=TestSurefire#testcaseFirst test“ (-pl module)
	// This command will execute only single test case method i.e. testcaseFirst().

	// ok := true
	logfile := "maven-test.log"
	// testName := tc.ClassName + "." + tc.Name
	resultsPath := path //+ "/build/test-results/test/TEST-" + tc.ClassName + ".xml"

	className := tc.ClassName[strings.LastIndex(tc.ClassName, ".")+1:]
	testName := tc.Name[strings.LastIndex(tc.Name, ".")+1:]

	log.Println(">>>------------------------------------------------ maven testcase", path, className, testName)
	fmt.Println(">>>------------------------------------------------ maven testcase", path, className, testName)

	var cmd *exec.Cmd
	var cmdStr string
	param := "-Dtest=" + className + "#" + testName
	if module != "" {
		cmdStr = "mvn test  -pl " + module + " " + param
		fmt.Println(cmdStr)

		cmd = exec.Command("mvn", "test", "-pl", module, param)
		resultsPath += "/" + module
	} else {
		cmdStr = "bash mvn test " + param
		fmt.Println(cmdStr)
		cmd = exec.Command("bash", "mvn", "test", param)
	}

	resultsPath += "/target/surefire-reports/TEST-" + tc.ClassName + ".xml"
	fmt.Println("resultsPath: ", resultsPath)
	cmd.Dir = path

	var output []byte
	var err error

	mr := &models.Run{
		MeasurementID: measurementID,
		Type:          "maven",
		TestCaseID:    tc.ID,
	}
	models.CreateRun(db, mr)

	err = cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	pid := cmd.Process.Pid

	stop := make(chan bool)
	go func() {
		// LOG_FILE := "/tmp/gorepodownloader_log"
		// // open log file
		// logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		// if err != nil {
		// 	log.Panic(err)
		// 	return
		// }
		// defer logFile.Close()
		// log.SetOutput(logFile)

		// optional: log date-time, filename, and line number
		// log.SetFlags(log.Lshortfile | log.LstdFlags)

		// log.Println("measurementID: ", measurementID)

		perfMetrics := []PerfMetrics{}
		for {
			select {
			case <-stop:
				// //save
				for _, perfMetric := range perfMetrics {
					saveMetrics(db, mr.ID, perfMetric)
				}
				return
			default:
				perfMetric, err := MonitorProcess(pid)
				if err == nil {
					perfMetrics = append(perfMetrics, perfMetric)
					// saveMetrics(db, mr.ID, perfMetric)

				}
				// log.Println(perfMetric)

			}
		}
	}()

	err = cmd.Wait()

	stop <- true

	if err != nil {
		fmt.Printf("maven test failed with %s\n", err.Error())
		log.Printf("Command finished with error: %s", err.Error())
	}

	// fmt.Printf("Mvn test out:\n%s\n", string(output))
	// log.Printf("gradle test out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		// ok = false
		fmt.Println("ERROR writing logfile: ", err.Error())
		panic(err)
	}
	suite := ParseMavenTestResults(resultsPath)
	for _, test := range suite.TestCases {
		if test.Name == tc.Name {
			// fmt.Printf("  %s\n", test.Name)
			// fmt.Printf("time:  %s\n", test.Time)
			// t, _ := strconv.ParseFloat(test.Time, 32)
			// fmt.Printf("float: %f\n", t)
			dur, errD := time.ParseDuration(test.Time + "s")
			if errD != nil {
				fmt.Println("ERROR parsing test time to duration: ", errD.Error())
			}
			// fmt.Printf("duration: %s\n", dur)
			mr.TestCaseTime = dur
			err := models.SaveRun(db, mr)
			if err != nil {
				fmt.Println("ERROR saving run: ", err.Error())
			}

			// if test.Error != nil {
			// 	fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
			// } else {
			// 	fmt.Printf("    %s\n", test.Status)
			// }
			// mr.TestCaseTime = test.Duration

		}
	}
	// resultsPath := path + "/build/test-results/test/TEST-" + tc.ClassName + ".xml"
	// // fmt.Println(resultsPath)
	// suites, err := junit.IngestFile(resultsPath)
	// if err != nil {
	// 	fmt.Printf("failed to ingest JUnit xml %v", err)
	// 	log.Fatalf("failed to ingest JUnit xml %v", err)
	// }
	// for _, suite := range suites {
	// 	// fmt.Println(suite.Name)

	// 	for _, test := range suite.Tests {
	// 		if test.Name == tc.Name {
	// 			fmt.Printf("  %s\n", test.Name)
	// 			mr.TestCaseTime = test.Duration
	// 			err := models.SaveRun(db, mr)
	// 			if err != nil {
	// 				fmt.Println("ERROR saving run: ", err.Error())
	// 			}

	// 			// if test.Error != nil {
	// 			// 	fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
	// 			// } else {
	// 			// 	fmt.Printf("    %s\n", test.Status)
	// 			// }
	// 			// mr.TestCaseTime = test.Duration

	// 		}
	// 	}
	// }

}
