package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	// "github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"gorm.io/gorm"
)

type Resource struct {
	Model
	RunID      uint `gorm:"index"`
	Run        Run
	Timestamp  time.Time
	CpuPercent float64
	MemPercent float32
	process.MemoryInfoStat
	process.IOCountersStat
	process.PageFaultsStat
	load.AvgStat
	mem.VirtualMemoryStat
	SwapMemoryStat

	// CPUInfo  []CPUInfo
	// CPUTimes []CPUTimes

	// DiskIOCounters map[string]disk.IOCountersStat
	// NetIOCounters []net.IOCountersStat
}

type SwapMemoryStat struct {
	SwapTotal       uint64
	SwapUsed        uint64
	SwapFree        uint64
	SwapUsedPercent float64
	Sin             uint64
	Sout            uint64
	PgIn            uint64
	PgOut           uint64
	PgFault         uint64
	PgMajFaults     uint64
}

type CPUInfo struct {
	Model
	ResourceID uint
	Resource   Resource
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	CPUModel   string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

type CPUTimes struct {
	Model
	ResourceID uint
	Resource   Resource
	// CPUID      uint
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}

type DiskIOCounters struct {
	Model
	ResourceID       uint
	Resource         Resource
	Device           string
	ReadCount        uint64 `json:"readCount"`
	MergedReadCount  uint64 `json:"mergedReadCount"`
	WriteCount       uint64 `json:"writeCount"`
	MergedWriteCount uint64 `json:"mergedWriteCount"`
	ReadBytes        uint64 `json:"readBytes"`
	WriteBytes       uint64 `json:"writeBytes"`
	ReadTime         uint64 `json:"readTime"`
	WriteTime        uint64 `json:"writeTime"`
	IopsInProgress   uint64 `json:"iopsInProgress"`
	IoTime           uint64 `json:"ioTime"`
	WeightedIO       uint64 `json:"weightedIO"`
	Name             string `json:"name"`
	SerialNumber     string `json:"serialNumber"`
	Label            string `json:"label"`
}

type NetIOCounters struct {
	Model
	ResourceID  uint
	Resource    Resource
	NICID       uint
	Name        string `json:"name"`        // interface name
	BytesSent   uint64 `json:"bytesSent"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv"` // number of packets received
	Errin       uint64 `json:"errin"`       // total number of errors while receiving
	Errout      uint64 `json:"errout"`      // total number of errors while sending
	Dropin      uint64 `json:"dropin"`      // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout"`     // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Fifoin      uint64 `json:"fifoin"`      // total number of FIFO buffers errors while receiving
	Fifoout     uint64 `json:"fifoout"`     // total number of FIFO buffers errors while sending

}

func (r *Resource) TableName() string {
	return "resources"
}

func (r *CPUTimes) TableName() string {
	return "cputimes"
}

func (r *DiskIOCounters) TableName() string {
	return "diskiocounters"
}

func (r *NetIOCounters) TableName() string {
	return "netiocounters"
}
func CreateResource(db *gorm.DB, mr *Resource) (uint, error) {
	err := db.Create(mr).Error
	if err != nil {
		return 0, err

	}
	return mr.ID, nil
}

func CreateCPUTimes(db *gorm.DB, ct *CPUTimes) (uint, error) {
	err := db.Create(ct).Error
	if err != nil {
		return 0, err
	}
	return ct.ID, nil
}

func CreateDiskIOCounters(db *gorm.DB, dic *DiskIOCounters) (uint, error) {
	err := db.Create(dic).Error
	if err != nil {
		return 0, err
	}
	return dic.ID, nil
}

func CreateNetIOCounters(db *gorm.DB, nic *NetIOCounters) (uint, error) {
	err := db.Create(nic).Error
	if err != nil {
		return 0, err
	}
	return nic.ID, nil
}

// type TestResources struct {
// 	Model
// 	TestID uint
// 	Test   TestCase
// 	Type   string
// 	Resources
// }

// func (r *TestResources) TableName() string {
// 	return "testresources"
// }

// func CreateTestResources(db *gorm.DB, tr *TestResources) (uint, error) {
// 	err := db.Create(tr).Error
// 	if err != nil {
// 		return 0, err
// 	}
// 	return tr.ID, nil
// }
