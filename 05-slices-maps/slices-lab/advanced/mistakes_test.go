package slicesadvanced

import "testing"

var mistake21ConvertResult []float32

func BenchmarkMistake21Convert_WithoutPreallocation(b *testing.B) {
	foos := make([]int, 1_000)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mistake21ConvertResult = convertEmptySlice(foos)
	}
}

func BenchmarkMistake21Convert_GivenCapacity(b *testing.B) {
	foos := make([]int, 1_000)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mistake21ConvertResult = convertGivenCapacity(foos)
	}
}

func BenchmarkMistake21Convert_GivenLength(b *testing.B) {
	foos := make([]int, 1_000)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mistake21ConvertResult = convertGivenLength(foos)
	}
}
