package magic

import (
	"fmt"
	"strings"
)

// getPngFirstWritableByte inspects the file header to find the index
// of the first writeable byte that won't destroy the file.
func getPngFirstWritableByte(src []byte) (int, error) {
	var headerStr strings.Builder
	// Build a string until 'IDAT' is present, the next byte should be writable
	for i, b := range src {
		headerStr.WriteString(string(b))
		if strings.Contains(headerStr.String(), "IDATx") {
			return i + 1000, nil
		}
	}
	return 0, fmt.Errorf("could not find IDATx marker for first writable byte")
}

// getJpegFirstWriteableByte inspects the file header to find the index
// of the first writeable byte that won't destroy the file.
func getJpegFirstWriteableByte(src []byte) (int, error) {
	var header []byte
	var nextWrite int
	for i, b := range src {
		header = append(header, b)
		// A JPEG Start Of Scan (SOS) marker are the exact bytes FF DA
		if header[i] == 0xDA && header[i-1] == 0xFF {
			nextWrite = i + 11
			// A Start Of Scan marker contains 10 bytes of metadata
			// The 11th byte should be writeable
			// Issue: some JPEG files contain more than one SOS marker
			// First could be image thumbnail?
			// Second is actual image?
		}
		return nextWrite, nil
	}
	return 0, fmt.Errorf("could not find SOS marker for first writeable byte")
}

// getBmpFirstWriteableByte inspects the file header to find the index
// of the first writeable byte that won't destroy the file
func getBmpFirstWriteableByte(src []byte) (int, error) {
	// Sum of 4 bytes starting at 0A (index 9)
	offset := int(src[9]) + int(src[10]) + int(src[11]) + int(src[12])
	return offset, fmt.Errorf("")
}
