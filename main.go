package main

import (
	"fmt"
	"time"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	_ "github.com/go-sql-driver/mysql"

	"go-repo-downloader/models"
)


func main(){
	fmt.Println("go-repo-downloader")
	url := "http://github.com/paulorfarah/refactoring-python-code"
	urlSplit := strings.Split(url, "/")
	//for k, v := range urlSplit {
	//	fmt.Printf("%s -> %s\n", k, v)
	//}
	//username := urlSplit[3]
	repoName := urlSplit[4]

	fmt.Println("git clone " + url)
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("git log")
	ref, err := r.Head()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ref)

	since := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)
	until := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---- commits ----")
	i := 0

	//https://github.com/go-git/go-git/blob/master/_examples/commit/main.go
	//Author: &object.Signature{
	//	Name:  "John Doe",
	//	Email: "john@doe.org",
	//	When:  time.Now(),
	//},

	db := models.GetDB()
	fmt.Println(db)
	platform, err := models.FindPlatformByName(db, "github") 
	if err != nil {
		fmt.Println("Create new record...")
		//fmt.Println(err)
		platform = &models.Platform{Name:"github"}
		models.CreatePlatform(db, platform)
	}
	fmt.Println(platform)

	//repository
	repository, err := models.FindRepositoryByName(db, repoName)
	if err != nil{
		fmt.Println("create new repo")
		fmt.Println(err)
		repository = &models.Repository{PlatformFK: platform.ID, Name: repoName}
		models.CreateRepository(db, repository)
	}
	fmt.Println(repository)

	//commit
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Printf("----- commIt %v -----\n", strconv.Itoa(i))
		fmt.Println(c.Hash)
		fmt.Println(c.Author.Email)
		fmt.Println(c.Committer)
		fmt.Println(c.Message)
		//Author
		author, err := models.FindAccountByEmail(db, c.Author.Email)
		if err != nil {
			fmt.Println("create new author...")
			fmt.Println(err)
			author = &models.Account{Email:c.Author.Email, Name: c.Author.Name}
			models.CreateAccount(db, author)
		}
		//Committer
		committer, err := models.FindAccountByEmail(db, c.Committer.Email)
		if err != nil {
			fmt.Println("create new committer...")
			fmt.Println(err)
			committer = &models.Account{Email:c.Committer.Email, Name: c.Committer.Name}
			models.CreateAccount(db, committer)
		}


		commit, err := models.FindCommitByHash(db, c.Hash.String())
		if err != nil {
			fmt.Println("create new commit")
			fmt.Println(err)
			parent, errj := json.Marshal(c.ParentHashes)
			if errj != nil {
				fmt.Println(errj)
			}
			commit = &models.Commit{CommitHash: c.Hash.String(),
						RepositoryFK: repository.ID,
						TreeHash:c.TreeHash.String(),
						ParentHashes:parent,
						Author:author.ID,
						AuthorDate:c.Author.When,
						Committer:committer.ID,
						CommitterDate:c.Committer.When,
						Subject:c.Message}
			models.CreateCommit(db, commit)
		}
		i = i + 1

		return nil
	})
	if err != nil{
		fmt.Println(err)
	}
}
