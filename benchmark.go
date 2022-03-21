package main

func BenchmarkMethod(packageName, commitHash string) {
	//java -javaagent:/home/usuario/eclipse-workspace/perfrt-profiler/target/perfrt-profiler-0.0.1-SNAPSHOT.jar -jar junit-platform-console-standalone-1.8.2.jar -cp target/test-classes/:target/classes -c org.apache.commons.io.ByteOrderParserTest -m errors

	// ~/eclipse-workspace/maven-project $
	// java -javaagent:/home/usuario/eclipse-workspace/perfrt-profiler/target/
	// perfrt-profiler-0.0.1-SNAPSHOT.jar=com.github.paulorfarah.mavenproject,
	// cc2ed3975de05d3a6f9616807b44f974425e0e74 -jar
	// /home/usuario/Downloads/junit-platform-console-standalone-1.8.2.jar
	// -cp target/test-classes/:target/classes -m com.github.paulorfarah.mavenproject.AppTest#testAppHasAGreeting

	// localpath, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// 	fmt.Println("error getting current path: ", err.Error())
	// }

	// out, err := exec.Command(
	// 	"java", "-javaagent:"+localpath+"/perfrt-profiler-0.0.1-SNAPSHOT.jar="+packageName+","+commitHash,
	// 	+localpath+"/junit-platform-console-standalone-1.8.2.jar", "-cp", ".:target/test-classes/:target/classes", "-m", arg[0],
	// ).Output()
}
