package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

func generateRandoopTests(repoDir, file string) ([]string, bool) {
	fmt.Println("Generating Randoop tests for " + file + "...")
	dir, pack := parseProjectPath(file)
	if dir != "" {
		dir += string(os.PathSeparator)
	}
	path := strings.Split(pack, ".java")[0]
	randoopJar := "$RANDOOP_JAR"
	cpSep := ":"
	if runtime.GOOS == "windows" {
		randoopJar = "%RANDOOP_JAR%"
		cpSep = ";"
	} else {
		// clean temporary files to avoid Too many links error
		cmdClean := exec.Command("bash", "-c", "find", "/tmp/", "-name", "\"*\"", "-print0|", "xargs", "-0", "rm", "-rf")
		cmdClean.Run()
	}
	classpath := repoDir + string(os.PathSeparator) + dir + "target" + string(os.PathSeparator) + "classes" + cpSep
	classpath += GetMavenDependenciesClasspath(repoDir)
	className := strings.ReplaceAll(path, string(os.PathSeparator), ".")

	randoopStr := "java -classpath " + classpath + cpSep + randoopJar + " randoop.main.Main gentests --testclass=" + className + " > " + className + ".txt"
	fmt.Println(randoopStr)
	cmdRandoop := exec.Command("bash", "-c", randoopStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdRandoop.Stdout = &out
	cmdRandoop.Stderr = &stderr
	err := cmdRandoop.Run()
	if err != nil {
		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot run randoop gentests (" + fmt.Sprint(err) + "): " + stderr.String())
		fmt.Println(out)
		return []string{}, false
	}
	return readRandoopGentestResults(className + ".txt"), true
}

func compileRandoopTests(repoDir string, testfiles []string) bool {
	// javac -classpath .:$JUNITPATH ErrorTest*.java RegressionTest*.java -sourcepath .:path/to/files/under/test/
	//javac -cp /mnt/sda4/go-work/src/github.com/paulorfarah/junit4/target/classes:/mnt/sda4/downloads/junit-4.13.2.jar:/mnt/sda4/downloads/hamcrest-core-1.3.jar:. RegressionTest2.java

	for _, file := range testfiles {
		dir, pack := parseProjectPath(file)
		if dir != "" {
			dir += string(os.PathSeparator)
		}
		path := strings.Split(pack, ".java")[0]
		junitJar := "$JUNITPATH"
		cpSep := ":"
		if runtime.GOOS == "windows" {
			junitJar = "%JUNITPATH%"
			cpSep = ";"
		}
		// classpath := "/mnt/sda4/downloads/junit-4.13.2.jar:/mnt/sda4/downloads/hamcrest-core-1.3.jar" + cpSep
		classpath := repoDir + string(os.PathSeparator) + dir + "target" + string(os.PathSeparator) + "classes" + cpSep
		classpath += GetMavenDependenciesClasspath(repoDir)
		className := strings.ReplaceAll(path, "/", ".")

		// clean temporary files to avoid Too many links error
		cmdClean := exec.Command("bash", "-c", "find", "/tmp/", "-name", "\"*\"", "-print0|", "xargs", "-0", "rm", "-rf")
		cmdClean.Run()

		randoopStr := "javac -cp " + classpath + cpSep + junitJar + className + " > " + className + "_comp.txt"
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> javac")
		fmt.Println(randoop)
		cmdRandoop := exec.Command("bash", "-c", randoopStr)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmdRandoop.Stdout = &out
		cmdRandoop.Stderr = &stderr
		err := cmdRandoop.Run()
		if err != nil {
			fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot compile randoop tests (" + fmt.Sprint(err) + "): " + stderr.String())
			fmt.Println(out)
			// return false
		}
	}
	return true
}

func runRandoopTests(testfiles []string) bool {
	// java -classpath .:$JUNITPATH:myclasspath org.junit.runner.JUnitCore RegressionTest
	// java -cp .:/usr/share/java/junit.jar org.junit.runner.JUnitCore [test class name]
	return true
}

func parseProjectPath(file string) (string, string) {
	dir := ""
	pack := ""
	paths := strings.Split(file, string(os.PathSeparator)+"src"+string(os.PathSeparator)+"main"+string(os.PathSeparator)+"java"+string(os.PathSeparator))
	if len(paths) > 1 {
		dir = paths[0]
		pack = paths[1]
	} else if len(paths) == 1 {
		if strings.HasPrefix(file, "src/main/java/") {
			//commons-io
			pack = strings.TrimLeft(file, "/src/main/java/")
		} else if strings.HasPrefix(file, "src/conf/") {
			pack = strings.TrimLeft(file, "/src/conf/")
		} else if strings.HasPrefix(file, "src/examples/") {
			pack = strings.TrimLeft(file, "/src/examples/")
		} else if strings.HasPrefix(file, "src/java/") {
			pack = strings.TrimLeft(file, "/src/java/")
		} else if strings.HasPrefix(file, "src/test/") {
			pack = strings.TrimLeft(file, "/src/test/")
		} else {
			fmt.Println("**************************** filefrom: " + file)
			paths = strings.Split(file, "src/")
			dir = paths[0]
			pack = paths[1]
		}
	}
	return dir, pack
}
func readRandoopGentestResults(path string) []string {
	// logfile := "randoop-gentest.log"
	// f, err := os.Open(path + string(os.PathSeparator) + logfile)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("[>>ERROR]: There has been an error openning randoop-gentest log file: ", err.Error())
		fmt.Println("log file: " + path)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var files []string
	for scanner.Scan() {
		row := scanner.Bytes()
		fmt.Println(string(row))
		elements := strings.Split(string(row), " ")
		if len(elements) > 9 {
			if bytes.Equal(row[:12], []byte("Created file")) {
				aux := strings.Split(string(row), " ")

				f := aux[1]
				files = append(files, f)

			}
		}
	}
	return files
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
