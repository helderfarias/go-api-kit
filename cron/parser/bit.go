package parser

import "math"

// bits return a bit field where all the bits less significant than or equal to the value bit is set to 1's.
func bits(value uint8) uint64 {
	return math.MaxUint64 >> (64 - value)
}

// field returns the bit field where the bit number corresponding to the value is set to 1.
func bit(value uint8) uint64 {
	return 1 << (value - 1)
}

// rangeField returns a bit field where the bits higher than or equal to from and lower than or equal to to are set to 1's.
func rangeField(from uint8, to uint8) uint64 {
	if from == to {
		return bit(from)
	}
	return bits(to) ^ bits(from-1)
}

// stepField returns a bit field where the bits, including from and to, are set to 1's at every step within the range..
func stepField(from uint8, to uint8, step uint8) uint64 {
	value := uint64(0)
	for i := from; i <= to; i = i + step {
		value |= bit(i)
	}
	return value
}

func listField(values []uint8) uint64 {
	var value uint64
	for _, val := range values {
		value |= bit(val)
	}
	return value
}
