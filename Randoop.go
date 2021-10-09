package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go-repo-downloader/models"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// func CollectRandoopMetricsByChange(repoDir, repoName, prevCommit, fileFrom, currCommit, fileTo string, changeID uint) {
// 	//checkout previous commit
// 	// var okB bool
// 	// var okA bool
// 	// var metricsB [5]string
// 	// var metricsA [5]string

// 	fmt.Println("prevCommit: " + prevCommit + "fileFrom: " + fileFrom)
// 	err := Checkout(repoName, prevCommit)
// 	if err == nil {
// 		okB, metricsB = runRandoop(repoDir, fileFrom)
// 	}

// 	//checkout current commit
// 	fmt.Println("current Commit: " + currCommit + "fileTo: " + fileTo)
// 	err = Checkout(repoName, prevCommit)
// 	if err == nil {
// 		okA, metricsA = runRandoop(repoDir, fileTo)
// 	}

// 	// if okB == true && okA == true {
// 	// 	db := models.GetDB()
// 	// 	rm := &models.RandoopMetrics{ChangeID: changeID,
// 	// 		NMEBefore:  metricsB[0],
// 	// 		EMEBefore:  metricsB[1],
// 	// 		AETNBefore: metricsB[2],
// 	// 		AETEBefore: metricsB[3],
// 	// 		AMUBefore:  metricsB[4],
// 	// 		NMEAfter:   metricsA[0],
// 	// 		EMEAfter:   metricsA[1],
// 	// 		AETNAfter:  metricsA[2],
// 	// 		AETEAfter:  metricsA[3],
// 	// 		AMUAfter:   metricsA[4],
// 	// 		NMEDiff:    0,
// 	// 		EMEDiff:    0,
// 	// 		AETNDiff:   0,
// 	// 		AETEDiff:   0,
// 	// 		AMUDiff:    0,
// 	// 		NMEPerc:    0,
// 	// 		EMEPerc:    0,
// 	// 		AETNPerc:   0,
// 	// 		AETEPerc:   0,
// 	// 		AMUPerc:    0}
// 	// 	models.CreateRandoopMetrics(db, rm)
// 	// }
// }

// func CollectRandoopMetricsByAllClasses(repoDir, repoName, prevCommit, fileFrom, currCommit, fileTo string, measurement models.Measurement) {
// 	//checkout previous commit
// 	var okB bool
// 	var okA bool
// 	var metricsB [5]string
// 	var metricsA [5]string

// 	fmt.Println("prevCommit: " + prevCommit + "fileFrom: " + fileFrom)
// 	err := Checkout(repoName, prevCommit)
// 	if err == nil {
// 		okB, metricsB = runRandoop(repoDir, fileFrom)
// 	}

// 	//checkout current commit
// 	fmt.Println("current Commit: " + currCommit + "fileTo: " + fileTo)
// 	err = Checkout(repoName, prevCommit)
// 	if err == nil {
// 		okA, metricsA = runRandoop(repoDir, fileTo)
// 	}

// 	if okB == true && okA == true {
// 		db := models.GetDB()
// 		rm := &models.RandoopMetrics{ChangeID: changeID,
// 			NMEBefore:  metricsB[0],
// 			EMEBefore:  metricsB[1],
// 			AETNBefore: metricsB[2],
// 			AETEBefore: metricsB[3],
// 			AMUBefore:  metricsB[4],
// 			NMEAfter:   metricsA[0],
// 			EMEAfter:   metricsA[1],
// 			AETNAfter:  metricsA[2],
// 			AETEAfter:  metricsA[3],
// 			AMUAfter:   metricsA[4],
// 			NMEDiff:    0,
// 			EMEDiff:    0,
// 			AETNDiff:   0,
// 			AETEDiff:   0,
// 			AMUDiff:    0,
// 			NMEPerc:    0,
// 			EMEPerc:    0,
// 			AETNPerc:   0,
// 			AETEPerc:   0,
// 			go biAMUPerc:    0}
// 		models.CreateRandoopMetrics(db, rm)
// 	}
// }

func generateRandoopTests(db *gorm.DB, dir, classpath, cpSep, randoopJar, envRandoopJar, className string, measurement models.Measurement, commitID uint) bool {
	log.Println("------------------------------------------------ generating Randoop tests for " + className + "...")

	// fmt.Println("dir: ", dir)
	// dir:  /home/farah/go-work/src/github.com/paulorfarah/repos/gradle-project-example/app/

	randoopStr := "java -classpath " + classpath + cpSep + randoopJar + cpSep + envRandoopJar + " randoop.main.Main gentests --testclass=" + className + " --time-limit=5 > " + getDirectory() + "/gentest/" + className + ".txt"
	log.Println(randoopStr)
	fmt.Println(randoopStr)
	cmdRandoop := exec.Command("bash", "-c", randoopStr)
	cmdRandoop.Dir = dir

	output, err := cmdRandoop.CombinedOutput()
	log.Printf("test generation out:\n%s\n", string(output))
	// fmt.Printf("test generation out:\n%s\n", string(output))
	if err != nil {
		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot generate randoop tests (" + fmt.Sprint(err.Error()) + "): ")
		log.Println(string(output))

		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot generate randoop tests (" + fmt.Sprint(err.Error()) + "): ")
		fmt.Println(string(output))
		return false
	}
	saveRandoopTestCases(db, dir, className, measurement, commitID)
	return readRandoopGentestResults("gentest/" + className + ".txt")
}

func compileRandoopTests(source, output, classpath, cpSep string) bool {
	log.Println("------------------------------------------------ compile randoop tests")
	fmt.Println("------------------------------------------------ compile randoop tests")

	junitJar := "$JUNITPATH"
	dt := time.Now().Format("01-02-2006_150405")
	randoopStr := "javac -cp " + classpath + cpSep + os.ExpandEnv(junitJar) + " -d " + output + " RegressionTest*.java > " + getDirectory() + "/compilation/RandoopRegressionTest_compilation_" + dt + ".txt" //RegressionTest_compilation.txt"
	log.Println(randoopStr)
	fmt.Println(randoopStr)
	cmdRandoop := exec.Command("bash", "-c", randoopStr)
	cmdRandoop.Dir = source
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdRandoop.Stdout = &out
	cmdRandoop.Stderr = &stderr
	err := cmdRandoop.Run()
	if err != nil {
		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot compile randoop tests (" + fmt.Sprint(err) + "): " + stderr.String())
		log.Println(out)

		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot compile randoop tests (" + fmt.Sprint(err) + "): " + stderr.String())
		fmt.Println(out)
		return false
	}
	// }
	return true
}

func runRandoopTests(dir, classpath, cpSep string) (float64, int, []PerfMetrics, bool) {
	log.Println("------------------------------------------------ run randoop tests")
	fmt.Println("------------------------------------------------ run randoop tests")

	junitJar := "${JUNITPATH}"
	dt := time.Now().Format("01-02-2006_150405")
	logfile := getDirectory() + "/run/RandoopRegressionTest_run_" + dt + ".txt"
	junitStr := "java -javaagent:" + getDirectory() + "/jacoco-0.8.6/jacocoagent.jar -cp ." + cpSep + classpath + cpSep + junitJar + " org.junit.runner.JUnitCore RegressionTest > " + logfile

	// java -jar jacoco-0.8.6/lib/jacococli.jar report jacoco.exec --classfiles classes --sourcefiles src --csv <file>

	log.Println(junitStr)
	fmt.Println(junitStr)
	cmdRandoop := exec.Command("bash", "-c", junitStr)
	cmdRandoop.Dir = dir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdRandoop.Stdout = &out
	cmdRandoop.Stderr = &stderr
	// err := cmdRandoop.Run()
	err := cmdRandoop.Start()
	if err != nil {
		log.Fatal(err)
	}
	pid := cmdRandoop.Process.Pid

	stop := make(chan bool)
	perfMetrics := []PerfMetrics{}
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				perfMetric, err := MonitorProcess(pid)
				if err == nil {
					perfMetrics = append(perfMetrics, perfMetric)
				}

			}
		}
	}()

	err = cmdRandoop.Wait()
	log.Printf("Command finished with error: %v", err)
	stop <- true

	if err != nil {
		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot execute randoop tests (" + err.Error() + "): " + stderr.String())
		log.Println(out)

		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot execute randoop tests (" + err.Error() + "): " + stderr.String())
		fmt.Println(out)
		return float64(0.0), 0, []PerfMetrics{}, false
	}

	testTime, numTests, ok := readRandoopTestResults(logfile)
	return testTime, numTests, perfMetrics, ok
}

func readPackage(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pack := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "package") {
			pack = strings.Replace(line, "package ", "", 1)
			pack = strings.Replace(pack, ";", "", 1)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("pack: ", pack)
	return pack
}

func parseProjectPath(file string) (string, string) {
	dir := ""
	pack := ""
	paths := strings.Split(file, string(os.PathSeparator)+"src"+string(os.PathSeparator)+"main"+string(os.PathSeparator)+"java"+string(os.PathSeparator))
	if len(paths) > 1 {
		dir = paths[0]
		pack = paths[1]
	} else {
		paths = strings.Split(file, string(os.PathSeparator)+"src"+string(os.PathSeparator)+"test"+string(os.PathSeparator)+"java"+string(os.PathSeparator))
		if len(paths) > 1 {
			dir = paths[0]
			pack = paths[1]
		} else if len(paths) == 1 {
			if strings.Contains(file, "/java/") {
				//commons-io
				pack = strings.TrimLeft(file, "/src/main/java/")
			} else if strings.Contains(file, "src/conf/") {
				pack = strings.TrimLeft(file, "/src/conf/")
			} else if strings.Contains(file, "src/examples/") {
				pack = strings.TrimLeft(file, "/src/examples/")
			} else if strings.Contains(file, "src/java/") {
				pack = strings.TrimLeft(file, "/src/java/")
			} else if strings.Contains(file, "/src/test/java/") {
				//junit4
				parts := strings.Split(file, "/src/test/java/")
				dir = parts[0] //+ "/src/test/java"
				pack = parts[1]

			} else if strings.Contains(file, "src/test/") {
				pack = strings.TrimLeft(file, "/src/test/")
			} else if strings.Contains(file, "core/src/test/") {
				pack = strings.TrimLeft(file, "/core/src/test/")
			} else {
				fmt.Println("Error in parseProjectPath, path not in list -  filefrom: " + file)
				// paths = strings.Split(file, "src/")
				// dir = paths[0]
				// pack = paths[1]
				pack = readPackage(file)
				packTmp := strings.ReplaceAll(pack, ".", "/")
				dir = strings.Split(file, packTmp)[0]
				fmt.Println("###################### parse project path: ")
				fmt.Println("file: ", file)
				fmt.Println("pack: ", pack)
				fmt.Println("packTmp: ", packTmp)
				fmt.Println("dir: ", dir)
				fmt.Println("######################")
			}
		}
	}
	return dir, pack
}
func readRandoopGentestResults(path string) bool {
	log.Println("readRandoopGentestResults: " + path)
	ok := false
	f, err := os.Open(path)
	if err != nil {
		log.Println("[>>ERROR]: There has been an error openning randoop-gentest log file: ", err.Error())
		log.Println("log file: " + path)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		row := scanner.Bytes()
		if len(string(row)) > 12 {
			if bytes.Equal(row[:12], []byte("Created file")) {
				ok = true
			} else if bytes.Contains(row, []byte("No regression tests to output.")) {
				ok = false

			}
		}
	}
	return ok
}

func readRandoopTestResults(path string) (float64, int, bool) {
	fmt.Println("readRandoopTestResults: " + path)
	ok := false
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("[>>ERROR]: There has been an error openning randoop-gentest log file: ", err.Error())
		fmt.Println("log file: " + path)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	testTime := float64(0.0)
	numTests := 0
	for scanner.Scan() {
		row := scanner.Bytes()
		if len(string(row)) > 5 {
			if bytes.Equal(row[:5], []byte("Time:")) {
				aux := strings.Split(string(row), " ")
				testTime, _ = strconv.ParseFloat(aux[1], 64)
			} else if bytes.Equal(row[:4], []byte("OK (")) {
				ok = true
				aux := strings.Split(string(row), " ")
				numTests, _ = strconv.Atoi(aux[1][1:])
			}
		}
	}
	return testTime, numTests, ok
}

func parseResult(line []byte, metric string) string {
	size := len(metric)
	res := ""
	if len(line) > len(metric) {
		if bytes.Equal(line[:size], []byte(metric)) {
			res = strings.Trim(string(line[size:]), " ")
		}
	}
	return res
}

func deleteOldRandoopTests() bool {
	dirname, err := os.Getwd()
	if err != nil {
		fmt.Println(">>>> ERROR: Cannot get local directory: " + err.Error())
	}
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println("Error openning dir to delete java files: " + err.Error())
		return false
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading java files: " + err.Error())
		return false
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".java" || filepath.Ext(file.Name()) == ".class" {
				fmt.Println("deleting file: " + file.Name())
				os.Remove(file.Name())
			}
		}
	}
	return true
}

func getDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}

func saveRandoopTestCases(db *gorm.DB, dir, className string, measurement models.Measurement, commitID uint) {
	//read generated test cases
	randoopTests, err := walkMatch(dir, "RegressionTest*")
	if err != nil {
		fmt.Println("Error reading randoop regresion tests: ", err.Error())
	}
	for _, test := range randoopTests {
		//rename
		testSuite, err := models.FindFileByEndsWithNameAndCommit(db, test, commitID)
		if err != nil {
			fmt.Println("Error trying to find file by name : ", err.Error())
		}
		//save
		mr := &models.TestCase{
			MeasurementID: measurement.ID,
			Type:          "randoop",
			ClassName:     className,
			CommitID:      commitID,
			// Duration:      test.Duration,
			TestSuiteID: testSuite.ID,
			Name:        test,
			// Status:        string(test.Status),
			// Error:         errorMsg,
			// Message:       test.Message,
			// SystemErr:     string(test.SystemErr),
			// SystemOut:     string(test.SystemOut),
		}
		_, errTC := models.CreateTestCase(db, mr)
		if errTC != nil {
			fmt.Println(errTC.Error())
		}
	}
}

func walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match("RegressionTest*", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
