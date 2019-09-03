package process

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"lsb_encoder/pkg/encoders"
	"os"
	"path/filepath"
)

// EncodeSrcFile does...
func EncodeSrcFile(conf EncodeConfig) error {
	var msg string
	var loadedImage image.Image
	if conf.MsgSrc == "stdin" {
		// Pull the message from Stdin
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Error reading from Stdin: (%v)", err)
		}
		msg = string(bytes)
	} else if conf.MsgSrc != "text" {
		// Read the message from the filepath
		bytes, err := ioutil.ReadFile(conf.MsgSrc)
		if err != nil {
			return fmt.Errorf("Error reading Message input file: (%v)", err)
		}
		msg = string(bytes)
	} else {
		// Accept the message passed via flag
		msg = conf.Msg
	}
	// Apply any Pre Encoding to the secret message
	if len(conf.PreEnc) != 0 {
		conf.Msg = encoders.ApplyPreEncoding(msg, conf.PreEnc)
	}
	// Read the Source file
	sourceFile, err := os.Open(conf.Src)
	// source, err := ioutil.ReadFile(conf.Src)
	if err != nil {
		return fmt.Errorf("Error reading Source file: (%v)", err)
	}
	defer sourceFile.Close()
	// reader = base64.NewDecoder(base64.StdEncoding, reader.(io.Reader))
	_, format, err := image.Decode(sourceFile)
	if err != nil {
		return fmt.Errorf("Error decoding source file: (%v)", err)
	}
	// Reset the file's reader to beginning
	sourceFile.Seek(0, 0)
	if format == "png" {
		loadedImage, err = png.Decode(sourceFile)
		if err != nil {
			return fmt.Errorf("Error decoding PNG file: (%v)", err)
		}
		embedded, err := EmbedMsgInImage(conf.Msg, format, loadedImage)
		if err != nil {
			return fmt.Errorf("Error embedding message in file: (%v)", err)
		}
		newFile, err := os.Create(filepath.Join(conf.Out, "output."+format))
		if err != nil {
			return fmt.Errorf("Error creating output file: (%v)", err)
		}
		defer newFile.Close()
		err = png.Encode(newFile, embedded)
		if err != nil {
			return fmt.Errorf("Error encoding new PNG image: (%v)", err)
		}
		return nil
	} else if format == "jpeg" {
		loadedImage, err = jpeg.Decode(sourceFile)
		if err != nil {
			return fmt.Errorf("Error decoding JPEG file: (%v)", err)
		}
		embedded, err := EmbedMsgInImage(conf.Msg, format, loadedImage)
		if err != nil {
			return fmt.Errorf("Error embedding message in file: (%v)", err)
		}
		newFile, err := os.Create(filepath.Join(conf.Out, "output."+format))
		if err != nil {
			return fmt.Errorf("Error creating output file: (%v)", err)
		}
		defer newFile.Close()
		err = jpeg.Encode(newFile, embedded, &jpeg.Options{Quality: 100})
		if err != nil {
			return fmt.Errorf("Error encoding new JPEG image: (%v)", err)
		}
		return nil
	} else if format == "gif" {
		loadedGIF, err := gif.DecodeAll(sourceFile)
		if err != nil {
			return fmt.Errorf("Error decoding GIF file: (%v)", err)
		}
		embedded, err := EmbedMsgInGIF(conf.Msg, format, loadedGIF)
		if err != nil {
			return fmt.Errorf("Error embedding message in file: (%v)", err)
		}
		newFile, err := os.Create(filepath.Join(conf.Out, "output."+format))
		if err != nil {
			return fmt.Errorf("Error creating output file: (%v)", err)
		}
		defer newFile.Close()
		err = gif.Encode(newFile, embedded, &gif.Options{})
	} else if format == "bmp" {
		// Do something else?
	} else {
		// ?
	}
	return nil
}

// DecodeSrcFile does...
func DecodeSrcFile(conf DecodeConfig) error {
	return nil
}
