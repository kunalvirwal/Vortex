package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/types"
)

type Scheduler struct {
	count       int // -1 for infinite
	interval    time.Duration
	task        func()
	stopTrigger chan bool
}

type SchedulerList struct {
	Mu       sync.RWMutex
	Schedule map[string]*Scheduler
}

var BackOffResetSchedulers = &SchedulerList{
	Schedule: make(map[string]*Scheduler),
}

// creates a new scheduler for a task
func NewScheduler(count int, interval time.Duration, task func()) *Scheduler {
	return &Scheduler{
		count:       count,
		interval:    interval,
		task:        task,
		stopTrigger: make(chan bool),
	}
}

// runs the scheduler asynchronously
func (s *Scheduler) RunAsync() {
	go s.Run()
}

// runs the scheduler task
func (s *Scheduler) Run() {

	var i int = 0
	ticker := time.NewTicker(s.interval)
	select {
	case <-ticker.C:
		i++
		if i <= s.count || s.count < 0 {
			s.task()
		} else {
			s.Terminate()
		}
	case <-s.stopTrigger:
		ticker.Stop()
		return
	}
}

// stops the scheduler
func (s *Scheduler) Terminate() {
	s.stopTrigger <- true
}

func ResetCrashbackOff(containerID string) {

	exists := false
	var Vcontainer *types.ContainerConfig
	state.VortexContainers.Mu.RLock()
	for _, vcontainers := range state.VortexContainers.List {
		if vcontainers.ID == containerID {
			exists = true
			Vcontainer = vcontainers
			break
		}
	}
	state.VortexContainers.Mu.RUnlock()

	if exists && !Vcontainer.CrashData.IsInCrashLoop {
		fmt.Println("Resetting CrashBackoff for ", containerID)
		Vcontainer.CrashData.Mu.Lock()
		Vcontainer.CrashData.CrashCount = 0
		Vcontainer.CrashData.CurrentBackoffDuration = time.Duration(0)
		Vcontainer.CrashData.Mu.Unlock()
	}

}

func DelScheduler(containerID string) {
	BackOffResetSchedulers.Mu.Lock()
	if BackOffResetSchedulers.Schedule[containerID] != nil {
		BackOffResetSchedulers.Schedule[containerID].Terminate()
		delete(BackOffResetSchedulers.Schedule, containerID)
	}

	BackOffResetSchedulers.Mu.Unlock()
}
