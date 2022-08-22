package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	// grmon "github.com/bcicen/grmon/agent"
	"net/http"
	_ "net/http/pprof"

	// grmon "github.com/bcicen/grmon/agent"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	_ "github.com/go-sql-driver/mysql"

	"perfrt/models"
)

func main() {
	fmt.Println("starting perfrt")

	go func() {
		http.ListenAndServe(":1234", nil)
	}()

	logFile, err := os.OpenFile("perfrt.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	// log.Println("starting...")
	url, ok := os.LookupEnv("repository")

	// fmt.Println("############## url: ", url)
	if ok {
		urlSplit := strings.Split(url, "/")
		repoName := urlSplit[4]
		repoDir := getParentDirectory() + "/repos/" + repoName
		// fmt.Println("repoDir: " + repoDir)
		// log.Println("repoDir: " + repoDir)

		// fmt.Println("git clone " + url)
		// log.Println("git clone " + url)

		repo, err := cloneRepository(url, repoDir)
		if err == nil {
			createDirs()
			db := models.GetDB()

			// platform
			platform, err := models.FindPlatformByName(db, "github")
			if err != nil {
				// log.Println("Create new platform: " + "github")
				platform = &models.Platform{Name: "github"}
				platformID, err := models.CreatePlatform(db, platform)
				if err != nil {
					fmt.Println("ERROR creating github platform: ", err.Error())

				}
				platform.ID = platformID
			}

			//search representative repositories

			//save repository in db
			repository, err := models.FindRepositoryByName(db, repoName)
			if err != nil {
				// log.Println("create new repo: " + repoName)
				repository = &models.Repository{PlatformID: platform.ID, Name: repoName}
				models.CreateRepository(db, repository)
			}

			// measurement
			// var measurement *models.Measurement
			measurement := &models.Measurement{RepositoryID: repository.ID}
			runsEnv, ok := os.LookupEnv("runs")
			// fmt.Println("############## runs: ", runsEnv)
			runs, err := strconv.Atoi(runsEnv)
			if err != nil {
				ok = false
			}
			if ok {
				measurement.Runs = runs
			} else {
				fmt.Println("ATTENTION: Number of runs not set, running with value 1!!!", "runs")
				measurement.Runs = 1
			}

			tcTimeOut := 3600
			tcTimeOutStr, ok := os.LookupEnv("testcase_timeout")
			if ok {
				tcTimeOut, err = strconv.Atoi(tcTimeOutStr)
				if err != nil {
					log.Println("WARNING: invalid testcase_timeout setting, using 1 hour.")
					tcTimeOut = 3600
				}
			} else {
				log.Println("WARNING: testcase_timeout setting not found, using 1 hour.")
			}
			// log.Println("monitoring time: ", tcTimeOut)
			measurement.TestcaseTimeout = time.Duration(tcTimeOut)

			var monitoringTime time.Duration
			monitoringTimeStr, ok := os.LookupEnv("monitoring_time")
			if ok {
				monitoringTime, err = time.ParseDuration(monitoringTimeStr + "s")
				if err != nil {
					log.Println("Error parsing monitoringTimeFlt: ", err)
					monitoringTime, _ = time.ParseDuration("1s")
				}
			} else {
				log.Println("WARNING: environment variable monitoring_time not found, using 1s!")
				monitoringTime, _ = time.ParseDuration("1s")
			}
			// log.Println("monitoring time: ", monitoringTime)
			measurement.MonitoringTime = monitoringTime

			models.CreateMeasurement(db, measurement)

			//issues
			//repository.Issues()
			// lastIssue, err := models.FindIssueByRepository(db, repository.ID)
			// fmt.Println("issue: ", lastIssue)
			// if err != nil {
			// 	fmt.Println("create new issues")
			// 	fmt.Println(err)
			// 	//issue = &models.Issue{Repository:repository.ID, Number: 1}
			// 	//models.CreateIssue(db, issue)
			// }
			// allIssues := models.GetIssues(lastIssue)
			// fmt.Println("issues...", allIssues)
			// for _, i := range allIssues {
			// 	fmt.Println("########################################", i.Title)
			// }

			packName := os.Getenv("package")
			releases, errRel := ReadListFromFile(".releases_" + packName)
			if errRel != nil {
				log.Println("Error reading file of releases: ", errRel)
			}

			//branches
			// branchCounter := 0
			// branches, _ := repo.Branches()
			// for {
			// 	branch, err := branches.Next()
			// 	if err != nil {
			// 		if err == io.EOF {
			// 			//Finished branch
			// 			break
			// 		} else {
			// 			fmt.Println("main/master branches not found.")
			// 			log.Fatal(err)
			// 		}
			// 	}
			// 	branchCounter++
			// 	fmt.Println("branch -->: ", branch.Name())

			// 	//commits
			// 	var cCommit *object.Commit
			// 	cCommit = nil
			// 	var cTree *object.Tree
			// 	cTree = nil

			//filter by dates
			// since := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)
			// until := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

			// commits, err := repo.Log(&git.LogOptions{From: branch.Hash()}) //, Since: &since, Until: &until})
			// if err != nil {
			// 	log.Println("Error in git log: " + err.Error())
			// }
			// defer commits.Close()
			// fmt.Println("---- commits ----")
			// i := 0
			// err = commits.ForEach(func(pCommit *object.Commit) error {
			// 	fmt.Println(pCommit.Hash.String())

			// ref, err := repo.Head()
			// if err != nil {
			// 	log.Println("ERROR retrieving the commit being pointed by HEAD")
			// 	fmt.Println("ERROR retrieving the commit being pointed by HEAD")
			// }

			for _, release := range releases {
				cCommit, err := repo.CommitObject(plumbing.NewHash(release))
				if err != nil {
					log.Println("ERROR retrieving commit: ", release, " (", err, ")")
					fmt.Println("ERROR retrieving commit: ", release, " (", err, ")")
					return
				}
				cTree, err := cCommit.Tree()
				if err != nil {
					log.Println("ERROR retrieving Tree: ", release, " (", err, ")")
					fmt.Println("ERROR retrieving Tree: ", release, " (", err, ")")
					return
				}
				// fmt.Println("*********************************************************************commDate: ", commDate)
				// fmt.Printf("\n----- commit %v: %v -----\n", strconv.Itoa(i), currCommit.Message)
				// fmt.Printf("###>  commit: %s <###\n", cCommit.Hash)
				//fmt.Println(currCommit.Author.Email)
				//fmt.Println(currCommit.Committer)
				//fmt.Println(currCommit.Message)
				//fmt.Printf("\nfile: %v\n", cs.Name)

				// pTree, err := pCommit.Tree()
				// if err != nil {
				// 	return err
				// }
				// if cTree != nil {
				// 	changes, err := pTree.Diff(cTree)
				// 	if err != nil {
				// 		return err
				// 	}
				//Author
				author, err := models.FindAccountByEmail(db, cCommit.Author.Email)
				if err != nil {
					// log.Println("create new author: " + cCommit.Author.Name)
					author = &models.Account{Email: cCommit.Author.Email, Name: cCommit.Author.Name}
					models.CreateAccount(db, author)
				}
				//Committer
				committer, err := models.FindAccountByEmail(db, cCommit.Committer.Email)
				if err != nil {
					// log.Println("create new committer: " + cCommit.Committer.Name)
					committer = &models.Account{Email: cCommit.Committer.Email, Name: cCommit.Committer.Name}
					models.CreateAccount(db, committer)
				}

				//Commit
				var commitId uint
				commit, err := models.FindCommitByHash(db, cCommit.Hash.String())
				if err != nil {
					// log.Println("create new commit: " + cCommit.Hash.String())
					// fmt.Println("#  create new commit: " + cCommit.Hash.String())
					// parent, errj := json.Marshal(currCommit.ParentHashes)
					// if errj != nil {
					// 	log.Println("Error Marshalling parent hashes: " + errj.Error())
					// }
					commit = &models.Commit{CommitHash: cCommit.Hash.String(),
						PreviousCommitHash: "", //pCommit.Hash.String(),
						RepositoryID:       repository.ID,
						// TreeHash:           cCommit.TreeHash.String(),
						// ParentHashes:       parent,
						AuthorID:      author.ID,
						AuthorDate:    cCommit.Author.When,
						CommitterID:   committer.ID,
						CommitterDate: cCommit.Committer.When,
						Subject:       cCommit.Message,
						Branch:        "", //branch.Name().String()
					}
					commitId, err = models.CreateCommit(db, commit)
					if err != nil {
						fmt.Printf("Error creating new commit %s\n", err.Error())
					}

				}
				// fmt.Println("commitId: ", commitId)
				// fmt.Println("commit.ID: ", commit.ID)

				//files
				// currTree.Files().ForEach(func(f *object.File) error {
				cTree.Files().ForEach(func(f *object.File) error {
					// contents := ""
					// if !(strings.HasSuffix(f.Name, ".class") || strings.HasSuffix(f.Name, ".jar")) {
					// 	contents, _ = f.Contents()
					// }
					isBin, _ := f.IsBinary()
					lines, _ := f.Lines()
					// fmt.Printf("%d	%s\n", commit.ID, f.Name)

					ls := []models.FileLine{}
					for _, l := range lines {
						ls = append(ls, models.FileLine{Line: l})
					}
					// fmt.Printf("Commit: %d - %s\n", commit.ID, commit.CommitHash)

					fl := &models.File{
						CommitID: commitId,
						Hash:     f.Hash.String(),
						Name:     f.Name,
						Size:     f.Size,
						// Contents: contents,
						IsBinary: isBin,
						// Lines:      ls,
						HasChanged: false}
					models.CreateFile(db, fl)
					return nil
				})

				//changes

				// for _, change := range changes {
				// 	// fmt.Println("change.From.Name: ", change.From.Name)
				// 	// fmt.Println("change.To.Name: ", change.To.Name)
				// 	// fmt.Println(change.Action())
				// 	// fmt.Println(change.Files())
				// 	// fmt.Println("------------------- start")
				// 	// fmt.Println(change.Patch())
				// 	// fmt.Println("-------------------")

				// 	// fmt.Printf("(change) file: %s - commit: %d\n", change.From.Name, commit.ID)

				// 	var fileFrom *models.File
				// 	var fileFromID uint
				// 	var fileTo *models.File
				// 	var fileToID uint
				// 	if len(change.From.Name) > 0 {
				// 		fileFrom, err = models.FindFileByEndsWithNameAndCommit(db, change.From.Name, commitId)
				// 		// fmt.Println("File From: ", fileFrom.ID, fileFrom.Name)
				// 		if err != nil {
				// 			log.Println("Cannot find file: " + change.From.Name)
				// 			log.Println(err.Error())
				// 			fmt.Println("Cannot find file: " + change.From.Name)
				// 			fmt.Println(err.Error())
				// 			fileFrom = nil
				// 			fileFromID = 0
				// 		} else {
				// 			fileFromID = fileFrom.ID
				// 		}
				// 	} else {
				// 		fileFrom = nil
				// 		fileFromID = 0
				// 	}

				// 	// FileTo
				// 	var err2 error
				// 	if len(change.To.Name) > 0 {
				// 		fileTo, err2 = models.FindFileByEndsWithNameAndCommit(db, change.To.Name, commit.ID)
				// 		if err2 != nil {
				// 			log.Println("Cannot find fileTo: " + change.To.Name)
				// 			log.Println(err.Error())
				// 			fmt.Println("Cannot find fileTo: " + change.To.Name)
				// 			fmt.Println(err.Error())
				// 			fileTo = nil
				// 			fileToID = 0
				// 		} else {
				// 			fileToID = fileTo.ID
				// 		}
				// 		// fmt.Println("File To: ", fileTo.ID, fileTo.Name)
				// 	} else {
				// 		fileTo = nil
				// 		fileToID = 0
				// 	}

				// 	act, _ := change.Action()
				// 	patch, _ := change.Patch()

				// 	ch := &models.Change{
				// 		// ChangeHash:
				// 		FileFromID: fileFromID,
				// 		FileToID:   fileToID,
				// 		Action:     act.String(),
				// 		Patch:      patch.String(),
				// 	}
				// 	models.CreateChange(db, ch)
				// }
				// var wg sync.WaitGroup
				// wg.Add(8)
				Measure(db, *measurement, repoDir, *repository, commit, cCommit)
				runtime.GC()
				// log.Println("#GoRoutines: ", runtime.NumGoroutine())
				// fmt.Println("#GoRoutines: ", runtime.NumGoroutine())
				// fmt.Println("finished Measure")
				// wg.Wait()
				// fmt.Println("finished wait group")

				//codeanalysis.Understand(cs.Name)
				// models.BarChart()
				// boxplot := charts.BoxplotExamples{}
				// boxplot.Examples()
				// }

				// cCommit = pCommit
				// cTree, _ = pCommit.Tree()

				// i = i + 1

				// }
				// return nil
				// })

				// if err != nil {
				// 	log.Println("Error iterating over commits: " + err.Error())
				// }
			}
			// deprecated
			// models.GetRandoopMetrics()
		} else {
			log.Println("Cannot get repository")
			fmt.Println("ERROR: cannot get repository!")
		}
	} else {
		log.Println("Cannot get url from .env file")
		fmt.Println("ERROR: Cannot get url from .env file")
	}
}

func isRelease(releases []string, hash string) bool {
	isRelease := false
	for _, release := range releases {
		if hash == release {
			isRelease = true
			break
		}
	}
	return isRelease
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("ERROR: cannot get local directory: %s\n", err.Error())
		log.Fatal("ERROR: cannot get local directory: ", err.Error())
	}
	// fmt.Println(dir)
	dir = strings.Replace(dir, "\\", "/", -1)
	// fmt.Println(dir)
	dir = substr(dir, 0, strings.LastIndex(dir, "/"))
	// fmt.Println(dir)
	return dir
}

func createDirs() {

	_, errd := os.Stat("coverage")
	if os.IsNotExist(errd) {
		err := os.Mkdir("coverage", 0755)
		if err != nil {
			fmt.Printf("ERROR: cannot create directory coverage: %s\n", err.Error())
			log.Fatal("ERROR: cannot create directory coverage: ", err.Error())
		}
	}

	_, errd = os.Stat("gentest")
	if os.IsNotExist(errd) {
		err := os.Mkdir("gentest", 0755)
		if err != nil {
			fmt.Printf("ERROR: cannot create directory gentest: %s\n", err.Error())
			log.Fatal("ERROR: cannot create directory gentest: ", err.Error())
		}
	}

	_, errd = os.Stat("compilation")
	if os.IsNotExist(errd) {
		err := os.Mkdir("compilation", 0755)
		if err != nil {
			fmt.Printf("ERROR: cannot create directory compilation: %s\n", err.Error())
			log.Fatal("ERROR: cannot create directory compilation: ", err.Error())
		}
	}

	_, errd = os.Stat("run")
	if os.IsNotExist(errd) {
		err := os.Mkdir("run", 0755)
		if err != nil {
			fmt.Printf("ERROR: cannot create directory run: %s\n", err.Error())
			log.Fatal("ERROR: cannot create directory run: ", err.Error())
		}
	}

	_, errd = os.Stat("profiler")
	if os.IsNotExist(errd) {
		err := os.Mkdir("profiler", 0755)
		if err != nil {
			fmt.Printf("ERROR: cannot create directory profiler: %s\n", err.Error())
			log.Fatal("ERROR: cannot create directory profiler: ", err.Error())
		}
	}

}

func ReadListFromFile(filename string) ([]string, error) {
	// reads file tcignore into memory
	// and returns a slice of testcases to be ignored
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ReadTCIgnoreMap(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tcignoreMap := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tcignoreMap[scanner.Text()] = struct{}{}
	}
	return tcignoreMap, scanner.Err()
}
