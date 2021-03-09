package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecuteRandoop(repoDir, repoName, prevCommit, fileFrom, currCommit, fileTo string) {
	//checkout previous commit
	fmt.Println("prevCommit: " + prevCommit + "fileFrom: " + fileFrom)
	err := checkout(repoName, prevCommit)

	if err == nil {
		fmt.Println("randoop")
		//fmt.Printf("java -classpath .;%RANDOOP_JAR% randoop.main.Main gentests --classlist=myclasses.txt")

		paths := strings.Split(fileFrom, "/src/main/java/")
		if len(paths) == 1 {
			paths = strings.Split(fileFrom, "/src/")
		}
		if len(paths) > 1 {
			fmt.Println(paths)
			path := strings.Split(paths[1], ".java")[0]

			// classpath := "D:" + string(os.PathSeparator) + "eclipse-workspace2" + string(os.PathSeparator) + "randoop" + string(os.PathSeparator) + "pdfbox" + string(os.PathSeparator) + paths[0] + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes" + string(os.PathSeparator)
			classpath := repoDir + string(os.PathSeparator) + paths[0] + string(os.PathSeparator) + "target" + string(os.PathSeparator) + "classes;" // + string(os.PathSeparator)
			fmt.Println("classpath: " + classpath)
			classpath += GetMavenDependenciesClasspath(repoDir)
			className := strings.ReplaceAll(path, "/", ".")

			cmd := exec.Command("java", "-classpath", classpath+";D:\\Download\\randoop-4.2.5\\randoop-all-4.2.5.jar", "randoop.main.Main", "gentests", "--testclass="+className)
			// dir := ".." + string(os.PathSeparator) + "repos" + string(os.PathSeparator) + repoName
			// cmd.Dir = dir
			err := cmd.Start()
			if err != nil {
				path, _ := os.Getwd()
				fmt.Println("currentdir: " + path)
				fmt.Println("START dir: " + dir)
				fmt.Printf("java -classpath " + classpath + ";D:\\Download\\randoop-4.2.5\\randoop-all-4.2.5.jar randoop.main.Main gentests --testclass=" + className)
				fmt.Println("\nCannot run randoop gentests: ", err)
			}
			err = cmd.Wait()
			if err != nil {
				fmt.Println("START dir:" + dir)
				fmt.Printf("java -classpath " + classpath + ";D:\\Download\\randoop-4.2.5\\randoop-all-4.2.5.jar randoop.main.Main gentests --testclass=" + className)
				fmt.Println("\nCannot run randoop gentests: ", err)
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
	}
}
