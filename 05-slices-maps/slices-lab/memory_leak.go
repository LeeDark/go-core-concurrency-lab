package sliceslab

// TakeSmallPart returns a slice containing the first 10 bytes from a 100 MB byte slice.
func TakeSmallPart() []byte {
	big := make([]byte, 100<<20) // 100 MB
	for i := range big {
		big[i] = byte(i)
	}

	return big[:10]
}

// TakeSmallPartSafeCopy copies a small portion of a larger byte slice into a new slice and returns it.
func TakeSmallPartSafeCopy() []byte {
	big := make([]byte, 100<<20)
	for i := range big {
		big[i] = byte(i)
	}

	small := make([]byte, 10)
	copy(small, big[:10])

	return small
}

// TakeSmallPartSafeAppend safely extracts a small portion of a large slice and returns a new slice to avoid memory overhead.
func TakeSmallPartSafeAppend() []byte {
	big := make([]byte, 100<<20)
	for i := range big {
		big[i] = byte(i)
	}

	return append([]byte(nil), big[:10]...)
}
