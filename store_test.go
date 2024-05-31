package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store, errCr := NewKVStore[string](_pathShards)
	require.NoError(t, errCr)
	require.NotNil(t, store)

	n := 7000

	for i := range n {
		value := strconv.Itoa(i + 1)

		require.NoError(t,
			store.Set(
				value,
				value,
			),
		)
	}

	store.FlushMemoryData()

	for i := range n {
		value := strconv.Itoa(i + 1)

		reconstructedValue, errGet := store.Get(
			value,
		)
		require.NoError(t, errGet, i)
		require.Equal(t,
			strconv.Itoa(i+1),
			reconstructedValue,
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkStoreSet-16    	   10000	    880386 ns/op	  120862 B/op	    2551 allocs/op
func BenchmarkStoreSet(b *testing.B) {
	store, errCr := NewKVStore[string](_pathShards)
	require.NoError(b, errCr)
	require.NotNil(b, store)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value := strconv.Itoa(i)

		store.Set(
			value,
			value,
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkStoreParallelSet-16    	   10000	    189949 ns/op	    7974 B/op	     204 allocs/op
func BenchmarkStoreParallelSet(b *testing.B) {
	store, errCr := NewKVStore[string](_pathShards)
	require.NoError(b, errCr)
	require.NotNil(b, store)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for i := 0; pb.Next(); i++ {
				value := strconv.Itoa(i)

				require.NoError(b,
					store.Set(
						value,
						value,
					),
				)
			}
		},
	)
}

func BenchmarkStoreGet(b *testing.B) {
	store, errCr := NewKVStore[string](_pathShards)
	require.NoError(b, errCr)
	require.NotNil(b, store)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value := strconv.Itoa(i)

		require.NoError(b,
			store.Set(
				value,
				value,
			),
		)
	}

	store.FlushMemoryData()

	for i := 0; i < b.N; i++ {
		value := strconv.Itoa(i)

		reconstructedValue, errGet := store.Get(
			value,
		)
		require.NoError(b, errGet)
		require.Equal(b, value, reconstructedValue)
	}
}
