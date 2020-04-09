package main

import (
	"fmt"
	"time"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func main(){
	fmt.Println("go-repo-downloader")
	fmt.Println("git clone https://github.com/paulorfarah/refactoring-python-code")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/paulorfarah/refactoring-python-code",
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
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)

		return nil
	})
	if err != nil{
		fmt.Println(err)
	}
}
