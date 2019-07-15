package parser

import (
	"fmt"
	"regexp"
)

// This spec tries to adhere to the 4th Berkely Distribution of the crontab
// manual (man 5 crontab) dated 19 April 2010.

// Regular expression strings
const (
	startExp               = `^`
	endExp                 = `$`
	everyExp               = `\*`
	singleOrDoubleDigitExp = `([\d]{1,2})`
	aliasExp               = `([[:alpha:]]{3})`
	stepExp                = `/` + singleOrDoubleDigitExp
	numberRangeExp         = singleOrDoubleDigitExp + `-` + singleOrDoubleDigitExp
	listExp                = singleOrDoubleDigitExp + `(?:,\s*` + singleOrDoubleDigitExp + `)*`
	nameExp                = `@[[:alpha:]]+`
)

var (
	every               = regexp.MustCompile(startExp + everyExp + endExp)
	step                = regexp.MustCompile(startExp + stepExp + endExp)
	everyStep           = regexp.MustCompile(startExp + everyExp + stepExp + endExp)
	singleOrDoubleDigit = regexp.MustCompile(startExp + singleOrDoubleDigitExp + endExp)
	alias               = regexp.MustCompile(startExp + aliasExp + endExp)
	numberRange         = regexp.MustCompile(startExp + numberRangeExp + endExp)
	list                = regexp.MustCompile(startExp + listExp + endExp)
	rangeStep           = regexp.MustCompile(startExp + numberRangeExp + stepExp + endExp)
	name                = regexp.MustCompile(startExp + nameExp + endExp)
)

// Days and months can be specified with named aliases such as "mon", "jan", etc.
type aliases map[string]uint8

// Every field has a minimum and maximum value and possibly aliases.
type fieldSpec struct {
	Min     uint8
	Max     uint8
	Aliases aliases
}

// Dealias returns the value aliased value by the given alias. Error returned if the field has no such alias or no aliases.
func (f *fieldSpec) Dealias(alias string) (uint8, error) {
	if f.Aliases == nil {
		return 0, fmt.Errorf("field has no aliases")
	}
	if number, ok := f.Aliases[alias]; !ok {
		return 0, fmt.Errorf(`"%v" is not a valid alias`, alias)
	} else {
		return number, nil
	}
}

// InRange returns a boolean indicating if the given number lies in the range of the minimum and maximum value of the field spec.
func (f *fieldSpec) InRange(number uint8) bool {
	if number < f.Min || number > f.Max {
		return false
	} else {
		return true
	}
}

func (f *fieldSpec) String() string {
	return fmt.Sprintf("min %v, max %v, aliases %+v", f.Min, f.Max, f.Aliases)
}

type fields struct {
	minute *fieldSpec
	hour   *fieldSpec
	dom    *fieldSpec
	month  *fieldSpec
	dow    *fieldSpec
}

var spec = &fields{
	minute: &fieldSpec{0, 59, nil},
	hour:   &fieldSpec{0, 23, nil},
	dom:    &fieldSpec{1, 31, nil},
	month: &fieldSpec{1, 12,
		aliases{
			"jan": 1,
			"feb": 2,
			"mar": 3,
			"apr": 4,
			"may": 5,
			"jun": 6,
			"jul": 7,
			"aug": 8,
			"sep": 9,
			"okt": 10,
			"nov": 11,
			"des": 12,
		},
	},
	dow: &fieldSpec{0, 7,
		aliases{
			"sun": 0,
			"mon": 1,
			"tue": 2,
			"wed": 3,
			"thu": 4,
			"fri": 5,
			"sat": 6,
			// "sun": 7,
		},
	},
}

// Common cron expressions can be specified using names
var names = map[string]string{
	"@yearly":   "0 0 1 1 *", // 1st day in the 1st month at midnight
	"@annually": "@yearly",
	"@monthly":  "0 0 1 * *", // 1st day of every month at midnight
	"@weekly":   "0 0 * * 0", // Every sunday at midnight
	"@daily":    "0 0 * * *", // Every day at noon
	"@midnight": "@daily",
	"@hourly":   "0 * * * *", // Every hour
}
