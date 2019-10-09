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
	stdIn   = flag.Bool("stdin", false, "Passing the message through stdin (ex: pipe command)")
	msgFile = flag.String("msgfile", "", "The path to a text file containing the message to be endcoded")
	bitOpt  = flag.Int("bitopt", 1, "Option for which LSBs to embed: 1 = Last Bit Only, 2 = 2-3-3/R-G-B method (default)")
	decode  = flag.Bool("decode", false, "Decoding mode")
	rot13   = flag.Bool("rot13", false, "Apply Rot13 pre encoding to the message before embedding")
	base16  = flag.Bool("base16", false, "Apply Base16 pre encoding to the message before embedding")
	base32  = flag.Bool("base32", false, "Apply Base32 pre encoding to the message before embedding")
	base64  = flag.Bool("base64", false, "Apply Base64 pre encoding to the message before embedding")
	base85  = flag.Bool("base85", false, "Apply Base85 pre encoding to the message before embedding")
	complex = flag.String("complex", "", "A comma separated list(no spaces) of encoding types, applied in the order they appear (limit 5)")
	help    = flag.Bool("help", false, "Print out help text")
)

func init() {
	flag.StringVar(srcFile, "s", "", "Path to the source file to be messed with")
	flag.StringVar(outDir, "o", "", "Path to the output directory when finished")
	flag.StringVar(text, "t", "", "The text string to encode into the file")
	flag.BoolVar(stdIn, "i", false, "Passing the message through stdin (ex: pipe command)")
	flag.StringVar(msgFile, "m", "", "The path to a text file containing the message to be encoded")
	flag.IntVar(bitOpt, "b", 2, "Option for which LSBs to embed: 1 = Last Bit Only, 2 = 2-3-3/R-G-B method (default)")
	flag.BoolVar(decode, "d", false, "Decoding mode")
	flag.BoolVar(rot13, "r13", false, "Apply Rot13 pre encoding to the message before embedding")
	flag.BoolVar(base16, "b16", false, "Apply Base16 pre encoding to the message before embedding")
	flag.BoolVar(base32, "b32", false, "Apply Base32 pre encoding to the message before embedding")
	flag.BoolVar(base64, "b64", false, "Apply Base64 pre encoding to the message before embedding")
	flag.BoolVar(base85, "b85", false, "Apply Base85 pre encoding to the message before embedding")
	flag.StringVar(complex, "c", "", "A comma separated list(no spaces) of encoding types, applied in the order they appear (limit 5)")
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
		- b # or --bitopt Option for which LSBs to embed: 1 = Last Bit Only, 2 = 2-3-3/R-G-B method (default)
		- d # or --decode Extract a message from an already embedded file
		- r13 # or --rot13 Apply Rot13 pre encoding to the message before embedding
		- b16 # or --base16 Apply Base16 pre encoding to the message before embedding
		- b32 # or --base32 Apply Base32 pre encoding to the message before embedding
		- b64 # or --base64 Apply Base64 pre encoding to the message before embedding
		- b85 # or --base85 Apply Base85 pre encoding to the message before embedding
		- c # or --complex A Comma separated list(no space) of encoding types, applied in the order they appear (limit 5)
		- h # or --help Print out help text
	
		# Example commands
		# Simple
		go run ./cmd/lsb_encoder/ \
			--srcfile ~/Desktop/Pics/kitty_cat.jpeg \
			--outdir ~/Desktop/Pics -base64 \
			--text "Kitty Cat"
	
		go run ./cmd/lsb_encoder/ --decode --b64 \
			-s ~/Desktop/Pics/output_jpeg.png \
			-o ~/Desktop/Pics \
	
		# Fancy
		go run ./cmd/lsb_encoder/ \
			-s ~/Desktop/Pics/funny_cat.gif \
			-o ~/Desktop/Pics \
			--complex "b16,b32,b64,b85" \
			--msgfile ~/Downloads/lorem_ipsum_paragraph.txt
	
		go run ./cmd/lsb_encoder/ --decode \
			-s ~/Desktop/Pics/output.gif \
			-o ~/Desktop/Pics \
			--complex "b85,b64,b32,b16"

		# Even Fancier
  	# embed a message in a small image file, like my_avatar.png
  	go run ./cmd/lsb_encoder/ \
    	-s ~/Desktop/Pics/my_avatar.png \
    	-o ~/Desktop/Pics/Output \
    	--text "Shhhh, don't tell anyone this is hidden in my avatar."
  
  	# embed the output from above in a wallpaper
  	go run ./cmd/lsb_encoder/ \
    	-s ~/Desktop/Pics/really_cool_wallpaper.jpeg \
    	-o ~/Desktop/Pics/Output \
    	--msgfile ~/Desktop/Pics/Output/output.png
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
		BitOpt:      *bitOpt,
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
