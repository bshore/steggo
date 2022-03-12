package process

// import (
// 	"fmt"
// 	"image"
// 	"image/gif"
// 	"image/jpeg"
// 	"image/png"
// 	"os"
// 	"path/filepath"

// 	"golang.org/x/image/bmp"
// )

// // Extract extracts the Secret's data and writes it to the output directory
// func Extract(secret *Secret) error {
// 	var loadedImage image.Image
// 	// Read the Source file
// 	sourceFile, err := os.Open(secret.SourcePath)
// 	if err != nil {
// 		return fmt.Errorf("Error reading Source file: (%v)", err)
// 	}
// 	defer sourceFile.Close()
// 	_, format, err := image.Decode(sourceFile)
// 	if err != nil {
// 		return fmt.Errorf("Error decoding source file: (%v)", err)
// 	}
// 	// Reset the file's reader to beginning
// 	sourceFile.Seek(0, 0)
// 	if format == "png" {
// 		// =====
// 		// Handle PNG
// 		// ====
// 		loadedImage, err = png.Decode(sourceFile)
// 		if err != nil {
// 			return fmt.Errorf("Error decoding PNG file: (%v)", err)
// 		}
// 		extracted, err := ExtractMsgFromImage(secret, loadedImage)
// 		if err != nil {
// 			return fmt.Errorf("Error extracting message from file: (%v)", err)
// 		}
// 		err = WriteFile(extracted.Message, extracted.OutputDir, extracted.DataHeader.Type)
// 		if err != nil {
// 			return fmt.Errorf("Error creating output file: (%v)", err)
// 		}
// 		return nil
// 	} else if format == "jpeg" {
// 		// =====
// 		// Handle JPEG
// 		// ====
// 		loadedImage, err = jpeg.Decode(sourceFile)
// 		if err != nil {
// 			return fmt.Errorf("Error decoding JPEG file: (%v)", err)
// 		}
// 		extracted, err := ExtractMsgFromImage(secret, loadedImage)
// 		if err != nil {
// 			return fmt.Errorf("Error extracting message from file: (%v)", err)
// 		}
// 		err = WriteFile(extracted.Message, extracted.OutputDir, extracted.DataHeader.Type)
// 		if err != nil {
// 			return fmt.Errorf("Error creating output file: (%v)", err)
// 		}
// 		return nil
// 	} else if format == "bmp" {
// 		// =====
// 		// Handle BMP
// 		// ====
// 		loadedImage, err = bmp.Decode(sourceFile)
// 		if err != nil {
// 			return fmt.Errorf("Error decoding BMP file: (%v)", err)
// 		}
// 		extracted, err := ExtractMsgFromImage(secret, loadedImage)
// 		if err != nil {
// 			return fmt.Errorf("Error extracting message from file: (%v)", err)
// 		}
// 		err = WriteFile(extracted.Message, extracted.OutputDir, extracted.DataHeader.Type)
// 		if err != nil {
// 			return fmt.Errorf("Error creating output file: (%v)", err)
// 		}
// 		return nil
// 	} else if format == "gif" {
// 		// =====
// 		// Handle GIF
// 		// ====
// 		loadedGIF, err := gif.DecodeAll(sourceFile)
// 		if err != nil {
// 			return fmt.Errorf("Error decoding GIF file: (%v)", err)
// 		}
// 		extracted, err := ExtractMsgFromGif(secret, loadedGIF)
// 		if err != nil {
// 			return fmt.Errorf("Error extracting message from file: (%v)", err)
// 		}
// 		err = WriteFile(extracted.Message, extracted.OutputDir, extracted.DataHeader.Type)
// 		if err != nil {
// 			return fmt.Errorf("Error creating output file: (%v)", err)
// 		}
// 		return nil
// 	}
// 	// =====
// 	// Handle ???
// 	// ====
// 	return fmt.Errorf("Unsupported source file format: %v", format)
// }
