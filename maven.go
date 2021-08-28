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

	"github.com/jinzhu/gorm"
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
		fmt.Printf("cmd.Run() failed with %s\n", err)
		log.Printf("Compilation out:\n%s\n", string(output))
		fmt.Printf("Compilation out:\n%s\n", string(output))
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

func MvnTest(db *gorm.DB, path string, measurementID uint) ([]MvnTestResult, bool) {
	ok := true
	logfile := "maven-test.log"

	log.Println("------------------------------------------------ mvn test")
	fmt.Println("------------------------------------------------ mvn test")
	cmd := exec.Command("mvn", "-Drat.skip=true", "-Djacoco.destFile=jacoco.exec", "clean", "org.jacoco:jacoco-maven-plugin:0.7.8:prepare-agent", "test")
	cmd.Dir = path

	var output []byte
	var err error

	// output, err = cmd.CombinedOutput()
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	pid := cmd.Process.Pid

	stop := make(chan bool)
	go func() {
		perfMetrics := []PerfMetrics{}
		for {
			select {
			case <-stop:
				//save
				for _, perfMetric := range perfMetrics {
					mr := &models.MeasurementResources{
						MeasurementID: measurementID,
						Type:          "maven",
						Resources: models.Resources{
							CpuPercent:        perfMetric.CpuPercent,
							MemPercent:        perfMetric.MemoryPercent,
							MemoryInfoStat:    *perfMetric.MemoryInfo,
							IOCountersStat:    *perfMetric.IOCounters,
							PageFaultsStat:    *perfMetric.PageFaults,
							AvgStat:           *perfMetric.Load,
							VirtualMemoryStat: *perfMetric.VirtualMemory,
							SwapMemory:        *perfMetric.SwapMemory,
							// CPUTime:           perfMetric.CPUTime,
							// DiskIOCounters:    perfMetric.DiskIOCounters,
							// NetIOCounters:     perfMetric.NetIOCounters,
						},
					}
					models.CreateMeasurementResources(db, mr)
					for _, cpuTime := range perfMetric.CPUTimes {
						models.CreateCPUTimes(db, &models.CPUTimes{
							MeasurementResourcesID: mr.ID,
							CPU:                    cpuTime.CPU,
							User:                   cpuTime.User,
							System:                 cpuTime.System,
							Idle:                   cpuTime.Idle,
							Nice:                   cpuTime.Nice,
							Iowait:                 cpuTime.Iowait,
							Irq:                    cpuTime.Irq,
							Softirq:                cpuTime.Softirq,
							Steal:                  cpuTime.Steal,
							Guest:                  cpuTime.Guest,
							GuestNice:              cpuTime.GuestNice,
						})
					}

					for i, diskIOCounter := range perfMetric.DiskIOCounters {
						models.CreateDiskIOCounters(db, &models.DiskIOCounters{
							MeasurementResourcesID: mr.ID,
							Device:                 i,
							ReadCount:              diskIOCounter.ReadCount,
							MergedReadCount:        diskIOCounter.MergedReadCount,
							WriteCount:             diskIOCounter.WriteCount,
							MergedWriteCount:       diskIOCounter.MergedWriteCount,
							ReadBytes:              diskIOCounter.ReadBytes,
							WriteBytes:             diskIOCounter.WriteBytes,
							ReadTime:               diskIOCounter.ReadTime,
							WriteTime:              diskIOCounter.WriteTime,
							IopsInProgress:         diskIOCounter.IopsInProgress,
							IoTime:                 diskIOCounter.IoTime,
							WeightedIO:             diskIOCounter.WeightedIO,
							Name:                   diskIOCounter.Name,
							SerialNumber:           diskIOCounter.SerialNumber,
							Label:                  diskIOCounter.Label,
						})
					}

					for i, netIOCounter := range perfMetric.NetIOCounters {
						models.CreateNetIOCounters(db, &models.NetIOCounters{
							MeasurementResourcesID: mr.ID,
							NICID:                  uint(i),
							Name:                   netIOCounter.Name,
							BytesSent:              netIOCounter.BytesSent,
							BytesRecv:              netIOCounter.BytesRecv,
							PacketsSent:            netIOCounter.PacketsSent,
							PacketsRecv:            netIOCounter.PacketsRecv,
							Errin:                  netIOCounter.Errin,
							Errout:                 netIOCounter.Errout,
							Dropin:                 netIOCounter.Dropin,
							Dropout:                netIOCounter.Dropout,
							Fifoin:                 netIOCounter.Fifoin,
							Fifoout:                netIOCounter.Fifoout,
						})
					}
				}
				return
			default:
				perfMetric, err := MonitorProcess(pid)
				if err == nil {
					perfMetrics = append(perfMetrics, perfMetric)
				}

			}
		}
	}()

	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
	stop <- true

	if err != nil {
		fmt.Printf("mvn -Drat.skip=true test failed with %s\n", err)
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
		fmt.Printf("Compilation out:\n%s\n", string(output))
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
