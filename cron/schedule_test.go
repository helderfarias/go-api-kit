package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSchedule(t *testing.T) {
	task := NewSchedule()

	assert.NotNil(t, task)
}

func TestCreateScheduleWithEvery(t *testing.T) {
	result := ""

	NewSchedule(Every("1s")).Run(func() {
		result = "ok"
	})

	assert.NotNil(t, "ok", result)
}
