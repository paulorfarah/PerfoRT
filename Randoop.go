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
		fmt.Println("randoop fileFrom: " + fileFrom)
		//fmt.Printf("java -classpath .;%RANDOOP_JAR% randoop.main.Main gentests --classlist=myclasses.txt")

		dir := ""
		pack := ""

		paths := strings.Split(fileFrom, "/src/main/java/")
		if len(paths) > 1 {
			dir = paths[0]
			pack = paths[1]
		} else if len(paths) == 1 {
			if strings.HasPrefix(fileFrom, "src/main/java/") {
				//commons-io
				pack = strings.TrimLeft(fileFrom, "/src/main/java/")
			} else if strings.HasPrefix(fileFrom, "src/conf/") {
				pack = strings.TrimLeft(fileFrom, "/src/conf/")
			} else if strings.HasPrefix(fileFrom, "src/examples/") {
				pack = strings.TrimLeft(fileFrom, "/src/examples/")
			} else if strings.HasPrefix(fileFrom, "src/java/") {
				pack = strings.TrimLeft(fileFrom, "/src/java/")
			} else if strings.HasPrefix(fileFrom, "src/test/") {
				pack = strings.TrimLeft(fileFrom, "/src/test/")
			} else {
				fmt.Println("**************************** filefrom: " + fileFrom)
				paths = strings.Split(fileFrom, "/src/")
				dir = paths[0]
				pack = paths[1]
			}
		}

		path := strings.Split(pack, ".java")[0]
		if dir != "" {
			dir += string(os.PathSeparator)
		}
		classpath := repoDir + string(os.PathSeparator) + dir + "target" + string(os.PathSeparator) + "classes;" // + string(os.PathSeparator)
		classpath += GetMavenDependenciesClasspath(repoDir)
		className := strings.ReplaceAll(path, "/", ".")
		cmd := exec.Command("java", "-classpath", classpath+";$RANDOOP_JAR", "randoop.main.Main", "gentests", "--testclass="+className)
		fmt.Printf("java -classpath " + classpath + ";$RANDOOP_JAR randoop.main.Main gentests --testclass=" + className)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("\n[>>ERROR]: Cannot run randoop gentests: ", err)
		} else {
			fmt.Println("\n [>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>SUCCESS]: Randoop executed successully!")
			fmt.Println(out)

		}

		// //checkout current commit
		// fmt.Println(currCommit, fileTo)
		// fmt.Printf("git --git-dir=repos"+string(os.PathSeparator)+"%v"+string(os.PathSeparator)+".git --work-tree=repos"+string(os.PathSeparator)+"%v checkout %s\n", repoName, repoName, currCommit)
		// _, err = exec.Command("git", "--git-dir=repos"+string(os.PathSeparator)+repoName+string(os.PathSeparator)+".git", "--work-tree=repos"+string(os.PathSeparator)+repoName, "checkout", currCommit).Output()
		// if err != nil {
		// 	fmt.Println("\nCannot run git checkout: ", err)
		// }
		// fmt.Printf("java -classpath .;%RANDOOP_JAR% randoop.main.Main gentests --classlist=myclasses.txt")
		// }
	}
}
