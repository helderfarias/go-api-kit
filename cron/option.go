package cron

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Option struct {
	Expr      string
	Recover   func(next Scheduled) Scheduled
	StackSize int
}

type Options func(o *Option)

func Expr(expr string) Options {
	return func(o *Option) {
		o.Expr = expr
	}
}

func Every(t string) Options {
	return func(o *Option) {
		o.Expr = fmt.Sprintf("@every %s", t)
	}
}

func Minutes() Options {
	return func(o *Option) {
		o.Expr = "@every 60s"
	}
}

func Seconds() Options {
	return func(o *Option) {
		o.Expr = "@every 1s"
	}
}

func Hourly() Options {
	return func(o *Option) {
		o.Expr = "@hourly"
	}
}

func Daily() Options {
	return func(o *Option) {
		o.Expr = "@daily"
	}
}

func Weekly() Options {
	return func(o *Option) {
		o.Expr = "@weekly"
	}
}

func Monthly() Options {
	return func(o *Option) {
		o.Expr = "@monthly"
	}
}

func Yearly() Options {
	return func(o *Option) {
		o.Expr = "@yearly"
	}
}

func Recover() Options {
	return func(o *Option) {
		o.Recover = func(next Scheduled) Scheduled {
			return func() {
				defer func() {
					if r := recover(); r != nil {
						err, ok := r.(error)
						if !ok {
							err = fmt.Errorf("%v", r)
						}

						stack := make([]byte, o.StackSize)
						length := runtime.Stack(stack, true)
						logrus.Errorf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					}
				}()

				next()
			}
		}
	}
}
