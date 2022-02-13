package process

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"lsb_encoder/pkg/encoders"
// 	"os"
// 	"path/filepath"
// 	"strings"
// )

// /*
// 	This file contains the Secret struct definition and functions related to it.
// */

// // Secret stores information about the input data, and the data itself
// type Secret struct {
// 	// SourcePath holds the path to the source image.
// 	// This is the image the secret will be embedded into.
// 	SourcePath string
// 	// OutputDir is where to put the output image once the
// 	// secret has been embedded inside it.
// 	OutputDir string
// 	// SecretPath holds the path to the 'secret' file to be
// 	// embedded inside the image.
// 	SecretPath string
// 	// Type holds the type of the input.
// 	// Since the input data might be a message or another
// 	// file, here is where a file extension
// 	// or 'text' would be stored.
// 	Type string
// 	// Size holds how large the input is.
// 	// The input needs to be able to fit inside
// 	// the source file and serves as a stopping point
// 	// during extraction.
// 	Size int
// 	// PreEncoding holds a slice of encoding methods applied.
// 	// Extraction functions need to be aware of pre encoding
// 	// so the input data can be fully decoded.
// 	PreEncoding []encoders.EncType
// 	// DataHeader holds information about the secret that assists
// 	// during extraction.
// 	// Header is included at the beginning of Data as a key/value store
// 	DataHeader Header
// 	// Data holds the contents of what is being embedded into a file.
// 	Data []byte
// 	// Message holds the contents of what had been extracted froma  file.
// 	Message []byte
// }

// // FormatSecretData formats the secret's header & secret itself into the Secret.Data field
// func FormatSecretData(msg string) error {
// 	header, err := json.Marshal(s.DataHeader)
// 	if err != nil {
// 		return err
// 	}
// 	msgBytes := []byte(string(header) + msg)
// 	var bitArr []byte
// 	for _, b := range msgBytes {
// 		// Get bit values in a group of 2-3-3 (R-G-B)
// 		// sevenEight uses & 131 to set the 128 bit, so embedding knows to zero out
// 		// the last 2 bits of a color value, instead of zeroing out the last 3 bits
// 		sevenEight := (b >> 6) & 131 // shifts bb------ to ------bb and gets last 2 bits value
// 		fourFiveSix := (b >> 3) & 7  // shifts --bbb--- to -----bbb and gets last 3 bits value
// 		oneTwoThree := b & 7         // just gets -----bbb last 3 bits value
// 		bitArr = append(bitArr, uint8(sevenEight), uint8(fourFiveSix), uint8(oneTwoThree))
// 	}
// 	s.Data = bitArr
// 	s.Message = msgBytes
// 	return nil
// }

// // ParseEmbedSecret takes Flags info and turns it into a Secret struct for embedding
// func ParseEmbedSecret(f *Flags) (*Secret, error) {
// 	s := &Secret{}
// 	var msg string
// 	var pre []encoders.EncType
// 	var enc []string
// 	var header Header
// 	header.BitOpt = f.BitOpt
// 	// Check for the source of the secret message/file
// 	if f.Stdin {
// 		header.Type = "stdin"
// 		s.Type = "stdin"
// 		// Pull the message from Stdin
// 		bytes, err := ioutil.ReadAll(os.Stdin)
// 		if err != nil {
// 			return s, fmt.Errorf("Error reading from Stdin: (%v)", err)
// 		}
// 		msg = string(bytes)
// 	} else if f.MessageFile != "" {
// 		s.SecretPath = f.MessageFile
// 		// Read the message from the filepath
// 		bytes, err := ioutil.ReadFile(f.MessageFile)
// 		if err != nil {
// 			return s, fmt.Errorf("Error reading Message input file: (%v)", err)
// 		}
// 		msg = string(bytes)
// 		header.Type = filepath.Ext(f.MessageFile)
// 	} else if f.Text != "" {
// 		header.Type = ".txt"
// 		s.Type = "txt"
// 		msg = f.Text
// 	}
// 	// Check for Pre Encoding
// 	if f.Complex != "" {
// 		typs := strings.Split(f.Complex, ",")
// 		if len(typs) > 5 {
// 			return s, fmt.Errorf("Too many --complex pre-encoders, maximum of 5")
// 		}
// 		for _, typ := range typs {
// 			encoder, err := encoders.EncTypeFromString(typ)
// 			if err != nil {
// 				return s, err
// 			}
// 			pre = append(pre, encoder)
// 			enc = append(enc, encoder.String())
// 		}
// 	} else {
// 		if f.Rot13 {
// 			pre = append(pre, encoders.R13)
// 			enc = append(enc, encoders.R13.String())
// 		}
// 		if f.Base16 {
// 			pre = append(pre, encoders.B16)
// 			enc = append(enc, encoders.B16.String())
// 		}
// 		if f.Base32 {
// 			pre = append(pre, encoders.B32)
// 			enc = append(enc, encoders.B32.String())
// 		}
// 		if f.Base64 {
// 			pre = append(pre, encoders.B64)
// 			enc = append(enc, encoders.B64.String())
// 		}
// 		if f.Base85 {
// 			pre = append(pre, encoders.B85)
// 			enc = append(enc, encoders.B85.String())
// 		}
// 	}
// 	if len(enc) != 0 {
// 		header.Enc = enc
// 	}
// 	s.PreEncoding = pre
// 	// Apply any Pre Encoding to the secret message
// 	if len(s.PreEncoding) != 0 {
// 		msg = encoders.ApplyPreEncoding(msg, s.PreEncoding)
// 	}
// 	// Size of the secret is length of msg * 4
// 	// because each byte (8) is split into pairs of 2
// 	header.Size = len(msg)
// 	s.DataHeader = header
// 	err := s.FormatSecretData(msg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s.SourcePath = f.SrcFile
// 	s.OutputDir = f.OutputDir
// 	return s, nil
// }

// // ParseExtractSecret takes Flags info and turns it into a Secret struct for extracting
// func ParseExtractSecret(f *Flags) (*Secret, error) {
// 	s := &Secret{}
// 	s.SourcePath = f.SrcFile
// 	s.OutputDir = f.OutputDir
// 	// The rest is figured out later?
// 	return s, nil
// }
