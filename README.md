# bit

Package bit provides a bit array implementation and some utility functions.

## Installation
```sh
go get github.com/robskie/bit
```

## Example

A bit array can be used to compactly store integer values in exchange for slower
operations. For example, suppose you have an input that varies from 0 to 1000,
and you want to store it in an array without wasting too much bits. You
can use an array of ```uint16``` to store it, but if you use a bit array, you
will save more than one third of those precious bit space. Here's how:

```go
import "github.com/robskie/bit"

bitsize := 10
array := bit.NewArray(0)

// Store the input to a bit array
for _, v := range input {
  array.Add(v, bitsize)
}

// Iterate through the array
numElem := array.Len() / bitsize
for i := 0; i < numElem; i++ {
  value := array.Get(i*bitsize, bitsize)

  // Do something useful with value
}
```

## API Reference

Godoc documentation can be found [here](https://godoc.org/github.com/robskie/bit).

## Benchmarks

I used a Core i5 running at 2.3GHz for these benchmarks. I used different bit
sizes in measuring the running time for bit array Add and Get methods as shown
in BenchmarkAddXX and BenchmarkGetXX where XX is the bit size. BenchmarkAddRand
and BenchmarkGetRand uses random bit sizes. Both BenchmarkGetXX and
BenchmarkGetRand measures the running time of Get using consecutive indices
while BenchmarkGetRandIdx measures the running time of Get using random indices
and bit sizes.

Here are the results by running ```go test github.com/robskie/bit -bench=.*```
from terminal.

```
BenchmarkAdd7           100000000           12.2 ns/op
BenchmarkAdd15          100000000           14.6 ns/op
BenchmarkAdd31          100000000           19.7 ns/op
BenchmarkAdd63          50000000            24.7 ns/op
BenchmarkAddRand        100000000           25.8 ns/op
BenchmarkInsertRandIdx  20000000            63.1 ns/op
BenchmarkGet7           100000000           22.5 ns/op
BenchmarkGet15          100000000           23.4 ns/op
BenchmarkGet31          50000000            24.3 ns/op
BenchmarkGet63          50000000            24.7 ns/op
BenchmarkGetRand        50000000            29.2 ns/op
BenchmarkGetRandIdx     30000000            53.1 ns/op
BenchmarkMSBIndex       300000000            4.96 ns/op
BenchmarkPopCount       500000000            3.51 ns/op
BenchmarkRank           200000000            6.69 ns/op
BenchmarkSelect         100000000           13.3 ns/op
```
