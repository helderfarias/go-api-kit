package cron

import (
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	option Option
}

type Scheduled func()

func NewSchedule(opts ...Options) *Schedule {
	option := Option{
		StackSize: 4 << 10, // 4 KB
	}

	all := []Options{}
	all = append(all, Expr("* * * * *"))
	all = append(all, Recover())
	for _, i := range opts {
		all = append(all, i)
	}

	for _, o := range all {
		o(&option)
	}

	return &Schedule{option: option}
}

func (s *Schedule) Run(task Scheduled) {
	target := task

	if s.option.Recover != nil {
		target = s.option.Recover(target)
	}

	c := cron.New()
	c.AddFunc(s.option.Expr, target)
	c.Start()
}
