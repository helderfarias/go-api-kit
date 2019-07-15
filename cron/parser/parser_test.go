package parser

import (
	"testing"
)

func TestParseSingleOrDoubleDigit(t *testing.T) {
	actual, err := parseSingleOrDoubleDigit("1", spec.minute)
	expected := uint64(2)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseSingleOrDoubleDigit("2", spec.minute)
	expected = uint64(4)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseSingleOrDoubleDigit("3", spec.minute)
	expected = uint64(8)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseSingleOrDoubleDigit("3", spec.dom)
	expected = 4
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseSingleOrDoubleDigit("7", spec.dow)
	expected = 1
	ok(t, err)
	equals(t, expected, actual)
}

func TestParseEvery(t *testing.T) {
	actual, err := parseEvery(spec.dow)
	ok(t, err)
	expected := bits(offset(spec.dow, spec.dow.Max-1))
	equals(t, expected, actual)

	actual, err = parseEvery(spec.dom)
	ok(t, err)
	expected = bits(spec.dom.Max)
	equals(t, expected, actual)

	actual, err = parseEvery(spec.minute)
	ok(t, err)
	expected = bits(offset(spec.minute, spec.minute.Max))
	equals(t, expected, actual)

	actual, err = parseEvery(spec.hour)
	ok(t, err)
	expected = bits(offset(spec.hour, spec.hour.Max))
	equals(t, expected, actual)
}

func TestParseEveryStep(t *testing.T) {
	actual, err := parseEveryStep("*/1", spec.minute)
	expected := bits(60)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseEveryStep("*/2", spec.minute)
	expected = uint64(0)
	for i := uint8(1); i <= 60; i += 2 {
		expected |= bit(i)
	}
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseEveryStep("*/3", spec.month)
	expected = uint64(0)
	for i := uint8(1); i <= spec.month.Max; i += 3 {
		expected |= bit(i)
	}
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseEveryStep("*/3", spec.dow)
	expected = uint64(spec.dow.Min)
	for i := uint8(1); i <= spec.dow.Max; i += 3 {
		expected |= bit(i)
	}
	ok(t, err)
	equals(t, expected, actual)
}

func TestParseRange(t *testing.T) {
	actual, err := parseRange("0-6", spec.dow)
	expected := bits(7)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRange("1-6", spec.dow)
	expected = bits(7) ^ bit(1)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRange("2-6", spec.dow)
	expected = bits(7) ^ bits(2)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRange("1-7", spec.dow)
	expected = bits(7)
	ok(t, err)
	equals(t, expected, actual)
}

func TestParseRangeStep(t *testing.T) {
	actual, err := parseRangeStep("2-5/2", spec.minute)
	expected := bit(3) + bit(5)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRangeStep("2-5/3", spec.minute)
	expected = bit(3) + bit(6) // 100100
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRangeStep("2-5/3", spec.dom)
	expected = bit(2) + bit(5) // 10010
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRangeStep("3-6/2", spec.dow) // 0101000
	expected = bit(4) + bit(6)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRangeStep("4-7/2", spec.dow) // 01010000
	expected = bit(5) + bit(7)
	ok(t, err)
	equals(t, expected, actual)

	actual, err = parseRangeStep("3-7/2", spec.dow) // 0101001
	expected = bit(1) + bit(4) + bit(6)
	ok(t, err)
	equals(t, expected, actual)
}

func TestParseList(t *testing.T) {
	actual, err := parseList("1, 2, 3", spec.dom)
	expected := bit(3) + bit(2) + bit(1)
	ok(t, err)
	equals(t, expected, actual)
}

func TestParseAlias(t *testing.T) {
	actual, err := parseAlias("sun", spec.dow)
	expected := bit(1)
	ok(t, err)
	equals(t, expected, actual)
}
