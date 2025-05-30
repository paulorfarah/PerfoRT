package main

import (
	"PerfoRT/models"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func JacocoTestCoverage(db *gorm.DB, repoDir, javaVer, testtype, buildTool string, measurementID, commitID uint) error {
	filename := "coverage/" + strings.ReplaceAll(repoDir, "/", "_") + "-" + strconv.Itoa(int(commitID)) + ".csv"

	classpath := ""
	switch buildTool {
	case "maven":
		classpath += repoDir + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes"
	case "gradle":
		classpath += repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" //+ cpSep + repoDir + string(os.PathSeparator) + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "main"
	}

	jacoco_exec := repoDir + "/jacoco.exec"

	jacocoStr := javaVer + "/bin/java -jar jacoco-0.8.6/jacococli.jar report " + jacoco_exec + " --classfiles " + classpath + " --sourcefiles " + repoDir + " --csv " + filename

	cmd := exec.Command("bash", "-c", jacocoStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ERROR]: Cannot execute JaCoCo coverage (" + err.Error() + "): " + stderr.String())
		log.Print(jacocoStr)
		log.Println(out.String())

		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ERROR]: Cannot execute JaCoCo coverage (" + err.Error() + "): " + stderr.String())
		fmt.Print(jacocoStr)
		fmt.Println(out.String())
	}

	err = saveCoverage(db, filename, testtype, measurementID, commitID)
	if err != nil {
		log.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ERROR]: Cannot save JaCoCo coverage: " + err.Error())
		log.Println(out.String())

		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ERROR]: Cannot save JaCoCo coverage: " + err.Error())
		fmt.Println(out.String())
	}

	return err
}

func saveCoverage(db *gorm.DB, filename string, testType string, measurementID, commitID uint) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	header := true
	for _, line := range lines {
		if header {
			header = false
		} else {
			im, err := strconv.Atoi(line[3])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Instruction Missed of JaCoCo coverage report for " + testType + "! Value: " + line[3] + ", Error: " + err.Error())
				im = 0
			}
			ic, err := strconv.Atoi(line[4])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Instruction Covered of JaCoCo coverage report for " + testType + "! Value: " + line[4] + ", Error: " + err.Error())
				ic = 0
			}
			bm, err := strconv.Atoi(line[5])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Branch Missed of JaCoCo coverage report for " + testType + "! Value: " + line[5] + ", Error: " + err.Error())
				bm = 0
			}
			bc, err := strconv.Atoi(line[6])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Branch Covered of JaCoCo coverage report for " + testType + "! Value: " + line[6] + ", Error: " + err.Error())
				bc = 0
			}
			lm, err := strconv.Atoi(line[7])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Line Missed of JaCoCo coverage report for " + testType + "! Value: " + line[7] + ", Error: " + err.Error())
				lm = 0
			}
			lc, err := strconv.Atoi(line[8])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Line Covered of JaCoCo coverage report for " + testType + "! Value: " + line[8] + ", Error: " + err.Error())
				lc = 0
			}
			cm, err := strconv.Atoi(line[9])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Complexity Missed of JaCoCo coverage report for " + testType + "! Value: " + line[9] + ", Error: " + err.Error())
				cm = 0
			}
			cc, err := strconv.Atoi(line[10])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Complexity Covered of JaCoCo coverage report for " + testType + "! Value: " + line[10] + ", Error: " + err.Error())
				cc = 0
			}
			mm, err := strconv.Atoi(line[11])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Method Missed of JaCoCo coverage report for " + testType + "! Value: " + line[11] + ", Error: " + err.Error())
				mm = 0
			}
			mc, err := strconv.Atoi(line[12])
			if err != nil {
				log.Println("[>>ERROR]: Error reading Method Covered of JaCoCo coverage report for " + testType + "! Value: " + line[12] + ", Error: " + err.Error())
				mc = 0
			}

			cov := &models.Coverage{
				MeasurementID:      measurementID,
				CommitID:           commitID,
				Type:               testType,
				Group:              line[0],
				Package:            line[1],
				Class:              line[2],
				InstructionMissed:  im,
				InstructionCovered: ic,
				BranchMissed:       bm,
				BranchCovered:      bc,
				LineMissed:         lm,
				LineCovered:        lc,
				ComplexityMissed:   cm,
				ComplexityCovered:  cc,
				MethodMissed:       mm,
				MethodCovered:      mc,
			}
			models.CreateCoverage(db, cov)
		}
	}

	return nil
}
