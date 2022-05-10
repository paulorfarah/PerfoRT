package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"perfrt/models"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/mitchellh/mapstructure"
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
			result.ForEach(func(key, value gjson.Result) bool {
				fmt.Println(value.String())
				fmt.Println("-----")
				return true // keep iterating
			})
			// gabs
			// jsonParsed, err := gabs.ParseJSON(jsonJFR)
			// if err != nil {
			// 	fmt.Println("Error parsing JFR json file: ", err)
			// }
			// events, err := jsonParsed.S("recording", "events").ChildrenMap()
			// if err != nil {
			// 	fmt.Println("ERROR parsing JFR json: ", err.Error())
			// }
			// fmt.Printf("%v\n", events)
			jfrMap := make(map[time.Time]models.Jvm)
			// for _, eventMap := range events {
			result.ForEach(func(key, eventMap gjson.Result) bool {
				// fmt.Printf("%v\n", eventMap.Data())
				var event Event
				err := mapstructure.Decode(eventMap, &event)
				if err != nil {
					panic(err)
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
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						tCpuLoad := &models.ThreadCPULoad{
							ThreadCPULoadOsName:       osName,
							ThreadCPULoadOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							ThreadCPULoadJavaName:     event.Values["eventThread.javaName"].(string),
							ThreadCPULoadJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
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
					}

				case "jdk.ThreadStart":
					osName, ok := event.Values["eventThread.osName"].(string)
					if ok {

						threadStart := &models.ThreadStart{
							ThreadStartOsName:                   osName,
							ThreadStartOsThreadId:               event.Values["eventThread.osThreadId"].(int),
							ThreadStartJavaName:                 event.Values["eventThread.javaName"].(string),
							ThreadStartJavaThreadId:             event.Values["eventThread.javaThreadId"].(int),
							ThreadStartParentThreadosName:       event.Values["parentThread.osName"].(string),
							ThreadStartParentThreadOsThreadId:   event.Values["parentThread.osThreadId"].(int),
							ThreadStartParentThreadJavaName:     event.Values["parentThread.javaName"].(string),
							ThreadStartParentThreadJavaThreadId: event.Values["parentThread.javaThreadId"].(int),
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
					}

				case "jdk.ThreadEnd":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						threadEnd := &models.ThreadEnd{
							ThreadEndOsName:       osName,
							ThreadEndOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							ThreadEndJavaName:     event.Values["eventThread.javaName"].(string),
							ThreadEndJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
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
					}

				case "jdk.ThreadSleep":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						threadSleep := &models.ThreadSleep{
							ThreadSleepDuration:     event.Values["duration"].(float64),
							ThreadSleepOsName:       osName,
							ThreadSleepOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							ThreadSleepJavaName:     event.Values["eventThread.javaName"].(string),
							ThreadSleepJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
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
					}

				case "jdk.ThreadPark":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						threadPark := &models.ThreadPark{
							ThreadParkDuration:     event.Values["duration"].(float64),
							ThreadParkOsName:       osName,
							ThreadParkOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							ThreadParkJavaName:     event.Values["eventThread.javaName"].(string),
							ThreadParkJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
							ThreadParkParkedClass:  event.Values["parkedClass.name"].(string),
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
					}

				case "jdk.JavaErrorThrow":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						javaErrorThrow := &models.JavaErrorThrow{
							JavaErrorThrowDuration:     event.Values["duration"].(float64),
							JavaErrorThrowOsName:       osName,
							JavaErrorThrowOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							JavaErrorThrowJavaName:     event.Values["eventThread.javaName"].(string),
							JavaErrorThrowJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
							JavaErrorThrowMessage:      event.Values["message"].(string),
							JavaErrorThrowThrownClass:  event.Values["thrownClass.name"].(string),
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
					}

				case "jdk.JavaExceptionThrow":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						javaExceptionThrow := &models.JavaExceptionThrow{
							JavaExceptionThrowDuration:     event.Values["duration"].(float64),
							JavaExceptionThrowOsName:       osName,
							JavaExceptionThrowOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							JavaExceptionThrowJavaName:     event.Values["eventThread.javaName"].(string),
							JavaExceptionThrowJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
							JavaExceptionThrowMessage:      event.Values["message"].(string),
							JavaExceptionThrowThrownClass:  event.Values["thrownClass.name"].(string),
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
					}

				case "jdk.JavaMonitorEnter":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						javaMonitorEnter := &models.JavaMonitorEnter{
							JavaMonitorEnterDuration:     event.Values["duration"].(float64),
							JavaMonitorEnterOsName:       osName,
							JavaMonitorEnterOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							JavaMonitorEnterJavaName:     event.Values["eventThread.javaName"].(string),
							JavaMonitorEnterJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
							JavaMonitorEnterMonitorClass: event.Values["monitorClass.name"].(string),
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
					}

				case "jdk.JavaMonitorWait":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						javaMonitorWait := &models.JavaMonitorWait{
							JavaMonitorWaitDuration:     event.Values["duration"].(float64),
							JavaMonitorWaitOsName:       osName,
							JavaMonitorWaitOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							JavaMonitorWaitJavaName:     event.Values["eventThread.javaName"].(string),
							JavaMonitorWaitJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
							JavaMonitorWaitMonitorClass: event.Values["monitorClass.name"].(string),
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
					}
				case "jdk.OldObjectSample":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						oldObjectSample := &models.OldObjectSample{
							OldObjectSampleDuration:           event.Values["duration"].(float64),
							OldObjectSampleOsName:             osName,
							OldObjectSampleOsThreadId:         event.Values["eventThread.osThreadId"].(int),
							OldObjectSampleJavaName:           event.Values["eventThread.javaName"].(string),
							OldObjectSampleJavaThreadId:       event.Values["eventThread.javaThreadId"].(int),
							OldObjectSampleAllocationTime:     event.Values["allocationTime"].(float64),
							OldObjectSampleLastKnownHeapUsage: event.Values["lastKnownHeapUsage"].(float64),
							OldObjectSampleObject:             event.Values["object.type.name"].(string),
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
					}

				case "jdk.ClassLoaderStatistics":
					e := fmt.Sprintf("%s", event)
					fmt.Println(e)
					jsonParsed, err := gabs.ParseJSON([]byte(e))
					name, nameOk := jsonParsed.Path("classLoader.name").Data().(string)
					// name, nameOk := event.Values["classLoader.name"].(string)
					fmt.Println("ClassLoader.name: ", name)
					if nameOk {

						parent, parentOk := event.Values["parentClassLoader.name"].(string)
						if !parentOk {
							parent = ""
						}
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
					}

				case "jdk.ObjectAllocationInNewTLAB":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {

						objectAllocationInNewTLAB := &models.ObjectAllocationInNewTLAB{

							ObjectAllocationInNewTLABOsName:         osName,
							ObjectAllocationInNewTLABOsThreadId:     event.Values["eventThread.osThreadId"].(int),
							ObjectAllocationInNewTLABJavaName:       event.Values["eventThread.javaName"].(string),
							ObjectAllocationInNewTLABJavaThreadId:   event.Values["eventThread.javaThreadId"].(int),
							ObjectAllocationInNewTLABObjectClass:    event.Values["objectClass.name"].(string),
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
					}

				case "jdk.ObjectAllocationOutsideTLAB":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						objectAllocationOutsideTLAB := &models.ObjectAllocationOutsideTLAB{
							ObjectAllocationOutsideTLABOsName:         osName,
							ObjectAllocationOutsideTLABOsThreadId:     event.Values["eventThread.osThreadId"].(int),
							ObjectAllocationOutsideTLABJavaName:       event.Values["eventThread.javaName"].(string),
							ObjectAllocationOutsideTLABJavaThreadId:   event.Values["eventThread.javaThreadId"].(int),
							ObjectAllocationOutsideTLABObjectClass:    event.Values["objectClass.name"].(string),
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
					}

				case "jdk.GCPhasePause":
					osName, osNameOk := event.Values["eventThread.osName"].(string)
					if osNameOk {
						gcPhasePause := &models.GCPhasePause{
							GCPhasePauseDuration:     event.Values["duration"].(float64),
							GCPhasePauseOsName:       osName,
							GCPhasePauseOsThreadId:   event.Values["eventThread.osThreadId"].(int),
							GCPhasePauseJavaName:     event.Values["eventThread.javaName"].(string),
							GCPhasePauseJavaThreadId: event.Values["eventThread.javaThreadId"].(int),
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
