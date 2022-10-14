package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"perform/models"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/joshdk/go-junit"
	"gorm.io/gorm"
)

func Measure(db *gorm.DB, measurement models.Measurement, repoDir string, repository models.Repository, commit *models.Commit, currCommit *object.Commit) { //}, wg *sync.WaitGroup) {
	// defer wg.Done()
	// dt := time.Now()
	// fmt.Println(currCommit.Hash.String() + " - " + dt.String())

	// src := ".." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repository.Name
	// dst := "copy" + string(os.PathSeparator) + repository.Name

	// err := CopyDirectory(src, dst)
	// if err != nil {
	// 	fmt.Println("Error copying commit directory: ", err.Error())
	// 	log.Println("Error copying commit directory: ", err.Error())
	// }

	err := Checkout(repository.Name, currCommit.Hash.String())
	if err != nil {
		fmt.Println("Error checkout commit " + currCommit.Hash.String() + " " + err.Error())
		log.Println("Error checkout commit " + currCommit.Hash.String() + " " + err.Error())
	} else {

		switch buildtool := checkBuildTool(repoDir); buildtool {
		case "":
			fmt.Println("ATTENTION: Maven or Gradle files not found in ", repoDir)
		case "maven":
			// projectModules := getProjectModules(repoDir)
			// if len(projectModules) == 0 {
			// buildPath := repoDir + string(os.PathSeparator)
			javaVer := getMavenJavaVersion(repoDir)
			log.Println("Java version: ", javaVer)
			commit.JavaVersion = javaVer
			db.Save(commit)

			// ok := MvnInstall(repoDir, javaVer)
			MvnInstall(repoDir, javaVer)
			// ok := MvnCompile(repoDir)
			// log.Println("MvnCompile ok: ", ok)
			// if ok {
			MeasureMavenTests(db, repoDir, javaVer, *commit, measurement)
			// JacocoTestCoverage(db, repoDir, "maven", "maven", measurement.ID, commit.ID)
			// mavenClasspath := GetMavenDependenciesClasspath(repoDir)
			// for _, file := range listJavaFiles(repoDir) {
			// 	MeasureRandoopTests(db, repoDir, file, "maven", mavenClasspath, commitID, measurement)
			// }
			// JacocoTestCoverage(db, repoDir, "randoop", "maven", measurement.ID)
			// }
			// } else {
			// 	MvnInstall(repoDir)
			// 	ok := MvnCompile(repoDir)
			// 	if ok {
			// 		for _, projectPath := range projectModules {
			// 			// buildPath := repoDir + string(os.PathSeparator) + projectPath
			// 			MeasureMavenTests(db, repoDir, projectPath, commitID, measurement)
			// 			// (db, repoDir, "maven", "maven", measurement.ID)
			// 			// mavenClasspath := GetMavenDependenciesClasspath(repoDir)
			// 			// for _, file := range listJavaFiles(repoDir) {
			// 			// 	MeasureRandoopTests(db, repoDir, file, "maven", mavenClasspath, commitID, measurement)
			// 			// }
			// 			// JacocoTestCoverage(db, repoDir, "randoop", "maven", measurement.ID)

			// 		}
			// 	}
			// } else {
			// 	log.Println("ERROR:  Compilation failed!")
			// }
		case "gradle":
			projectPaths := getProjectPaths(repoDir)
			if len(projectPaths) == 0 {
				buildPath := repoDir + string(os.PathSeparator)
				ok := GradleBuild(buildPath)
				if ok {
					MeasureGradleTests(db, buildPath, *commit, measurement)
					// JacocoTestCoverage(db, buildPath, "gradle", "gradle", measurement.ID, commitID)
					// gradleClasspath := GetGradleDependenciesClasspath(buildPath)
					// for _, file := range listJavaFiles(buildPath) {
					// 	MeasureRandoopTests(db, buildPath, file, "gradle", gradleClasspath, commitID, measurement)
					// }
					// JacocoTestCoverage(db, buildPath, "randoop", "gradle", measurement.ID)
				}

			} else {
				for _, projectPath := range projectPaths {
					buildPath := repoDir + string(os.PathSeparator) + projectPath
					ok := GradleBuild(buildPath)
					if ok {
						MeasureGradleTests(db, buildPath, *commit, measurement)
						// JacocoTestCoverage(db, buildPath, "gradle", "gradle", measurement.ID)
						// gradleClasspath := GetGradleDependenciesClasspath(buildPath)
						// for _, file := range listJavaFiles(buildPath) {
						// 	MeasureRandoopTests(db, buildPath, file, "gradle", gradleClasspath, commitID, measurement)
						// }
						// JacocoTestCoverage(db, buildPath, "randoop", "gradle", measurement.ID)
					}
				}
			}
		default:
			fmt.Println("Error: Build Tool not recognized!")
			log.Println("Error: Build Tool not recognized!")
			log.Fatal("Error: Build Tool not recognized!")
		}
	}
}

func MeasureMavenTests(db *gorm.DB, repoDir, javaVer string, commit models.Commit, measurement models.Measurement) {

	projectModules := getProjectModules(repoDir)
	// fmt.Println("modules: ", projectModules)
	path := repoDir
	packName := os.Getenv("package")

	minTestTime := 0.0
	minTestTimeStr, ok := os.LookupEnv("min_test_time")
	if ok {
		minTestTime, _ = strconv.ParseFloat(minTestTimeStr, 32)
	}

	profiler := "/perform-tracer-1.11.jar"
	if javaVer != "" {
		if strings.Contains(javaVer, "8") {
			profiler = "/perform-tracer-1.8.jar"
		}
	}

	//set environment variable to activate profiler during testcases execution
	localpath, errPath := os.Getwd()
	if errPath != nil {
		log.Println(errPath)
		fmt.Println("error getting current path: ", errPath.Error())
	}

	mavenClasspath := GetMavenDependenciesClasspath(repoDir, javaVer)
	// log.Println("- junit testcase: ", path, tc.ClassName, testName)
	// fmt.Println("- junit testcase: ", path, tc.ClassName, testName)

	// read testcase ignore file
	tcignoreMap, errIgn := ReadTCIgnoreMap(".tcignore/" + packName)
	if errIgn != nil {
		log.Println("Error reading lfist of ignored testcases: ", errIgn)
	}

	for _, module := range projectModules {
		// fmt.Printf("\n*** Module: %s\n\n", module)
		// log.Printf("\n*** Module: %s\n\n", module)
		if module != "" {
			path = repoDir + "/" + module
		}

		MvnTest(db, path, javaVer, measurement.ID, commit.ID)
		JacocoTestCoverage(db, path, javaVer, "maven", "maven", measurement.ID, commit.ID)
		files, err := ioutil.ReadDir(path + "/target/surefire-reports/")

		if err != nil {
			log.Printf("cannot find surefire results in path: %s - %s\n", path+"/target/surefire-reports/", err.Error())
			fmt.Printf("cannot find surefire results in path: %s - %s\n", path+"/target/surefire-reports/", err.Error())
		} else {

			localClasspath := ".:target/test-classes/:target/classes:"
			if module != "" {
				localClasspath += module + "/target/test-classes/:" + module + "/target/classes/:"
			}

			for _, file := range files {
				// log.Println("test file: ", file.Name())
				if !file.IsDir() {
					suites := ParseMavenTestResults(path + "/target/surefire-reports/" + file.Name())
					// maxGoroutines, err := strconv.Atoi(os.Getenv("threads"))
					// if err != nil {
					// 	fmt.Println("ATTENTION: Error reading number of threads from .env file, using 1")
					// 	maxGoroutines = 1
					// }
					// guard := make(chan struct{}, maxGoroutines)
					// count := -1
					for _, test := range suites.TestCases {
						testName := test.ClassName + "#" + test.Name
						// log.Println("testcase:", testName)
						_, ignore := tcignoreMap[testName]
						if !ignore {
							if tcTime, err := strconv.ParseFloat(test.Time, 32); err == nil {
								if tcTime >= minTestTime {
									// guard <- struct{}{} // would block if guard channel is already filled
									// go func(n int) {
									classname := strings.Replace(test.ClassName, ".", "/", -1)
									filename := classname + ".java"
									testSuite, errF := models.FindFileByEndsWithNameAndCommit(db, filename, commit.ID)
									if errF != nil {
										fmt.Println("error finding file: ", test.ClassName, commit.CommitHash)
									}
									tc := &models.TestCase{
										Type:      "maven",
										ClassName: test.ClassName,
										FileID:    testSuite.ID,
										Name:      test.Name,
									}
									_, errTC := models.CreateTestCase(db, tc)
									if errTC != nil {
										fmt.Println("Error creating test case: ", errTC.Error())
									}
									// RunMavenTestCase(db, repoDir, module, tc, measurement.ID, commit)
									RunJUnitTestCase(db, repoDir, module, javaVer, tc, measurement, commit, packName, profiler, localpath, mavenClasspath, localClasspath)
									// <-guard
									// count++
									// }(count)
								} else {
									log.Printf("Testcase %s#%s was not executed because its time %s is lower than the mininum test time threshold (%f).\n", test.ClassName, test.Name, test.Time, minTestTime)
								}
							} else {
								log.Println("ERROR: Invalid testcase time!!! " + test.ClassName + "#" + test.Name + " - " + test.Time)
							}
						}
					}

				}
			}
		}
	}
}

func MeasureGradleTests(db *gorm.DB, repoDir string, commit models.Commit, measurement models.Measurement) {
	ok := GradleTest(db, repoDir, measurement.ID)
	if ok {
		// JacocoTestCoverage(db, repoDir, "gradle", "gradle", measurement.ID, commit.ID)
		// read tests xml file
		// fmt.Printf("repoDir gradle tests: %s\n", repoDir)
		suites, err := junit.IngestDir(repoDir + "/build/test-results/test/")
		if err != nil {
			log.Fatalf("failed to ingest JUnit xml %v", err)
		}
		// fmt.Println("suites: ", suites)
		for _, suite := range suites {
			// fmt.Println(suite.Name)
			// fmt.Printf("%s\n", suite.Tests)
			for _, test := range suite.Tests {
				// fmt.Println(test.Classname + ".java")
				dt := time.Now()
				fmt.Printf("  %s %s\n", test.Name, dt.String())
				// if test.Error != nil {
				// 	fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
				// } else {
				// 	fmt.Printf("    %s %f\n", test.Status, test.Duration.Seconds())
				// }
				classname := strings.Replace(test.Classname, ".", "/", -1)
				filename := classname + ".java"
				// fmt.Println(filename)
				testSuite, errF := models.FindFileByEndsWithNameAndCommit(db, filename, commit.ID)
				if errF != nil {
					fmt.Println("error finding file: ", test.Classname, commit.CommitHash)
				}
				// fmt.Println("testSuite: ", testSuite)

				// errorMsg := ""
				// if test.Error != nil {
				// 	errorMsg = test.Error.Error()
				// }
				tc := &models.TestCase{
					Type:      "gradle",
					ClassName: test.Classname,
					// Duration :      test.Duration,
					FileID: testSuite.ID,
					Name:   test.Name,
					// Status:        string(test.Status),
					// Error:         errorMsg,
					// Message:       test.Message,
					// SystemErr:     string(test.SystemErr),
					// SystemOut:     string(test.SystemOut),
				}
				_, errTC := models.CreateTestCase(db, tc)
				if errTC != nil {
					fmt.Println("Error creating test case: ", errTC.Error())
				}

				//gradle test --test "com.cloudhadoop.emp.SuiteTest.testTestCaseName"
				RunGradleTestCase(db, repoDir, tc, measurement.ID)

			}
		}

	} else {
		log.Println("********************** CRITICAL ERROR ***************")
		log.Println("successAfter is false measuring maven tests")
	}
}

func MeasureRandoopTests(db *gorm.DB, repoDir, file, buildTool, buildToolClasspath string, commitID uint, measurement models.Measurement) {
	//java -classpath ${RANDOOP_JAR} randoop.main.Main gentests --classlist=myclasses.txt --time-limit=60
	//Randoop prints out is the name of the JUnit files containing the tests it generated

	dir, pack := parseProjectPath(file)
	if dir != "" {
		dir += string(os.PathSeparator)
	}

	//create gentest dir in project to save log files of randoop generation phase
	fmt.Println("create gentest dir in project to save log files of randoop generation phase: ", dir+"gentest")
	_, errd := os.Stat(dir + "gentest")
	if os.IsNotExist(errd) {
		err := os.Mkdir(dir+"gentest", 0755)
		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
		} else {
			fmt.Println("gentest directory created...")
		}
	}

	path := strings.Split(pack, ".java")[0]
	// fmt.Println("path: ", path)
	randoopJar := "${RANDOOP_JAR}"
	cpSep := ":"
	if runtime.GOOS == "windows" {
		randoopJar = "%RANDOOP_JAR%"
		cpSep = ";"
	}

	envRandoopJar := os.Getenv("RANDOOP_JAR")
	// remove old tests
	// deleteOldRandoopTests()

	classpath := ""
	switch buildTool {
	case "maven":
		classpath += dir + "target" + string(os.PathSeparator) + "classes" + cpSep
	case "gradle":
		classpath += dir + "build" + string(os.PathSeparator) + "classes" + cpSep +
			dir + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "main" + cpSep +
			dir + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "test"
	}
	classpath += buildToolClasspath
	className := strings.ReplaceAll(path, string(os.PathSeparator), ".")

	fmt.Println("------------------------------------------------ Generating Randoop tests for " + file + "...")

	// gradle-project-example/app/src/test/java
	dirSourceTest := dir + "src" + string(os.PathSeparator) + "test" + string(os.PathSeparator) + "java" + string(os.PathSeparator)
	okGen := generateRandoopTests(db, dirSourceTest, classpath, cpSep, randoopJar, envRandoopJar, className, measurement, commitID)

	// Compile and run the tests. (The classpath should include the code under test, the generated tests, and JUnit files junit.jar and hamcrest-core.jar. Classes in java.util.* are always on the Java classpath, so the myclasspath part is not needed in this particular example, but it is shown because you will usually need to supply it.)
	// export JUNITPATH=.../junit.jar:.../hamcrest-core.jar
	// javac -classpath .:$JUNITPATH ErrorTest*.java RegressionTest*.java -sourcepath .:path/to/files/under/test/
	// java -classpath .:$JUNITPATH:myclasspath org.junit.runner.JUnitCore ErrorTest
	// java -classpath .:$JUNITPATH:myclasspath org.junit.runner.JUnitCore RegressionTest

	if okGen {
		dirClassTest := dir + "build" + string(os.PathSeparator) + "classes" + string(os.PathSeparator) + "java" + string(os.PathSeparator) + "test"
		okComp := compileRandoopTests(dirSourceTest, dirClassTest, classpath, cpSep)
		if okComp {
			// _, _, perfMetrics, okTest := runRandoopTests(dirSourceTest, classpath, cpSep)
			_, _, resources, okTest := runRandoopTests(dirSourceTest, classpath, cpSep)
			filename, errF := models.FindFileByEndsWithNameAndCommit(db, path, commitID)
			if errF != nil {
				fmt.Println(errF.Error())
			}
			if okTest {
				r := &models.TestCase{
					Type:      "randoop",
					ClassName: file,

					// Duration:  time.Duration(testTime * float64(time.Second)),
					// TestSuiteID: testSuite.ID,
					FileTargetID: filename.ID,
					Name:         file,
					// Status:    string(test.Status),
					// Error:     errorMsg,
					// Message:   test.Message,
					// SystemErr: string(test.SystemErr),
					// SystemOut: string(test.SystemOut),

					// Duration:  testTime,
					// TestsRun:  numTests,
					// Failures:    failures,
					// Errors:      errors,
					// Skipped:     skipped,
					// TimeElapsed: testTime
				}
				testID, err := models.CreateTestCase(db, r)
				if err != nil {
					log.Println("Error creating randoop: " + err.Error())
					fmt.Println("Error creating randoop: " + err.Error())
				} else {
					rr := &models.Run{
						TestCaseID: testID,
						Type:       "randoop",
					}
					models.CreateRun(db, rr)

					for i := range resources {
						resources[i].RunID = rr.ID
					}
					db.CreateInBatches(resources, 1000)

					// ATTENTION: deprecates perfmetric changed by resources implemented above
					// for _, perfMetric := range perfMetrics {

					// 	resource := models.Resource{
					// 		RunID:          rr.ID,
					// 		CpuPercent:     perfMetric.CpuPercent,
					// 		MemPercent:     perfMetric.MemoryPercent,
					// 		MemoryInfoStat: *perfMetric.MemoryInfo,
					// 		IOCountersStat: *perfMetric.IOCounters,
					// 		PageFaultsStat: *perfMetric.PageFaults,
					// 		// AvgStat:           *perfMetric.Load,
					// 		// VirtualMemoryStat: *perfMetric.VirtualMemory,
					// 		// SwapMemoryStat:    *perfMetric.SwapMemory,
					// 		// CPUTime:           perfMetric.CPUTime,
					// 		// DiskIOCounters:    perfMetric.DiskIOCounters,
					// 		// NetIOCounters:     perfMetric.NetIOCounters,
					// 	}
					// 	_, err = models.CreateResource(db, &resource)
					// 	if err != nil {
					// 		fmt.Println("ERROR creating randoop resource: ", err.Error())
					// 	}
					// 	// for _, cpuTime := range perfMetric.CPUTimes {
					// 	// 	models.CreateCPUTimes(db, &models.CPUTimes{
					// 	// 		ResourceID: resource.ID,
					// 	// 		CPU:        cpuTime.CPU,
					// 	// 		User:       cpuTime.User,
					// 	// 		System:     cpuTime.System,
					// 	// 		Idle:       cpuTime.Idle,
					// 	// 		Nice:       cpuTime.Nice,
					// 	// 		Iowait:     cpuTime.Iowait,
					// 	// 		Irq:        cpuTime.Irq,
					// 	// 		Softirq:    cpuTime.Softirq,
					// 	// 		Steal:      cpuTime.Steal,
					// 	// 		Guest:      cpuTime.Guest,
					// 	// 		GuestNice:  cpuTime.GuestNice,
					// 	// 	})

					// 	// }

					// 	// for i, diskIOCounter := range perfMetric.DiskIOCounters {
					// 	// 	models.CreateDiskIOCounters(db, &models.DiskIOCounters{
					// 	// 		ResourceID:       resource.ID,
					// 	// 		Device:           i,
					// 	// 		ReadCount:        diskIOCounter.ReadCount,
					// 	// 		MergedReadCount:  diskIOCounter.MergedReadCount,
					// 	// 		WriteCount:       diskIOCounter.WriteCount,
					// 	// 		MergedWriteCount: diskIOCounter.MergedWriteCount,
					// 	// 		ReadBytes:        diskIOCounter.ReadBytes,
					// 	// 		WriteBytes:       diskIOCounter.WriteBytes,
					// 	// 		ReadTime:         diskIOCounter.ReadTime,
					// 	// 		WriteTime:        diskIOCounter.WriteTime,
					// 	// 		IopsInProgress:   diskIOCounter.IopsInProgress,
					// 	// 		IoTime:           diskIOCounter.IoTime,
					// 	// 		WeightedIO:       diskIOCounter.WeightedIO,
					// 	// 		Name:             diskIOCounter.Name,
					// 	// 		SerialNumber:     diskIOCounter.SerialNumber,
					// 	// 		Label:            diskIOCounter.Label,
					// 	// 	})
					// 	// }

					// 	// for i, netIOCounter := range perfMetric.NetIOCounters {
					// 	// 	models.CreateNetIOCounters(db, &models.NetIOCounters{
					// 	// 		ResourceID:  resource.ID,
					// 	// 		NICID:       uint(i),
					// 	// 		Name:        netIOCounter.Name,
					// 	// 		BytesSent:   netIOCounter.BytesSent,
					// 	// 		BytesRecv:   netIOCounter.BytesRecv,
					// 	// 		PacketsSent: netIOCounter.PacketsSent,
					// 	// 		PacketsRecv: netIOCounter.PacketsRecv,
					// 	// 		Errin:       netIOCounter.Errin,
					// 	// 		Errout:      netIOCounter.Errout,
					// 	// 		Dropin:      netIOCounter.Dropin,
					// 	// 		Dropout:     netIOCounter.Dropout,
					// 	// 		Fifoin:      netIOCounter.Fifoin,
					// 	// 		Fifoout:     netIOCounter.Fifoout,
					// 	// 	})
					// 	// }
					// }
				}
			}

		}

	}

	// CollectRandoopMetrics(repoDir, repository.Name, commit.PreviousCommitHash, change.From.Name, commit.CommitHash, change.To.Name, changeObj.ID)
}

// func MeasureChanges(db *gorm.DB, repoDir string, repository models.Repository, commit models.Commit, changes object.Changes) {
// 	//randoop
// 	for _, change := range changes {
// 		// fmt.Println(change.From.Name)
// 		// fmt.Println(change.To.Name)
// 		// fmt.Println(change.Action())
// 		// fmt.Println(change.Files())
// 		// fmt.Println("------------------- start")
// 		// fmt.Println(change.Patch())

// 		patch, _ := change.Patch()
// 		diff, _ := diffparser.Parse(patch.String())

// 		//files
// 		count := 0
// 		for _, file := range diff.Files {
// 			// fmt.Println("************************** file: ", file)

// 			sc := fmt.Sprintf("%d", count)

// 			fNew, _ := os.Create("results/" + commit.CommitHash + "f" + sc + "_new.java")
// 			defer fNew.Close()

// 			fOld, _ := os.Create("results/" + commit.CommitHash + "f" + sc + "_old.java")
// 			defer fOld.Close()

// 			// //hunks
// 			for _, hunk := range file.Hunks {
// 				for _, l := range hunk.NewRange.Lines {
// 					fNew.WriteString(l.Content + "\n")
// 				}
// 				for _, l := range hunk.OrigRange.Lines {
// 					fOld.WriteString(l.Content + "\n")
// 				}
// 			}
// 			count++

// 		}

// 		hasher := sha1.New()
// 		patch, err := change.Patch()
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		hasher.Write([]byte(patch.String()))
// 		changeSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
// 		//fmt.Println(changeSha)
// 		//	id := fmt.Sprintf("%s",currCommit.ID)
// 		//	fmt.Printf("*************  %s\n", id)
// 		_, err = models.FindChangeByHash(db, changeSha, commit.ID)
// 		if err != nil {
// 			fmt.Println("new change")
// 			fmt.Println(err)
// 			action, err := change.Action()
// 			if err != nil {
// 				fmt.Println(err.Error()) //return err
// 			}
// 			changeObj := &models.Change{CommitID: commit.ID, ChangeHash: changeSha, FileFrom: change.From.Name, FileTo: change.To.Name, Action: action.String(), Patch: patch.String()}
// 			models.CreateChange(db, changeObj)

// 			//call randoop
// 			fmt.Println(change.From.Name)
// 			if action.String() == "Modify" &&
// 				strings.Contains(change.From.Name, ".java") &&
// 				strings.Contains(change.To.Name, ".java") &&
// 				!strings.HasPrefix(change.From.Name, "src/test/") &&
// 				!strings.HasPrefix(change.From.Name, "src/test/") {
// 				// CollectRandoopMetrics(repoDir, repository.Name, commit.PreviousCommitHash, change.From.Name, commit.CommitHash, change.To.Name, changeObj.ID)
// 			}
// 		} else {
// 			fmt.Println("change already exists in database...")
// 		}
// 	}
// }

func listJavaFiles(repoDir string) []string {
	var files []string
	err := filepath.Walk(repoDir, visit(&files))
	if err != nil {
		log.Println("ERROR listing java files: ", err)
	}
	return files
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(path) == ".java" {
			*files = append(*files, path)
		}

		return nil
	}
}

// exists returns whether the given file or directory exists
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkBuildTool(repoDir string) string {
	gradleExists, err := fileExists(repoDir + "/" + "settings.gradle")
	if err != nil {
		fmt.Println("ERROR looking for settings.gradle...")
	}
	if gradleExists {
		return "gradle"
	}

	gradleExists, err = fileExists(repoDir + "/" + "build.gradle")
	if err != nil {
		fmt.Println("ERROR looking for settings.gradle...")
	}
	if gradleExists {
		fmt.Println("build.gradle found...")
		return "gradle"
	}

	pomExists, err := fileExists(repoDir + "/" + "pom.xml")
	if err != nil {
		fmt.Println("ERROR looking for pom.xml...")
	}
	if pomExists {
		return "maven"
	}
	return ""

}

func getProjectModules(repoDir string) []string {
	var includes []string

	includes = append(includes, "")
	pomPath := repoDir + "/pom.xml"
	// parsedPom, err := gopom.Parse(pomPath)
	parsedPom, err := ParsePom(pomPath)
	if err != nil {
		log.Printf("unable to unmarshal pom file. Reason: %s\n", err)
	}
	if err != nil {
		fmt.Printf("unable to unmarshal pom file. File: %s Reason: %s\n", repoDir+"/pom.xml", err)
		log.Printf("Unable to unmarshal pom file. File: %s Reason: %s\n", repoDir+"/pom.xml", err)
	}

	for _, m := range parsedPom.Modules {
		includes = append(includes, m)
	}
	for _, p := range parsedPom.Profiles {
		for _, m := range p.Modules {
			includes = append(includes, m)
		}
	}
	return includes
}

func getMavenJavaVersion(repoDir string) string {
	var version string

	// Java version:  /usr/lib/jvm/java-${java.version}.0-openjdk-amd64

	pomPath := repoDir + "/pom.xml"
	// fmt.Println("pomPath: ", pomPath)
	// parsedPom, err := gopom.Parse(pomPath)
	parsedPom, err := ParsePom(pomPath)
	if err != nil {
		log.Printf("unable to unmarshal pom file. Reason: %s\n", err)
	}

	for k, v := range parsedPom.Properties.Entries {
		if k == "maven.compiler.source" || k == "compileSource" || k == "java.version" || k == "java_source_version" {
			version = strings.Replace(v, "1.", "", 1)
			break
		}
	}

	if version == "" {
		//search in plugins
		// var pomPath2 string = repoDir + "/pom.xml"
		// parsedPom2, err := 	.Parse(pomPath2)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		for _, plug := range parsedPom.Build.BuildBase.Plugins {
			// fmt.Println(plug.ArtifactID)
			// fmt.Println(plug.Configuration.Source)
			// fmt.Println(plug.Configuration.Target)
			if plug.ArtifactID == "maven-compiler-plugin" {
				// fmt.Println("maven-compiler-plugin: ", plug.Version)
				version = strings.Replace(plug.Version, "1.", "", 1)
				break
			}
		}
	}
	if version == "" {
		b, err := os.ReadFile(pomPath)
		if err != nil {
			fmt.Print(err)
		}

		str := string(b)
		vnum := between(str, "<source>", "</source>")
		version = strings.Replace(vnum, "1.", "", 1)
	}

	if version == "" {
		version = "11"
	}
	fmt.Println("JAVA: ", version)
	version = "/usr/lib/jvm/java-" + version + "-openjdk-amd64"
	return version
}

func getProjectPaths(repoDir string) []string {
	var includes []string
	file, err := os.Open(repoDir + "/settings.gradle")
	if err != nil {
		log.Printf("cannot find settings.gradle file: %s\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`\(\'(.*?)\'\)`)
	for scanner.Scan() {
		str1 := scanner.Text()
		// fmt.Println(str1)
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		if strings.Contains(str1, "include('") {
			submatchall := re.FindAllString(str1, -1)
			for _, element := range submatchall {
				element = strings.Trim(element, "('")
				element = strings.Trim(element, "')")
				includes = append(includes, element)
			}
		}
	}

	return includes
}

func CopyDirectory(srcDir, dest string) error {
	deleteDir(dest)
	fmt.Println("Copying directory")
	fmt.Println(srcDir)
	fmt.Println(dest)
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func deleteDir(dir string) error {
	fmt.Println("deleting directory: " + dir)
	log.Println("deleting directory: " + dir)
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
