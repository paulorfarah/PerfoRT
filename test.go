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

func RunJUnitTestCase(db *gorm.DB, repoDir, module, javaVer string, tc *models.TestCase, measurement models.Measurement, commit models.Commit, packageName, profiler, localpath, mavenClasspath, localClasspath string) {
	//java -javaagent:/home/usuario/go-work/src/github.com/paulorfarah/PerfoRT/PerfoRT-profiler-0.0.1-SNAPSHOT.jar=com.github.paulorfarah.mavenproject.,8df83daaa39f3e341f4057f4ae329edd425a2c7b,181 -jar /home/usuario/go-work/src/github.com/paulorfarah/PerfoRT/junit-platform-console-standalone-1.8.2.jar -cp .:target/test-classes/:target/classes -m com.github.paulorfarah.mavenproject.AppTest#testAppHasAGreeting

	// var wg sync.WaitGroup
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
						fmt.Println("saved resources... ", len(resources))
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

		//    java -javaagent:/home/usuario/go-work/src/github.com/paulorfarah/PerfoRT/PerfoRT-tracer-0.0.1-SNAPSHOT.jar=
		//    com.github.paulorfarah.mavenproject.,f07a8a880ab962572ff8fb013958afd55e4f282a,11
		//    -XX:StartFlightRecording:maxsize=200M,name=sized,dumponexit=true,filename=/home/usuario/teste5.jfr,
		// settings=/home/usuario/Downloads/PerfoRT.jfc
		// -jar /home/usuario/go-work/src/github.com/paulorfarah/PerfoRT/junit-platform-console-standalone-1.8.2.jar -cp .:target/test-classes/:target/classes -m com.github.paulorfarah.mavenproject.AppTest#testAppHasAGreeting

		jfrFilename := "/jfr/PerfoRT" + strconv.Itoa(int(run.ID)) + ".jfr"

		strJunitTC := javaVer + "/bin/java -javaagent:" + localpath + profiler + "=" + packageName + "," + commit.CommitHash + "," + strconv.Itoa(int(run.ID)) +
			" -XX:StartFlightRecording:maxsize=200M,dumponexit=true,filename=" + localpath + jfrFilename + ",settings=" + localpath + "/jfr/PerfoRT.jfc" +
			" -jar " +
			localpath + "/junit-platform-console-standalone-1.8.2.jar -cp " + localClasspath + mavenClasspath + " -m " + tc.ClassName + "#" + testName
		log.Println()
		log.Println(strJunitTC)

		// https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738
		cmd = exec.CommandContext(ctx, javaVer+"/bin/java", "-javaagent:"+localpath+profiler+"="+packageName+","+commit.CommitHash+","+strconv.Itoa(int(run.ID)),
			"-XX:StartFlightRecording:maxsize=200M,dumponexit=true,filename="+localpath+jfrFilename+",settings="+localpath+"/jfr/PerfoRT.jfc",
			"-jar", localpath+"/junit-platform-console-standalone-1.8.2.jar", "-cp", localClasspath+mavenClasspath, "-m", tc.ClassName+"#"+testName)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		cmd.Dir = repoDir
		// log.Println("path: ", path)
		// wg.Wait()

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
