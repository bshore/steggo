package process

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"os"
	"strconv"
)

// ExtractMsgFromImage takes an Image that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromImage(secret *Secret, file image.Image) (*Secret, error) {
	// var err error
	var size int
	var headBytes, msgBytes []byte
	var headerFound bool
	var header Header
	var poo int
	bounds := file.Bounds()
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			poo++
			if poo == 2 {
				os.Exit(0)
			}
			r, g, b, a := file.At(x, y).RGBA()
			if headerFound {
				if len(msgBytes) < size {
					msgbyte := extract16BitColor(r, g, b, a)
					msgBytes = append(msgBytes, msgbyte)
				}
				continue
			} else {
				// Build up headBytes until it can be Unmarshaled
				headbyte := extract16BitColor(r, g, b, a)
				headBytes = append(headBytes, headbyte)
				err := json.Unmarshal(headBytes, &header)
				if err == nil {
					headerFound = true
					fmt.Println("Header found! ", string(headBytes))
					i, err := strconv.ParseInt(header.Size, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("Error parsing secret size from header: (%v)", err)
					}
					size = int(i)
				}
			}
		}
	}
	fmt.Println("Header: ", headBytes[:20])
	secret.DataHeader = header
	secret.Message = msgBytes
	return secret, nil
}

// ExtractMsgFromGif takes a GIF that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromGif(secret *Secret, file *gif.GIF) (string, error) {
	return "", nil
}

func extract8BitColor(vals []uint32, msg string) (string, error) {
	return "", nil
}

func extract16BitColor(r, g, b, a uint32) byte {
	// Get last 2 bits of each color to reconstruct a message byte
	rBits := r & 3
	gBits := g & 3
	bBits := b & 3
	aBits := a & 3

	var newByte uint32
	// Assign last 2 bits and shift left for each color pixel
	newByte = newByte | (rBits << 0) // ------00
	newByte = newByte << 2           // ----00--
	newByte = newByte | (gBits << 0) // ----0000
	newByte = newByte << 2           // --0000--
	newByte = newByte | (bBits << 0) // --000000
	newByte = newByte << 2           // 000000--
	newByte = newByte | (aBits << 0) // 00000000
	fmt.Printf("%08b\n", uint8(newByte))
	return uint8(newByte)

	// var bits []string
	// 	colorVal := strconv.FormatUint(uint64(v), 2)
	// 	colorBits := strings.Split(ZeroPadLeft(colorVal, 16), "")
	// 	twoBit := colorBits[len(colorBits)-2]
	// 	oneBit := colorBits[len(colorBits)-1]
	// 	bits = append(bits, twoBit, oneBit)

	// msgStr := strings.Join(bits, "")
	// char, err := strconv.ParseUint(msgStr, 2, 8)
	// fmt.Println(char, msgStr)
	// if err != nil {
	// 	return 0x00, err
	// }
	// return uint8(char), nil
}
