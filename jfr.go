package main

import (
	"PerfoRT/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

//https://docs.oracle.com/en/java/javase/16/docs/specs/man/jcmd.html
// https://medium.com/@nate510/dynamic-json-umarshalling-in-go-88095561d6a0

type Recording struct {
	Events []Event `json:"events"`
}

type Event struct {
	Type   string                 `json:"type"`
	Values map[string]interface{} `json:"values"`
}

func SaveJFRMetrics(db *gorm.DB, runID uint, tcID uint) {
	// log.Println("****************** SaveJFRMetrics", time.Now())
	// generate json
	jfrFilename := "jfr/PerfoRT" + strconv.Itoa(int(runID)) + ".jfr"
	jsonFilename := "jfr/PerfoRT" + strconv.Itoa(int(runID)) + ".json"
	if _, err := os.Stat(jfrFilename); err == nil {
		// file PerfoRT.jfr exists
		// log.Println("- jfr print --json " + jfrFilename + " > " + jsonFilename)
		cmd := exec.Command("bash", "-c", "jfr print --json "+jfrFilename+" > "+jsonFilename)
		err := cmd.Run()
		if err != nil {
			log.Println("-> Error converting jfr to json in "+jfrFilename, err.Error())
		}

		// parse json file
		// Open our jsonFile
		jsonFile, err := os.Open(jsonFilename)
		// defer the closing of our jsonFile so that we can parse it later on

		// if we os.Open returns an error then handle it
		if err != nil {
			log.Println("-> Error opening jfr json file: ", err.Error())
		} else {
			defer jsonFile.Close()
			// fmt.Println("Successfully Opened PerfoRT.json")

			// read our opened jsonFile as a byte array.
			jsonJFR, _ := ioutil.ReadAll(jsonFile)

			// GJson
			result := gjson.Get(string(jsonJFR), "recording.events")
			jfrMap := make(map[time.Time]models.Jvm)
			result.ForEach(func(key, eventMap gjson.Result) bool {
				var event Event

				// fmt.Println(eventMap)
				// fmt.Println("-----")
				var result map[string]interface{}
				jsonFile := fmt.Sprintf("%s", eventMap)
				if err := json.Unmarshal([]byte(jsonFile), &result); err != nil {
					log.Println("### ERROR unmarshalling json file: ", err)
					log.Println(jsonFile)
					fmt.Println("### ERROR unmarshalling json file: ", err)
					fmt.Println(jsonFile)
				} else {
					values := result["values"]
					if v, ok := values.(map[string]interface{}); ok {
						// fmt.Printf(" %s", v["startTime"])
						event.Type = result["type"].(string)
						event.Values = v

					} else {
						log.Printf("record not a map[string]interface{}: %v\n", values)
						// panic("record not a map[string]interface{}: ")

					}

					// fmt.Printf("%#v\n", event.Values)
					layout := "2006-01-02T15:04:05.000000000-07:00"
					startTimeStr, ok := event.Values["startTime"].(string)
					if ok {
						indStartTime := strings.Index(startTimeStr[20:29], "-")

						if indStartTime != -1 {
							// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> starttime1: ", str)
							var auxStartime string
							for i := indStartTime; i < 9; i++ {
								auxStartime += "0"
							}
							startTimeStr = startTimeStr[:20+indStartTime] + auxStartime + startTimeStr[20+indStartTime:]
							// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> starttime2: ", str)
						}
						t, err := time.Parse(layout, startTimeStr)

						if err != nil {
							log.Println("ERROR parsing startTime: ", err)
						}
						// fmt.Println("->->->->->->->->->->->->-> Event type: "+event.Type+" ", t)
						// log.Println("->->->->->->->->->->->->-> Event type: "+event.Type+" ", t)
						switch event.Type {
						case "jdk.CPULoad":
							cpuLoad := &models.CPULoad{
								JvmUser:      event.Values["jvmUser"].(float64),
								JvmSystem:    event.Values["jvmSystem"].(float64),
								MachineTotal: event.Values["machineTotal"].(float64),
							}
							// fmt.Printf("%v\n", cpuLoad)
							// save
							if val, ok := jfrMap[t]; ok {
								val.CPULoad = *cpuLoad
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:       models.Run{},
									StartTime: &t,
									CPULoad:   *cpuLoad,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ThreadCPULoad":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							tCpuLoad := &models.ThreadCPULoad{
								ThreadCPULoadOsName:       osName,
								ThreadCPULoadOsThreadId:   osThreadId,
								ThreadCPULoadJavaName:     javaName,
								ThreadCPULoadJavaThreadId: javaThreadId,
								ThreadCPULoadUser:         event.Values["user"].(float64),
								ThreadCPULoadSystem:       event.Values["system"].(float64),
							}

							if val, ok := jfrMap[t]; ok {
								val.ThreadCPULoad = *tCpuLoad
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:           models.Run{},
									StartTime:     &t,
									ThreadCPULoad: *tCpuLoad,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ThreadStart":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							osNameParent, osThreadIdParent, javaNameParent, javaThreadIdParent := getEventParentThread(event)
							threadStart := &models.ThreadStart{
								ThreadStartOsName:                   osName,
								ThreadStartOsThreadId:               osThreadId,
								ThreadStartJavaName:                 javaName,
								ThreadStartJavaThreadId:             javaThreadId,
								ThreadStartParentThreadosName:       osNameParent,
								ThreadStartParentThreadOsThreadId:   osThreadIdParent,
								ThreadStartParentThreadJavaName:     javaNameParent,
								ThreadStartParentThreadJavaThreadId: javaThreadIdParent,
							}
							if val, ok := jfrMap[t]; ok {
								val.ThreadStart = *threadStart
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:         models.Run{},
									StartTime:   &t,
									ThreadStart: *threadStart,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ThreadEnd":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							threadEnd := &models.ThreadEnd{
								ThreadEndOsName:       osName,
								ThreadEndOsThreadId:   osThreadId,
								ThreadEndJavaName:     javaName,
								ThreadEndJavaThreadId: javaThreadId,
							}
							if val, ok := jfrMap[t]; ok {
								val.ThreadEnd = *threadEnd
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:       models.Run{},
									StartTime: &t,
									ThreadEnd: *threadEnd,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ThreadSleep":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							// var duration time.Time
							// if event.Values["duration"] != nil {
							// 	// auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	// auxDur = strings.ReplaceAll(auxDur, "S", "")
							// 	// duration, err = strconv.ParseFloat(auxDur, 64)
							// 	// if err != nil {
							// 	// 	log.Println("ERROR in jdk.ThreadSleep: cannot parse duration value. ", err)
							// 	// }
							// 	var errParse error
							// 	auxDur := event.Values["duration"].(string)
							// 	auxDur = strings.ReplaceAll(auxDur, "PT", "")
							// 	if strings.Contains(auxDur, "M") {
							// 		// time := "29M25.974594914"
							// 		duration, errParse = time.Parse("04M05.999999999", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxDur, "S") {
							// 		// time := "0.05S"
							// 		duration, errParse = time.Parse("5.999999999S", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxDur)
							// 	}
							// }

							// // "time": "PT30S"
							// var tsTime time.Time
							// if event.Values["time"] != nil {
							// 	// auxTime := strings.ReplaceAll(event.Values["time"].(string), "PT", "")
							// 	// auxTime = strings.ReplaceAll(auxTime, "S", "")
							// 	// tsTime, err = strconv.ParseFloat(auxTime, 64)
							// 	// if err != nil {
							// 	// 	log.Println("ERROR in jdk.ThreadSleep: cannot parse time value. ", err)
							// 	// }
							// 	var errParse error
							// 	auxT := event.Values["time"].(string)
							// 	auxT = strings.ReplaceAll(auxT, "PT", "")
							// 	if strings.Contains(auxT, "M") {
							// 		// time := "29M25.974594914"
							// 		tsTime, errParse = time.Parse("04M05.999999999", auxT)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark tsTime: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxT, "S") {
							// 		// time := "0.05S"
							// 		tsTime, errParse = time.Parse("5.999999999S", auxT)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark tsTime: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass tsTime: ", auxT)
							// 	}
							// }

							duration := ParseDuration(event.Values["duration"].(string))
							tsTime := ParseDuration(event.Values["time"].(string))
							threadSleep := &models.ThreadSleep{
								ThreadSleepDuration:     duration,
								ThreadSleepOsName:       osName,
								ThreadSleepOsThreadId:   osThreadId,
								ThreadSleepJavaName:     javaName,
								ThreadSleepJavaThreadId: javaThreadId,
								ThreadSleepTime:         tsTime,
							}
							if val, ok := jfrMap[t]; ok {
								val.ThreadSleep = *threadSleep
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:         models.Run{},
									StartTime:   &t,
									ThreadSleep: *threadSleep,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ThreadPark":
							// 2022/08/26 23:46:04 ERROR in jdk.ThreadPark: cannot parse timeOut value.  strconv.ParseFloat: parsing "29M25.974594914": invalid syntax
							// 2022/08/26 23:46:04 ERROR parsing ThreadPark.until timestamp:  strconv.ParseFloat: parsing "29M25.974594914": invalid syntax
							// 2022/08/26 23:46:04 ERROR parsing ThreadPark.until timestamp:  <nil>

							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							parkedClass := getClass(event, "parkedClass")
							// var duration time.Time
							// var timeOut time.Time

							// if event.Values["duration"] != nil {
							// 	var errParse error
							// 	auxDur := event.Values["duration"].(string)
							// 	auxDur = strings.ReplaceAll(auxDur, "PT", "")
							// 	if strings.Contains(auxDur, "M") {
							// 		// time := "29M25.974594914"
							// 		duration, errParse = time.Parse("04M05.999999999", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxDur, "S") {
							// 		// time := "0.05S"
							// 		duration, errParse = time.Parse("5.999999999S", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxDur)
							// 	}
							// }
							// if event.Values["timeout"] != nil {
							// 	//"timeout": "PT0.05S",
							// 	// "29M25.974594914"
							// 	var errParse error
							// 	auxTO := strings.ReplaceAll(event.Values["timeout"].(string), "PT", "")
							// 	if strings.Contains(auxTO, "M") {
							// 		// time := "29M25.974594914"
							// 		duration, errParse = time.Parse("04M05.999999999", auxTO)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark timeout: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxTO, "S") {
							// 		// time := "0.05S"
							// 		duration, errParse = time.Parse("5.999999999S", auxTO)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark timeout: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxTO)
							// 	}
							// }
							// var until time.Time
							// // var errUntil error
							// if event.Values["until"] != nil {
							// 	// event.Values["until"].(float64),"until": "-999999999-01-01T00:00+18:00",

							// 	layout := "2006-01-02T15:04:05.000000000-07:00"
							// 	untilStr, okUntil := event.Values["until"].(string)

							// 	if okUntil {
							// 		// until, errUntil = time.Parse(layout, untilStr)
							// 		until, _ = time.Parse(layout, untilStr)

							// 		// if errUntil != nil {
							// 		// 	log.Println("ERROR parsing ThreadPark.until timestamp: ", errUntil)
							// 		// }

							// 	}
							// }
							duration := ParseDuration(event.Values["duration"].(string))
							timeOut := ParseDuration(event.Values["timeout"].(string))
							until := ParseDuration(event.Values["until"].(string))
							threadPark := &models.ThreadPark{
								ThreadParkDuration:     duration,
								ThreadParkOsName:       osName,
								ThreadParkOsThreadId:   osThreadId,
								ThreadParkJavaName:     javaName,
								ThreadParkJavaThreadId: javaThreadId,
								ThreadParkParkedClass:  parkedClass,
								ThreadParkTimeout:      timeOut,
								ThreadParkUntil:        until,
							}
							if val, ok := jfrMap[t]; ok {
								val.ThreadPark = *threadPark
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:        models.Run{},
									StartTime:  &t,
									ThreadPark: *threadPark,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaErrorThrow":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							thrownClass := getClass(event, "thrownClass")
							// var duration time.Time
							// if event.Values["duration"] != nil {
							// 	// auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	// auxDur = strings.ReplaceAll(auxDur, "S", "")
							// 	// duration, err = strconv.ParseFloat(auxDur, 64)
							// 	// if err != nil {
							// 	// 	log.Println("ERROR in jdk.JavaErrorThrow: cannot parse duration value. ", err)
							// 	// }
							// 	var errParse error
							// 	auxDur := event.Values["duration"].(string)
							// 	auxDur = strings.ReplaceAll(auxDur, "PT", "")
							// 	if strings.Contains(auxDur, "M") {
							// 		// time := "29M25.974594914"
							// 		duration, errParse = time.Parse("04M05.999999999", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxDur, "S") {
							// 		// time := "0.05S"
							// 		duration, errParse = time.Parse("5.999999999S", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxDur)
							// 	}
							// }
							duration := ParseDuration(event.Values["duration"].(string))

							message := ""
							if event.Values["message"] != nil {
								message = event.Values["message"].(string)
							}

							javaErrorThrow := &models.JavaErrorThrow{
								JavaErrorThrowDuration:     duration,
								JavaErrorThrowOsName:       osName,
								JavaErrorThrowOsThreadId:   osThreadId,
								JavaErrorThrowJavaName:     javaName,
								JavaErrorThrowJavaThreadId: javaThreadId,
								JavaErrorThrowMessage:      message,
								JavaErrorThrowThrownClass:  thrownClass,
							}
							if val, ok := jfrMap[t]; ok {
								val.JavaErrorThrow = *javaErrorThrow
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:            models.Run{},
									StartTime:      &t,
									JavaErrorThrow: *javaErrorThrow,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaExceptionThrow":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							thrownClass := getClass(event, "thrownClass")
							// var duration time.Time
							// if event.Values["duration"] != nil {
							// 	// auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	// auxDur = strings.ReplaceAll(auxDur, "S", "")
							// 	// duration, err = strconv.ParseFloat(auxDur, 64)
							// 	// if err != nil {
							// 	// 	log.Println("ERROR in jdk.JavaExceptionThrow: cannot parse duration value. ", err)
							// 	// }
							// 	var errParse error
							// 	auxDur := event.Values["duration"].(string)
							// 	auxDur = strings.ReplaceAll(auxDur, "PT", "")
							// 	if strings.Contains(auxDur, "M") {
							// 		// time := "29M25.974594914"
							// 		duration, errParse = time.Parse("04M05.999999999", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxDur, "S") {
							// 		// time := "0.05S"
							// 		duration, errParse = time.Parse("5.999999999S", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark duration: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxDur)
							// 	}
							// }
							duration := ParseDuration(event.Values["duration"].(string))
							javaExceptionThrow := &models.JavaExceptionThrow{
								JavaExceptionThrowDuration:     duration,
								JavaExceptionThrowOsName:       osName,
								JavaExceptionThrowOsThreadId:   osThreadId,
								JavaExceptionThrowJavaName:     javaName,
								JavaExceptionThrowJavaThreadId: javaThreadId,
								JavaExceptionThrowMessage:      event.Values["message"].(string),
								JavaExceptionThrowThrownClass:  thrownClass,
							}
							if val, ok := jfrMap[t]; ok {
								val.JavaExceptionThrow = *javaExceptionThrow
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:                models.Run{},
									StartTime:          &t,
									JavaExceptionThrow: *javaExceptionThrow,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaMonitorEnter":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							monitorClass := getClass(event, "monitorClass")
							// var duration float64
							// if event.Values["duration"] != nil {
							// 	auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	auxDur = strings.ReplaceAll(auxDur, "S", "")
							// 	duration, err = strconv.ParseFloat(auxDur, 64)
							// 	if err != nil {
							// 		log.Println("ERROR in jdk.JavaMonitorEnter: cannot parse duration value. ", err)
							// 	}
							// }
							duration := ParseDuration(event.Values["duration"].(string))
							javaMonitorEnter := &models.JavaMonitorEnter{
								JavaMonitorEnterDuration:     duration,
								JavaMonitorEnterOsName:       osName,
								JavaMonitorEnterOsThreadId:   osThreadId,
								JavaMonitorEnterJavaName:     javaName,
								JavaMonitorEnterJavaThreadId: javaThreadId,
								JavaMonitorEnterMonitorClass: monitorClass,
							}
							if val, ok := jfrMap[t]; ok {
								val.JavaMonitorEnter = *javaMonitorEnter
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:              models.Run{},
									StartTime:        &t,
									JavaMonitorEnter: *javaMonitorEnter,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaMonitorWait":
							// JavaMonitorWai
							// 2022/08/27 19:36:56 ERROR in jdk.JavaMonitorWait: cannot parse duration value.  strconv.ParseFloat: parsing "1M0.000118777": invalid syntax
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							monitorClass := getClass(event, "monitorClass")

							// var duration time.Time
							// if event.Values["duration"] != nil {
							// 	var errParse error
							// 	auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	if strings.Contains(auxDur, "M") {
							// 		// Error parsing JavaMonitorWait duration:  parsing time "1M0.000088702S" as "04M05.999999999": cannot parse "1M0.000088702S" as "04"
							// 		duration, errParse = time.Parse("94M05.999999999", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing JavaMonitorWait duration: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxDur, "S") {
							// 		//parsing time                  "0.100194234S": extra text: "S"0.100193907S
							// 		duration, errParse = time.Parse("5.999999999S", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing JavaMonitorWait duration: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxDur)
							// 	}
							// }
							duration := ParseDuration(event.Values["duration"].(string))

							// var timeOut time.Time
							// if event.Values["timeOut"] != nil {
							// 	// auxtimeOut := strings.ReplaceAll(event.Values["timeOut"].(string), "PT", "")
							// 	// auxtimeOut = strings.ReplaceAll(auxtimeOut, "S", "")
							// 	// timeOut, err = strconv.ParseFloat(auxtimeOut, 64)
							// 	// if err != nil {
							// 	// 	log.Println("ERROR in jdk.JavaMonitorWait: cannot parse timeOut value. ", err)
							// 	// }
							// 	var errParse error
							// 	auxTO := event.Values["timeOut"].(string)
							// 	auxTO = strings.ReplaceAll(auxTO, "PT", "")
							// 	if strings.Contains(auxTO, "M") {
							// 		// time := "29M25.974594914"
							// 		timeOut, errParse = time.Parse("04M05.999999999", auxTO)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark timeOut: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxTO, "S") {
							// 		// time := "0.05S"
							// 		timeOut, errParse = time.Parse("5.999999999S", auxTO)
							// 		if errParse != nil {
							// 			log.Println("Error parsing ThreadPark timeOut: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass timeOut: ", auxTO)
							// 	}
							// }

							var timeOut time.Duration
							if event.Values["timeOut"] != nil {
								timeOut = ParseDuration(event.Values["timeOut"].(string))
							}
							javaMonitorWait := &models.JavaMonitorWait{
								JavaMonitorWaitDuration:     duration,
								JavaMonitorWaitOsName:       osName,
								JavaMonitorWaitOsThreadId:   osThreadId,
								JavaMonitorWaitJavaName:     javaName,
								JavaMonitorWaitJavaThreadId: javaThreadId,
								JavaMonitorWaitMonitorClass: monitorClass,
								JavaMonitorWaitTimeout:      timeOut,
								JavaMonitorWaitTimedOut:     event.Values["timedOut"].(bool),
							}
							if val, ok := jfrMap[t]; ok {
								val.JavaMonitorWait = *javaMonitorWait
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:             models.Run{},
									StartTime:       &t,
									JavaMonitorWait: *javaMonitorWait,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.OldObjectSample":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							objectType := getObjectType(event)
							// var duration time.Time
							// if event.Values["duration"] != nil {
							// 	// auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	// auxDur = strings.ReplaceAll(auxDur, "S", "")
							// 	// duration, err = strconv.ParseFloat(auxDur, 64)
							// 	// if err != nil {
							// 	// 	log.Println("ERROR in jdk.OldObjectSample: cannot parse duration value. ", err)
							// 	// }
							// 	var errParse error
							// 	auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	if strings.Contains(auxDur, "M") {
							// 		// Error parsing JavaMonitorWait duration:  parsing time "1M0.000088702S" as "04M05.999999999": cannot parse "1M0.000088702S" as "04"
							// 		duration, errParse = time.Parse("94M05.999999999", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing JavaMonitorWait duration: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxDur, "S") {
							// 		//parsing time                  "0.100194234S": extra text: "S"0.100193907S
							// 		duration, errParse = time.Parse("5.999999999S", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing JavaMonitorWait duration: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxDur)
							// 	}
							// }
							duration := ParseDuration(event.Values["duration"].(string))

							// Error 1292: Incorrect datetime value: '0000-00-00' for column 'old_object_sample_allocation_time' at row 1

							oldObjectSample := &models.OldObjectSample{
								OldObjectSampleDuration:     duration,
								OldObjectSampleOsName:       osName,
								OldObjectSampleOsThreadId:   osThreadId,
								OldObjectSampleJavaName:     javaName,
								OldObjectSampleJavaThreadId: javaThreadId,
								// OldObjectSampleAllocationTime:     at,
								OldObjectSampleLastKnownHeapUsage: event.Values["lastKnownHeapUsage"].(float64),
								OldObjectSampleObject:             objectType,
								OldObjectSampleArrayElements:      event.Values["arrayElements"].(float64),
							}

							// "allocationTime": "2022-05-22T18:48:37.932136923-07:00",
							//    2006-01-02T15:04:05.000000000-07:00
							// ERROR in jdk.OldObjectSample: cannot parse allocationTime value.  parsing time "2022-08-27T18:44:26.019938-03:00" as "2006-01-02T15:04:05.000000000-07:00": cannot parse ":00" as ".000000000"
							var at time.Time
							atStr := event.Values["allocationTime"].(string)
							layoutAT := "2006-01-02T15:04:05.999999-07:00"
							at, err = time.Parse(layoutAT, atStr)
							if err != nil {
								log.Println("ERROR in jdk.OldObjectSample: cannot parse allocationTime value. ", err)
							}

							if !at.IsZero() {
								oldObjectSample.OldObjectSampleAllocationTime.Time = at
								oldObjectSample.OldObjectSampleAllocationTime.Valid = true
							} else {
								oldObjectSample.OldObjectSampleAllocationTime.Time = time.Time{}
								oldObjectSample.OldObjectSampleAllocationTime.Valid = false
							}

							if val, ok := jfrMap[t]; ok {
								val.OldObjectSample = *oldObjectSample
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:             models.Run{},
									StartTime:       &t,
									OldObjectSample: *oldObjectSample,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ClassLoaderStatistics":
							name, parent := getClassLoader(event)
							aCStr := event.Values["anonymousClassCount"]
							aCSStr := event.Values["anonymousChunkSize"]
							aBSStr := event.Values["anonymousBlockSize"]

							var aC float64
							var aCS float64
							var aBS float64
							if aCStr != nil {
								aC = aCStr.(float64)
							}
							if aCSStr != nil {
								aCS = aCSStr.(float64)
							}
							if aBSStr != nil {
								aBS = aBSStr.(float64)
							}
							classLoaderStatistics := &models.ClassLoaderStatistics{
								ClassLoader:         name,
								ParentClassLoader:   parent,
								ClassLoaderData:     event.Values["classLoaderData"].(float64),
								ClassCount:          event.Values["classCount"].(float64),
								ChunkSize:           event.Values["chunkSize"].(float64),
								BlockSize:           event.Values["blockSize"].(float64),
								AnonymousClassCount: aC,
								AnonymousChunkSize:  aCS,
								AnonymousBlockSize:  aBS,
							}
							if val, ok := jfrMap[t]; ok {
								val.ClassLoaderStatistics = *classLoaderStatistics
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:                   models.Run{},
									StartTime:             &t,
									ClassLoaderStatistics: *classLoaderStatistics,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ObjectAllocationInNewTLAB":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							objectClass := getClass(event, "objectClass")
							objectAllocationInNewTLAB := &models.ObjectAllocationInNewTLAB{

								ObjectAllocationInNewTLABOsName:         osName,
								ObjectAllocationInNewTLABOsThreadId:     osThreadId,
								ObjectAllocationInNewTLABJavaName:       javaName,
								ObjectAllocationInNewTLABJavaThreadId:   javaThreadId,
								ObjectAllocationInNewTLABObjectClass:    objectClass,
								ObjectAllocationInNewTLABAllocationSize: event.Values["allocationSize"].(float64),
								ObjectAllocationInNewTLABTlabSize:       event.Values["tlabSize"].(float64),
							}
							if val, ok := jfrMap[t]; ok {
								val.ObjectAllocationInNewTLAB = *objectAllocationInNewTLAB
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:                       models.Run{},
									StartTime:                 &t,
									ObjectAllocationInNewTLAB: *objectAllocationInNewTLAB,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ObjectAllocationOutsideTLAB":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							objectClass := getClass(event, "objectClass")
							objectAllocationOutsideTLAB := &models.ObjectAllocationOutsideTLAB{
								ObjectAllocationOutsideTLABOsName:         osName,
								ObjectAllocationOutsideTLABOsThreadId:     osThreadId,
								ObjectAllocationOutsideTLABJavaName:       javaName,
								ObjectAllocationOutsideTLABJavaThreadId:   javaThreadId,
								ObjectAllocationOutsideTLABObjectClass:    objectClass,
								ObjectAllocationOutsideTLABAllocationSize: event.Values["allocationSize"].(float64),
							}
							if val, ok := jfrMap[t]; ok {
								val.ObjectAllocationOutsideTLAB = *objectAllocationOutsideTLAB
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID: runID,
									// Run:                         models.Run{},
									StartTime:                   &t,
									ObjectAllocationOutsideTLAB: *objectAllocationOutsideTLAB,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.GCPhasePause":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							// var duration time.Time
							// if event.Values["duration"] != nil {
							// 	// auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	// auxDur = strings.ReplaceAll(auxDur, "S", "")
							// 	// duration, err = strconv.ParseFloat(auxDur, 64)
							// 	// if err != nil {
							// 	// 	log.Println("ERROR in jdk.GCPhasePause: cannot parse duration value. ", err)
							// 	// }
							// 	var errParse error
							// 	auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
							// 	if strings.Contains(auxDur, "M") {
							// 		// Error parsing JavaMonitorWait duration:  parsing time "1M0.000088702S" as "04M05.999999999": cannot parse "1M0.000088702S" as "04"
							// 		duration, errParse = time.Parse("94M05.999999999", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing JavaMonitorWait duration: ", errParse)
							// 		}
							// 	} else if strings.Contains(auxDur, "S") {
							// 		//parsing time                  "0.100194234S": extra text: "S"0.100193907S
							// 		duration, errParse = time.Parse("5.999999999S", auxDur)
							// 		if errParse != nil {
							// 			log.Println("Error parsing JavaMonitorWait duration: ", errParse)
							// 		}
							// 	} else {
							// 		log.Println("ATTENTION: Cannot parse parkedClass duration: ", auxDur)
							// 	}
							// }
							duration := ParseDuration(event.Values["duration"].(string))
							gcPhasePause := &models.GCPhasePause{
								GCPhasePauseDuration:     duration,
								GCPhasePauseOsName:       osName,
								GCPhasePauseOsThreadId:   osThreadId,
								GCPhasePauseJavaName:     javaName,
								GCPhasePauseJavaThreadId: javaThreadId,
								GcId:                     event.Values["gcId"].(float64),
								GCPhasePauseName:         event.Values["name"].(string),
							}
							if val, ok := jfrMap[t]; ok {
								val.GCPhasePause = *gcPhasePause
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID:        runID,
									StartTime:    &t,
									GCPhasePause: *gcPhasePause,
								}
								jfrMap[t] = *jvm

							}
						}
						//end ok starttime
					} else {
						log.Println("Cannot convert interface: interface {} is not string")
					}
				}

				// e := event.Data().(map[string]interface{})
				// for key, value := range events {
				// 	if key == "type" {
				// 		switch value {
				// 		case "jdk.CPULoad":
				// 			var result CPULoad
				// 			err := Decode(input, &result)
				// 			if err != nil {
				// 				panic(err)
				// 			}

				// 			fmt.Printf("%#v", result)
				// 		}
				// 	}
				// }

				// 	var result map[string]interface{}
				// 	err := mapstructure.Decode(events, &result)
				// 	if err != nil {
				// 		panic(err)
				// 	}
				// 	for key, value := range result {
				// 		fmt.Printf("%s: %#v\n", key, value)
				// 	}
				// 	// for key, value := range events {
				// 	// 	if key == "type" {
				// 	// 		switch value {
				// 	// 		case "jdk.CPULoad":
				// 	// 			var result CPULoad
				// 	// 			err := Decode(input, &result)
				// 	// 			if err != nil {
				// 	// 				panic(err)
				// 	// 			}

				// 	// 			fmt.Printf("%#v", result)
				// 	// 		}
				// 	// 	}
				// 	// 	fmt.Printf("%s - %s\n", key, value)
				// 	// }
				return true // keep iterating
			})

			// for _, jvm := range jfrMap {

			// 	_, err = models.CreateJvm(db, &jvm)
			// 	if err != nil {
			// 		log.Printf("Error saving jvm: %s\n", err.Error())
			// 		log.Println("jvm.StartTime: ", jvm.StartTime)
			// 		log.Println("ThreadSleepDuration: ", jvm.ThreadSleepDuration)
			// 		log.Println("ThreadParkDuration: ", jvm.ThreadParkDuration)
			// 		log.Println("ThreadParkTimeout: ", jvm.ThreadParkTimeout)
			// 		log.Println("ThreadParkUntil: ", jvm.ThreadParkUntil)
			// 		log.Println("JavaErrorThrowDuration: ", jvm.JavaErrorThrowDuration)
			// 		log.Println("JavaExceptionThrowDuration: ", jvm.JavaExceptionThrowDuration)
			// 		log.Println("JavaMonitorWaitDuration: ", jvm.JavaMonitorWaitDuration)
			// 		log.Println("JavaMonitorWaitTimeout: ", jvm.JavaMonitorWaitTimeout)
			// 		log.Println("OldObjectSampleDuration: ", jvm.OldObjectSampleDuration)
			// 		log.Println("GCPhasePauseDuration: ", jvm.GCPhasePauseDuration)
			// 		log.Println("-----")
			// 	}
			// }
			jfrValues := maps.Values(jfrMap)
			// fmt.Println("jfrValues: ", len(jfrValues))
			db.CreateInBatches(jfrValues, 1000)

		}

	} else {
		log.Println("ATTENTION!!!!!!! JFR file not found: " + jfrFilename + "!!!")
	}

}

func getClassLoader(event Event) (string, string) {
	var name string
	classloader := event.Values["classLoader"]
	if n, ok := classloader.(map[string]interface{}); ok {
		// fmt.Println(n)
		if n["name"] != nil {
			name = n["name"].(string)
		}
	}

	var parent string
	parentClassloader := event.Values["parentClassLoader"]
	if n, ok := parentClassloader.(map[string]interface{}); ok {
		if n["name"] != nil {
			parent = n["name"].(string)
		}
	}
	return name, parent
}

func getEventThread(event Event) (string, float64, string, float64) {
	var osName string
	eventThread := event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osName"] != nil {
			osName = e["osName"].(string)
		}
	}

	var osThreadId float64
	eventThread = event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osThreadId"] != nil {
			osThreadId = e["osThreadId"].(float64)
		}
	}

	var javaName string
	eventThread = event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaName"] != nil {
			javaName = e["javaName"].(string)
		}
	}

	var javaThreadId float64
	eventThread = event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaThreadId"] != nil {
			osThreadId = e["javaThreadId"].(float64)
		}
	}
	return osName, osThreadId, javaName, javaThreadId
}

func getEventParentThread(event Event) (string, float64, string, float64) {
	var osName string
	eventThread := event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osName"] != nil {
			osName = e["osName"].(string)
		}
	}

	var osThreadId float64
	eventThread = event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osThreadId"] != nil {
			osThreadId = e["osThreadId"].(float64)
		}
	}

	var javaName string
	eventThread = event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaName"] != nil {
			javaName = e["javaName"].(string)
		}
	}

	var javaThreadId float64
	eventThread = event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaThreadId"] != nil {
			osThreadId = e["javaThreadId"].(float64)
		}
	}
	return osName, osThreadId, javaName, javaThreadId
}

func getClass(event Event, classType string) string {
	var name string
	class := event.Values[classType]
	if e, ok := class.(map[string]interface{}); ok {
		if e["name"] != nil {
			name = e["name"].(string)
		}
	}
	return name
}

func getObjectType(event Event) string {
	var name string
	obj := event.Values["object"]
	if e, ok := obj.(map[string]interface{}); ok {
		if e["type"] != nil {
			oType := e["type"]
			if t, ok := oType.(map[string]interface{}); ok {
				if t["name"] != nil {
					name = t["name"].(string)
				}
			}
		}
	}
	return name
}

func ParseDuration(str string) time.Duration {
	// ParseDuration converts a ISO8601 duration into a time.Duration
	dur := time.Duration(0)

	var durationRegex = regexp.MustCompile(`P([\d\.]+Y)?([\d\.]+M)?([\d\.]+D)?T?([\d\.]+H)?([\d\.]+M)?([\d\.]+?S)?`)
	matches := durationRegex.FindStringSubmatch(str)
	if len(matches) > 0 {

		years := parseDurationPart(matches[1], time.Hour*24*365)
		months := parseDurationPart(matches[2], time.Hour*24*30)
		days := parseDurationPart(matches[3], time.Hour*24)
		hours := parseDurationPart(matches[4], time.Hour)
		minutes := parseDurationPart(matches[5], time.Second*60)
		seconds := parseDurationPart(matches[6], time.Second)
		dur = time.Duration(years + months + days + hours + minutes + seconds)
	}
	return dur
}

func parseDurationPart(value string, unit time.Duration) time.Duration {
	if len(value) != 0 {
		if parsed, err := strconv.ParseFloat(value[:len(value)-1], 64); err == nil {
			return time.Duration(float64(unit) * parsed)
		}
	}
	return 0
}
