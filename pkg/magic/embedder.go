package magic

import (
	"fmt"
	"strconv"
	"strings"
)

// EmbedMsgInPNG embeds the message into the PNG file's bytes
func EmbedMsgInPNG(bits [][]string, file []byte) ([]byte, error) {
	firstWB, err := getPngFirstWritableByte(file)
	if err != nil {
		return nil, err
	}
	var bitsIndex int
	// Copy the file up to the point of just before first writeable byte
	embeddedFile := file[:firstWB-1]
	for i := firstWB; i < len(file)-firstWB; i++ {
		if bitsIndex+1 < len(bits) {
			newByte, err := writeBitsToFile(bits[bitsIndex], file[i])
			if err != nil {
				return nil, err
			}
			bitsIndex++
			embeddedFile = append(embeddedFile, newByte)
		} else {
			embeddedFile = append(embeddedFile, file[i:]...)
			return embeddedFile, nil
		}
	}
	return nil, fmt.Errorf("Secret message is too long, it won't fit into the Source file")
}

// EmbedMsgInJPEG embeds the message into the JPEG file's bytes
func EmbedMsgInJPEG(bits [][]string, file []byte) ([]byte, error) {
	firstWB, err := getJpegFirstWriteableByte(file)
	if err != nil {
		return nil, err
	}
	_ = firstWB
	return nil, nil
}

// EmbedMsgInBMP embeds the message into the BMP file's bytes
func EmbedMsgInBMP(bits [][]string, file []byte) ([]byte, error) {
	firstWB, err := getBmpFirstWriteableByte(file)
	if err != nil {
		return nil, err
	}
	_ = firstWB
	return nil, nil
}

// EmbedMsgInGIF embeds the message into the GIF file's bytes
func EmbedMsgInGIF(bits [][]string, file []byte) ([]byte, error) {
	writeableBytes, err := getGifWriteableBytes(file)
	if err != nil {
		return nil, err
	}
	_ = writeableBytes
	return nil, nil
}

func writeBitsToFile(bits []string, fileByte byte) (byte, error) {
	byteVal := strconv.FormatInt(int64(fileByte), 2)
	fileBits := strings.Split(ZeroPadLeft(byteVal), "")
	// Encode the message bits into the least significant bits
	// of the fileBits
	fileBits[6] = bits[0]
	fileBits[7] = bits[1]
	byteStr := strings.Join(fileBits, "")
	newByte, err := strconv.ParseInt(byteStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return byte(newByte), nil
}
