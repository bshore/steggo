package main

import (
	"flag"
	"fmt"
	"log"
	"lsb_encoder/pkg/process"
	"os"
	"path/filepath"
)

var (
	srcFile = flag.String("srcfile", "", "Path to the source file to be messed with")
	outFile = flag.String("outfile", "", "Path to the output directory when finished")
	text    = flag.String("text", "", "The text string to encode into the file")
	stdIn   = flag.Bool("stdin", false, "T/F flag if passing through concat file and pipe | ")
	msgFile = flag.String("msgfile", "", "The path to a text file containing the message to be endcoded")
	decode  = flag.Bool("decode", false, "T/F flag for decoding mode")
	rot13   = flag.Bool("rot13", false, "T/F flag to encode the message into Rot13 before writing to the output file")
	base16  = flag.Bool("base16", false, "T/F flag to encode the message into Base16 before writing to the output file")
	base32  = flag.Bool("base32", false, "T/F flag to encode the message into Base32 before writing to the output file")
	base64  = flag.Bool("base64", false, "T/F flag to encode the message into Base64 before writing to the output file")
	base85  = flag.Bool("base85", false, "T/F flag to encode the message into asciiBase85 before writing to the output file")
	complex = flag.String("complex", "", "A comma separated list of encoding types, applied in the order they appear")
)

func init() {
	flag.StringVar(srcFile, "s", "", "Path to the source file to be messed with")
	flag.StringVar(outFile, "o", "", "Path to the output directory when finished")
	flag.StringVar(text, "t", "", "The text string to encode into the file")
	flag.BoolVar(stdIn, "i", false, "T/F flag if passing through concat file and pipe | ")
	flag.StringVar(msgFile, "m", "", "The path to a text file containing the message to be encoded")
	flag.BoolVar(decode, "d", false, "T/F flag for decoding mode")
	flag.BoolVar(rot13, "r13", false, "T/F flag to encode the message into Rot13 before writing to the output file")
	flag.BoolVar(base16, "b16", false, "T/F flag to encode the message into Base16 before writing to the output file")
	flag.BoolVar(base32, "b32", false, "T/F flag to encode the message into Base32 before writing to the output file")
	flag.BoolVar(base64, "b64", false, "T/F flag to encode the message into Base64 before writing to the output file")
	flag.BoolVar(base85, "b85", false, "T/F flag to encode the message into asciiBase85 before writing to the output file")
	flag.StringVar(complex, "c", "", "A comma separated list of encoding types, applied in the order they appear")
}

func parseFlags() (*process.Flags, []error) {
	flag.Parse()
	var errs []error
	srcFilePath, err := filepath.Abs(*srcFile)
	if err != nil {
		errs = append(errs, err)
	}
	outFilePath, err := filepath.Abs(*outFile)
	if err != nil {
		errs = append(errs, err)
	}
	msgFilePath, err := filepath.Abs(*msgFile)
	if err != nil {
		errs = append(errs, err)
	}
	// If attempting to encode without a source message
	if !*decode && (msgFilePath == "") && (*text == "") {
		errs = append(errs, fmt.Errorf("need a text source to encode"))
	}
	return &process.Flags{
		SrcFile:     srcFilePath,
		OutputFile:  outFilePath,
		Text:        *text,
		MessageFile: msgFilePath,
		Decode:      *decode,
		Rot13:       *rot13,
		Base16:      *base16,
		Base32:      *base32,
		Base64:      *base64,
		Base85:      *base85,
		Complex:     *complex,
	}, errs
}

func main() {
	flags, errs := parseFlags()
	if len(errs) != 0 {
		for _, err := range errs {
			log.Println(err)
		}
		os.Exit(1)
	}
	if flags.Decode {
		// Call decode
	}
	encodeConfig, err := flags.ToEncConf()
	if err != nil {
		panic(err)
	}
	// Call encode
}
