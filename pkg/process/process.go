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

	"golang.org/x/image/bmp"
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
	} else if conf.MsgSrc != "text" && conf.MsgSrc != "" {
		// Read the message from the filepath
		bytes, err := ioutil.ReadFile(conf.MsgSrc)
		if err != nil {
			return fmt.Errorf("Error reading Message input file: (%v)", err)
		}
		msg = string(bytes)
	} else if conf.MsgSrc == "text" {
		// Accept the message passed via flag
		msg = conf.Msg
	} else {
		return fmt.Errorf("Error determining MsgSrc: %v", conf.MsgSrc)
	}
	conf.Msg = msg
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

	// Do all the work, can totally be cleaned up & refactored
	// Maybe format becomes an iota enum so this can be a type switch vs if, else if, etc...
	// and each format type has Decode, Embed, and Encode methods?
	// Is that cleaner?
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
		err = gif.EncodeAll(newFile, embedded)
		if err != nil {
			return err
		}
	} else if format == "bmp" {
		loadedImage, err := bmp.Decode(sourceFile)
		if err != nil {
			return fmt.Errorf("Error decoding BMP file: (%v)", err)
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
		err = bmp.Encode(newFile, embedded)
		if err != nil {
			return fmt.Errorf("Error encoding new JPEG image: (%v)", err)
		}
		return nil
	} else {
		return fmt.Errorf("Unrecognized file format: %v", format)
	}
	return nil
}

// DecodeSrcFile does...
func DecodeSrcFile(conf DecodeConfig) error {
	return nil
}
