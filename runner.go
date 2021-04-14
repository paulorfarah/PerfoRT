package main

import (
	// "bytes"
	// "fmt"
	// "go-repo-downloader/models"
	// "log"
	// "os"
	// "os/exec"
	// "strings"
	// "github.com/wcharczuk/go-chart/v2"
	// "github.com/wcharczuk/go-chart/v2/drawing"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
)

func main2() {
	// plotRandoopResults()

	// db := models.GetDB()
	// models.GetRandoopMetrics()
	// getReleaseList("junit-team", "junit4")
	// Measure()
	// listJavaFiles()
}

// func listJavaFiles() []string {
// 	var files []string

// 	root := "/mnt/sda4/go-work/src/github.com/paulorfarah/junit4" //"D:\\go-work\\src\\github.com\\paulorfarah\\repos\\junit4"
// 	err := filepath.Walk(root, visit(&files))
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, file := range files {
// 		fmt.Println(file)
// 	}
// 	return files
// }

// func visit(files *[]string) filepath.WalkFunc {
// 	return func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if filepath.Ext(path) == ".java" {
// 			*files = append(*files, path)
// 		}

// 		return nil
// 	}
// }

// func Measure() {
// 	repoName := "junit4"
// 	repoPath := "/mnt/sda4/go-work/src/github.com/paulorfarah/junit4"

// 	releases := []string{"05fe2a64f59127c02135be22f416e91260d6ede6", "1b683f4ec07bcfa40149f086d32240f805487e66", "038f7518fc1018b26df608e3e5dce6db4611be29", "c0bdd7d4312862dbc6e1a8430cf75024a18158c9", "17d340a7d2661f0a0c7e284b49cd70f5a4495d6b", "fc3813fba9e2250ddd96877d01f2f694127edb80", "69424956c3c0d1f983cc2d489bcd7bebbf8b67a9", "6e1d2e6ebbc484af60c06bd26cc55349b352e49e", "64155f8a9babcfcf4263cf4d08253a1556e75481", "ceaafcafc7d7ac8d80509ab046acc7a472ddb515", "8d63bc65ace8c36eb3b189528c06e201f2caa535", "0550afa689cc5c7860d4378b27ecd38d35c05570", "c2e4d911fadfbd64444fb285342a8f1b72336169", "88c28a42a6fb7dc462c4bc504189a76a815fc265", "45a44647e7306262162e1346b750c3209019f2e1", "61f06547599bb6b98bca99d5bc457eb20bc17cab", "7ec443add70809418d2bbe1314cd4744742d854d", "5a9eb89e15ca86a3db6a4e21f5c9f94f9ab8fb60", "751f75986b11336ac8310d73c89003b0b09ecb92", "a30e87b6ac67f14a42b97d427bb1c8c6ba18cd87", "832bb8322f2ca09af52769a0198b276269b53988", "5a3a326096cf65a58272ee89a5ef1c164cfd9d33", "0c45278e830fdd1fc752a0eb1a3b25a3395d3e0e", "7aac4b19d359285041ccb51d575235339a1a8be0", "7aac4b19d359285041ccb51d575235339a1a8be0", "a8629da96207e1ce71ead9ba9f85bc324f09bcab", "b5e9885854a0d594451800b9127eb50afb645433", "a0f0ee1b3f72d9361eb09b3a25156c69a748aa47"}

// 	db := models.GetDB()
// 	platform, err := models.FindPlatformByName(db, "github")
// 	if err != nil {
// 		fmt.Println("Create new record...")
// 		//fmt.Println(err)
// 		platform = &models.Platform{Name: "github"}
// 		models.CreatePlatform(db, platform)
// 	}
// 	repository, err := models.FindRepositoryByName(db, repoName)
// 	if err != nil {
// 		fmt.Println("create new repo")
// 		fmt.Println(err)
// 		repository = &models.Repository{PlatformID: platform.ID, Name: repoName}
// 		models.CreateRepository(db, repository)
// 	}
// 	measurement := &models.Measurement{RepositoryID: repository.ID}
// 	models.CreateMeasurement(db, measurement)

// 	//before
// 	hash := releases[0]
// 	err = Checkout(repoName, hash)
// 	if err != nil {
// 		fmt.Println("Error in Before checkout " + hash + " " + err.Error())
// 	}

// 	MvnCompile(repoPath)
// 	successBefore, testResultsBefore := MvnTest(repoPath)

// 	for i := 1; i < len(releases); i++ {

// 		//after
// 		hash = releases[i]
// 		err = Checkout(repoName, hash)
// 		if err != nil {
// 			fmt.Println("Error in After checkout " + hash + " " + err.Error())
// 		}

// 		MvnCompile(repoPath)
// 		successAfter, testResultsAfter := MvnTest(repoPath)

// 		fmt.Println(successBefore, successAfter)
// 		fmt.Println("len: ", len(testResultsBefore), len(testResultsAfter))
// 		if successBefore && successAfter == true {
// 			if len(testResultsBefore) == len(testResultsAfter) {
// 				for ind := range testResultsBefore {
// 					fmt.Println(testResultsBefore[ind].ClassName, testResultsAfter[ind].ClassName)
// 					if testResultsBefore[ind].ClassName == testResultsAfter[ind].ClassName {
// 						mr := &models.MeasurementResults{MeasurementID: measurement.ID,
// 							Type:              byte('R'),
// 							ClassName:         testResultsBefore[ind].ClassName,
// 							TestsRunBefore:    testResultsBefore[ind].TestsRun,
// 							FailuresBefore:    testResultsBefore[ind].Failures,
// 							ErrorsBefore:      testResultsBefore[ind].Errors,
// 							SkippedBefore:     testResultsBefore[ind].Skipped,
// 							TimeElapsedBefore: testResultsBefore[ind].TimeElapsed,
// 							TestsRunAfter:     testResultsAfter[ind].TestsRun,
// 							FailuresAfter:     testResultsAfter[ind].Failures,
// 							ErrorsAfter:       testResultsAfter[ind].Errors,
// 							SkippedAfter:      testResultsAfter[ind].Skipped,
// 							TimeElapsedAfter:  testResultsAfter[ind].TimeElapsed}
// 						models.CreateMeasurementResults(db, mr)
// 					} else {
// 						fmt.Println("********************** CRITICAL ERROR ***************")
// 						fmt.Println("Class name of tests before and after are different, not considering this result")
// 					}
// 				}
// 			} else {
// 				fmt.Println("********************** CRITICAL ERROR ***************")
// 				fmt.Println("size of tests before and after are different, not considering these results")
// 			}
// 		} else {
// 			fmt.Println("********************** CRITICAL ERROR ***************")
// 			fmt.Println("size of tests before and after are different, not considering these results")
// 		}

// 		successBefore, testResultsBefore = successAfter, testResultsAfter
// 	}
// }

func getReleaseList(owner, repo string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{Page: 2, PerPage: 10}
	releases, rsp, err := client.Repositories.ListReleases(ctx, owner, repo, opt)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n%+v\n", releases)
	// for _, r := range releases {
	// 	fmt.Println("--------------------------")
	// 	fmt.Println(*r.Name)
	// 	fmt.Println(r.PublishedAt.String())
	// 	fmt.Println(*r.URL)

	// }
	fmt.Printf("\n%+v\n", rsp)
}

func randoop() {
	fmt.Println("teste")
	// cmd := exec.Command("java", "-classpath", "/mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR", "randoop.main.Main", "gentests", "--testclass=testproject.Test")
	// script := CreateRandoopScript("testproject.Test")
	// cmd := exec.Command("bash " + script)
	c := "java -classpath /mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR randoop.main.Main gentests --testclass=testproject.Test > testproject.Test.txt"
	cmd := exec.Command("bash", "-c", c)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("\n[>>ERROR]: Cannot run randoop gentests (" + fmt.Sprint(err) + "): " + stderr.String())
		fmt.Println(out)
	} else {
		fmt.Println("\n [>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>SUCCESS]: Randoop executed successully!")
		fmt.Println(out.String())
		// fmt.Println(ReadRandoopResults("testproject.Test.txt"))

	}
}

func CreateRandoopScript(class string) string {
	fn := strings.ReplaceAll(class, ".", "_") + ".sh"
	// Create new file
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = os.Chmod(fn, 0700)
	if err != nil {
		log.Fatal(err)
	}

	c := "java -classpath /mnt/sda4/go-work/src/github.com/paulorfarah/repos/TestProject/target/classes:$RANDOOP_JAR randoop.main.Main gentests --testclass=" + class + " > " + class + ".txt"
	_, err2 := f.WriteString(c)

	if err2 != nil {
		log.Fatal(err2)
	}

	return fn
}
