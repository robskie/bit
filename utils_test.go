package bit

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMSBIndex(t *testing.T) {
	for i := uint(0); i < 64; i++ {
		if !assert.EqualValues(t, i, MSBIndex(1<<i)) {
			break
		}
	}
}

func TestPopCount(t *testing.T) {
	for i := uint(0); i < 64; i++ {
		v := uint64(1<<i) - 1
		if !assert.EqualValues(t, i, PopCount(v)) {
			break
		}
	}
	assert.EqualValues(t, 64, PopCount(^uint64(0)))
}

func TestRank(t *testing.T) {
	for i := 0; i < 64; i++ {
		v := ^uint64(0)
		if !assert.Equal(t, i+1, Rank(v, i)) {
			break
		}
	}
}

func TestSelect(t *testing.T) {
	for i := 0; i < 64; i++ {
		v := ^uint64(0)
		if !assert.Equal(t, i, Select(v, i+1)) {
			break
		}
	}
}

func BenchmarkMSBIndex(b *testing.B) {
	val := make([]uint64, b.N)
	for i := range val {
		val[i] = uint64(rand.Int63())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MSBIndex(val[i])
	}
}

func BenchmarkPopCount(b *testing.B) {
	val := make([]uint64, b.N)
	for i := range val {
		val[i] = uint64(rand.Int63())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PopCount(val[i])
	}
}

func BenchmarkRank(b *testing.B) {
	idx := make([]int, b.N)
	for i := range idx {
		idx[i] = rand.Intn(64)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Rank(^uint64(0), idx[i])
	}
}

func BenchmarkSelect(b *testing.B) {
	rank := make([]int, b.N)
	for i := range rank {
		rank[i] = rand.Intn(64) + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Select(^uint64(0), rank[i])
	}
}
