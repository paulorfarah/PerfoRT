package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"perfrt/models"
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
	fmt.Println("****************** SaveJFRMetrics")
	// generate json
	if _, err := os.Stat("perfrt.jfr"); err == nil {
		// file perfrt.jfr exists
		fmt.Println("jfr print --json perfrt.jfr > perfrt.json")
		cmd := exec.Command("bash", "-c", "jfr print --json perfrt.jfr > perfrt.json")
		err := cmd.Run()
		if err != nil {
			fmt.Println("->->->->->->->->->->->->-> error executing jfr: ", err.Error())
			// log.Fatal(err)
		}

		// parse json file
		// Open our jsonFile
		jsonFile, err := os.Open("perfrt.json")
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println("->->->->->->->->->->->->-> error opening jfr json file: ", err.Error())
		} else {
			fmt.Println("Successfully Opened perfrt.json")

			// read our opened jsonFile as a byte array.
			jsonJFR, _ := ioutil.ReadAll(jsonFile)

			// GJson
			result := gjson.Get(string(jsonJFR), "recording.events")
			jfrMap := make(map[time.Time]models.Jvm)
			result.ForEach(func(key, eventMap gjson.Result) bool {
				var event Event

				fmt.Println(eventMap)
				fmt.Println("-----")
				var result map[string]interface{}
				jsonFile := fmt.Sprintf("%s", eventMap)
				if err := json.Unmarshal([]byte(jsonFile), &result); err != nil {
					panic(err)
				}
				values := result["values"]
				if v, ok := values.(map[string]interface{}); ok {
					// fmt.Printf(" %s", v["startTime"])
					event.Type = result["type"].(string)
					event.Values = v

				} else {
					fmt.Printf("record not a map[string]interface{}: %v\n", values)
					panic("record not a map[string]interface{}: ")
				}

				// fmt.Printf("%#v\n", event.Values)
				layout := "2006-01-02T15:04:05.000000000-07:00"
				str := event.Values["startTime"].(string)
				t, err := time.Parse(layout, str)

				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("->->->->->->->->->->->->-> Event type: "+event.Type+" ", t)
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
							StartTime: t,
							CPULoad:   *cpuLoad,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
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
							StartTime:     t,
							ThreadCPULoad: *tCpuLoad,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
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
							StartTime:   t,
							ThreadStart: *threadStart,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
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
							StartTime: t,
							ThreadEnd: *threadEnd,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.ThreadSleep":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					threadSleep := &models.ThreadSleep{
						ThreadSleepDuration:     event.Values["duration"].(float64),
						ThreadSleepOsName:       osName,
						ThreadSleepOsThreadId:   osThreadId,
						ThreadSleepJavaName:     javaName,
						ThreadSleepJavaThreadId: javaThreadId,
						ThreadSleepTime:         event.Values["time"].(float64),
					}
					if val, ok := jfrMap[t]; ok {
						val.ThreadSleep = *threadSleep
						jfrMap[t] = val
					} else {
						jvm := &models.Jvm{
							RunID:       measurementID,
							Run:         models.Run{},
							StartTime:   t,
							ThreadSleep: *threadSleep,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.ThreadPark":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					parkedClass := getClass(event, "parkedClass")
					threadPark := &models.ThreadPark{
						ThreadParkDuration:     event.Values["duration"].(float64),
						ThreadParkOsName:       osName,
						ThreadParkOsThreadId:   osThreadId,
						ThreadParkJavaName:     javaName,
						ThreadParkJavaThreadId: javaThreadId,
						ThreadParkParkedClass:  parkedClass,
						ThreadParkTimeout:      event.Values["timeout"].(float64),
						ThreadParkUntil:        event.Values["until"].(float64),
					}
					if val, ok := jfrMap[t]; ok {
						val.ThreadPark = *threadPark
						jfrMap[t] = val
					} else {
						jvm := &models.Jvm{
							RunID:      measurementID,
							Run:        models.Run{},
							StartTime:  t,
							ThreadPark: *threadPark,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.JavaErrorThrow":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					thrownClass := getClass(event, "thrownClass")
					javaErrorThrow := &models.JavaErrorThrow{
						JavaErrorThrowDuration:     event.Values["duration"].(float64),
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
							StartTime:      t,
							JavaErrorThrow: *javaErrorThrow,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.JavaExceptionThrow":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					thrownClass := getClass(event, "thrownClass")
					javaExceptionThrow := &models.JavaExceptionThrow{
						JavaExceptionThrowDuration:     event.Values["duration"].(float64),
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
							StartTime:          t,
							JavaExceptionThrow: *javaExceptionThrow,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.JavaMonitorEnter":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					monitorClass := getClass(event, "monitorClass")
					javaMonitorEnter := &models.JavaMonitorEnter{
						JavaMonitorEnterDuration:     event.Values["duration"].(float64),
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
							StartTime:        t,
							JavaMonitorEnter: *javaMonitorEnter,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.JavaMonitorWait":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					monitorClass := getClass(event, "monitorClass")
					javaMonitorWait := &models.JavaMonitorWait{
						JavaMonitorWaitDuration:     event.Values["duration"].(float64),
						JavaMonitorWaitOsName:       osName,
						JavaMonitorWaitOsThreadId:   osThreadId,
						JavaMonitorWaitJavaName:     javaName,
						JavaMonitorWaitJavaThreadId: javaThreadId,
						JavaMonitorWaitMonitorClass: monitorClass,
						JavaMonitorWaitTimeout:      event.Values["timeOut"].(float64),
						JavaMonitorWaitTimedOut:     event.Values["timedOut"].(bool),
					}
					if val, ok := jfrMap[t]; ok {
						val.JavaMonitorWait = *javaMonitorWait
						jfrMap[t] = val
					} else {
						jvm := &models.Jvm{
							RunID:           measurementID,
							Run:             models.Run{},
							StartTime:       t,
							JavaMonitorWait: *javaMonitorWait,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.OldObjectSample":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					objectType := getObjectType(event)
					oldObjectSample := &models.OldObjectSample{
						OldObjectSampleDuration:           event.Values["duration"].(float64),
						OldObjectSampleOsName:             osName,
						OldObjectSampleOsThreadId:         osThreadId,
						OldObjectSampleJavaName:           javaName,
						OldObjectSampleJavaThreadId:       javaThreadId,
						OldObjectSampleAllocationTime:     event.Values["allocationTime"].(float64),
						OldObjectSampleLastKnownHeapUsage: event.Values["lastKnownHeapUsage"].(float64),
						OldObjectSampleObject:             objectType,
						OldObjectSampleArrayElements:      event.Values["arrayElements"].(int),
					}
					if val, ok := jfrMap[t]; ok {
						val.OldObjectSample = *oldObjectSample
						jfrMap[t] = val
					} else {
						jvm := &models.Jvm{
							RunID:           measurementID,
							Run:             models.Run{},
							StartTime:       t,
							OldObjectSample: *oldObjectSample,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
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
							StartTime:             t,
							ClassLoaderStatistics: *classLoaderStatistics,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}
					// }

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
							StartTime:                 t,
							ObjectAllocationInNewTLAB: *objectAllocationInNewTLAB,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
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
							StartTime:                   t,
							ObjectAllocationOutsideTLAB: *objectAllocationOutsideTLAB,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.GCPhasePause":
					osName, osThreadId, javaName, javaThreadId := getEventThread(event)
					gcPhasePause := &models.GCPhasePause{
						GCPhasePauseDuration:     event.Values["duration"].(float64),
						GCPhasePauseOsName:       osName,
						GCPhasePauseOsThreadId:   osThreadId,
						GCPhasePauseJavaName:     javaName,
						GCPhasePauseJavaThreadId: javaThreadId,
						GcId:                     event.Values["gcId"].(int),
						GCPhasePauseName:         event.Values["name"].(string),
					}
					if val, ok := jfrMap[t]; ok {
						val.GCPhasePause = *gcPhasePause
						jfrMap[t] = val
					} else {
						jvm := &models.Jvm{
							RunID:        measurementID,
							Run:          models.Run{},
							StartTime:    t,
							GCPhasePause: *gcPhasePause,
						}
						jfrMap[t] = *jvm

					}
					// _, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
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
				models.CreateJvm(db, &jvm)
			}

		}

	} else {
		fmt.Println("!!!!!!! JFR file not found!!!")
	}

}

func getClassLoader(event Event) (string, string) {
	var name string
	classloader := event.Values["classLoader"]
	if n, ok := classloader.(map[string]interface{}); ok {
		fmt.Println(n)
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

func getEventThread(event Event) (string, int, string, int) {
	var osName string
	eventThread := event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osName"] != nil {
			osName = e["osName"].(string)
		}
	}

	var osThreadId int
	eventThread = event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osThreadId"] != nil {
			osThreadId = e["osThreadId"].(int)
		}
	}

	var javaName string
	eventThread = event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaName"] != nil {
			javaName = e["javaName"].(string)
		}
	}

	var javaThreadId int
	eventThread = event.Values["eventThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaThreadId"] != nil {
			osThreadId = e["javaThreadId"].(int)
		}
	}
	return osName, osThreadId, javaName, javaThreadId
}

func getEventParentThread(event Event) (string, int, string, int) {
	var osName string
	eventThread := event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osName"] != nil {
			osName = e["osName"].(string)
		}
	}

	var osThreadId int
	eventThread = event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["osThreadId"] != nil {
			osThreadId = e["osThreadId"].(int)
		}
	}

	var javaName string
	eventThread = event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaName"] != nil {
			javaName = e["javaName"].(string)
		}
	}

	var javaThreadId int
	eventThread = event.Values["parentThread"]
	if e, ok := eventThread.(map[string]interface{}); ok {
		if e["javaThreadId"] != nil {
			osThreadId = e["javaThreadId"].(int)
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
