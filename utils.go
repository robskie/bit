package bit

var logTable [256]int8

func init() {
	logTable[0] = -1
	for i := 2; i < 256; i++ {
		logTable[i] = 1 + logTable[i/2]
	}
}

// MSBIndex returns the index of the most significant set bit.
// This is equivalent to log base 2 and returns -1 when v is 0.
func MSBIndex(v uint64) int {
	// From https://graphics.stanford.edu/~seander/bithacks.html

	r := int8(0)
	if tt := v >> 56; tt > 0 {
		r = 56 + logTable[tt]
	} else if tt = v >> 48; tt > 0 {
		r = 48 + logTable[tt]
	} else if tt = v >> 40; tt > 0 {
		r = 40 + logTable[tt]
	} else if tt = v >> 32; tt > 0 {
		r = 32 + logTable[tt]
	} else if tt := v >> 24; tt > 0 {
		r = 24 + logTable[tt]
	} else if tt := v >> 16; tt > 0 {
		r = 16 + logTable[tt]
	} else if tt = v >> 8; tt > 0 {
		r = 8 + logTable[tt]
	} else {
		r = logTable[v]
	}

	return int(r)
}

// PopCount counts the number of
// set bits in the given integer.
func PopCount(v uint64) int {
	// From http://vigna.di.unimi.it/ftp/papers/Broadword.pdf
	v -= (v & 0xAAAAAAAAAAAAAAAA) >> 1
	v = (v & 0x3333333333333333) + ((v >> 2) & 0x3333333333333333)
	v = (v + (v >> 4)) & 0x0F0F0F0F0F0F0F0F

	return int(v * 0x0101010101010101 >> 56)
}

// Size returns the minimum number of bits
// required to represent the given integer.
func Size(v uint64) int {
	if v == 0 {
		return 1
	}
	return MSBIndex(v) + 1
}
