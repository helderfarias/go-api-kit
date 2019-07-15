package parser

import (
	"math"
	"testing"
)

func TestStepField(t *testing.T) {
	actual := stepField(1, 2, 1)
	expected := bit(2) + bit(1)
	equals(t, expected, actual)

	// 11100
	actual = stepField(3, 5, 1)
	expected = bit(5) + bit(4) + bit(3)

	equals(t, expected, actual)

	actual = stepField(3, 5, 2)
	expected = bit(5) + bit(3)
	equals(t, expected, actual)

	actual = stepField(3, 5, 3)
	expected = bit(3)
	equals(t, expected, actual)
}

func TestRangeFieldFunction(t *testing.T) {
	// It should return just bit value if from equals to
	actual := rangeField(2, 2)
	expected := bit(2)
	equals(t, expected, actual)

	// e.g. 2-3 should produce binary 5 (110)
	actual = rangeField(2, 3)
	expected = uint64(6)
	equals(t, expected, actual)

	// e.g. 3-6 should produce binary 55 (111100)
	actual = rangeField(3, 6)
	expected = bits(6) ^ (bit(3) - 1)
	// 111111 ^ 111
	equals(t, expected, actual)
	//1000000 -> 111111 ^ 11
	expected = (bit(7) - 1) ^ bits(3-1)
	equals(t, expected, actual)

}

func TestBitsFunction(t *testing.T) {
	// 1111111
	actual := bits(7)
	expected := uint64(127)
	equals(t, expected, actual)
	// 111
	actual = bits(3)
	expected = uint64(7)
	equals(t, expected, actual)
}

func TestBitFunction(t *testing.T) {
	actual := bit(2)
	expected := uint64(2)
	equals(t, expected, actual)
	// binary 8 (1000)
	actual = bit(4)
	expected = uint64(8)
	equals(t, expected, actual)

	actual = uint64(math.Pow(2, 63.0))
	expected = uint64(1) << 63
	equals(t, expected, actual)

	// bits(n)+1 should equal bit(n+1)  and
	actual = bit(7)
	expected = bits(6) + 1
	equals(t, expected, actual)
	// vica versa
	actual = bit(7) - 1
	expected = bits(6)
	equals(t, expected, actual)
}
