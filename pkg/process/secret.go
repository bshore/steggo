package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lsb_encoder/pkg/encoders"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
	This file contains the Secret struct definition and functions related to it.
*/

// Header is a prefix to to identify information used during extraction.
type Header struct {
	// Size is the Secret.Size but stored as a JSON string.
	Size string `json:"size"`
	// Type is the Secret.Type but stored as a JSON string.
	Type string `json:"type"`
	// Enc is Secret.PreEncoding but stored as a JSON string.
	Enc []string `json:"enc"`
}

// Secret stores information about the input data, and the data itself
type Secret struct {
	// SourcePath holds the path to the source image.
	// This is the image the secret will be embedded into.
	SourcePath string
	// OutputDir is where to put the output image once the
	// secret has been embedded inside it.
	OutputDir string
	// SecretPath holds the path to the 'secret' file to be
	// embedded inside the image.
	SecretPath string
	// Type holds the type of the input.
	// Since the input data might be a message or another
	// file, here is where a file extension
	// or 'text' would be stored.
	Type string
	// Size holds how large the input is.
	// The input needs to be able to fit inside
	// the source file and serves as a stopping point
	// during extraction.
	Size int
	// PreEncoding holds a slice of encoding methods applied.
	// Extraction functions need to be aware of pre encoding
	// so the input data can be fully decoded.
	PreEncoding []encoders.EncType
	// DataHeader holds information about the secret that assists
	// during extraction.
	// Header is included at the beginning of Data as a key/value store
	DataHeader Header
	// Data holds the contents of what is being embedded into a file.
	Data [][]string
	// Message holds the contents of what had been extracted froma  file.
	Message []byte
}

// FormatSecretData formats the secret's header & secret itself into the Secret.Data field
func (s *Secret) FormatSecretData(msg string) error {
	header, err := json.Marshal(s.DataHeader)
	if err != nil {
		return err
	}
	msgBytes := []byte(string(header) + msg)
	var bitArr [][]string
	for _, b := range msgBytes {
		binStr := strconv.FormatInt(int64(b), 2)
		bits := strings.Split(ZeroPadLeft(binStr, 8), "")
		for i, bit := range bits {
			// Current iteration is odd ?
			if i%2 != 0 && i != 0 {
				// Grab bits in pairs (01, 23, 45, 67)
				// ------1- two's bit
				two := bit
				// -------1 one's bit
				one := bits[i-1]
				bitArr = append(bitArr, []string{two, one})
			}
		}
	}
	s.Size = len(bitArr)
	s.Data = bitArr
	return nil
}

// ParseEmbedSecret takes Flags info and turns it into a Secret struct for embedding
func ParseEmbedSecret(f *Flags) (*Secret, error) {
	s := &Secret{}
	var msg string
	var pre []encoders.EncType
	var enc []string
	var header Header
	// Check for the source of the secret message/file
	if f.Stdin {
		header.Type = "stdin"
		s.Type = "stdin"
		// Pull the message from Stdin
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return s, fmt.Errorf("Error reading from Stdin: (%v)", err)
		}
		msg = string(bytes)
	} else if f.MessageFile != "" {
		s.SecretPath = f.MessageFile
		// Read the message from the filepath
		bytes, err := ioutil.ReadFile(f.MessageFile)
		if err != nil {
			return s, fmt.Errorf("Error reading Message input file: (%v)", err)
		}
		msg = string(bytes)
		header.Type = filepath.Ext(f.MessageFile)
	} else if f.Text != "" {
		header.Type = "text"
		s.Type = "text"
		msg = f.Text
	}
	// Check for Pre Encoding
	if f.Complex != "" {
		typs := strings.Split(f.Complex, ",")
		if len(typs) > 5 {
			return s, fmt.Errorf("Too many --complex pre-encoders, maximum of 5")
		}
		for _, typ := range typs {
			encoder, err := encoders.EncTypeFromString(typ)
			if err != nil {
				return s, err
			}
			pre = append(pre, encoder)
			enc = append(enc, encoder.String())
		}
	} else {
		if f.Rot13 {
			pre = append(pre, encoders.R13)
			enc = append(enc, encoders.R13.String())
		}
		if f.Base16 {
			pre = append(pre, encoders.B16)
			enc = append(enc, encoders.B16.String())
		}
		if f.Base32 {
			pre = append(pre, encoders.B32)
			enc = append(enc, encoders.B32.String())
		}
		if f.Base64 {
			pre = append(pre, encoders.B64)
			enc = append(enc, encoders.B64.String())
		}
		if f.Base85 {
			pre = append(pre, encoders.B85)
			enc = append(enc, encoders.B85.String())
		}
	}
	if len(enc) != 0 {
		header.Enc = enc
	}
	s.PreEncoding = pre
	// Apply any Pre Encoding to the secret message
	if len(s.PreEncoding) != 0 {
		msg = encoders.ApplyPreEncoding(msg, s.PreEncoding)
	}
	// Size of the secret is length of msg * 4
	// because each byte (8) is split into pairs of 2
	header.Size = strconv.FormatInt(int64((len(msg) * 4)), 10)
	s.DataHeader = header
	err := s.FormatSecretData(msg)
	if err != nil {
		return nil, err
	}
	s.SourcePath = f.SrcFile
	s.OutputDir = f.OutputDir
	return s, nil
}

// ParseExtractSecret takes Flags info and turns it into a Secret struct for extracting
func ParseExtractSecret(f *Flags) (*Secret, error) {
	s := &Secret{}
	s.SourcePath = f.SrcFile
	s.OutputDir = f.OutputDir
	// The rest is figured out later?
	return s, nil
}

// ReconstructMessage takes bit pairs and reconstructs a message
// back into its source string
func ReconstructMessage(bitArr [][]string) string {
	return ""
}