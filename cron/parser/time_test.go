package parser

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	if DaysIn(time.January, 2017) != 31 {
		t.Error("Expected january to have 31 days")
	}

	if DaysIn(time.February, 2017) != 28 {
		t.Error("Expected february 2017 to have 28 days")
	}
}
