package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	//	fdiff "github.com/go-git/go-git/v5/plumbing/format/diff"

	_ "github.com/go-sql-driver/mysql"

	"go-repo-downloader/models"
	//"go-repo-downloader/codeanalysis"

	"github.com/waigani/diffparser"
)

// type FileStat struct {
// 	Name     string
// 	Addition int
// 	Deletion int
// }

func main2() {
	fmt.Println("go-repo-downloader")

	url := "https://github.com/apache/commons-io" //"https://github.com/eclipse/jgit" (cant compile) //"https://github.com/apache/pdfbox" (svnexit
	//  "https://github.com/paulorfarah/TestProject"
	urlSplit := strings.Split(url, "/")
	//for k, v := range urlSplit {
	//	fmt.Printf("%s -> %s\n", k, v)
	//}
	//username := urlSplit[3]
	repoName := urlSplit[4]
	// repoDir := ".." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName
	repoDir := getParentDirectory() + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName
	fmt.Println("repoDir: " + repoDir)

	// fmt.Println("git clone " + url)
	r, err := cloneRepository(url, repoDir)

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

		db := models.GetDB()
		platform, err := models.FindPlatformByName(db, "github")
		if err != nil {
			fmt.Println("Create new record...")
			//fmt.Println(err)
			platform = &models.Platform{Name: "github"}
			models.CreatePlatform(db, platform)
		}

		//search representative repositories

		//save repository in db
		repository, err := models.FindRepositoryByName(db, repoName)
		if err != nil {
			fmt.Println("create new repo")
			fmt.Println(err)
			repository = &models.Repository{PlatformID: platform.ID, Name: repoName}
			models.CreateRepository(db, repository)
		}

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
		branches, _ := r.Branches()
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

			commits, err := r.Log(&git.LogOptions{From: branch.Hash()}) //, Since: &since, Until: &until})
			if err != nil {
				fmt.Println(err)
			}
			defer commits.Close()
			//		fmt.Println("---- commits ----")
			i := 0

			err = commits.ForEach(func(currCommit *object.Commit) error {

				if prevCommit != nil {
					fmt.Printf("\n----- commit %v: %v -----\n", strconv.Itoa(i), currCommit.Message)
					//fmt.Println(currCommit.Hash)
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
						if err != nil {
							return err
						}
						//Author
						author, err := models.FindAccountByEmail(db, currCommit.Author.Email)
						if err != nil {
							fmt.Println("create new author...")
							fmt.Println(err)
							author = &models.Account{Email: currCommit.Author.Email, Name: currCommit.Author.Name}
							models.CreateAccount(db, author)
						}
						//Committer
						committer, err := models.FindAccountByEmail(db, currCommit.Committer.Email)
						if err != nil {
							fmt.Println("create new committer...")
							fmt.Println(err)
							committer = &models.Account{Email: currCommit.Committer.Email, Name: currCommit.Committer.Name}
							models.CreateAccount(db, committer)
						}

						//Commit
						commit, err := models.FindCommitByHash(db, currCommit.Hash.String())
						if err != nil {
							fmt.Println("create new commit")
							fmt.Println(err)
							parent, errj := json.Marshal(currCommit.ParentHashes)
							if errj != nil {
								fmt.Println(errj)
							}
							commit = &models.Commit{CommitHash: currCommit.Hash.String(),
								PreviousCommitHash: prevCommit.Hash.String(),
								RepositoryID:       repository.ID,
								TreeHash:           currCommit.TreeHash.String(),
								ParentHashes:       parent,
								Author:             author.ID,
								AuthorDate:         currCommit.Author.When,
								Committer:          committer.ID,
								CommitterDate:      currCommit.Committer.When,
								Subject:            currCommit.Message,
								Branch:             fmt.Sprintf("%s", branch.Name())}
							models.CreateCommit(db, commit)
						}

						// Changes
						for _, change := range changes {
							// fmt.Println(change.From.Name)
							// fmt.Println(change.To.Name)
							// fmt.Println(change.Action())
							// fmt.Println(change.Files())
							// fmt.Println("------------------- start")
							// fmt.Println(change.Patch())

							patch, _ := change.Patch()
							diff, _ := diffparser.Parse(patch.String())

							//files
							count := 0
							for _, file := range diff.Files {
								// fmt.Println("************************** file: ", file)

								sc := fmt.Sprintf("%d", count)

								fNew, _ := os.Create("results/" + currCommit.Hash.String() + "f" + sc + "_new.java")
								defer fNew.Close()

								fOld, _ := os.Create("results/" + currCommit.Hash.String() + "f" + sc + "_old.java")
								defer fOld.Close()

								// //hunks
								for _, hunk := range file.Hunks {
									for _, l := range hunk.NewRange.Lines {
										fNew.WriteString(l.Content + "\n")
									}
									for _, l := range hunk.OrigRange.Lines {
										fOld.WriteString(l.Content + "\n")
									}
								}
								count++

							}

							hasher := sha1.New()
							patch, err := change.Patch()
							if err != nil {
								return err
							}
							hasher.Write([]byte(patch.String()))
							changeSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
							//fmt.Println(changeSha)
							//	id := fmt.Sprintf("%s",currCommit.ID)
							//	fmt.Printf("*************  %s\n", id)
							changeObj, err := models.FindChangeByHash(db, changeSha, commit.ID)
							if err != nil {
								fmt.Println("new change")
								fmt.Println(err)
								action, err := change.Action()
								if err != nil {
									return err
								}
								changeObj = &models.Change{CommitID: commit.ID, ChangeHash: changeSha, FileFrom: change.From.Name, FileTo: change.To.Name, Action: action.String(), Patch: patch.String()}
								models.CreateChange(db, changeObj)

								//call randoop
								fmt.Println(change.From.Name)
								if action.String() == "Modify" &&
									strings.Contains(change.From.Name, ".java") &&
									strings.Contains(change.To.Name, ".java") &&
									!strings.HasPrefix(change.From.Name, "src/test/") &&
									!strings.HasPrefix(change.From.Name, "src/test/") {
									CollectRandoopMetrics(repoDir, repoName, prevCommit.Hash.String(), change.From.Name, currCommit.Hash.String(), change.To.Name, changeObj.ID)
								}
							} else {
								fmt.Println("change already exists in database...")
							}
						}
						//codeanalysis.Understand(cs.Name)
					}
				}
				prevCommit = currCommit
				prevTree, _ = currCommit.Tree()

				i = i + 1

				return nil
			})
			if err != nil {
				fmt.Println(err)
			}

		}
	} else {
		fmt.Println("Cannot get repository")
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
	dir = strings.Replace(dir, "\\", "/", -1)

	return substr(dir, 0, strings.LastIndex(dir, "/"))
}
