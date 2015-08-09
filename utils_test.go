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
