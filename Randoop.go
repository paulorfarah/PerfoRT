package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go-repo-downloader/models"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CollectRandoopMetrics(repoDir, repoName, prevCommit, fileFrom, currCommit, fileTo string, changeID uint) {
	//checkout previous commit
	var okB bool
	var okA bool
	var metricsB [5]string
	var metricsA [5]string

	fmt.Println("prevCommit: " + prevCommit + "fileFrom: " + fileFrom)
	err := checkout(repoName, prevCommit)
	if err == nil {
		okB, metricsB = runRandoop(repoDir, fileFrom)
	}

	//checkout current commit
	fmt.Println("current Commit: " + currCommit + "fileTo: " + fileTo)
	err = checkout(repoName, prevCommit)
	if err == nil {
		okA, metricsA = runRandoop(repoDir, fileTo)
	}

	if okB == true && okA == true {
		db := models.GetDB()
		rm := &models.RandoopMetrics{ChangeID: changeID,
			NMEBefore:  metricsB[0],
			EMEBefore:  metricsB[1],
			AETNBefore: metricsB[2],
			AETEBefore: metricsB[3],
			AMUBefore:  metricsB[4],
			NMEAfter:   metricsA[0],
			EMEAfter:   metricsA[1],
			AETNAfter:  metricsA[2],
			AETEAfter:  metricsA[3],
			AMUAfter:   metricsA[4],
			NMEDiff:    0,
			EMEDiff:    0,
			AETNDiff:   0,
			AETEDiff:   0,
			AMUDiff:    0,
			NMEPerc:    0,
			EMEPerc:    0,
			AETNPerc:   0,
			AETEPerc:   0,
			AMUPerc:    0}
		models.CreateRandoopMetrics(db, rm)
		// if s, err := strconv.ParseFloat(metricsB[0], 64); err == nil {
		// 	fmt.Printf("%T, %v\n", s, s)
		// }

	}
}

func runRandoop(repoDir, file string) (bool, [5]string) {
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
	}
	classpath := repoDir + string(os.PathSeparator) + dir + "target" + string(os.PathSeparator) + "classes" // + string(os.PathSeparator)
	classpath += GetMavenDependenciesClasspath(repoDir)
	className := strings.ReplaceAll(path, "/", ".")

	// fmt.Printf("java -classpath " + classpath + cpSep + randoopJar + " randoop.main.Main gentests --testclass=" + className + "\n")
	c := "java -classpath " + classpath + cpSep + randoopJar + " randoop.main.Main gentests --testclass=" + className + " > " + className + ".txt"
	fmt.Println(c)
	cmd := exec.Command("bash", "-c", c)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>CRITICAL ERROR]: Cannot run randoop gentests (" + fmt.Sprint(err) + "): " + stderr.String())
		fmt.Println(out)
		return false, [5]string{}
	}
	return true, readRandoopResults(className + ".txt")
}

func parseProjectPath(file string) (string, string) {
	dir := ""
	pack := ""
	paths := strings.Split(file, "/src/main/java/")
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

func readRandoopResults(fn string) [5]string {
	var res, nme, eme, aetn, aete, amu string
	f, err := os.Open(fn)
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		res = parseResult(line, "Normal method executions:")
		if res != "" {
			nme = res
		}
		res = parseResult(line, "Exceptional method executions:")
		if res != "" {
			eme = res
		}
		res = parseResult(line, "Average method execution time (normal termination):")
		if res != "" {
			aetn = res
		}
		res = parseResult(line, "Average method execution time (exceptional termination):")
		if res != "" {
			aete = res
		}
		res = parseResult(line, "Approximate memory usage")
		if res != "" {
			amu = res
		}
	}
	return [5]string{nme, eme, aetn, aete, amu}
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
