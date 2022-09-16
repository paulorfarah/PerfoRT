package main

import (
	"perfrt/models"
	"time"

	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// func MonitorResources(db *gorm.DB, measurementID uint) {
// 	go cpuMonitor(db, measurementID)
// 	// go loadMonitor()
// 	// go memMonitor()
// }
// type PerfMetrics struct {
// 	Timestamp     time.Time
// 	CpuPercent    float64
// 	MemoryPercent float32
// 	IOCounters    *process.IOCountersStat
// 	MemoryInfo    *process.MemoryInfoStat
// 	PageFaults    *process.PageFaultsStat
// 	// Load           *load.AvgStat
// 	// CPUTimes       []cpu.TimesStat
// 	// VirtualMemory  *mem.VirtualMemoryStat
// 	// SwapMemory     *models.SwapMemoryStat
// 	// DiskIOCounters map[string]disk.IOCountersStat
// 	// NetIOCounters  []net.IOCountersStat
// }

func MonitorProcess(pid int, runID uint) (models.Resource, error) {

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		// log.Println("Error CPU Percent: ", err.Error())
		return models.Resource{}, err
	}
	cp, err := p.CPUPercent()
	if err != nil {
		return models.Resource{}, err
	}
	m, err := p.MemoryPercent()
	if err != nil {
		return models.Resource{}, err
	}
	io, err := p.IOCounters()
	if err != nil {
		return models.Resource{}, err
	}

	mi, err := p.MemoryInfo()
	if err != nil {
		return models.Resource{}, err
	}
	pf, err := p.PageFaults()
	if err != nil {
		return models.Resource{}, err
	}

	l, _ := load.Avg()

	vm, _ := mem.VirtualMemory()

	swap, _ := mem.SwapMemory()
	sm := &models.SwapMemoryStat{
		SwapTotal:       swap.Total,
		SwapUsed:        swap.Used,
		SwapFree:        swap.Free,
		SwapUsedPercent: swap.UsedPercent,
		Sin:             swap.Sin,
		Sout:            swap.Sout,
		PgIn:            swap.PgIn,
		PgOut:           swap.PgOut,
		PgFault:         swap.PgFault,
		PgMajFaults:     swap.PgMajFault,
	}

	// var cpusInfo []models.CPUInfo
	// cpuInfo, _ := cpu.Info()

	// for _, cpuI := range cpuInfo {
	// 	c := models.CPUInfo{
	// 		CPU:        cpuI.CPU,
	// 		VendorID:   cpuI.VendorID,
	// 		Family:     cpuI.Family,
	// 		CPUModel:   cpuI.Model,
	// 		Stepping:   cpuI.Stepping,
	// 		PhysicalID: cpuI.PhysicalID,
	// 		CoreID:     cpuI.CoreID,
	// 		Cores:      cpuI.Cores,
	// 		ModelName:  cpuI.ModelName,
	// 		Mhz:        cpuI.Mhz,
	// 		CacheSize:  cpuI.CacheSize,
	// 		Flags:      cpuI.Flags,
	// 		Microcode:  cpuI.Microcode,
	// 	}
	// 	cpusInfo = append(cpusInfo, c)
	// }

	// var cpusTimes []models.CPUTimes
	// cpuTimes, _ := cpu.Times(false)
	// for _, cpuT := range cpuTimes {
	// 	c := models.CPUTimes{
	// 		CPU:       cpuT.CPU,
	// 		User:      cpuT.User,
	// 		System:    cpuT.System,
	// 		Idle:      cpuT.Idle,
	// 		Nice:      cpuT.Nice,
	// 		Iowait:    cpuT.Iowait,
	// 		Irq:       cpuT.Irq,
	// 		Softirq:   cpuT.Softirq,
	// 		Steal:     cpuT.Steal,
	// 		Guest:     cpuT.Guest,
	// 		GuestNice: cpuT.GuestNice,
	// 	}
	// 	cpusTimes = append(cpusTimes, c)
	// }

	resource := &models.Resource{
		RunID:             runID,
		Timestamp:         time.Now(),
		CpuPercent:        cp,
		MemPercent:        m,
		MemoryInfoStat:    *mi,
		IOCountersStat:    *io,
		PageFaultsStat:    *pf,
		AvgStat:           *l,
		VirtualMemoryStat: *vm,
		SwapMemoryStat:    *sm,
		// CPUInfo:           cpusInfo,
		// CPUTimes:          cpusTimes,
	}

	return *resource, nil
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

// func saveMetrics(db *gorm.DB, measurementID uint, perfMetric PerfMetrics) {
// 	resource := &models.Resource{
// 		RunID:          measurementID,
// 		Timestamp:      perfMetric.Timestamp,
// 		CpuPercent:     perfMetric.CpuPercent,
// 		MemPercent:     perfMetric.MemoryPercent,
// 		MemoryInfoStat: *perfMetric.MemoryInfo,
// 		IOCountersStat: *perfMetric.IOCounters,
// 		PageFaultsStat: *perfMetric.PageFaults,
// 		// AvgStat:           *perfMetric.Load,
// 		// VirtualMemoryStat: *perfMetric.VirtualMemory,
// 		// SwapMemoryStat:    *perfMetric.SwapMemory,
// 		// DiskIOCounters:    perfMetric.DiskIOCounters,
// 		// NetIOCounters:     perfMetric.NetIOCounters,

// 	}
// 	_, err := models.CreateResource(db, resource)
// 	if err != nil {
// 		fmt.Printf("Error saving resource: %s\n", err.Error())
// 	}

// 	// for _, cpuTime := range perfMetric.CPUTimes {
// 	// 	models.CreateCPUTimes(db, &models.CPUTimes{
// 	// 		ResourceID: resource.ID,
// 	// 		CPU:        cpuTime.CPU,
// 	// 		User:       cpuTime.User,
// 	// 		System:     cpuTime.System,
// 	// 		Idle:       cpuTime.Idle,
// 	// 		Nice:       cpuTime.Nice,
// 	// 		Iowait:     cpuTime.Iowait,
// 	// 		Irq:        cpuTime.Irq,
// 	// 		Softirq:    cpuTime.Softirq,
// 	// 		Steal:      cpuTime.Steal,
// 	// 		Guest:      cpuTime.Guest,
// 	// 		GuestNice:  cpuTime.GuestNice,
// 	// 	})
// 	// }

// 	// for i, diskIOCounter := range perfMetric.DiskIOCounters {
// 	// 	models.CreateDiskIOCounters(db, &models.DiskIOCounters{
// 	// 		ResourceID:       resource.ID,
// 	// 		Device:           i,
// 	// 		ReadCount:        diskIOCounter.ReadCount,
// 	// 		MergedReadCount:  diskIOCounter.MergedReadCount,
// 	// 		WriteCount:       diskIOCounter.WriteCount,
// 	// 		MergedWriteCount: diskIOCounter.MergedWriteCount,
// 	// 		ReadBytes:        diskIOCounter.ReadBytes,
// 	// 		WriteBytes:       diskIOCounter.WriteBytes,
// 	// 		ReadTime:         diskIOCounter.ReadTime,
// 	// 		WriteTime:        diskIOCounter.WriteTime,
// 	// 		IopsInProgress:   diskIOCounter.IopsInProgress,
// 	// 		IoTime:           diskIOCounter.IoTime,
// 	// 		WeightedIO:       diskIOCounter.WeightedIO,
// 	// 		Name:             diskIOCounter.Name,
// 	// 		SerialNumber:     diskIOCounter.SerialNumber,
// 	// 		Label:            diskIOCounter.Label,
// 	// 	})
// 	// }

// 	// for i, netIOCounter := range perfMetric.NetIOCounters {
// 	// 	models.CreateNetIOCounters(db, &models.NetIOCounters{
// 	// 		ResourceID:  resource.ID,
// 	// 		NICID:       uint(i),
// 	// 		Name:        netIOCounter.Name,
// 	// 		BytesSent:   netIOCounter.BytesSent,
// 	// 		BytesRecv:   netIOCounter.BytesRecv,
// 	// 		PacketsSent: netIOCounter.PacketsSent,
// 	// 		PacketsRecv: netIOCounter.PacketsRecv,
// 	// 		Errin:       netIOCounter.Errin,
// 	// 		Errout:      netIOCounter.Errout,
// 	// 		Dropin:      netIOCounter.Dropin,
// 	// 		Dropout:     netIOCounter.Dropout,
// 	// 		Fifoin:      netIOCounter.Fifoin,
// 	// 		Fifoout:     netIOCounter.Fifoout,
// 	// 	})
// 	// }
// }
