package magic

// Block definitions found here: http://www.onicos.com/staff/iz/formats/gif.html
//
// In this use case, I'm just assuming the Local Color Table Flag is set
// and not counting those 0-3 bytes as writable

// getGifWriteableBytes inspects the file to find the indexes
// of writeable bytes that won't destroy the file.
// Since a GIF contains Frame Blocks, need to return a collection of
// startAt, endAt indexes of writeable bytes.
func getGifWriteableBytes(src []byte) ([][]int, error) {
	var build []byte
	var onBlock bool
	var blockStart int
	var writeableBytes [][]int
	for i, b := range src {
		build = append(build, b)
		if onBlock {
			// Read until Block Terminator ( 0x00 ) is followed by Image Separator ( 0x2C )
			if build[i-1] == 0x00 && build[i] == 0x2C {
				// Here is where I mean I'm just assuming the LCTF is set and assuming 3 bytes
				// 0x2C..<14 bytes block metadata>..<image data>..0x00..0x2C(next block)
				writeableBytes = append(writeableBytes, []int{blockStart + 14, i - 2})
				// End of a Block
				onBlock = false
			}
		}
		// Read until Block Terminator ( 0x00 ) is followed by Image Separator ( 0x2C )
		if build[i-1] == 0x00 && build[i] == 0x2C {
			// Beginning of a block
			onBlock = true
			blockStart = i
		}
	}
	return writeableBytes, nil
}
