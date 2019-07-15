package parser

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

const (
	colorRED    = "\033[31m"
	colorNORMAL = "\033[39m"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(colorRED)
		fmt.Printf("%s:%d:\n\n"+msg, append([]interface{}{filepath.Base(file), line}, v...)...)
		fmt.Printf(colorNORMAL)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(colorRED)
		fmt.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		fmt.Printf(colorNORMAL)
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		pc, file, line, _ := runtime.Caller(1)
		fn := runtime.FuncForPC(pc).Name()
		fmt.Printf(colorRED)
		fmt.Printf("%s:%d (%v)\n\n\texp: %+v\n\n\tgot: %+v\n\n", filepath.Base(file), line, fn, exp, act)
		fmt.Printf(colorNORMAL)
		tb.FailNow()
	}
}
