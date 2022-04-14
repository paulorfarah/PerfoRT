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
	fmt.Println("->->->->->->->->->->->->->->->->->->->->->->  JFR")
	// generate json

	if _, err := os.Stat("perfrt.jfr"); err == nil {
		// file perfrt.jfr exists
		fmt.Println("jfr print --json perfrt.jfr>perfrt.json")
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

			jsonParsed, err := gabs.ParseJSON(jsonJFR)
			if err != nil {
				fmt.Println("Error parsing JFR json file: ", err)
			}

			events, err := jsonParsed.S("recording", "events").Children()
			// fmt.Printf("%v\n", events)
			for _, eventMap := range events {
				// fmt.Printf("%v\n", eventMap.Data())
				var event Event
				err := mapstructure.Decode(eventMap.Data(), &event)
				if err != nil {
					panic(err)
				}

				// fmt.Printf("%#v\n", event.Values)
				switch event.Type {
				case "jdk.CPULoad":
					layout := "2006-01-02T15:04:05.000000000-07:00"
					str := event.Values["startTime"].(string)
					t, err := time.Parse(layout, str)

					if err != nil {
						fmt.Println(err)
					}
					cpuLoad := &models.CPULoad{StartTime: t,
						JvmUser:      event.Values["jvmUser"].(float64),
						JvmSystem:    event.Values["jvmSystem"].(float64),
						MachineTotal: event.Values["machineTotal"].(float64),
					}
					// fmt.Printf("%v\n", cpuLoad)
					// save
					jvm := &models.Jvm{
						RunID:   measurementID,
						CPULoad: *cpuLoad,
					}
					_, err = models.CreateJvm(db, jvm)
					if err != nil {
						fmt.Printf("Error saving resource: %s\n", err.Error())
					}

				case "jdk.":
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
			}

		}

	}

}
