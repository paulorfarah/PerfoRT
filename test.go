package main

import (
	"PerfoRT/models"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gorm.io/gorm"
)

func RunJUnitTestCase(db *gorm.DB, repoDir, module, javaHome string, tc *models.TestCase, measurement models.Measurement, commit models.Commit, packageName, profiler, localpath, mavenClasspath, localClasspath string, tracerClasspath string) {
	// var wg sync.WaitGroup

	// fmt.Println("######## " + tc.Name)
	testName := tc.Name[strings.LastIndex(tc.Name, ".")+1:]

	var cmd *exec.Cmd

	// log.Println("Number of runs: ", measurement.Runs)
	finish := make(chan bool)
	for runNumber := 0; runNumber < measurement.Runs; runNumber++ {
		// log.Println("#Run: ", runNumber)
		run := &models.Run{
			MeasurementID: measurement.ID,
			TestCaseID:    tc.ID,
			Type:          "junit",
			Number:        runNumber,
		}
		models.CreateRun(db, run)

		// prepare goroutine
		start := make(chan int)
		stop := make(chan bool)
		ctx, cancel := context.WithTimeout(context.Background(), measurement.TestcaseTimeout*time.Second)
		defer cancel()
		// wg.Add(1)
		go func() {
			active := false
			var pid int
			resources := []models.Resource{}
			for {
				select {
				case pid = <-start:
					// log.Println("+++ start monitoring... ", time.Now())
					active = true
					resources = []models.Resource{}
					// wg.Done()
				case <-stop:
					log.Println("### stop monitoring... ", time.Now())
					active = false
					//save
					// log.Println("saving resources: ", len(resources))
					db.CreateInBatches(resources, 1000)
					log.Println("saved resources... ", len(resources))
					SaveJFRMetrics(db, run.ID, tc.ID)
					// log.Println("saved jvm...")
				case <-ctx.Done():
					if active {
						// fmt.Println("!!! time out monitoring... ", time.Now())
						log.Println("!!! time out monitoring... ", time.Now())
						active = false
						// errKill := cmd.Process.Kill()
						// if errKill != nil {
						// 	fmt.Println("Error killing process: ", errKill)
						// 	log.Println("Error killing process: ", errKill)
						// }
						// fmt.Println("Testcase monitoring timed out: ", tc.ClassName, "#", tc.Name)
						// log.Println("Testcase monitoring timed out", tc.ClassName, "#", tc.Name)

						db.CreateInBatches(resources, 1000)
						// fmt.Println("saved resources... ", len(resources))
						SaveJFRMetrics(db, run.ID, tc.ID)
						// fmt.Println("saved jvm...")
					}

				case <-finish:
					// fmt.Println(">>> finished monitoring... ", time.Now())
					return
				default:
					if active {
						resource, err := MonitorProcess(pid, run.ID)
						if err == nil {
							resources = append(resources, resource)
						}
						time.Sleep(measurement.MonitoringTime)
					}
				}
			}
		}()

		jfrFilename := "/jfr/PerfoRT" + strconv.Itoa(int(run.ID)) + ".jfr"

		// // strJunitTC := javaHome + "/bin/java -javaagent:" + localpath + profiler + "=" + packageName + "," + commit.CommitHash + "," + strconv.Itoa(int(run.ID)) +
		// // 	" -XX:StartFlightRecording:maxsize=200M,dumponexit=true,filename=" + localpath + jfrFilename + ",settings=" + localpath + "/jfr/PerfoRT.jfc" +
		// // 	" -jar " +
		// // 	localpath + "/junit-platform-console-standalone-1.8.2.jar -cp " + localClasspath + mavenClasspath + " -m " + tc.ClassName + "#" + testName

		// strJunitTC := javaHome + "/bin/java -javaagent:" + localpath + "/aspectjweaver-1.9.24.jar -Dpackage.name=" + packageName + " -Dhash=" + commit.CommitHash + " -Drun_id=" + strconv.Itoa(int(run.ID)) +
		// 	" -XX:StartFlightRecording:maxsize=200M,dumponexit=true,filename=" + localpath + jfrFilename + ",settings=" + localpath + "/jfr/PerfoRT.jfc" +
		// 	" -cp .:/home/farah/eclipse-workspace/method-timing-agent-maven/target/method-timing-agent-1.0-SNAPSHOT.jar" +
		// 	" -jar " +
		// 	localpath + "/junit-platform-console-standalone-1.8.2.jar -cp " + localClasspath + mavenClasspath + ":/home/farah/eclipse-workspace/method-timing-agent-maven/target/method-timing-agent-1.0-SNAPSHOT.jar -m " + tc.ClassName + "#" + testName

		// log.Println()
		// log.Println(strJunitTC)

		// // https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738
		// // cmd = exec.CommandContext(ctx, javaHome+"/bin/java", "-javaagent:"+localpath+profiler+"="+packageName+","+commit.CommitHash+","+strconv.Itoa(int(run.ID)),
		// // 	"-XX:StartFlightRecording:maxsize=200M,dumponexit=true,filename="+localpath+jfrFilename+",settings="+localpath+"/jfr/PerfoRT.jfc",
		// // 	"-jar", localpath+"/junit-platform-console-standalone-1.8.2.jar", "-cp", localClasspath+mavenClasspath, "-m", tc.ClassName+"#"+testName)

		// cmd = exec.CommandContext(ctx, javaHome+"/bin/java", "-javaagent:"+localpath+"/aspectjweaver-1.9.24.jar", "-Dpackage.name="+packageName, "-Dhash="+commit.CommitHash, "-Drun_id="+strconv.Itoa(int(run.ID)),
		// 	"-XX:StartFlightRecording:maxsize=200M,dumponexit=true,filename="+localpath+jfrFilename+",settings="+localpath+"/jfr/PerfoRT.jfc",
		// 	"-jar", localpath+"/junit-platform-console-standalone-1.8.2.jar", "-cp", localClasspath+mavenClasspath+":/home/farah/eclipse-workspace/method-timing-agent-maven/target/method-timing-agent-1.0-SNAPSHOT.jar", "-m", tc.ClassName+"#"+testName)

		// var outb, errb bytes.Buffer
		// cmd.Stdout = &outb
		// cmd.Stderr = &errb
		// cmd.Dir = repoDir
		// // log.Println("path: ", path)
		// // wg.Wait()

		//
		// Config vars
		javaAgent := localpath + "/aspectjweaver-1.9.24.jar"
		junitPlatform := localpath + "/junit-platform-console-standalone-1.12.2.jar"
		jfrSettings := localpath + "/jfr/PerfoRT.jfc"
		jfrFilename = localpath + jfrFilename

		runID := strconv.Itoa(int(run.ID))

		// Build classpath
		classpath := fmt.Sprintf(".:"+repoDir+"/target/test-classes:"+repoDir+"/target/classes:%s:%s:%s:%s",
			profiler, junitPlatform, mavenClasspath, tracerClasspath)

		// Build the full Java command
		cmdArgs := []string{
			"-javaagent:" + javaAgent,
			"-Dpackage.name=" + packageName,
			"-Dhash=" + commit.CommitHash,
			"-Drun_id=" + runID,
			"-XX:StartFlightRecording=maxsize=200M,filename=" + jfrFilename + ",settings=" + jfrSettings,
			"-cp", classpath,
			"org.junit.platform.console.ConsoleLauncher",
			"--select-method=" + tc.ClassName + "#" + testName,
		}

		// Print command (for debugging)
		log.Println("Running:", javaHome+"/bin/java", cmdArgs)

		//////////////////////////////
		// Build the full Java command string
		cmdStr := fmt.Sprintf(`%s/bin/java -javaagent:%s -Dpackage.name=%s -Dhash=%s -Drun_id=%s -XX:StartFlightRecording=maxsize=200M,filename=%s,settings=%s -cp "%s" org.junit.platform.console.ConsoleLauncher --select-method=%s#%s`,
			javaHome, javaAgent, packageName, commit.CommitHash, runID, jfrFilename, jfrSettings, classpath, tc.ClassName, testName)

		// Path do script
		scriptPath := repoDir + "/" + ".sh"

		// Criar script .sh
		scriptContent := fmt.Sprintf("#!/bin/bash\n\n%s\n", cmdStr)
		err2 := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
		if err2 != nil {
			log.Fatalf("Erro ao criar script: %v", err2)
		}

		// Executar script via bash
		var outb, errb bytes.Buffer
		cmd2 := exec.CommandContext(ctx, "bash", scriptPath)
		cmd2.Stdout = &outb
		cmd2.Stderr = &errb
		cmd2.Dir = repoDir

		err2 = cmd2.Run()
		if err2 != nil {
			log.Printf("Erro ao executar script: %v", err2)
		}

		// Exibir outputs
		log.Println("STDOUT:", outb.String())
		log.Println("STDERR:", errb.String())
		////////////////////////////////////////////

		// Run command
		// cmd := exec.Command(javaHome+"/bin/java", cmdArgs...)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr

		// Set working directory if needed
		// cmd.Dir = "/home/farah/eclipse-workspace/method-timing-agent-maven"
		// cmd.Dir = repoDir

		// var outb, errb bytes.Buffer
		cmd = exec.CommandContext(ctx, javaHome+"/bin/java", cmdArgs...)
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		cmd.Dir = repoDir

		err := cmd.Start()
		if cmd.Process != nil {
			start <- cmd.Process.Pid
			if err != nil {
				log.Println("Error setting JUnit process id: ", err)
				// log.Fatal(err)
			}

			err = cmd.Wait()
			// fmt.Println(ctx.Err())
			stop <- true
			if ctx.Err() == context.DeadlineExceeded {
				// log.Println("DeadlineExceeded...")
				// fmt.Println("DeadlineExceeded...")
				models.SetTestCaseError(db, tc)
				cancel()
				// return
			}

			if err != nil {
				process, err := os.FindProcess(int(cmd.Process.Pid))
				if err != nil {
					fmt.Printf("Failed to find process: %s\n", err)
				} else {
					errPid := process.Signal(syscall.Signal(0))
					// fmt.Printf("process.Signal on pid %d returned: %v\n", cmd.Process.Pid, errPid)
					resPid := fmt.Sprintf("%v", errPid)
					if resPid != "os: process already finished" {
						fmt.Printf("junit test failed with %s\n", err.Error())
						log.Printf("Command finished with error: %s", err.Error())
					}
				}
			}
		}
	}
	finish <- true
	db.Model(&models.Method{}).Where("Finished = ?", false).Update("Finished", true)
	runtime.GC()
}
