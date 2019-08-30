package process

import (
	"fmt"
	"io/ioutil"
	"lsb_encoder/pkg/encoders"
	"os"
	"path/filepath"
	"strings"
)

// EncodeSrcFile does...
func EncodeSrcFile(conf EncodeConfig) error {
	var msg string
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
	source, err := ioutil.ReadFile(conf.Src)
	if err != nil {
		return fmt.Errorf("Error reading Source file: (%v)", err)
	}
	split := strings.Split(filepath.Base(conf.Src), ".")
	ext := split[len(split)-1]
	embedded, err := EmbedMsgInFile(conf.Msg, ext, source)
	if err != nil {
		return err
	}
	err = WriteEmbeddedFile(embedded, conf.Out, ext)
	if err != nil {
		return err
	}
	return nil
}

// DecodeSrcFile does...
func DecodeSrcFile(conf DecodeConfig) error {
	return nil
}
