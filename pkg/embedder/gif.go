package embedder

// import (
// 	"fmt"
// 	"image/gif"
// 	"io"
// 	"os"
// )

// func ProcessGIF(msg, dest string, src io.Reader) error {
// 	loadedImage, err := gif.Decode(src)
// 	if err != nil {
// 		return fmt.Errorf("error decoding PNG file: %v", err)
// 	}
// 	embedded, err := EmbedMsgInGIF(msg, loadedImage)
// 	if err != nil {
// 		return fmt.Errorf("error embedding message in file: %v", err)
// 	}
// 	newFile, err := os.Create(dest)
// 	if err != nil {
// 		return fmt.Errorf("error creating output file: %v", err)
// 	}
// 	defer newFile.Close()

// 	err = gif.EncodeAll(newFile, embedded)
// 	if err != nil {
// 		return fmt.Errorf("error encoding new PNG image: %v", err)
// 	}
// 	return nil
// }
