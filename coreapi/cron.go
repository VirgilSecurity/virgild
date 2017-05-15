package coreapi

import (
	"time"

	"gopkg.in/robfig/cron.v2"
)

type scheduler struct {
	c *cron.Cron
	l Logger
}

func (s *scheduler) Add(period time.Duration, f func() error, name string) {
	s.c.Schedule(&timeSchedule{period}, makeJob(f, name, s.l))
}

type timeSchedule struct {
	d time.Duration
}

func (ts *timeSchedule) Next(t time.Time) time.Time {
	return t.Add(ts.d)
}

func makeJob(f func() error, name string, l Logger) cron.Job {
	job := func() {
		// l.Info("Cron.Job(%s).START", name)
		err := f()

		if err != nil {
			l.Err("Cron.Job(%s): %+v", name, err)
		}

		// l.Info("Cron.Job(%s).FINISH", name)
	}
	return cron.FuncJob(job)
}
