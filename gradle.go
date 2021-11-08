package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go-repo-downloader/models"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/joshdk/go-junit"
	"gorm.io/gorm"
)

type GradleTestResult struct {
	ClassName   string
	TestsRun    int
	Failures    int
	Errors      int
	Skipped     int
	TimeElapsed float64
}

func GetGradleDependenciesClasspath(path string) string {
	logfile := "gradle-classpath.log"

	fmt.Println("gradle dependencies")
	cmd := exec.Command("gradle", "dependencies")
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

	return getGradleClasspath(path)
}

func getGradleClasspath(path string) string {
	found := false
	classpath := ""
	logfile := "gradle-classpath.log"
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

func GradleBuild(path string) bool {
	logfile := "gradle-compiler.log"

	fmt.Println("------------------------------------------------ gradle build")
	cmd := exec.Command("gradle", "build")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("gradle build failed with %s\n", err)
		fmt.Printf("gradle build failed with %s\n", err)
		log.Printf("gradle build Compilation out:\n%s\n", string(output))
		fmt.Printf("gradle build Compilation out:\n%s\n", string(output))
		return false
	}
	log.Printf("gradle build Compilation out:\n%s\n", string(output))
	fmt.Printf("gradle build Compilation out:\n%s\n", string(output))
	// fmt.Printf("Compilation out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		log.Printf(err.Error())
		fmt.Printf(err.Error())
		return false
	}
	if strings.Contains(string(output), "BUILD SUCCESSFUL") {
		fmt.Println("BUILD SUCCESSFUL")
		return true
	} else {
		return false
	}
}

func GradleTest(db *gorm.DB, path string, measurementID uint) bool {
	ok := true
	logfile := "gradle-test.log"

	log.Println("------------------------------------------------ gradle test")
	fmt.Println("------------------------------------------------ gradle test")
	cmd := exec.Command("gradle", "test")
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
	// 					Type:          "gradle",
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
	log.Printf("Command finished with error: %v", err)
	// stop <- true

	if err != nil {
		fmt.Printf("gradle test failed with %s\n", err)
	}

	// fmt.Printf("Mvn test out:\n%s\n", string(output))
	log.Printf("gradle test out:\n%s\n", string(output))
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

func readGradleTestResults(path string) []MvnTestResult {
	logfile := "gradle-test.log"
	f, err := os.Open(path + string(os.PathSeparator) + logfile)
	if err != nil {
		log.Println("[>>ERROR]: There has been an error running gradle test: ", err.Error())
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

func RunGradleTestCase(db *gorm.DB, path string, tc *models.TestCase, measurementID uint) {
	// # Executes a single specified test in SomeTestClass
	// gradle test --tests SomeTestClass.someSpecificMethod
	// fmt.Println("TC: ", tc.ID)

	// ok := true
	logfile := "gradle-test.log"
	testName := tc.ClassName + "." + tc.Name

	log.Println(">>>------------------------------------------------ gradle testcase", path, testName)
	fmt.Println(">>>------------------------------------------------ gradle testcase", path, testName)
	fmt.Printf("gradle test --rerun-tasks --tests %s (dir: %s)\n", testName, path)
	cmd := exec.Command("gradle", "test", "--rerun-tasks", "--tests", testName)
	cmd.Dir = path

	var output []byte
	var err error

	mr := &models.Run{
		MeasurementID: measurementID,
		Type:          "gradle",
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
		fmt.Printf("gradle test failed with %s\n", err.Error())
		log.Printf("Command finished with error: %s", err.Error())
	}

	// fmt.Printf("Mvn test out:\n%s\n", string(output))
	// log.Printf("gradle test out:\n%s\n", string(output))
	err = ioutil.WriteFile(path+string(os.PathSeparator)+logfile, []byte(output), 0644)
	if err != nil {
		// ok = false
		panic(err)
	}

	resultsPath := path + "/build/test-results/test/TEST-" + tc.ClassName + ".xml"
	// fmt.Println(resultsPath)
	suites, err := junit.IngestFile(resultsPath)
	if err != nil {
		log.Fatalf("failed to ingest JUnit xml %v", err)
	}
	for _, suite := range suites {
		// fmt.Println(suite.Name)

		for _, test := range suite.Tests {
			if test.Name == tc.Name {
				fmt.Printf("  %s\n", test.Name)
				mr.TestCaseTime = test.Duration
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
	}
}

func saveMetrics(db *gorm.DB, measurementID uint, perfMetric PerfMetrics) {
	resource := &models.Resource{
		RunID:             measurementID,
		CpuPercent:        perfMetric.CpuPercent,
		MemPercent:        perfMetric.MemoryPercent,
		MemoryInfoStat:    *perfMetric.MemoryInfo,
		IOCountersStat:    *perfMetric.IOCounters,
		PageFaultsStat:    *perfMetric.PageFaults,
		AvgStat:           *perfMetric.Load,
		VirtualMemoryStat: *perfMetric.VirtualMemory,
		SwapMemoryStat:    *perfMetric.SwapMemory,
		// DiskIOCounters:    perfMetric.DiskIOCounters,
		// NetIOCounters:     perfMetric.NetIOCounters,

	}
	_, err := models.CreateResource(db, resource)
	if err != nil {
		fmt.Printf("Error saving resource: %s\n", err.Error())
	}

	for _, cpuTime := range perfMetric.CPUTimes {
		models.CreateCPUTimes(db, &models.CPUTimes{
			ResourceID: resource.ID,
			CPU:        cpuTime.CPU,
			User:       cpuTime.User,
			System:     cpuTime.System,
			Idle:       cpuTime.Idle,
			Nice:       cpuTime.Nice,
			Iowait:     cpuTime.Iowait,
			Irq:        cpuTime.Irq,
			Softirq:    cpuTime.Softirq,
			Steal:      cpuTime.Steal,
			Guest:      cpuTime.Guest,
			GuestNice:  cpuTime.GuestNice,
		})
	}

	for i, diskIOCounter := range perfMetric.DiskIOCounters {
		models.CreateDiskIOCounters(db, &models.DiskIOCounters{
			ResourceID:       resource.ID,
			Device:           i,
			ReadCount:        diskIOCounter.ReadCount,
			MergedReadCount:  diskIOCounter.MergedReadCount,
			WriteCount:       diskIOCounter.WriteCount,
			MergedWriteCount: diskIOCounter.MergedWriteCount,
			ReadBytes:        diskIOCounter.ReadBytes,
			WriteBytes:       diskIOCounter.WriteBytes,
			ReadTime:         diskIOCounter.ReadTime,
			WriteTime:        diskIOCounter.WriteTime,
			IopsInProgress:   diskIOCounter.IopsInProgress,
			IoTime:           diskIOCounter.IoTime,
			WeightedIO:       diskIOCounter.WeightedIO,
			Name:             diskIOCounter.Name,
			SerialNumber:     diskIOCounter.SerialNumber,
			Label:            diskIOCounter.Label,
		})
	}

	for i, netIOCounter := range perfMetric.NetIOCounters {
		models.CreateNetIOCounters(db, &models.NetIOCounters{
			ResourceID:  resource.ID,
			NICID:       uint(i),
			Name:        netIOCounter.Name,
			BytesSent:   netIOCounter.BytesSent,
			BytesRecv:   netIOCounter.BytesRecv,
			PacketsSent: netIOCounter.PacketsSent,
			PacketsRecv: netIOCounter.PacketsRecv,
			Errin:       netIOCounter.Errin,
			Errout:      netIOCounter.Errout,
			Dropin:      netIOCounter.Dropin,
			Dropout:     netIOCounter.Dropout,
			Fifoin:      netIOCounter.Fifoin,
			Fifoout:     netIOCounter.Fifoout,
		})
	}
}
