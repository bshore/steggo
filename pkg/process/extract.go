package process

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"strconv"
	"strings"
)

// ExtractMsgFromImage takes an Image that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromImage(secret *Secret, file image.Image) (*Secret, error) {
	var err error
	var size int
	var headStr, msgStr string
	var headerFound bool
	var header Header
	bounds := file.Bounds()
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := file.At(x, y).RGBA()
			if headerFound {
				if len(msgStr) < size {
					msgStr, err = extract16BitColor([]uint32{r, g, b, a}, msgStr)
				}
				continue
			} else {
				// Build up headStr until it can be Unmarshaled
				headStr, err = extract16BitColor([]uint32{r, g, b, a}, headStr)
				if err != nil {
					return nil, fmt.Errorf("Error extracting from pixel: (%v)", err)
				}
				fmt.Println(headStr)
				err = json.Unmarshal([]byte(headStr), header)
				if err == nil {
					headerFound = true
					i, err := strconv.ParseInt(header.Size, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("Error parsing secret size from header: (%v)", err)
					}
					size = int(i)
				}
			}
		}
	}
	secret.DataHeader = header
	secret.Message = []byte(msgStr)
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

func extract16BitColor(vals []uint32, msg string) (string, error) {
	var bits []string
	for _, v := range vals {
		colorVal := strconv.FormatUint(uint64(v), 2)
		colorBits := strings.Split(ZeroPadLeft(colorVal, 16), "")
		bits = append(bits, colorBits[len(colorBits)-2])
		bits = append(bits, colorBits[len(colorBits)-1])
	}
	msgStr := strings.Join(bits, "")
	char, err := strconv.ParseInt(msgStr, 2, 64)
	if err != nil {
		return "", err
	}
	// I have no idea how to handle this...
	msg += strconv.QuoteRune(int32(char))
	return msg, nil
}
