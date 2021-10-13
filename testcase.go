package main

import (
	"fmt"
	"log"
	"os/exec"
)

// func measureTestCases(xmlFile string) {
// 	suites, err := junit.IngestDir(xmlFile)
// 	if err != nil {
// 		log.Fatalf("failed to ingest JUnit xml %v", err)
// 	}
// 	// fmt.Println("suites: ", suites)
// 	for _, suite := range suites {
// 		// fmt.Println(suite.Name)
// 		for _, test := range suite.Tests {
// 			// fmt.Println(test.Classname + ".java")
// 			// fmt.Printf("  %s\n", test.Name)
// 			// if test.Error != nil {
// 			// 	fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
// 			// } else {
// 			// 	fmt.Printf("    %s %f\n", test.Status, test.Duration.Seconds())
// 			// }
// 			classname := strings.Replace(test.Classname, ".", "/", -1)
// 			filename := classname + ".java"
// 			// fmt.Println(filename)
// 			testSuite, errF := models.FindFileByEndsWithNameAndCommit(db, filename, commitID)
// 			if errF != nil {
// 				fmt.Println("error finding file: ", test.Classname, commitID)
// 			}
// 			// fmt.Println("testSuite: ", testSuite)

// 			// errorMsg := ""
// 			// if test.Error != nil {
// 			// 	errorMsg = test.Error.Error()
// 			// }
// 			mr := &models.TestCase{
// 				MeasurementID: measurement.ID,
// 				Type:          "gradle",
// 				ClassName:     test.Classname,
// 				CommitID:      commitID,
// 				// Duration:      test.Duration,
// 				TestSuiteID: testSuite.ID,
// 				Name:        test.Name,
// 				// Status:        string(test.Status),
// 				// Error:         errorMsg,
// 				// Message:       test.Message,
// 				// SystemErr:     string(test.SystemErr),
// 				// SystemOut:     string(test.SystemOut),
// 			}
// 			_, errTC := models.CreateTestCase(db, mr)
// 			if errTC != nil {
// 				fmt.Println(errTC.Error())
// 			}

// 		}
// 	}
// }

func GradleTestCase(path, testCase string) {
	// # Executes a single specified test in SomeTestClass
	// gradle test --tests SomeTestClass.someSpecificMethod

	// ok := true
	// logfile := "gradle-test.log"

	log.Println(">>>------------------------------------------------ gradle testcase", path, testCase)
	fmt.Println(">>>------------------------------------------------ gradle testcase", path, testCase)
	cmd := exec.Command("gradle", "test", "--tests", testCase)
	cmd.Dir = path

	var output []byte
	var err error

	// output, err = cmd.CombinedOutput()
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("gradle test out:\n%s\n", string(output))
}
