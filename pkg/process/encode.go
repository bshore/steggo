package process

import (
	"io/ioutil"
	"os"
)

// EncodeSrcFile does...
func EncodeSrcFile(conf EncodeConfig) error {
	var msg string
	if conf.MsgSrc == "stdin" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		msg = string(bytes)
	} else if conf.MsgSrc != "file" {
		bytes, err := ioutil.ReadFile(conf.Msg)
		if err != nil {
			return err
		}
		msg = string(bytes)
	} else {
		msg = conf.Msg
	}
	if len(conf.PreEnc) != 0 {
		conf.Msg = encoders.ApplyPreEncoding(msg, conf.PreEnc)
	}
}
