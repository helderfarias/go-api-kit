package main

import (
	"log"
	"runtime"

	"github.com/helderfarias/go-api-kit/cron"
)

func main() {
	cron.NewSchedule("0 * * * *").Run(func() {
		log.Println("ok")
	})

	runtime.Goexit()
}
