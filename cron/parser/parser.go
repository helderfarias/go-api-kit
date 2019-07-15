package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(string) Schedule
}

type parser struct{}

var defaultParser = new(parser)

func NewParse(expression string) (Schedule, error) {
	return defaultParser.Parse(expression)
}

// Parses a cron expression to a struct which represents the cron expression string as bitfields.
// Since the maximum number of bits needed is 60 (for minutes) an uint64 can be (and is) used to represent each bitfield.
//
// A bit in position 1 means (according to the the spec): 0th minute, 0 hour, 1st, January, Sunday.
// A bit in position 3 means (according to the spec): 2nd minute, 2nd hour, 3 dom, March, Tuesday.
// A bit in position ...
func (p *parser) Parse(expression string) (Schedule, error) {
	fields := strings.Fields(expression)
	nrOfFields := len(fields)

	if nrOfFields != 1 && nrOfFields != 5 {
		return nil, fmt.Errorf("number of fields expected to be either 1 or 5, got %d", nrOfFields)
	}

	if nrOfFields == 1 {
		return p.parseNamedExpression(fields[0])
	} else {
		minute, err := parseField(fields[0], spec.minute)
		hour, err := parseField(fields[1], spec.hour)
		dom, err := parseField(fields[2], spec.dom)
		month, err := parseField(fields[3], spec.month)
		dow, err := parseField(fields[4], spec.dow)

		// return last error
		if err != nil {
			return nil, err
		}

		return &bitCron{minute, hour, dom, month, dow}, nil
	}
}

// Parse named expressions like @yearly, @daily, etc..
func (p *parser) parseNamedExpression(value string) (Schedule, error) {
	if expression, ok := names[value]; ok {
		return p.Parse(expression)
	} else {
		return nil, fmt.Errorf("no such named cron expression")
	}
}

// offset increments the value by 1 if the spec minimum for the field is 0,
func offset(fieldSpec *fieldSpec, value uint8) uint8 {
	if fieldSpec.Min == 0 {
		return value + 1
	} else {
		return value
	}
}

func parseEvery(fieldSpec *fieldSpec) (uint64, error) {
	min := offset(fieldSpec, fieldSpec.Min)
	max := offset(fieldSpec, fieldSpec.Max)
	if fieldSpec == spec.dow {
		max -= 1
	}
	return rangeField(min, max), nil
}

func parseSingleOrDoubleDigit(value string, fieldSpec *fieldSpec) (uint64, error) {
	parsedNumber, err := strconv.ParseUint(value, 10, 8)
	num := uint8(parsedNumber)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	if !fieldSpec.InRange(num) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", num, fieldSpec.Min, fieldSpec.Max)
	}
	if fieldSpec == spec.dow && num == spec.dow.Max {
		num = 0 // wrap sunday (when valued as 7) to 0
	}
	return bit(offset(fieldSpec, num)), nil
}

func parseEveryStep(value string, fieldSpec *fieldSpec) (uint64, error) {
	sub := everyStep.FindStringSubmatch(value)
	step, err := strconv.ParseUint(sub[1], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	from := offset(fieldSpec, fieldSpec.Min)
	to := offset(fieldSpec, fieldSpec.Max)
	if fieldSpec == spec.dow {
		to -= 1 // shrink dow max to 0 (because sunday can be both 7 and 0)
	}
	return stepField(from, to, uint8(step)), nil
}

func parseRange(value string, fieldSpec *fieldSpec) (uint64, error) {
	sub := numberRange.FindStringSubmatch(value)
	parsedFrom, err := strconv.ParseUint(sub[1], 10, 8)
	parsedTo, err := strconv.ParseUint(sub[2], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	from := uint8(parsedFrom)
	to := uint8(parsedTo)
	if !fieldSpec.InRange(from) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", from, fieldSpec.Min, fieldSpec.Max)
	}
	if !fieldSpec.InRange(to) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", to, fieldSpec.Min, fieldSpec.Max)
	}
	if from > to {
		return 0, fmt.Errorf("from larger than to in range %d-%d", from, to)
	}
	if to == fieldSpec.Max && fieldSpec == spec.dow {
		// wrap the last value in range to the beginning (because sunday can be 7)
		return 1 + rangeField(offset(fieldSpec, from), offset(fieldSpec, to-1)), nil
	} else {
		return rangeField(offset(fieldSpec, from), offset(fieldSpec, to)), nil
	}
}

func parseRangeStep(value string, fieldSpec *fieldSpec) (uint64, error) {
	sub := rangeStep.FindStringSubmatch(value)
	from, err := strconv.ParseUint(sub[1], 10, 8)
	to, err := strconv.ParseUint(sub[2], 10, 8)
	step, err := strconv.ParseUint(sub[3], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	if !fieldSpec.InRange(uint8(from)) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", from, fieldSpec.Min, fieldSpec.Max)
	}
	if !fieldSpec.InRange(uint8(to)) {
		return 0, fmt.Errorf("expected %d in the range %d-%d", to, fieldSpec.Min, fieldSpec.Max)
	}
	bitField := stepField(offset(fieldSpec, uint8(from)), offset(fieldSpec, uint8(to)), uint8(step))
	if fieldSpec.Max == uint8(to) && fieldSpec == spec.dow {
		for i := int(fieldSpec.Max); i > 0; i -= int(step) {
			if i == int(from) {
				return 1 + bitField - bit(fieldSpec.Max+1), nil
			}
		}
	}
	return bitField, nil
}

func parseAlias(alias string, fieldSpec *fieldSpec) (uint64, error) {
	number, err := fieldSpec.Dealias(alias)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	return bit(offset(fieldSpec, number)), nil
}

// NB: Doesn't allow ranges in the list
func parseList(value string, fieldSpec *fieldSpec) (uint64, error) {
	strValues := strings.Split(value, ",")
	values := make([]uint8, len(strValues))
	for i := range strValues {
		val, err := strconv.ParseUint(strings.TrimSpace(strValues[i]), 10, 8)
		if err != nil {
			return 0, fmt.Errorf("%v", err)
		}
		if !fieldSpec.InRange(uint8(val)) {
			return 0, fmt.Errorf("expected %d in the range %d-%d", val, fieldSpec.Min, fieldSpec.Max)
		}
		if fieldSpec == spec.dow && uint8(val) == spec.dow.Max {
			val = 0 // wrap sunday (when valued as 7) to 0
		}
		values[i] = offset(fieldSpec, uint8(val))
	}
	return listField(values), nil
}

// parseField parses any field of a cron expression
func parseField(value string, fieldSpec *fieldSpec) (uint64, error) {
	switch {
	case every.MatchString(value):
		return parseEvery(fieldSpec)

	case singleOrDoubleDigit.MatchString(value):
		return parseSingleOrDoubleDigit(value, fieldSpec)

	case everyStep.MatchString(value):
		return parseEveryStep(value, fieldSpec)

	case numberRange.MatchString(value):
		return parseRange(value, fieldSpec)

	case rangeStep.MatchString(value):
		return parseRangeStep(value, fieldSpec)

	case alias.MatchString(value):
		return parseAlias(value, fieldSpec)

	case list.MatchString(value):
		return parseList(value, fieldSpec)

	default:
		return 0, fmt.Errorf("field %v does not match any pattern", value)
	}
}
