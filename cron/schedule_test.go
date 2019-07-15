package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedule(t *testing.T) {
	task := NewSchedule()

	task.Run(func() {})

	assert.NotNil(t, task)
}
