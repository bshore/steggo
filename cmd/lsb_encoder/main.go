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
	outDir  = flag.String("outdir", "", "Path to the output directory when finished")
	text    = flag.String("text", "", "The text string to encode into the file")
	stdIn   = flag.Bool("stdin", false, "T/F flag if passing the message through stdin (ex: pipe command)")
	msgFile = flag.String("msgfile", "", "The path to a text file containing the message to be endcoded")
	decode  = flag.Bool("decode", false, "T/F flag for decoding mode")
	rot13   = flag.Bool("rot13", false, "T/F flag to encode the message into Rot13 before writing to the output file")
	base16  = flag.Bool("base16", false, "T/F flag to encode the message into Base16 before writing to the output file")
	base32  = flag.Bool("base32", false, "T/F flag to encode the message into Base32 before writing to the output file")
	base64  = flag.Bool("base64", false, "T/F flag to encode the message into Base64 before writing to the output file")
	base85  = flag.Bool("base85", false, "T/F flag to encode the message into asciiBase85 before writing to the output file")
	complex = flag.String("complex", "", "A comma separated list(no spaces) of encoding types, applied in the order they appear")
	help    = flag.Bool("help", false, "Print out help text")
)

func init() {
	flag.StringVar(srcFile, "s", "", "Path to the source file to be messed with")
	flag.StringVar(outDir, "o", "", "Path to the output directory when finished")
	flag.StringVar(text, "t", "", "The text string to encode into the file")
	flag.BoolVar(stdIn, "i", false, "T/F flag if passing the message through stdin (ex: pipe command)")
	flag.StringVar(msgFile, "m", "", "The path to a text file containing the message to be encoded")
	flag.BoolVar(decode, "d", false, "T/F flag for decoding mode")
	flag.BoolVar(rot13, "r13", false, "T/F flag to encode the message into Rot13 before writing to the output file")
	flag.BoolVar(base16, "b16", false, "T/F flag to encode the message into Base16 before writing to the output file")
	flag.BoolVar(base32, "b32", false, "T/F flag to encode the message into Base32 before writing to the output file")
	flag.BoolVar(base64, "b64", false, "T/F flag to encode the message into Base64 before writing to the output file")
	flag.BoolVar(base85, "b85", false, "T/F flag to encode the message into asciiBase85 before writing to the output file")
	flag.StringVar(complex, "c", "", "A comma separated list(no spaces) of encoding types, applied in the order they appear")
	flag.BoolVar(help, "h", false, "Print out help text")
}

func parseFlags() (*process.Flags, []error) {
	flag.Parse()
	if *help {
		fmt.Println(`
		Flag Options:
		- s # or --srcfile /path/to/input/source.png (.gif, .bmp, or .jpeg)
		- o # or --outdir /path/to/output/ Directory to save output.png (.gif, .bmp, or .jpeg)
		- t # or --text "The Secret Message to embed"
		- m # or --msgfile /path/to/secret_message.txt (can be anything)
		- i # or --stdin The secret message to embed comes from stdin (ex: pipe command)
		- d # or --decode Extract a message from an already embedded file
		- r13 # or --rot13 (Apply Rot13 pre encoding to the message before embedding)
		- b16 # or --base16 (Apply Base16 pre encoding to the message before embedding)
		- b32 # or --base32 (Apply Base32 pre encoding to the message before embedding)
		- b64 # or --base64 (Apply Base64 pre encoding to the message before embedding)
		- b85 # or --base85 (Apply Base85 pre encoding to the message before embedding)
		- c # or --complex (A comma separated list(no spaces) of encoding types, applied in the order they appear)
		- h # or --help (Print out this text block)
	
		# Example commands
		# Simple
		go run ./cmd/lsb_encoder/ \
			--srcfile ~/Desktop/Pics/kitty_cat.jpeg \
			--outdir ~/Desktop/Pics -base64 \
			--text "Kitty Cat"
	
		go run ./cmd/lsb_encoder/ --decode --b64 \
			-s ~/Desktop/Pics/output.jpeg \
			-o ~/Desktop/Pics \
	
		# Fancy
		go run ./cmd/lsb_encoder/ \
			-s ~/Desktop/Pics/funny_cat.gif \
			-o ~/Desktop/Pics \
			--complex "b16,b32,b64,b85" \
			--msgfile ~/Downloads/harry_potter_prisoner_of_azkaban.txt
	
		go run ./cmd/lsb_encoder/ --decode \
			-s ~/Desktop/Pics/output.gif \
			-o ~/Desktop/Pics \
			--complex "b85,b64,b32,b16"
		`)
		os.Exit(0)
	}
	var msgFilePath string
	var errs []error
	srcFilePath, err := filepath.Abs(*srcFile)
	if err != nil {
		errs = append(errs, err)
	}
	outFilePath, err := filepath.Abs(*outDir)
	if err != nil {
		errs = append(errs, err)
	}
	if *msgFile != "" {
		msgFilePath, err = filepath.Abs(*msgFile)
		if err != nil {
			errs = append(errs, err)
		}
	}
	// If attempting to encode without a source message
	if !*decode && (msgFilePath == "") && (*text == "") {
		errs = append(errs, fmt.Errorf("need a text source to encode"))
	}
	return &process.Flags{
		SrcFile:     srcFilePath,
		OutputDir:   outFilePath,
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
		extractSecret, err := process.ParseExtractSecret(flags)
		if err != nil {
			panic(err)
		}
		err = process.Extract(extractSecret)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}
	embedSecret, err := process.ParseEmbedSecret(flags)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = process.Embed(embedSecret)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
