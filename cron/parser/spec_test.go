package parser

import (
	"regexp"
	"strconv"
	"testing"
)

type expectation = map[bool][]string
type expectations = map[*regexp.Regexp]expectation

var samples = expectations{
	every: {
		true:  {"*"},
		false: {"**", "2", "a"},
	},
	singleOrDoubleDigit: {
		true:  {"1", "2", "10", "12"},
		false: {"a1", "123", "a", "."},
	},
	step: {
		true:  {"/2", "/28"},
		false: {"/a", "2", "/123"},
	},
	everyStep: {
		true:  {"*/2", "*/28"},
		false: {"*/a", "2", "/123"},
	},
	alias: {
		true:  {"abb", "man"},
		false: {"m", "mann", "123"},
	},
	list: {
		true:  {"1,2,3,4", "12,1,42,2"},
		false: {"a,1,2,3", "1,2,3,", "123,2"},
	},
	rangeStep: {
		true:  {"12-31/23", "1-2/2", "1-31/32"},
		false: {"*/2", "1,2,4", "123"},
	},
	numberRange: {
		true:  {"1-2", "21-2"},
		false: {"123-2", "5-123", "123-123"},
	},
	name: {
		true:  {"@daily", "@annually", "@midnight", "@hourly"},
		false: {"@every", "monday"},
	},
}

func matchSamples(t *testing.T, r *regexp.Regexp) {
	match := func(exp string) bool {
		return r.MatchString(exp)
	}
	for expectation := range samples[r] {
		for _, sample := range samples[r][expectation] {
			if !expectation == match(sample) {
				t.Error("Expected", expectation, sample)
			}
		}
	}
}

func TestEveryRegexp(t *testing.T) {
	matchSamples(t, every)
}

func TestSingleOrDoubleDigitRexexp(t *testing.T) {
	matchSamples(t, singleOrDoubleDigit)
}

func TestStepRegexp(t *testing.T) {
	matchSamples(t, step)
}

func TestEveryStepRegexp(t *testing.T) {
	matchSamples(t, everyStep)
}

func TestAliasRegexp(t *testing.T) {
	matchSamples(t, alias)
}

func TestRangeRegexp(t *testing.T) {
	matchSamples(t, numberRange)
}

func TestRangeStepRegexp(t *testing.T) {
	matchSamples(t, rangeStep)
}

func TestListRegexp(t *testing.T) {
	matchSamples(t, list)
}

func TestNameRegexp(t *testing.T) {
	matchSamples(t, numberRange)
}

func TestEveryStepSubmatch(t *testing.T) {
	expectedStep := uint64(32)
	sub := everyStep.FindStringSubmatch("*/32")
	if len(sub) != 2 {
		t.Fatal("expected length of submatches", 2, "got", len(sub))
	}
	step, err := strconv.ParseUint(sub[1], 10, 8)
	if err != nil {
		t.Error(err)
	}
	if step != expectedStep {
		t.Error("expected", expectedStep, "got", step)
	}
}

func TestRangeStepSubmatch(t *testing.T) {
	expectedFrom := uint8(1)
	expectedTo := uint8(31)
	expectedStep := uint8(32)

	s := "1-31/32"

	sub := rangeStep.FindStringSubmatch(s)
	if len(sub) != 4 {
		t.Fatal("expected length of submatch", 4, "got", len(sub), s)
	}

	from, err := strconv.ParseUint(sub[1], 10, 8)
	if err != nil {
		t.Error(err)
	} else if uint8(from) != expectedFrom {
		t.Error("expected", expectedFrom, "got", from)
	}

	to, err := strconv.ParseUint(sub[2], 10, 8)
	if err != nil {
		t.Error(err)
	} else if uint8(to) != expectedTo {
		t.Error("expected", expectedTo, "got", to)
	}

	step, err := strconv.ParseUint(sub[3], 10, 8)
	if err != nil {
		t.Error(err)
	} else if uint8(step) != expectedStep {
		t.Error("expected", expectedStep, "got", step)
	}
}

func TestValidNumber(t *testing.T) {
	if spec.minute.InRange(60) {
		t.Error("60 should not be a valid number for minute field")
	}
	if !spec.minute.InRange(0) {
		t.Error("0 should be a valid number")
	}
	if !spec.minute.InRange(59) {
		t.Error("60 not valid number for minute")
	}
}

func TestValidAlias(t *testing.T) {
	if _, err := spec.dow.Dealias("mond"); err == nil {
		t.Error("should error on invalid alias")
	}
	if _, err := spec.dow.Dealias("mon"); err != nil {
		t.Error("mon should be a valid alias")
	}
}
