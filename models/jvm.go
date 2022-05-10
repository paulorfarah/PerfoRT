package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type Jvm struct {
	Model
	RunID     uint `gorm:"not null"`
	Run       Run
	StartTime time.Time
	CPULoad
	ThreadCPULoad
	ThreadStart
	ThreadEnd
	ThreadSleep
	ThreadPark
	JavaErrorThrow
	JavaExceptionThrow
	JavaMonitorEnter
	JavaMonitorWait
	OldObjectSample
	LoadedClassCount   int
	UnloadedClassCount int
	ClassLoaderStatistics
	ObjectAllocationInNewTLAB
	ObjectAllocationOutsideTLAB
	GCPhasePause
}

type CPULoad struct {
	// StartTime    time.Time
	JvmUser      float64
	JvmSystem    float64
	MachineTotal float64
}

type ThreadCPULoad struct {
	ThreadCPULoadOsName       string
	ThreadCPULoadOsThreadId   int
	ThreadCPULoadJavaName     string
	ThreadCPULoadJavaThreadId int
	ThreadCPULoadUser         float64
	ThreadCPULoadSystem       float64
}
type ThreadStart struct { //ok
	ThreadStartOsName                   string
	ThreadStartOsThreadId               int
	ThreadStartJavaName                 string
	ThreadStartJavaThreadId             int
	ThreadStartParentThreadosName       string
	ThreadStartParentThreadOsThreadId   int
	ThreadStartParentThreadJavaName     string
	ThreadStartParentThreadJavaThreadId int
}

type ThreadEnd struct { //ok
	ThreadEndOsName       string
	ThreadEndOsThreadId   int
	ThreadEndJavaName     string
	ThreadEndJavaThreadId int
}

type ThreadSleep struct {
	ThreadSleepDuration     float64
	ThreadSleepOsName       string
	ThreadSleepOsThreadId   int
	ThreadSleepJavaName     string
	ThreadSleepJavaThreadId int
	ThreadSleepTime         float64
}

type ThreadPark struct {
	ThreadParkDuration     float64
	ThreadParkOsName       string
	ThreadParkOsThreadId   int
	ThreadParkJavaName     string
	ThreadParkJavaThreadId int
	ThreadParkParkedClass  string
	ThreadParkTimeout      float64
	ThreadParkUntil        float64
}

type JavaErrorThrow struct {
	JavaErrorThrowDuration     float64
	JavaErrorThrowOsName       string
	JavaErrorThrowOsThreadId   int
	JavaErrorThrowJavaName     string
	JavaErrorThrowJavaThreadId int
	JavaErrorThrowMessage      string
	JavaErrorThrowThrownClass  string
}

type JavaExceptionThrow struct {
	JavaExceptionThrowDuration     float64
	JavaExceptionThrowOsName       string
	JavaExceptionThrowOsThreadId   int
	JavaExceptionThrowJavaName     string
	JavaExceptionThrowJavaThreadId int
	JavaExceptionThrowMessage      string
	JavaExceptionThrowThrownClass  string
}

type JavaMonitorEnter struct {
	JavaMonitorEnterDuration     float64
	JavaMonitorEnterOsName       string
	JavaMonitorEnterOsThreadId   int
	JavaMonitorEnterJavaName     string
	JavaMonitorEnterJavaThreadId int
	JavaMonitorEnterMonitorClass string
}

type JavaMonitorWait struct {
	JavaMonitorWaitDuration     float64
	JavaMonitorWaitOsName       string
	JavaMonitorWaitOsThreadId   int
	JavaMonitorWaitJavaName     string
	JavaMonitorWaitJavaThreadId int
	JavaMonitorWaitMonitorClass string
	JavaMonitorWaitTimeout      float64 //	Maximum wait time
	JavaMonitorWaitTimedOut     bool    //Wait has been timed out
}

type OldObjectSample struct {
	OldObjectSampleDuration           float64
	OldObjectSampleOsName             string
	OldObjectSampleOsThreadId         int
	OldObjectSampleJavaName           string
	OldObjectSampleJavaThreadId       int
	OldObjectSampleAllocationTime     float64
	OldObjectSampleLastKnownHeapUsage float64
	OldObjectSampleObject             string
	OldObjectSampleArrayElements      int
}

type ClassLoaderStatistics struct {
	ClassLoader         string
	ParentClassLoader   string
	ClassLoaderData     float64
	ClassCount          float64
	ChunkSize           float64
	BlockSize           float64
	AnonymousClassCount float64
	AnonymousChunkSize  float64
	AnonymousBlockSize  float64
}

type ObjectAllocationInNewTLAB struct {
	ObjectAllocationInNewTLABOsName         string
	ObjectAllocationInNewTLABOsThreadId     int
	ObjectAllocationInNewTLABJavaName       string
	ObjectAllocationInNewTLABJavaThreadId   int
	ObjectAllocationInNewTLABObjectClass    string // Class of allocated object
	ObjectAllocationInNewTLABAllocationSize float64
	ObjectAllocationInNewTLABTlabSize       float64
}

type ObjectAllocationOutsideTLAB struct {
	ObjectAllocationOutsideTLABOsName         string
	ObjectAllocationOutsideTLABOsThreadId     int
	ObjectAllocationOutsideTLABJavaName       string
	ObjectAllocationOutsideTLABJavaThreadId   int
	ObjectAllocationOutsideTLABObjectClass    string // Class of allocated object
	ObjectAllocationOutsideTLABAllocationSize float64
}

type GCPhasePause struct {
	GCPhasePauseOsName       string
	GCPhasePauseOsThreadId   int
	GCPhasePauseJavaName     string
	GCPhasePauseJavaThreadId int
	GCPhasePauseDuration     float64 `json:"duration"`
	GcId                     int     `json:"gcId"`
	GCPhasePauseName         string  `json:"name"`
}

func (r *Jvm) TableName() string {
	return "jvms"
}

func CreateJvm(db *gorm.DB, j *Jvm) (uint, error) {
	err := db.Create(j).Error
	if err != nil {
		return 0, err

	}
	return j.ID, nil
}
