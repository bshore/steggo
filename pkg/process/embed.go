package process

import (
	"lsb_encoder/pkg/magic"
	"strings"
)

// EmbedMsgInFile takes the message string and embeds it
// in the source file's byte string using Least Significant Byte
func EmbedMsgInFile(msg, ext string, file []byte) ([]byte, error) {
	bitArr := magic.BreakupMessageBytes(msg)
	if strings.EqualFold(ext, "png") {
		newFile, err := magic.EmbedMsgInPNG(bitArr, file)
		if err != nil {
			return nil, err
		}
		return newFile, nil
	}
	// } else if (strings.EqualFold(ext, "jpg")) || (strings.EqualFold(ext, "jpeg")) {
	// 	newFile, err := magic.EmbedMsgInJPEG(bitArr, file)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return newFile, nil
	// } else if strings.EqualFold(ext, "bmp") {
	// 	newFile, err := magic.EmbedMsgInBMP(bitArr, file)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return newFile, nil
	// } else if strings.EqualFold(ext, "gif") {
	// 	newFile, err := magic.EmbedMsgInGIF(bitArr, file)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return newFile, nil
	// }
	return nil, nil
}
