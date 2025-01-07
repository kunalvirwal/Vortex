package utils

import "time"

type scheduler struct {
	count       int // -1 for infinite
	interval    time.Duration
	task        func()
	stopTrigger chan bool
}

// creates a new scheduler for a task
func NewScheduler(count int, interval time.Duration, task func()) *scheduler {
	return &scheduler{
		count:    count,
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

	var i int = 0
	ticker := time.NewTicker(s.interval)
	select {
	case <-ticker.C:
		i++
		if i <= s.count || s.count < 0 {
			s.task()
		}
	case <-s.stopTrigger:
		ticker.Stop()
		return
	}
}

// stops the scheduler
func (s *scheduler) Terminate() {
	s.stopTrigger <- true
}
