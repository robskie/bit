package bit

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	array := NewArray(0)

	// Simple case
	array.Add(0x7, 4)
	array.Add(0x5, 4)
	assert.Equal(t, 8, array.Len())
	assert.Equal(t, 1, len(array.bits))
	assert.EqualValues(t, 0x57, array.bits[0])

	// Bit array full then add
	// Resulting array in hex:
	// FF00000000000057 [63-0]
	// 00000000000000A5 [127-64]
	array.Add(0xFF<<48, 56)
	array.Add(0xA5, 8)
	assert.Equal(t, 2, len(array.bits))
	assert.Equal(t, 72, array.Len())
	assert.EqualValues(t, 0xA5, array.bits[1])

	// Bit array partially full then add
	// Resulting array in hex:
	// FF00000000000057 [63-0]
	// 00000000000000A5 [127-64]
	// 0000000000000057 [191-128]
	array.Add(0x57<<56, 64)
	assert.Equal(t, 136, array.Len())
	assert.Equal(t, 3, len(array.bits))
	assert.EqualValues(t, 0x57, array.bits[2])
}

func TestGet(t *testing.T) {
	array := NewArray(0)

	// Get bits spanning only one uint64 element
	array.Add(0xF5, 60)
	assert.EqualValues(t, 0xF5, array.Get(0, 8))

	// Get bits that spans 2 uint64 array element
	// Resulting array in hex:
	// 70000000000000F5 [63-0]
	// 000000000000000F [127-64]
	array.Add(0xF7, 8)
	assert.EqualValues(t, 0xF5, array.Get(0, 8))
	assert.EqualValues(t, 0xF7, array.Get(60, 8))
}

func TestInsert(t *testing.T) {
	array := NewArray(0)

	// Insert bits spanning only one uint64 element
	// Resulting array in hex:
	// 000000000000058F [63-0]
	array.Add(0xFF, 8)
	array.Insert(4, 0x58, 8)
	assert.Equal(t, 12, array.Len())
	assert.Equal(t, 1, len(array.bits))
	assert.EqualValues(t, 0x58F, array.bits[0])

	// Insert bits spanning 2 uint64 element
	// Resulting array in hex:
	// AF0000000000058F [63-0]
	// 0000000000000009 [127-64]
	array.Add(0xFF<<44, 52)
	array.Insert(60, 0x9A, 8)
	assert.Equal(t, 68, array.Len())
	assert.Equal(t, 2, len(array.bits))
	assert.EqualValues(t, 0xA, array.bits[0]>>60)
	assert.EqualValues(t, 0x9, array.bits[1])
}

func TestReset(t *testing.T) {
	array := NewArray(0)

	array.Add(0xFF, 64)
	array.Add(0xFF, 64)
	array.Reset()

	array.Add(0x55, 8)
	assert.EqualValues(t, 0x55, array.bits[0])
	assert.Equal(t, 8, array.Len())
}

func TestAddGetRandom(t *testing.T) {
	array := NewArray(63 * 1e6)
	tc := make([]uint64, 1e6)

	for i := range tc {
		r := uint64(rand.Int63())
		tc[i] = r
		array.Add(r, 63)
	}

	for i := range tc {
		if !assert.Equal(t, tc[i], array.Get(i*63, 63)) {
			break
		}
	}
}

func TestInsertGetRandom(t *testing.T) {
	array := NewArray(63 * 1e6)
	tc := make([]uint64, 1e6)

	// Populate array
	for _ = range tc {
		r := uint64(rand.Int63())
		array.Add(r, 63)
	}

	// Insert and overwrite bits
	for i := range tc {
		r := uint64(rand.Int63())
		tc[i] = r
		array.Insert(i*63, r, 63)
	}

	for i := range tc {
		if !assert.Equal(t, tc[i], array.Get(i*63, 63)) {
			break
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	array := NewArray(63 * 1e6)
	tc := make([]uint64, 1e6)

	for i := range tc {
		r := uint64(rand.Int63())
		tc[i] = r
		array.Add(r, 63)
	}

	data, _ := array.GobEncode()
	narray := NewArray(0)
	narray.GobDecode(data)

	for i := range tc {
		if !assert.Equal(t, tc[i], narray.Get(i*63, 63)) {
			break
		}
	}
}

func benchmarkAdd(size int, b *testing.B) {
	array := NewArray(0)
	for i := 0; i < b.N; i++ {
		array.Add(0xFF, size)
	}
}

func BenchmarkAdd7(b *testing.B)  { benchmarkAdd(7, b) }
func BenchmarkAdd15(b *testing.B) { benchmarkAdd(15, b) }
func BenchmarkAdd31(b *testing.B) { benchmarkAdd(31, b) }
func BenchmarkAdd63(b *testing.B) { benchmarkAdd(63, b) }

func BenchmarkAddRand(b *testing.B) {
	array := NewArray(0)

	// Create random sizes
	sz := make([]int, b.N)
	for i := range sz {
		sz[i] = rand.Intn(64) + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.Add(0xFF, sz[i])
	}
}

var bigArray *Array

func initBigArray() {
	if bigArray == nil {
		bigArray = NewArray(64 * 1e7)
		for i := 0; i < 1e7; i++ {
			bigArray.Add(0xFF, 64)
		}
	}
}

func BenchmarkInsertRandIdx(b *testing.B) {
	initBigArray()

	// Create random indices and sizes
	sz := make([]int, b.N)
	idx := make([]int, b.N)
	maxidx := bigArray.Len() - 64
	for i := range sz {
		sz[i] = rand.Intn(64) + 1
		idx[i] = rand.Intn(maxidx)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bigArray.Insert(idx[i], 0xFF, sz[i])
	}
}

func benchmarkGet(size int, b *testing.B) {
	initBigArray()

	idx := 0
	maxidx := bigArray.Len() - size

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bigArray.Get(idx%maxidx, size)
		idx += size
	}
}

func BenchmarkGet7(b *testing.B)  { benchmarkGet(7, b) }
func BenchmarkGet15(b *testing.B) { benchmarkGet(15, b) }
func BenchmarkGet31(b *testing.B) { benchmarkGet(31, b) }
func BenchmarkGet63(b *testing.B) { benchmarkGet(63, b) }

func BenchmarkGetRand(b *testing.B) {
	initBigArray()

	// Create random sizes
	sz := make([]int, b.N)
	maxidx := bigArray.Len() - 64
	for i := range sz {
		sz[i] = rand.Intn(64) + 1
	}
	b.ResetTimer()

	idx := 0
	for i := 0; i < b.N; i++ {
		bigArray.Get(idx%maxidx, sz[i])
		idx += sz[i]
	}
}

func BenchmarkGetRandIdx(b *testing.B) {
	initBigArray()

	// Create random indices and sizes
	sz := make([]int, b.N)
	idx := make([]int, b.N)
	maxidx := bigArray.Len() - 64
	for i := range sz {
		sz[i] = rand.Intn(64) + 1
		idx[i] = rand.Intn(maxidx)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bigArray.Get(idx[i], sz[i])
	}
}
