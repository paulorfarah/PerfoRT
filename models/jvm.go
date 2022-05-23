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
	ThreadCPULoadOsThreadId   float64
	ThreadCPULoadJavaName     string
	ThreadCPULoadJavaThreadId float64
	ThreadCPULoadUser         float64
	ThreadCPULoadSystem       float64
}
type ThreadStart struct { //ok
	ThreadStartOsName                   string
	ThreadStartOsThreadId               float64
	ThreadStartJavaName                 string
	ThreadStartJavaThreadId             float64
	ThreadStartParentThreadosName       string
	ThreadStartParentThreadOsThreadId   float64
	ThreadStartParentThreadJavaName     string
	ThreadStartParentThreadJavaThreadId float64
}

type ThreadEnd struct { //ok
	ThreadEndOsName       string
	ThreadEndOsThreadId   float64
	ThreadEndJavaName     string
	ThreadEndJavaThreadId float64
}

type ThreadSleep struct {
	ThreadSleepDuration     float64
	ThreadSleepOsName       string
	ThreadSleepOsThreadId   float64
	ThreadSleepJavaName     string
	ThreadSleepJavaThreadId float64
	ThreadSleepTime         float64
}

type ThreadPark struct {
	ThreadParkDuration     float64
	ThreadParkOsName       string
	ThreadParkOsThreadId   float64
	ThreadParkJavaName     string
	ThreadParkJavaThreadId float64
	ThreadParkParkedClass  string
	ThreadParkTimeout      float64
	ThreadParkUntil        float64
}

type JavaErrorThrow struct {
	JavaErrorThrowDuration     float64
	JavaErrorThrowOsName       string
	JavaErrorThrowOsThreadId   float64
	JavaErrorThrowJavaName     string
	JavaErrorThrowJavaThreadId float64
	JavaErrorThrowMessage      string
	JavaErrorThrowThrownClass  string
}

type JavaExceptionThrow struct {
	JavaExceptionThrowDuration     float64
	JavaExceptionThrowOsName       string
	JavaExceptionThrowOsThreadId   float64
	JavaExceptionThrowJavaName     string
	JavaExceptionThrowJavaThreadId float64
	JavaExceptionThrowMessage      string
	JavaExceptionThrowThrownClass  string
}

type JavaMonitorEnter struct {
	JavaMonitorEnterDuration     float64
	JavaMonitorEnterOsName       string
	JavaMonitorEnterOsThreadId   float64
	JavaMonitorEnterJavaName     string
	JavaMonitorEnterJavaThreadId float64
	JavaMonitorEnterMonitorClass string
}

type JavaMonitorWait struct {
	JavaMonitorWaitDuration     float64
	JavaMonitorWaitOsName       string
	JavaMonitorWaitOsThreadId   float64
	JavaMonitorWaitJavaName     string
	JavaMonitorWaitJavaThreadId float64
	JavaMonitorWaitMonitorClass string
	JavaMonitorWaitTimeout      float64 //	Maximum wait time
	JavaMonitorWaitTimedOut     bool    //Wait has been timed out
}

type OldObjectSample struct {
	OldObjectSampleDuration           float64
	OldObjectSampleOsName             string
	OldObjectSampleOsThreadId         float64
	OldObjectSampleJavaName           string
	OldObjectSampleJavaThreadId       float64
	OldObjectSampleAllocationTime     time.Time
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
	ObjectAllocationInNewTLABOsThreadId     float64
	ObjectAllocationInNewTLABJavaName       string
	ObjectAllocationInNewTLABJavaThreadId   float64
	ObjectAllocationInNewTLABObjectClass    string // Class of allocated object
	ObjectAllocationInNewTLABAllocationSize float64
	ObjectAllocationInNewTLABTlabSize       float64
}

type ObjectAllocationOutsideTLAB struct {
	ObjectAllocationOutsideTLABOsName         string
	ObjectAllocationOutsideTLABOsThreadId     float64
	ObjectAllocationOutsideTLABJavaName       string
	ObjectAllocationOutsideTLABJavaThreadId   float64
	ObjectAllocationOutsideTLABObjectClass    string // Class of allocated object
	ObjectAllocationOutsideTLABAllocationSize float64
}

type GCPhasePause struct {
	GCPhasePauseOsName       string
	GCPhasePauseOsThreadId   float64
	GCPhasePauseJavaName     string
	GCPhasePauseJavaThreadId float64
	GCPhasePauseDuration     float64 `json:"duration"`
	GcId                     float64 `json:"gcId"`
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
