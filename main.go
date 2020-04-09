package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
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
}
