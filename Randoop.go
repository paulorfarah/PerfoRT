package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecuteRandoop(repoName, prevCommit, fileFrom, currCommit, fileTo string) {
	//checkout previous commit
	fmt.Println(prevCommit, fileFrom)
	fmt.Println("checkout")
	fmt.Printf("git --git-dir=repos"+string(os.PathSeparator)+"%v"+string(os.PathSeparator)+".git --work-tree=repos"+string(os.PathSeparator)+"%v checkout %s\n", repoName, repoName, prevCommit)
	_, err := exec.Command("git", "--git-dir=repos"+string(os.PathSeparator)+repoName+string(os.PathSeparator)+".git", "--work-tree=repos"+string(os.PathSeparator)+repoName, "checkout", prevCommit).Output()
	if err != nil {
		fmt.Println("\nCannot run git checkout: ", err)
	}
	fmt.Println("randoop")
	//fmt.Printf("java -classpath .;%RANDOOP_JAR% randoop.main.Main gentests --classlist=myclasses.txt")

	paths := strings.Split(fileFrom, "/src/main/java/")
	path := strings.Split(paths[1], ".java")[0]

	classpath := "D:" + string(os.PathSeparator) + "eclipse-workspace2" + string(os.PathSeparator) + "randoop" + string(os.PathSeparator) + "pdfbox" + string(os.PathSeparator) + paths[0] + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes" + string(os.PathSeparator)
	classpath += GetMavenDependenciesClasspath("d:\\eclipse-workspace2\\randoop\\pdfbox")
	className := strings.ReplaceAll(path, "/", ".")
	// fmt.Printf("java -classpath " + classpath + ";D:\\Download\\randoop-4.2.5\\randoop-all-4.2.5.jar randoop.main.Main gentests --testclass=" + className)
	_, err = exec.Command("java", "-classpath", classpath+";D:\\Download\\randoop-4.2.5\\randoop-all-4.2.5.jar", "randoop.main.Main", "gentests", "--testclass="+className).Output()
	if err != nil {
		fmt.Println("\nCannot run git checkout: ", err)
	}

	// //checkout current commit
	// fmt.Println(currCommit, fileTo)
	// fmt.Printf("git --git-dir=repos"+string(os.PathSeparator)+"%v"+string(os.PathSeparator)+".git --work-tree=repos"+string(os.PathSeparator)+"%v checkout %s\n", repoName, repoName, currCommit)
	// _, err = exec.Command("git", "--git-dir=repos"+string(os.PathSeparator)+repoName+string(os.PathSeparator)+".git", "--work-tree=repos"+string(os.PathSeparator)+repoName, "checkout", currCommit).Output()
	// if err != nil {
	// 	fmt.Println("\nCannot run git checkout: ", err)
	// }
	// fmt.Printf("java -classpath .;%RANDOOP_JAR% randoop.main.Main gentests --classlist=myclasses.txt")
}
