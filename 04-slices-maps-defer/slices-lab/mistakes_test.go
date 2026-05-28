package sliceslab

import "testing"

var mistake21ConvertResult []float32

func BenchmarkMistake21Convert_EmptySlice(b *testing.B) {
	foos := []int{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mistake21ConvertResult = convertEmptySlice(foos)
	}
}

func BenchmarkMistake21Convert_GivenCapacity(b *testing.B) {
	foos := []int{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mistake21ConvertResult = convertGivenCapacity(foos)
	}
}

func BenchmarkMistake21Convert_GivenLength(b *testing.B) {
	foos := []int{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mistake21ConvertResult = convertGivenLength(foos)
	}
}
