package main

import (
	"fmt"
	"go-repo-downloader/models"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/jinzhu/gorm"
)

func Measure(db *gorm.DB, repoDir string, repository models.Repository, commitID uint, currCommit *object.Commit, changes object.Changes) {
	measurement := &models.Measurement{RepositoryID: repository.ID}
	models.CreateMeasurement(db, measurement)

	dt := time.Now()
	fmt.Println(currCommit.Hash.String() + " - " + dt.String())

	err := Checkout(repository.Name, currCommit.Hash.String())
	if err != nil {
		fmt.Println("Error checkout commit " + currCommit.Hash.String() + " " + err.Error())
		log.Println("Error checkout commit " + currCommit.Hash.String() + " " + err.Error())
	} else {

		ok := MvnCompile(repoDir)
		if ok {
			MeasureMavenTests(db, repoDir, commitID, *measurement)
			JacocoTestCoverage(db, repoDir, "maven", measurement.ID)
			mavenClasspath := GetMavenDependenciesClasspath(repoDir)
			for _, file := range listJavaFiles(repoDir) {
				MeasureRandoopTests(db, repoDir, file, mavenClasspath, commitID, *measurement)
			}
			JacocoTestCoverage(db, repoDir, "randoop", measurement.ID)
		}
	}
}

func MeasureMavenTests(db *gorm.DB, repoDir string, commitID uint, measurement models.Measurement) {
	testResults, ok := MvnTest(db, repoDir, measurement.ID)
	if ok {
		for ind := range testResults {
			mr := &models.Test{MeasurementID: measurement.ID,
				Type:        "maven",
				ClassName:   testResults[ind].ClassName,
				CommitID:    commitID,
				TestsRun:    testResults[ind].TestsRun,
				Failures:    testResults[ind].Failures,
				Errors:      testResults[ind].Errors,
				Skipped:     testResults[ind].Skipped,
				TimeElapsed: testResults[ind].TimeElapsed}
			models.CreateTest(db, mr)
		}
	} else {
		log.Println("********************** CRITICAL ERROR ***************")
		log.Println("successAfter is false measuring maven tests")
	}
}

func MeasureRandoopTests(db *gorm.DB, repoDir, file, mavenClasspath string, commitID uint, measurement models.Measurement) {
	//java -classpath ${RANDOOP_JAR} randoop.main.Main gentests --classlist=myclasses.txt --time-limit=60
	//Randoop prints out is the name of the JUnit files containing the tests it generated

	okGen := generateRandoopTests(repoDir, file, mavenClasspath)

	// Compile and run the tests. (The classpath should include the code under test, the generated tests, and JUnit files junit.jar and hamcrest-core.jar. Classes in java.util.* are always on the Java classpath, so the myclasspath part is not needed in this particular example, but it is shown because you will usually need to supply it.)
	// export JUNITPATH=.../junit.jar:.../hamcrest-core.jar
	// javac -classpath .:$JUNITPATH ErrorTest*.java RegressionTest*.java -sourcepath .:path/to/files/under/test/
	// java -classpath .:$JUNITPATH:myclasspath org.junit.runner.JUnitCore ErrorTest
	// java -classpath .:$JUNITPATH:myclasspath org.junit.runner.JUnitCore RegressionTest

	if okGen {
		okComp := compileRandoopTests(repoDir, mavenClasspath)
		if okComp {
			testTime, numTests, perfMetrics, okTest := runRandoopTests(repoDir)
			if okTest {
				r := &models.Test{MeasurementID: measurement.ID,
					Type:      "randoop",
					ClassName: file,
					CommitID:  commitID,
					TestsRun:  numTests,
					// Failures:    failures,
					// Errors:      errors,
					// Skipped:     skipped,
					TimeElapsed: testTime}
				testID, err := models.CreateTest(db, r)
				if err != nil {
					log.Println("Error creating randoop: " + err.Error())
					fmt.Println("Error creating randoop: " + err.Error())
				} else {
					for _, perfMetric := range perfMetrics {
						rr := &models.TestResources{
							TestID: testID,
							Type:   "randoop",
							Resources: models.Resources{
								Cpu:        perfMetric.Cpu,
								Mem:        perfMetric.Mem,
								ReadCount:  perfMetric.IO.ReadCount,
								WriteCount: perfMetric.IO.WriteCount,
								ReadBytes:  perfMetric.IO.ReadBytes,
								WriteBytes: perfMetric.IO.WriteBytes,
							},
						}
						models.CreateTestResources(db, rr)
					}
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
		panic(err)
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
