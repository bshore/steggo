package embedder

// import (
// 	"fmt"
// 	"io"
// 	"os"

// 	"golang.org/x/image/bmp"
// )

// func ProcessBMP(msg, dest string, src io.Reader) error {
// 	loadedImage, err := bmp.Decode(src)
// 	if err != nil {
// 		return fmt.Errorf("error decoding PNG file: %v", err)
// 	}
// 	embedded, err := EmbedMsgInImage(msg, loadedImage)
// 	if err != nil {
// 		return fmt.Errorf("error embedding message in file: %v", err)
// 	}
// 	newFile, err := os.Create(dest)
// 	if err != nil {
// 		return fmt.Errorf("error creating output file: %v", err)
// 	}
// 	defer newFile.Close()

// 	err = bmp.Encode(newFile, embedded)
// 	if err != nil {
// 		return fmt.Errorf("error encoding new PNG image: %v", err)
// 	}
// 	return nil
// }
