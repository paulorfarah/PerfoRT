package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	_ "github.com/go-sql-driver/mysql"

	"go-repo-downloader/models"
	"go-repo-downloader/models/charts"
)

func main() {
	fmt.Println("go-repo-downloader")

	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	log.Println("starting...")

	url := "https://github.com/apache/commons-math"
	// url := "https://github.com/paulorfarah/gradle-project-example"
	// url := "https://github.com/ReactiveX/RxJava"
	// url := "https://github.com/apache/pdfbox"
	// "https://github.com/dev9com/gradle-example"
	//"https://github.com/ReactiveX/RxJava"
	//  "https://github.com/zxing/zxing"
	//  "https://github.com/junit-team/junit4"
	//  "https://github.com/paulorfarah/TestProject"
	urlSplit := strings.Split(url, "/")
	//for k, v := range urlSplit {
	//	fmt.Printf("%s -> %s\n", k, v)
	//}
	//username := urlSplit[3]
	repoName := urlSplit[4]
	repoDir := getParentDirectory() + "/repos/" + repoName
	// repoDir := "/home/farah/go-work/src/github.com/paulorfarah/repos/" + repoName
	fmt.Println("repoDir: " + repoDir)
	log.Println("repoDir: " + repoDir)

	// fmt.Println("git clone " + url)
	log.Println("git clone " + url)

	repo, err := cloneRepository(url, repoDir)

	if err == nil {
		//	ref, err := r.Head()
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//	fmt.Println(ref)

		//https://github.com/go-git/go-git/blob/master/_examples/commit/main.go
		//Author: &object.Signature{
		//	Name:  "John Doe",
		//	Email: "john@doe.org",
		//	When:  time.Now(),
		//},

		createDirs()

		db := models.GetDB()

		// platform
		platform, err := models.FindPlatformByName(db, "github")
		if err != nil {
			log.Println("Create new platform: " + "github")
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
			log.Println("create new repo: " + repoName)
			repository = &models.Repository{PlatformID: platform.ID, Name: repoName}
			models.CreateRepository(db, repository)
		}

		// measurement
		measurement := &models.Measurement{RepositoryID: repository.ID}
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

		//branches
		branchCounter := 0
		branches, _ := repo.Branches()
		for {
			branch, err := branches.Next()
			if err != nil {
				if err == io.EOF {
					//Finished branch
					break
				} else {
					log.Fatal(err)
				}
			}
			branchCounter++
			//		fmt.Println("branch -->: ", branch.Name())

			//commits
			var prevCommit *object.Commit
			prevCommit = nil
			var prevTree *object.Tree
			prevTree = nil

			//filter by dates
			// since := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)
			// until := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

			commits, err := repo.Log(&git.LogOptions{From: branch.Hash()}) //, Since: &since, Until: &until})
			if err != nil {
				log.Println("Error in git log: " + err.Error())
			}
			defer commits.Close()
			//		fmt.Println("---- commits ----")
			i := 0

			err = commits.ForEach(func(currCommit *object.Commit) error {

				if prevCommit != nil {
					// fmt.Printf("\n----- commit %v: %v -----\n", strconv.Itoa(i), currCommit.Message)
					fmt.Printf("currCommit: %s\n", currCommit.Hash)
					//fmt.Println(currCommit.Author.Email)
					//fmt.Println(currCommit.Committer)
					//fmt.Println(currCommit.Message)
					//fmt.Printf("\nfile: %v\n", cs.Name)

					currTree, err := currCommit.Tree()
					if err != nil {
						return err
					}
					if prevTree != nil {
						changes, err := currTree.Diff(prevTree)
						// _, err := currTree.Diff(prevTree)
						if err != nil {
							return err
						}
						//Author
						author, err := models.FindAccountByEmail(db, currCommit.Author.Email)
						if err != nil {
							log.Println("create new author: " + currCommit.Author.Name)
							author = &models.Account{Email: currCommit.Author.Email, Name: currCommit.Author.Name}
							models.CreateAccount(db, author)
						}
						//Committer
						committer, err := models.FindAccountByEmail(db, currCommit.Committer.Email)
						if err != nil {
							log.Println("create new committer: " + currCommit.Committer.Name)
							committer = &models.Account{Email: currCommit.Committer.Email, Name: currCommit.Committer.Name}
							models.CreateAccount(db, committer)
						}

						//Commit
						commit, err := models.FindCommitByHash(db, currCommit.Hash.String())
						if err != nil {
							log.Println("create new commit: " + currCommit.Hash.String())
							fmt.Println("create new commit: " + currCommit.Hash.String())
							// parent, errj := json.Marshal(currCommit.ParentHashes)
							// if errj != nil {
							// 	log.Println("Error Marshalling parent hashes: " + errj.Error())
							// }
							commit = &models.Commit{CommitHash: currCommit.Hash.String(),
								PreviousCommitHash: prevCommit.Hash.String(),
								RepositoryID:       repository.ID,
								TreeHash:           currCommit.TreeHash.String(),
								// ParentHashes:       parent,
								AuthorID:      author.ID,
								AuthorDate:    currCommit.Author.When,
								CommitterID:   committer.ID,
								CommitterDate: currCommit.Committer.When,
								Subject:       currCommit.Message,
								Branch:        branch.Name().String()}
							_, err = models.CreateCommit(db, commit)
							if err != nil {
								fmt.Printf("Error creating new commit %s\n", err.Error())
							}

						}

						//files
						currTree.Files().ForEach(func(f *object.File) error {
							contents := ""
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
								CommitID: commit.ID,
								Hash:     f.Hash.String(),
								Name:     f.Name,
								Size:     f.Size,
								Contents: contents,
								IsBinary: isBin,
								// Lines:      ls,
								HasChanged: false}
							models.CreateFile(db, fl)
							return nil
						})

						//changes

						// changes.ForEach(func(change *object.Change) error {
						// 	fmt.Println(change.From.Name)
						// })
						for _, change := range changes {
							// fmt.Println(change.To.Name)
							// 		// fmt.Println(change.Action())
							// 		// fmt.Println(change.Files())
							// 		// fmt.Println("------------------- start")
							// 		// fmt.Println(change.Patch())

							// fmt.Printf("(change) file: %s - commit: %d\n", change.From.Name, commit.ID)
							fileFrom, err := models.FindFileByEndsWithNameAndCommit(db, change.From.Name, commit.ID)
							var fileTo *models.File
							if err != nil {
								log.Println("Cannot find file: " + change.From.Name)
								log.Println(err.Error())
								fmt.Println("Cannot find file: " + change.From.Name)
								fmt.Println(err.Error())
							} else {
								if change.From.Name == change.To.Name {
									fileTo = fileFrom
								} else {
									var err2 error
									fileTo, err2 = models.FindFileByEndsWithNameAndCommit(db, change.To.Name, commit.ID)
									if err2 != nil {
										log.Println("Cannot find file: " + change.From.Name)
										log.Println(err2.Error())
										fmt.Println("Cannot find file: " + change.From.Name)
										fmt.Println(err2.Error())
									}
								}
								act, _ := change.Action()
								patch, _ := change.Patch()

								ch := &models.Change{
									// ChangeHash:
									FileFromID: fileFrom.ID,
									FileToID:   fileTo.ID,
									Action:     act.String(),
									Patch:      patch.String(),
								}
								models.CreateChange(db, ch)
							}
						}
						// var wg sync.WaitGroup
						// wg.Add(8)
						Measure(db, *measurement, repoDir, *repository, commit.ID, currCommit)
						fmt.Println("finished Measure")
						// wg.Wait()
						fmt.Println("finished wait group")

						//codeanalysis.Understand(cs.Name)
						// models.BarChart()
						boxplot := charts.BoxplotExamples{}
						boxplot.Examples()
					}
				}
				prevCommit = currCommit
				prevTree, _ = currCommit.Tree()

				i = i + 1

				return nil
			})
			if err != nil {
				log.Println("Error iterating over commits: " + err.Error())
			}
		}
		// deprecated
		// models.GetRandoopMetrics()
	} else {
		log.Println("Cannot get repository")
	}
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
		log.Fatal(err)
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
			log.Fatal(err)
		}
	}

	_, errd = os.Stat("gentest")
	if os.IsNotExist(errd) {
		err := os.Mkdir("gentest", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, errd = os.Stat("compilation")
	if os.IsNotExist(errd) {
		err := os.Mkdir("compilation", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, errd = os.Stat("run")
	if os.IsNotExist(errd) {
		err := os.Mkdir("run", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

}
