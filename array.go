// Package bit provides a bit array implementation and some utility functions.
package bit

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

// Array represents a bit array.
type Array struct {
	bits   []uint64
	length int
}

// NewArray creates a new bit array
// with an initial bit capacity of n.
func NewArray(n int) *Array {
	if n < 0 {
		panic("bit: array size must be greater than or equal 0")
	}

	b := make([]uint64, 1, (n>>6)+1)
	return &Array{b, 0}
}

// Add appends the bits given its size to the array.
func (a *Array) Add(bits uint64, size int) {
	if size <= 0 || size > 64 {
		panic("bit: bit size must be in range [1,64]")
	}

	// Extend bits if necessary
	lenbits := len(a.bits)
	freespace := (lenbits << 6) - a.length
	overflow := size - freespace
	if overflow > 0 {
		a.bits = append(a.bits, 0)
	}

	// Append bits
	idx := lenbits - 1
	if freespace > 0 {
		a.bits[idx] |= bits << uint(a.length&63)
	}

	if overflow > 0 {
		a.bits[idx+1] |= bits >> uint(freespace)
	}

	// Increment size
	a.length += size
}

// Insert inserts bits to index idx overwriting its contents.
func (a *Array) Insert(idx int, bits uint64, size int) {
	if idx > a.length {
		panic("bit: index out of bounds")
	} else if size <= 0 || size > 64 {
		panic("bit: bit size must be in range [1,64]")
	}

	// Extend bits if necessary
	overflow := idx + size - a.length
	if overflow > 0 {
		lenbits := len(a.bits)
		freespace := (lenbits << 6) - a.length
		if overflow > freespace {
			a.bits = append(a.bits, 0)
		}

		a.length += overflow
	}

	bitIdx := idx & 63
	arrayIdx := idx >> 6
	lowBitSz := 64 - bitIdx
	overflow = size - lowBitSz

	ba := a.bits[arrayIdx]
	bb := bits << uint(bitIdx)
	mask := uint64(1<<uint(size)) - 1

	// Use bit twiddling hacks (merging bits)
	// https://graphics.stanford.edu/~seander/bithacks.html
	a.bits[arrayIdx] = ba ^ ((ba ^ bb) & (mask << uint(bitIdx)))

	if overflow > 0 {
		arrayIdx++
		ba = a.bits[arrayIdx]
		bb = bits >> uint(lowBitSz)

		a.bits[arrayIdx] = ba ^ ((ba ^ bb) & (mask >> uint(lowBitSz)))
	}
}

// Get returns the uint64 representation of
// bits starting from index idx given the bit size.
func (a *Array) Get(idx, size int) uint64 {
	if idx > a.length {
		panic("bit: index out of bounds")
	} else if size <= 0 || size > 64 {
		panic("bit: bit size must be in range [1,64]")
	}

	bitIdx := idx & 63
	arrayIdx := idx >> 6
	lowBitSz := 64 - bitIdx
	overflow := size - lowBitSz

	res := a.bits[arrayIdx] >> uint(bitIdx)
	if overflow > 0 {
		res |= a.bits[arrayIdx+1] << uint(lowBitSz)
	}

	return res & ((1 << uint(size)) - 1)
}

// Bits returns the underlying array.
func (a *Array) Bits() []uint64 {
	return a.bits
}

// Reset resets the array
// to its initial state.
func (a *Array) Reset() {
	for i := range a.bits {
		a.bits[i] = 0
	}

	a.bits = a.bits[0:1]
	a.length = 0
}

// Len returns the number
// of bits stored in the array.
func (a *Array) Len() int {
	return a.length
}

// Size returns the size
// of the array in bytes.
func (a *Array) Size() int {
	return len(a.bits) * 8
}

// String returns a hexadecimal
// string representation of the array.
func (a *Array) String() string {
	buf := new(bytes.Buffer)
	for i := len(a.bits) - 1; i >= 0; i-- {
		bits := fmt.Sprintf("%16X", a.bits[i])
		bits = strings.Replace(bits, " ", "0", -1)
		fmt.Fprintf(buf, "%s [%d-%d] ", bits, (i<<6)+63, i<<6)
	}

	return buf.String()
}

// GobEncode allows this array
// to be encoded into gob streams.
func (a *Array) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	enc.Encode(a.bits)
	enc.Encode(a.length)

	return buf.Bytes(), nil
}

// GobDecode allows this array
// to be decoded from gob streams.
func (a *Array) GobDecode(data []byte) error {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)

	dec.Decode(&a.bits)
	dec.Decode(&a.length)

	return nil
}
