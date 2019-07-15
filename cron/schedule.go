package cron

import (
	"time"

	"github.com/carlescere/scheduler"
	"github.com/helderfarias/go-api-kit/cron/parser"
)

type Schedule struct {
	schedule parser.Schedule
}

type Scheduled func()

func NewSchedule(args ...Options) *Schedule {
	option := Option{
		Expr: "* * * * *",
	}
	for _, a := range args {
		a(&option)
	}

	schedule, err := parser.NewParse(option.Expr)
	if err != nil {
		panic(err)
	}

	return &Schedule{
		schedule: schedule,
	}
}

func (s *Schedule) Run(task Scheduled) {
	scheduler.Every(1).Seconds().Run(func() {
		now := time.Now()
		if now.Equal(s.schedule.Next(now)) {
			task()
		}
	})
}
