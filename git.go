package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func cloneRepository(url, directory string) (*git.Repository, error) {
	// Clone the given repository to the given directory
	removeContents(directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	// r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
	// 	URL: url,
	// })
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Repository cloned to: " + directory)
	return r, err
}

func checkout(repoName, hash string) error {
	// fmt.Println("------------------------------------------------ checkout")
	// path, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("path: " + path)

	// fmt.Printf("git --git-dir=repos"+string(os.PathSeparator)+"%v"+string(os.PathSeparator)+".git --work-tree=repos"+string(os.PathSeparator)+"%v checkout %s\n", repoName, repoName, commit)
	// cmd := exec.Command("git", "--git-dir=repos"+string(os.PathSeparator)+repoName+string(os.PathSeparator)+".git", "--work-tree=repos"+string(os.PathSeparator)+repoName, "checkout", commit)

	dir := ".." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName
	cmd := exec.Command("git", "checkout", "-f", hash)
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		fmt.Println("START dir: " + dir)
		fmt.Printf("git checkout -f %s\n", hash)
		fmt.Println("\nCannot run git checkout: ", hash, err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("WAIT dir: " + dir)
		fmt.Printf("git checkout -f %s\n", hash)
		fmt.Println("\nCannot run git checkout: ", hash, err)
	}

	// fmt.Println("checkout output: " + string(out))
	// _, err = exec.Command("git", "--git-dir=repos"+string(os.PathSeparator)+repoName+string(os.PathSeparator)+".git", "--work-tree=repos"+string(os.PathSeparator)+repoName, "checkout", commit).Output()
	// if err != nil {
	// 	fmt.Println("\nCannot run git checkout: ", err)
	// }
	return err
}

func removeContents(dir string) error {
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
