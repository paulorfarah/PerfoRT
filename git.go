package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	billy "github.com/go-git/go-billy/v5"
	git "github.com/go-git/go-git/v5"
	memory "github.com/go-git/go-git/v5/storage/memory"
)

var storer *memory.Storage
var fs billy.Filesystem

func cloneRepository(url, dstDir string) (*git.Repository, error) {
	// Clone the given repository to the given directory

	removeContents(dstDir)

	//urlSplit := strings.Split(url, "/")
	//repoName := urlSplit[4]
	//fmt.Println("git clone " + url + " .." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName)
	//log.Println("git clone " + url + " .." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName)
	//cmd := exec.Command("git", "clone", url, ".."+string(os.PathSeparator)+"repos"+string(os.PathSeparator)+repoName)
	// Clone the repository
	repo, err := git.PlainClone(dstDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	fmt.Println("Repository cloned to: " + dstDir)
	return repo, err
}

func Checkout(repoName, hash, repoDir string) error {
	dt := time.Now()
	log.Println()
	log.Println("################################################ git checkout -f " + hash + " " + dt.String())
	fmt.Println()
	fmt.Println("################################################ checkout " + hash + " " + dt.String())

	//dir := ".." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName
	dir := repoDir
	cmd := exec.Command("git", "checkout", "-f", hash)
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		log.Println("[>>ERROR]: START dir: " + dir)
		log.Printf("git checkout -f %s\n", hash)
		log.Println("\nCannot run git checkout: ", hash, err)
		fmt.Println("\nCannot run git checkout: ", hash, err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Println("[>>ERROR]: WAIT dir: " + dir)
		log.Printf("git checkout -f %s\n", hash)
		log.Println("\nCannot run git checkout: ", hash, err)
		fmt.Println("\nCannot run git checkout: ", hash, err)
	} //else {
	//fmt.Println("\ncheckout successfull: ", hash)
	//}

	// fmt.Println("checkout output: " + string(out))
	// _, err = exec.Command("git", "--git-dir=repos"+string(os.PathSeparator)+repoName+string(os.PathSeparator)+".git", "--work-tree=repos"+string(os.PathSeparator)+repoName, "checkout", commit).Output()
	// if err != nil {
	// 	fmt.Println("\nCannot run git checkout: ", err)
	// }
	return err
}

func removeContents(dir string) error {
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
