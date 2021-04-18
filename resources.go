package main

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

// func MonitorResources(db *gorm.DB, measurementID uint) {
// 	go cpuMonitor(db, measurementID)
// 	// go loadMonitor()
// 	// go memMonitor()
// }
type PerfMetrics struct {
	Cpu float64
	Mem float32
	IO  *process.IOCountersStat
}

func MonitorProcess(pid int) (PerfMetrics, error) {

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("Error CPU Percent: ", err.Error())
	}
	cpu, err := p.CPUPercent()
	if err != nil {
		return PerfMetrics{}, err
	}
	mem, err := p.MemoryPercent()
	if err != nil {
		return PerfMetrics{}, err
	}
	io, err := p.IOCounters()
	if err != nil {
		return PerfMetrics{}, err
	}
	perfMetrics := &PerfMetrics{
		Cpu: cpu,
		Mem: mem,
		IO:  io,
	}
	return *perfMetrics, nil

}

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
