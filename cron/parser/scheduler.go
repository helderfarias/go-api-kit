package parser

import (
	"log"
	"time"
)

type Schedule interface {
	Next(now time.Time) time.Time
}

type bitCron struct {
	minute, hour, dom, month, dow uint64
}

// diff answers the question: How many times must the given value's bit field
// be shifted (and possibly wrapped) before this bit and the given field both are 1's.
func diff(fieldSpec *fieldSpec, value int, field uint64) int {
	b := bit(offset(fieldSpec, uint8(value)))
	msb := bit(offset(fieldSpec, fieldSpec.Max)) // most significant bit for field

	// r is the number of values in the allowed range for the specified field, e.g. 0-59 -> 60, 1-31 -> 31.
	// Since sunday can be both 0 and 7 we must not increment dow.
	r := int(fieldSpec.Max)
	if fieldSpec.Min == 0 && fieldSpec != spec.dow {
		r += 1
	}

	// search for next bit in in field by starting at bit b, counting the
	// number of shifts until a 1 bit is found. wraps around when bit b equals max bit for field.
	for i := 0; i < r; i++ {
		if b&field > 0 {
			return i
		}
		if b >= msb {
			b = 1 // wrap
		} else {
			b <<= 1 // shift
		}
	}
	// should never happen unless field or value equals 0 (which should never happen):
	// - a value of 0 should be offset to 1 above
	// - a field of 0 should never be an output from the parser
	log.Fatalf("diff should never be zero, %d, %b\n, %+v", value, field, fieldSpec)
	return 0
}

// Find the next time the cron expression should run (could possibly be now).
func (c *bitCron) Next(now time.Time) time.Time {
	t := now

	// month, dom, dow, day, hour and minute variables represent diffs in this function

	month := diff(spec.month, int(t.Month()), c.month)
	if month > 0 {
		t = time.Date(t.Year(), time.Month(int(t.Month())+month), 0, 0, 0, 0, 0, t.Location())
	}

	dom := diff(spec.dom, t.Day(), c.dom)
	// Adjust day of month because months can be 28, 29, 30 or 31 days long
	// and the diff function is blind to this fact (it sees 31).
	dom -= 31 - DaysIn(t.Month(), t.Year())

	dow := diff(spec.month, int(t.Weekday())+1, c.dow)

	// According to the crontab manual (man 5 crontab) if one of dom or dow equals every (*) the other one is used
	day := 0
	if c.dom == bits(spec.dom.Max) {
		day = dow
	} else if c.dow >= bits(spec.dow.Max) { // >= because every dow (*) yields 8 bits (because sunday can be both 0 and 7
		day = dom
	} else {
		//  or else the earliest of dom and dow is used when neither of them are equal to every (*)
		if dom < dow {
			day = dom
		} else {
			day = dow
		}
	}
	if day > 0 {
		t = time.Date(t.Year(), t.Month(), t.Day()+day, 0, 0, 0, 0, t.Location())
	}

	hour := diff(spec.hour, t.Hour(), c.hour)
	if hour > 0 {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+hour, 0, 0, 0, t.Location())
	}

	minute := diff(spec.minute, t.Minute(), c.minute)
	if minute > 0 {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+minute, 0, 0, t.Location())
	}

	return t
}
