package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go-repo-downloader/models"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func ParseProfilingClock(db *gorm.DB, commit models.Commit, testcase models.TestCase, output string) {
	output += "-wall.txt"
	e := os.Rename("profiler/wall.txt", output)
	if e != nil {
		fmt.Println("Error renaming profile file: ", e.Error())
		fmt.Println("FATAL: Check if MAVEN_OPTS is exported!")
		log.Fatal("Check if MAVEN_OPTS is exported!")
	}

	lines := 0
	f, err := os.Open(output)
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer f.Close()

	foundElement := false
	finished := false
	var stack []models.Method
	var duration string
	var ownDur time.Duration
	totalCalls := 0
	calls := 0
	firstInStack := true
	var prevMethod models.Method = models.Method{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines++
		line := scanner.Bytes()
		if lines == 1 {
			continue
		} else if lines == 2 {
			values := strings.Split(string(line), ":")
			totalCalls, e = strconv.Atoi(strings.Trim(values[1], " "))
			if e != nil {
				fmt.Println("Error converting total samples to integer: ", e.Error())
			}
		} else {
			if len(line) > 2 {
				if bytes.Equal(line[:3], []byte("---")) || bytes.Equal(line[:3], []byte("  -")) {
					if foundElement {
						//save stack
						fmt.Println("*********************************************************************** stack")
						fmt.Printf("line: %d ownDur: %v", lines, ownDur)
						lastMethod := stack[len(stack)-1]
						e := false
						if lastMethod.Name != testcase.Name {
							// check if name of last method is exactly the same name of testcase, they should be...
							e = true
						}
						// for _, m := range stack {
						for i := len(stack) - 2; i >= 0; i-- {
							m := stack[i]
							m.Error = e
							fmt.Printf("Name: %s, Duration: %d\n", m.Name, m.OwnDuration)
							if m.FileID >= 0 {
								method, errM := models.FindMethodByEndsWithNameAndFileAndTestcase(db, m.Name, m.FileID, testcase.ID)
								if errM != nil {
									log.Printf("Create new method: %s (%s)\n"+m.Name, errM.Error())
									if i == 0 {
										m.OwnDuration = ownDur
										m.TotalCalls = totalCalls
										m.OwnCalls = calls
										m.CallsPercent = float64(calls / totalCalls)
									}
									if !firstInStack {
										m.CallerID = &prevMethod.ID
										// err = models.SaveMethod(db, &prevMethod)
										// if err != nil {
										// 	fmt.Println("Error saving previous method: ", err.Error())
										// }
									}
									m.ID, err = models.CreateMethod(db, &m)
									if err != nil {
										fmt.Println("Error creating method: ", err.Error())
									}

									prevMethod = m
								} else {
									//update durations
									if i == 0 {
										method.OwnDuration += ownDur
										method.OwnCalls += calls
										method.TotalCalls = totalCalls
										method.CallsPercent = float64(float64(calls) / float64(totalCalls))
									}
									if !firstInStack {
										method.CallerID = &prevMethod.ID
										// err = models.SaveMethod(db, &prevMethod)
										// if err != nil {
										// 	fmt.Println("Error saving previous method: ", err.Error())
										// }
									}
									err = models.SaveMethod(db, method)
									if err != nil {
										fmt.Println("Error saving method: ", err.Error())
									}

									prevMethod = *method
								}
								if firstInStack {
									firstInStack = false
								}
							}
						}
					}
					if !bytes.Equal(line[:3], []byte("  -")) {
						values := strings.Split(string(line), " ")
						duration = values[1] + values[2]
						var errPD error
						ownDur, errPD = time.ParseDuration(duration)
						if errPD != nil {
							fmt.Println("Error parsing duration of method: ", errPD.Error())
						}
						calls, errPD = strconv.Atoi(values[4])
						if errPD != nil {
							fmt.Println("Error parsing calls of method: ", errPD.Error())
						}

						firstInStack = true
						foundElement = false
						stack = nil
						finished = false
						prevMethod = models.Method{}
					}

					// fmt.Println("duration: ", values[1])
					// fmt.Println("%: ", values[1])
					// fmt.Println("samples: ", values[1])
					// fmt.Println(values)
				} else {
					if !finished {
						call := strings.Split(string(line), "] ")
						// fmt.Println("call: ", call)
						if len(call) > 1 {
							elements := strings.Split(call[1], ".")
							// fmt.Println(elements)
							if len(elements) > 1 {
								element := elements[0]
								method := &models.Method{Name: elements[len(elements)-1], TestCaseID: testcase.ID}
								fmt.Printf(">>>>> method: %s File: %d\n", method.Name, method.FileID)
								// fmt.Println(elements)
								// fmt.Println(len(elements))
								for i := 1; i < len(elements)-1; i++ {
									// fmt.Println("i:", i)
									element += "." + elements[i]
								}

								// search package and class name
								file, err := models.FindFileByEndsWithNameAndCommit(db, element+".java", commit.ID)
								if err != nil {
									// fmt.Println("Error searching for profiled class: ", err.Error())
									// files = append(files, -1)
									if foundElement {
										finished = true
									}

								} else {
									// lastFound = foundElement
									// files = append(files, int(file.ID))
									fmt.Println("encontrou: ", file.ID)
									method.FileID = file.ID
									foundElement = true
									stack = append(stack, *method)
								}

							}
						}
					}

				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// fmt.Println(lines, requestA
}

func ParseProfilingAlloc(db *gorm.DB, commit models.Commit, testcase models.TestCase, output string) {
	output += "-alloc.txt"
	e := os.Rename("profiler/alloc.txt", output)
	if e != nil {
		fmt.Println("Error renaming profile file: ", e.Error())
		fmt.Println("FATAL: Check if MAVEN_OPTS is exported!")
		log.Fatal("Check if MAVEN_OPTS is exported!")
	}

	lines := 0
	f, err := os.Open(output)
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer f.Close()

	foundElement := false
	finished := false
	var stack []models.Method
	var size string
	var ownSize int
	totalCalls := 0
	calls := 0
	// firstInStack := true
	// var prevMethod models.Method = models.Method{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines++
		line := scanner.Bytes()
		if lines == 1 {
			// if strings.HasPrefix(string(line), "--- Execution") {
			continue
		} else if lines == 2 {
			values := strings.Split(string(line), ":")
			totalCalls, e = strconv.Atoi(strings.Trim(values[1], " "))
			if e != nil {
				fmt.Println("Error converting total samples to integer: ", e.Error())
			}
		} else {
			if len(line) > 2 {
				if bytes.Equal(line[:3], []byte("---")) || bytes.Equal(line[:3], []byte("  -")) {
					if foundElement {
						savePreviousStack(db, stack, ownSize, calls, totalCalls, testcase)
					}
					if !bytes.Equal(line[:3], []byte("  -")) {
						values := strings.Split(string(line), " ")
						size = values[1] //+ values[2]
						var errPD error
						ownSize, errPD = strconv.Atoi(size)
						if errPD != nil {
							fmt.Println("Error parsing size of method: ", errPD.Error())
						}
						calls, errPD = strconv.Atoi(values[4])
						if errPD != nil {
							fmt.Println("Error parsing calls of method: ", errPD.Error())
						}

						// firstInStack = true
						foundElement = false
						stack = nil
						finished = false
						// prevMethod = models.Method{}
					}

					// fmt.Println("size: ", values[1])
					// fmt.Println("%: ", values[1])
					// fmt.Println("samples: ", values[1])
					// fmt.Println(values)
				} else {
					if !finished {
						call := strings.Split(string(line), "] ")
						// fmt.Println("call: ", call)
						if len(call) > 1 {
							elements := strings.Split(call[1], ".")
							fmt.Println(elements)
							if len(elements) > 1 {
								element := elements[0]
								method := &models.Method{Name: elements[len(elements)-1], TestCaseID: testcase.ID}
								fmt.Printf(">>>>> method: %s File: %d\n", method.Name, method.FileID)
								// fmt.Println(elements)
								// fmt.Println(len(elements))
								for i := 1; i < len(elements)-1; i++ {
									// fmt.Println("i:", i)
									element += "." + elements[i]
								}

								// search package and class name
								file, err := models.FindFileByEndsWithNameAndCommit(db, element+".java", commit.ID)
								if err != nil {
									// fmt.Println("Error searching for profiled class: ", err.Error())
									// files = append(files, -1)
									if foundElement {
										finished = true
									}
								} else {
									// lastFound = foundElement
									// files = append(files, int(file.ID))
									fmt.Println("encontrou: ", file.ID)
									method.FileID = file.ID
									foundElement = true
									stack = append(stack, *method)
								}

							}
						}
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// fmt.Println(lines, requestA
}

func ParseJfrFile() {
	///home/farah/Downloads/jdk1.8.0_202/jre/bin/java -Dfile.encoding=UTF-8 -classpath /home/farah/eclipse-workspace/async-parser/bin:/home/farah/Downloads/async-profiler-2.5.1-linux-x64/build/async-profiler.jar:/home/farah/Downloads/async-profiler-2.5.1-linux-x64/build/converter.jar JfrParser /home/farah/go-work/src/github.com/paulorfarah/go-repo-downloader/profiler/13__AppTest_testAppHasAGreeting.jfr /home/farah/go-work/src/github.com/paulorfarah/go-repo-downloader/profiler/13__AppTest_testAppHasAGreeting-PARSED.txt
	fmt.Println("TODO")
}

func savePreviousStack(db *gorm.DB, stack []models.Method, ownSize, calls, totalCalls int, testcase models.TestCase) {

	//save previous stack
	fmt.Println("*********************************************************************** stack")
	// fmt.Printf("line: %d ownSize: %v", lines, ownSize)
	firstInStack := true
	var prevMethod models.Method = models.Method{}
	lastMethod := stack[len(stack)-1]
	e := false
	if lastMethod.Name != testcase.Name {
		// check if name of last method is exactly the same name of testcase, they should be...
		e = true
	}
	// for _, m := range stack {
	for i := len(stack) - 2; i >= 0; i-- {
		m := stack[i]
		m.Error = e
		fmt.Printf("Name: %s, Size: %d\n", m.Name, m.OwnSize)
		if m.FileID > 0 {
			method, errM := models.FindMethodByEndsWithNameAndFileAndTestcase(db, m.Name, m.FileID, testcase.ID)
			if errM != nil {
				log.Printf("Create new method: %s (%s)\n"+m.Name, errM.Error())
				if i == 0 {
					m.OwnSize = ownSize
					m.TotalAllocCalls = totalCalls
					m.AllocCalls = calls
					m.AllocCallsPercent = float64(calls / totalCalls)
				}
				if !firstInStack {
					m.CallerID = &prevMethod.ID
					// err = models.SaveMethod(db, &prevMethod)
					// if err != nil {
					// 	fmt.Println("Error saving previous method: ", err.Error())
					// }
				}
				var err error
				m.ID, err = models.CreateMethod(db, &m)
				if err != nil {
					fmt.Println("Error creating method: ", err.Error())
				}

				prevMethod = m
			} else {
				//update size
				if i == 0 {
					method.OwnSize += ownSize
					method.AllocCalls += calls
					method.TotalAllocCalls = totalCalls
					method.AllocCallsPercent = float64(float64(calls) / float64(totalCalls))
				}
				if !firstInStack {
					method.CallerID = &prevMethod.ID
					// err = models.SaveMethod(db, &prevMethod)
					// if err != nil {
					// 	fmt.Println("Error saving previous method: ", err.Error())
					// }
				}
				err := models.SaveMethod(db, method)
				if err != nil {
					fmt.Println("Error saving method: ", err.Error())
				}

				prevMethod = *method
			}
			if firstInStack {
				firstInStack = false
			}
		}
	}
}
