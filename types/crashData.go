package types

import (
	"sync"
	"time"
)

type CrashData struct {
	CrashCount             uint // reset with every modification to this service and reset of backoff duration
	LastCrashTime          time.Time
	CurrentBackoffDuration time.Duration
	IsInCrashLoop          bool
	CrashHistory           []CrashEvent
	Mu                     sync.RWMutex
}

type CrashEvent struct {
	Timestamp time.Time
	ExitCode  string
	Signal    string
	ExitError string
}
