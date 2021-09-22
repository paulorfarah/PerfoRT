package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
// 			AMUPerc:    0}
// 		models.CreateRandoopMetrics(db, rm)
// 	}
// }

func generateRandoopTests(classpath, cpSep, randoopJar, envRandoopJar, className string) bool {
	log.Println("------------------------------------------------ Generating Randoop tests for " + className + "...")
	// fmt.Println("------------------------------------------------ Generating Randoop tests for " + className + "...")
	// dir, pack := parseProjectPath(file)
	// if dir != "" {
	// 	dir += string(os.PathSeparator)
	// }

	// fmt.Println("=============================")
	// fmt.Println("dir: ", dir)
	// fmt.Println("pack: ", pack)

	// path := strings.Split(pack, ".java")[0]
	// // fmt.Println("path: ", path)
	// randoopJar := "${RANDOOP_JAR}"
	// cpSep := ":"
	// if runtime.GOOS == "windows" {
	// 	randoopJar = "%RANDOOP_JAR%"
	// 	cpSep = ";"
	// }

	// envRandoopJar := os.Getenv("RANDOOP_JAR")
	// // remove old tests
	// // deleteOldRandoopTests()

	// classpath := ""
	// switch buildTool {
	// case "maven":
	// 	classpath += dir + "target" + string(os.PathSeparator) + "classes" + cpSep
	// case "gradle":
	// 	classpath += dir + "build" + string(os.PathSeparator) + "classes" + cpSep + dir + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "main"
	// }
	// classpath += buildToolClasspath
	// className := strings.ReplaceAll(path, string(os.PathSeparator), ".")

	randoopStr := "java -classpath " + classpath + cpSep + randoopJar + cpSep + envRandoopJar + " randoop.main.Main gentests --testclass=" + className + " --time-limit=5 > gentest/" + className + ".txt"
	log.Println(randoopStr)
	fmt.Println(randoopStr)
	cmdRandoop := exec.Command("bash", "-c", randoopStr)

	output, err := cmdRandoop.CombinedOutput()
	if err != nil {
		log.Printf("GenerateRandoopTests failed with %s\n", err.Error())
		fmt.Printf("GenerateRandoopTests failed with %s\n", err.Error())
	}
	log.Printf("test generation out:\n%s\n", string(output))
	// fmt.Printf("test generation out:\n%s\n", string(output))
	if err != nil {
		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot run randoop gentests (" + fmt.Sprint(err.Error()) + "): ")
		log.Println(string(output))

		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot run randoop gentests (" + fmt.Sprint(err.Error()) + "): ")
		fmt.Println(string(output))
		return false
	}
	return readRandoopGentestResults("gentest/" + className + ".txt")
}

func compileRandoopTests(classpath, cpSep string) bool {
	log.Println("------------------------------------------------ compile randoop tests")
	fmt.Println("------------------------------------------------ compile randoop tests")
	// javac -classpath .:$JUNITPATH ErrorTest*.java RegressionTest*.java -sourcepath .:path/to/files/under/test/
	// javac -cp /mnt/sda4/go-work/src/github.com/paulorfarah/junit4/target/classes:/mnt/sda4/downloads/junit-4.13.2.jar:/mnt/sda4/downloads/hamcrest-core-1.3.jar:. RegressionTest2.java

	// gerado:
	// javac
	//-cp /mnt/sda4/go-work/src/github.com/paulorfarah/repos/junit4/ (repoDir) ok
	//    /mnt/sda4/go-work/ (dir)
	//    /target/classes
	//   :/users/farah/.m2/repository/org/hamcrest/hamcrest-core/1.3/hamcrest-core-1.3.jar:/users/farah/.m2/repository/org/hamcrest/hamcrest-library/1.3/hamcrest-library-1.3.jar:$JUNITPATHgithub.com.paulorfarah.go-repo-downloader.RegressionTest2

	// correto:
	//javac -cp /mnt/sda4/go-work/src/github.com/paulorfarah/repos/junit4/:/mnt/sda4/go-work/src/github.com/paulorfarah/repos/junit4/target/classes:/users/farah/.m2/repository/org/hamcrest/hamcrest-core/1.3/hamcrest-core-1.3.jar:/users/farah/.m2/repository/org/hamcrest/hamcrest-library/1.3/hamcrest-library-1.3.jar:$JUNITPATH RegressionTest.java -sourcepath /mnt/sda4/go-work/src/github.com/paulorfarah/repos/junit4/src/main/java

	// cpSep := ":"
	// if runtime.GOOS == "windows" {
	// 	junitJar = "%JUNITPATH%"
	// 	cpSep = ";"
	// }

	// classpath := ""
	// switch buildTool {
	// case "maven":
	// 	classpath += repoDir + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes" + cpSep
	// case "gradle":
	// 	classpath += repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" + cpSep + repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "main"
	// }
	// classpath += buildToolClasspath

	// className := strings.ReplaceAll(path, "/", ".")

	// clean temporary files to avoid Too many links error
	// cmdClean := exec.Command("bash", "-c", "find", "/tmp/", "-name", "\"*\"", "-print0|", "xargs", "-0", "rm", "-rf")
	// cmdClean.Run()

	junitJar := "$JUNITPATH"

	randoopStr := "javac -cp " + classpath + cpSep + os.ExpandEnv(junitJar) + " RegressionTest*.java > RegressionTest_compilation.txt"
	log.Println(randoopStr)
	fmt.Println(randoopStr)
	cmdRandoop := exec.Command("bash", "-c", randoopStr)
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

func runRandoopTests(classpath, cpSep string) (float64, int, []PerfMetrics, bool) {
	log.Println("------------------------------------------------ run randoop tests")
	fmt.Println("------------------------------------------------ run randoop tests")
	// java -classpath .:$JUNITPATH:myclasspath org.junit.runner.JUnitCore RegressionTest
	// java -cp .:/usr/share/java/junit.jar org.junit.runner.JUnitCore [test class name]
	// java -javaagent:jacoco-0.8.6/lib/jacocoagent.jar -cp junit-4.12.jar:hamcrest-core-1.3.jar:classes:test-classes org.junit.runner.JUnitCore CalculatorTest

	// cpSep := ":"
	// if runtime.GOOS == "windows" {
	// 	junitJar = "%JUNITPATH%"
	// 	cpSep = ";"
	// }

	// classpath := ""
	// // classpath := repoDir + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes"
	// switch buildTool {
	// case "maven":
	// 	classpath += repoDir + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes" + cpSep
	// case "gradle":
	// 	classpath += repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" + cpSep + repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "main"
	// }

	junitJar := "${JUNITPATH}"
	junitStr := "java -javaagent:jacoco-0.8.6/jacocoagent.jar -cp ." + cpSep + classpath + cpSep + junitJar + " org.junit.runner.JUnitCore RegressionTest > runRT.txt"

	// java -jar jacoco-0.8.6/lib/jacococli.jar report jacoco.exec --classfiles classes --sourcefiles src --csv <file>

	log.Println(junitStr)
	fmt.Println(junitStr)
	// fmt.Println(junitStr)
	cmdRandoop := exec.Command("bash", "-c", junitStr)
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

	testTime, numTests, ok := readRandoopTestResults("runRT.txt")
	return testTime, numTests, perfMetrics, ok
}

// func coverageRandoopTests(repoDir, buildTool string) {
// 	log.Println("------------------------------------------------ coverage randoop tests")
// 	fmt.Println("------------------------------------------------ coverage randoop tests")

// 	cpSep := ":"
// 	if runtime.GOOS == "windows" {
// 		cpSep = ";"
// 	}

// 	// classpath := repoDir + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes"
// 	classpath := ""
// 	switch buildTool {
// 	case "maven":
// 		classpath += repoDir + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes"
// 	case "gradle":
// 		classpath += repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" + cpSep + repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "main"
// 	}

// 	jacocoStr := "java -jar jacoco-0.8.6/jacococli.jar report jacoco.exec --classfiles " + classpath + " --sourcefiles " + repoDir + " --csv coverage/" + strings.ReplaceAll(repoDir, "/", "_") + ".csv"

// 	log.Println(jacocoStr)
// 	fmt.Println(jacocoStr)
// 	cmdRandoop := exec.Command("bash", "-c", jacocoStr)
// 	var out bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmdRandoop.Stdout = &out
// 	cmdRandoop.Stderr = &stderr
// 	err := cmdRandoop.Run()
// 	if err != nil {
// 		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot execute randoop tests (" + err.Error() + "): " + stderr.String())
// 		log.Println(out)

// 		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot execute randoop tests (" + err.Error() + "): " + stderr.String())
// 		fmt.Println(out)
// 	}

// 	// return readRandoopTestResults("runRT.txt")
// }

// Deprecated: changed to coverage.go
// func coverageRandoopTestsByFile(repoDir, src string) { //(map[string]int, bool) {
// 	log.Println("------------------------------------------------ coverage randoop tests")
// 	fmt.Println("------------------------------------------------ coverage randoop tests")
// 	// java -classpath .:$JUNITPATH:myclasspath org.junit.runner.JUnitCore RegressionTest
// 	// java -cp .:/usr/share/java/junit.jar org.junit.runner.JUnitCore [test class name]
// 	// java -javaagent:jacoco-0.8.6/lib/jacocoagent.jar -cp junit-4.12.jar:hamcrest-core-1.3.jar:classes:test-classes org.junit.runner.JUnitCore CalculatorTest

// 	// junitJar := "$JUNITPATH"
// 	// cpSep := ":"
// 	// if runtime.GOOS == "windows" {
// 	// 	junitJar = "%JUNITPATH%"
// 	// 	cpSep = ";"
// 	// }

// 	classpath := repoDir + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes"

// 	// junitStr := "java -javaagent:jacoco-0.8.6/lib/jacocoagent.jar -cp ." + cpSep + classpath + cpSep + junitJar + " org.junit.runner.JUnitCore RegressionTest > runRT.txt"
// 	jacocoStr := "java -jar jacoco-0.8.6/jacococli.jar report jacoco.exec --classfiles " + classpath + " --sourcefiles " + src + " --csv coverage/" + strings.ReplaceAll(src, "/", "_") + ".csv"

// 	// java -jar jacoco-0.8.6/lib/jacococli.jar report jacoco.exec --classfiles classes --sourcefiles src --csv <file>

// 	log.Println(jacocoStr)
// 	fmt.Println(jacocoStr)
// 	cmdRandoop := exec.Command("bash", "-c", jacocoStr)
// 	var out bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmdRandoop.Stdout = &out
// 	cmdRandoop.Stderr = &stderr
// 	err := cmdRandoop.Run()
// 	if err != nil {
// 		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot execute randoop tests (" + err.Error() + "): " + stderr.String())
// 		log.Println(out)

// 		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot execute randoop tests (" + err.Error() + "): " + stderr.String())
// 		fmt.Println(out)
// 	}

// 	// return readRandoopTestResults("runRT.txt")
// }

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
			if strings.Contains(file, "src/main/java/") {
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
	fmt.Println("readRandooTestResults: " + path)
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
