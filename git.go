package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	billy "github.com/go-git/go-billy/v5"
	memfs "github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	memory "github.com/go-git/go-git/v5/storage/memory"
)

var storer *memory.Storage
var fs billy.Filesystem

func cloneRepository(url, directory string) (*git.Repository, error) {
	// Clone the given repository to the given directory

	//removeContents(directory)

	urlSplit := strings.Split(url, "/")
	repoName := urlSplit[4]
	fmt.Println("git clone " + url + " .." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName)
	cmd := exec.Command("git", "clone", url, ".."+string(os.PathSeparator)+"repos"+string(os.PathSeparator)+repoName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("\n[>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ERROR]: Cannot clone repository: ", err, out)
		// return nil, err
	} else {
		fmt.Println("\n [>>SUCCESS]: repository cloned successully!")
		// fmt.Println(out)

	}

	// r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
	// 	URL: url,
	// })

	storer = memory.NewStorage()
	fs = memfs.New()
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println("Repository cloned to: " + directory)
	return r, err
}

func Checkout(repoName, hash string) error {
	// fmt.Println("################################################ checkout " + hash)
	// log.Println("################################################ checkout " + hash)

	dir := ".." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName
	cmd := exec.Command("git", "checkout", "-f", hash)
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		log.Println("[>>ERROR]: START dir: " + dir)
		log.Printf("git checkout -f %s\n", hash)
		log.Println("\nCannot run git checkout: ", hash, err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Println("[>>ERROR]: WAIT dir: " + dir)
		log.Printf("git checkout -f %s\n", hash)
		log.Println("\nCannot run git checkout: ", hash, err)
	}

	// fmt.Println("checkout output: " + string(out))
	// _, err = exec.Command("git", "--git-dir=repos"+string(os.PathSeparator)+repoName+string(os.PathSeparator)+".git", "--work-tree=repos"+string(os.PathSeparator)+repoName, "checkout", commit).Output()
	// if err != nil {
	// 	fmt.Println("\nCannot run git checkout: ", err)
	// }
	return err
}

func removeContents(dir string) error {
	fmt.Println("deleting directory: " + dir)
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
