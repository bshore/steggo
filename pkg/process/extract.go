package process

import (
	"fmt"
	"image"
)

// ExtractMsgFromImage takes an Image that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
func ExtractMsgFromImage(file image.Image) (*Header, []byte, error) {
	var headBytes, msgBytes []byte
	var headerFound bool
	var header = &Header{}
	bounds := file.Bounds()
	// For each vertical row
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// For each pixel in each row
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := file.At(x, y).RGBA()
			if headerFound {
				if len(msgBytes) < header.Size {
					msgbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
					msgBytes = append(msgBytes, msgbyte)
				}
				continue
			} else {
				// Build up headBytes until the full header is identified
				headbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
				headBytes = append(headBytes, headbyte)
				if header.Found(headBytes) {
					headerFound = true
				}
			}
		}
	}

	if !headerFound {
		return nil, nil, fmt.Errorf("failed to extract message")
	}

	return header, msgBytes, nil
}

// ExtractMsgFromGif takes a GIF that has had a message embedded
// inside it and extracts the message using Least Significant Bit(s)
// func ExtractMsgFromGif(secret *Secret, file *gif.GIF) (*Secret, error) {
// 	var err error
// 	var size int
// 	var headBits, msgBits []uint8
// 	var headBytes, headBitBytes, msgBytes, msgBitBytes []byte
// 	var headerFound bool
// 	var header Header
// 	// For each image frame
// 	for _, img := range file.Image {
// 		bounds := img.Bounds()
// 		// For each vertical row
// 		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 			// vor each pixel in each row
// 			for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 				r, g, b, _ := img.At(x, y).RGBA()
// 				if headerFound {
// 					if header.BitOpt == 2 {
// 						if len(msgBytes) < size {
// 							msgbyte := extractFromColor(uint8(r), uint8(g), uint8(b))
// 							msgBytes = append(msgBytes, msgbyte)
// 						}
// 					} else {
// 						if len(msgBytes) < size {
// 							newBits := extractBitFromColor(uint8(r), uint8(g), uint8(b))
// 							msgBits := append(msgBits, newBits...)
// 							if len(msgBitBytes) >= 8 {
// 								newByte := rebuildFromBits(msgBits[:8])
// 								msgBytes = append(msgBytes, newByte)
// 								msgBits = msgBits[8:]
// 							}
// 						}
// 					}
// 					continue
// 				} else {
// 					// Build up headBytes & headBits until it can be Unmarshaled
// 					headByte := extractFromColor(uint8(r), uint8(g), uint8(b))
// 					headBytes = append(headBytes, headByte)
// 					err = json.Unmarshal(headBytes, &header)
// 					if err == nil {
// 						headerFound = true
// 						size = header.Size
// 					}
// 					newBits := extractBitFromColor(uint8(r), uint8(g), uint8(b))
// 					headBits = append(headBits, newBits...)
// 					if len(headBits) >= 8 {
// 						newByte := rebuildFromBits(headBits[:8])
// 						headBitBytes = append(headBitBytes, newByte)
// 						err = json.Unmarshal(headBitBytes, &header)
// 						if err == nil {
// 							headerFound = true
// 						}
// 						headBits = headBits[8:]
// 					}
// 				}
// 			}
// 		}
// 	}
// 	secret.DataHeader = header
// 	secret.Message = msgBytes
// 	return secret, nil
// }
