package utils

import "time"

type scheduler struct {
	interval    time.Duration
	task        func()
	stopTrigger chan bool
}

// creates a new scheduler for a task
func NewScheduler(interval time.Duration, task func()) *scheduler {
	return &scheduler{
		interval: interval,
		task:     task,
	}
}

// runs the scheduler asynchronously
func (s *scheduler) RunAsync() {
	go s.Run()
}

// runs the scheduler task
func (s *scheduler) Run() {

	ticker := time.NewTicker(s.interval)
	select {
	case <-ticker.C:
		s.task()
	case <-s.stopTrigger:
		ticker.Stop()
		return
	}
}

// stops the scheduler
func (s *scheduler) Terminate() {
	s.stopTrigger <- true
}
