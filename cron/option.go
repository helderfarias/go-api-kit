package cron

type Option struct {
	Expr string
}

type Options func(o *Option)

func Cron(expr string) {
	return func(o *Option) {
		o.Expr = expr
	}
}
