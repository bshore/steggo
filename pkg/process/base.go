package process

import (
	"lsb_encoder/pkg/encoders"
	"strings"
)

// Flags holds the types of flags allowed by the script
type Flags struct {
	SrcFile     string
	OutputFile  string
	Text        string
	Stdin       bool
	MessageFile string
	Decode      bool
	Rot13       bool
	Base16      bool
	Base32      bool
	Base64      bool
	Base85      bool
	Complex     string
}

// ToEncConf takes Flags info and turns it into a config struct for Encoding
func (f Flags) ToEncConf() (*EncodeConfig, error) {
	var e EncodeConfig
	var p []encoders.EncType
	if f.Stdin {
		e.MsgSrc = "stdin"
	} else if f.MessageFile != "" {
		e.MsgSrc = f.MessageFile
	} else if f.Text != "" {
		e.MsgSrc = "text"
		e.Msg = f.Text
	}
	// Check Pre Encoding
	if f.Complex != "" {
		typs := strings.Split(f.Complex, ",")
		for _, typ := range typs {
			enc, err := encoders.EncTypeFromString(typ)
			if err != nil {
				return nil, err
			}
			p = append(p, enc)
		}
	} else {
		if f.Rot13 {
			p = append(p, encoders.R13)
		}
		if f.Base16 {
			p = append(p, encoders.B16)
		}
		if f.Base32 {
			p = append(p, encoders.B32)
		}
		if f.Base64 {
			p = append(p, encoders.B64)
		}
		if f.Base85 {
			p = append(p, encoders.B85)
		}
	}
	return &EncodeConfig{
		Src:    f.SrcFile,
		Out:    f.OutputFile,
		PreEnc: p,
	}, nil
}

// EncodeConfig stores config options for handling Encoding
type EncodeConfig struct {
	Src    string
	Out    string
	MsgSrc string
	Msg    string
	PreEnc []encoders.EncType
}
