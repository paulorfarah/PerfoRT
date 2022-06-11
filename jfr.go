package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"perfrt/models"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
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

// type JFRMetrics struct {
// 	Recording
// 	CPULoad models.CPULoad
// MemoryPercent  float32
// IOCounters     *process.IOCountersStat
// MemoryInfo     *process.MemoryInfoStat
// PageFaults     *process.PageFaultsStat
// Load           *load.AvgStat
// CPUTimes       []cpu.TimesStat
// VirtualMemory  *mem.VirtualMemoryStat
// SwapMemory     *models.SwapMemoryStat
// DiskIOCounters map[string]disk.IOCountersStat
// NetIOCounters  []net.IOCountersStat
// }

// func StartJFR(pid int) {
// 	cmd := exec.Command("jcmd", strconv.Itoa(pid), "JFR.start", "settings=perfrt.jfc")

// 	err := cmd.Run()

// 	if err != nil {
// 		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! ERROR starting JFR: ", err)
// 	}
// }

// func DumpJVM(pid int) (JFRMetrics, error) {

// p, err := process.NewProcess(int32(pid))
// if err != nil {
// 	fmt.Println("Error CPU Percent: ", err.Error())
// }
// cp, err := p.CPUPercent()
// if err != nil {
// 	return JFRMetrics{}, err
// }
// m, err := p.MemoryPercent()
// if err != nil {
// 	return JFRMetrics{}, err
// }
// io, err := p.IOCounters()
// if err != nil {
// 	return JFRMetrics{}, err
// }

// mi, err := p.MemoryInfo()
// if err != nil {
// 	return PerfMetrics{}, err
// }
// pf, err := p.PageFaults()
// if err != nil {
// 	return PerfMetrics{}, err
// }

// l, _ := load.Avg()
// ct, _ := cpu.Times(false)
// vm, _ := mem.VirtualMemory()

// swap, _ := mem.SwapMemory()
// sm := &models.SwapMemoryStat{
// 	SwapTotal:       swap.Total,
// 	SwapUsed:        swap.Used,
// 	SwapFree:        swap.Free,
// 	SwapUsedPercent: swap.UsedPercent,
// 	Sin:             swap.Sin,
// 	Sout:            swap.Sout,
// 	PgIn:            swap.PgIn,
// 	PgOut:           swap.PgOut,
// 	PgFault:         swap.PgFault,
// 	PgMajFaults:     swap.PgMajFault,
// }

// mem.SwapMemory()
// di, _ := disk.IOCounters()
// ni, _ := net.IOCounters(false)

// perfMetrics := &PerfMetrics{
// 	CpuPercent:     cp,
// 	MemoryPercent:  m,
// 	IOCounters:     io,
// 	MemoryInfo:     mi,
// 	PageFaults:     pf,
// 	Load:           l,
// 	CPUTimes:       ct,
// 	VirtualMemory:  vm,
// 	SwapMemory:     sm,
// 	DiskIOCounters: di,
// 	NetIOCounters:  ni,
// }
// return JFRMetrics{}, nil
// }

// func memMonitor() {
// 	v, _ := mem.VirtualMemory()

// 	// almost every return value is a struct
// 	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

// 	// convert to JSON. String() is also implemented
// 	fmt.Println(v)
// }

// func cpuMonitor(db *gorm.DB, measurementID uint) {
// 	percent, _ := cpu.Percent(time.Second, true)
// 	fmt.Printf("  User: %.2f\n", percent[cpu.CPUser])
// 	fmt.Printf("  Nice: %.2f\n", percent[cpu.CPNice])
// 	fmt.Printf("   Sys: %.2f\n", percent[cpu.CPSys])
// 	fmt.Printf("  Intr: %.2f\n", percent[cpu.CPIntr])
// 	fmt.Printf("  Idle: %.2f\n", percent[cpu.CPIdle])
// 	fmt.Printf("States: %.2f\n", percent[cpu.CPUStates])
// 	mr := &models.MavenResources{
// 		MeasurementID: measurementID,
// 		CPUser:        percent[cpu.CPUser],
// 		CPNice:        percent[cpu.CPNice],
// 		CPSys:         percent[cpu.CPSys],
// 		CPIntr:        percent[cpu.CPIntr],
// 		CPIdle:        percent[cpu.CPIdle],
// 		CPUStates:     percent[cpu.CPUStates],
// 	}
// 	models.CreateMavenResources(db, mr)
// }

// func loadMonitor() {
// 	load, _ := load.Avg()
// 	fmt.Println(" 1 min ave:", load.Load1)
// 	fmt.Println(" 5 min ave:", load.Load5)
// 	fmt.Println("15 min ave:", load.Load15)
// }

func SaveJFRMetrics(db *gorm.DB, measurementID uint, tcID uint) {
	// fmt.Println("****************** SaveJFRMetrics")
	// generate json
	if _, err := os.Stat("perfrt.jfr"); err == nil {
		// file perfrt.jfr exists
		log.Println("- jfr print --json perfrt.jfr > perfrt.json")
		cmd := exec.Command("bash", "-c", "jfr print --json perfrt.jfr > perfrt.json")
		err := cmd.Run()
		if err != nil {
			log.Println("-> Error executing jfr: ", err.Error())
		}

		// parse json file
		// Open our jsonFile
		jsonFile, err := os.Open("perfrt.json")
		// defer the closing of our jsonFile so that we can parse it later on

		// if we os.Open returns an error then handle it
		if err != nil {
			log.Println("-> Error opening jfr json file: ", err.Error())
		} else {
			defer jsonFile.Close()
			// fmt.Println("Successfully Opened perfrt.json")

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
									RunID:     measurementID,
									Run:       models.Run{},
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
									RunID:         measurementID,
									Run:           models.Run{},
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
									RunID:       measurementID,
									Run:         models.Run{},
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
									RunID:     measurementID,
									Run:       models.Run{},
									StartTime: &t,
									ThreadEnd: *threadEnd,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ThreadSleep":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							var duration float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.ThreadSleep: cannot parse duration value. ", err)
								}
							}

							// "time": "PT30S"
							var tsTime float64
							if event.Values["time"] != nil {
								auxTime := strings.ReplaceAll(event.Values["time"].(string), "PT", "")
								auxTime = strings.ReplaceAll(auxTime, "S", "")
								tsTime, err = strconv.ParseFloat(auxTime, 64)
								if err != nil {
									log.Println("ERROR in jdk.ThreadSleep: cannot parse time value. ", err)
								}
							}

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
									RunID:       measurementID,
									Run:         models.Run{},
									StartTime:   &t,
									ThreadSleep: *threadSleep,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ThreadPark":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							parkedClass := getClass(event, "parkedClass")
							var duration float64
							var timeOut float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.ThreadPark: cannot parse duration value. ", err)
								}
							}
							if event.Values["timeout"] != nil {
								//"timeout": "PT0.05S",
								auxTO := strings.ReplaceAll(event.Values["timeout"].(string), "PT", "")
								auxTO = strings.ReplaceAll(auxTO, "S", "")
								timeOut, err = strconv.ParseFloat(auxTO, 64)
								if err != nil {
									log.Println("ERROR in jdk.ThreadPark: cannot parse timeOut value. ", err)
								}
							}

							var until time.Time
							var errUntil error
							if event.Values["until"] != nil {
								// event.Values["until"].(float64),"until": "-999999999-01-01T00:00+18:00",

								layout := "2006-01-02T15:04:05.000000000-07:00"
								untilStr, okUntil := event.Values["until"].(string)
								if okUntil {
									until, errUntil = time.Parse(layout, untilStr)

									if errUntil != nil {
										log.Println("ERROR parsing ThreadPark.until timestamp: ", err)
									}

								}
							}
							threadPark := &models.ThreadPark{
								ThreadParkDuration:     duration,
								ThreadParkOsName:       osName,
								ThreadParkOsThreadId:   osThreadId,
								ThreadParkJavaName:     javaName,
								ThreadParkJavaThreadId: javaThreadId,
								ThreadParkParkedClass:  parkedClass,
								ThreadParkTimeout:      timeOut,
								ThreadParkUntil:        &until,
							}
							if val, ok := jfrMap[t]; ok {
								val.ThreadPark = *threadPark
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID:      measurementID,
									Run:        models.Run{},
									StartTime:  &t,
									ThreadPark: *threadPark,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaErrorThrow":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							thrownClass := getClass(event, "thrownClass")
							var duration float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.JavaErrorThrow: cannot parse duration value. ", err)
								}
							}
							javaErrorThrow := &models.JavaErrorThrow{
								JavaErrorThrowDuration:     duration,
								JavaErrorThrowOsName:       osName,
								JavaErrorThrowOsThreadId:   osThreadId,
								JavaErrorThrowJavaName:     javaName,
								JavaErrorThrowJavaThreadId: javaThreadId,
								JavaErrorThrowMessage:      event.Values["message"].(string),
								JavaErrorThrowThrownClass:  thrownClass,
							}
							if val, ok := jfrMap[t]; ok {
								val.JavaErrorThrow = *javaErrorThrow
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID:          measurementID,
									Run:            models.Run{},
									StartTime:      &t,
									JavaErrorThrow: *javaErrorThrow,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaExceptionThrow":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							thrownClass := getClass(event, "thrownClass")
							var duration float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.JavaExceptionThrow: cannot parse duration value. ", err)
								}
							}
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
									RunID:              measurementID,
									Run:                models.Run{},
									StartTime:          &t,
									JavaExceptionThrow: *javaExceptionThrow,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaMonitorEnter":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							monitorClass := getClass(event, "monitorClass")
							var duration float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.JavaMonitorEnter: cannot parse duration value. ", err)
								}
							}
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
									RunID:            measurementID,
									Run:              models.Run{},
									StartTime:        &t,
									JavaMonitorEnter: *javaMonitorEnter,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.JavaMonitorWait":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							monitorClass := getClass(event, "monitorClass")

							var duration float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.JavaMonitorWait: cannot parse duration value. ", err)
								}
							}

							var timeOut float64
							if event.Values["timeOut"] != nil {
								auxtimeOut := strings.ReplaceAll(event.Values["timeOut"].(string), "PT", "")
								auxtimeOut = strings.ReplaceAll(auxtimeOut, "S", "")
								timeOut, err = strconv.ParseFloat(auxtimeOut, 64)
								if err != nil {
									log.Println("ERROR in jdk.JavaMonitorWait: cannot parse timeOut value. ", err)
								}
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
									RunID:           measurementID,
									Run:             models.Run{},
									StartTime:       &t,
									JavaMonitorWait: *javaMonitorWait,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.OldObjectSample":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							objectType := getObjectType(event)
							var duration float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.OldObjectSample: cannot parse duration value. ", err)
								}
							}

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
							var at time.Time
							atStr := event.Values["allocationTime"].(string)
							at, err = time.Parse(layout, atStr)
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
									RunID:           measurementID,
									Run:             models.Run{},
									StartTime:       &t,
									OldObjectSample: *oldObjectSample,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.ClassLoaderStatistics":
							name, parent := getClassLoader(event)
							classLoaderStatistics := &models.ClassLoaderStatistics{
								ClassLoader:         name,
								ParentClassLoader:   parent,
								ClassLoaderData:     event.Values["classLoaderData"].(float64),
								ClassCount:          event.Values["classCount"].(float64),
								ChunkSize:           event.Values["chunkSize"].(float64),
								BlockSize:           event.Values["blockSize"].(float64),
								AnonymousClassCount: event.Values["anonymousClassCount"].(float64),
								AnonymousChunkSize:  event.Values["anonymousChunkSize"].(float64),
								AnonymousBlockSize:  event.Values["anonymousBlockSize"].(float64),
							}
							if val, ok := jfrMap[t]; ok {
								val.ClassLoaderStatistics = *classLoaderStatistics
								jfrMap[t] = val
							} else {
								jvm := &models.Jvm{
									RunID:                 measurementID,
									Run:                   models.Run{},
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
									RunID:                     measurementID,
									Run:                       models.Run{},
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
									RunID:                       measurementID,
									Run:                         models.Run{},
									StartTime:                   &t,
									ObjectAllocationOutsideTLAB: *objectAllocationOutsideTLAB,
								}
								jfrMap[t] = *jvm

							}

						case "jdk.GCPhasePause":
							osName, osThreadId, javaName, javaThreadId := getEventThread(event)
							var duration float64
							if event.Values["duration"] != nil {
								auxDur := strings.ReplaceAll(event.Values["duration"].(string), "PT", "")
								auxDur = strings.ReplaceAll(auxDur, "S", "")
								duration, err = strconv.ParseFloat(auxDur, 64)
								if err != nil {
									log.Println("ERROR in jdk.GCPhasePause: cannot parse duration value. ", err)
								}
							}
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
									RunID:        measurementID,
									Run:          models.Run{},
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
			for _, jvm := range jfrMap {
				_, err = models.CreateJvm(db, &jvm)
				if err != nil {
					log.Printf("Error saving jvm: %s\n", err.Error())
				}
			}

		}

	} else {
		log.Println("!!!!!!! JFR file not found!!!")
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
