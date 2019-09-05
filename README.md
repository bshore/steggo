# lsb_encoder
A Least Significant Bit(s) Steganography encoder/decoder. Takes a secret message and embeds it into an image file using LSB insertion. 

**Note:** The secret to embed needs to be approx < 20% the file size of the source image to safely embed the whole secret. I haven't gotten that totally pinned down yet.

### Install Dependencies
Visit the [Golang](https://golang.org/dl/) downloads page and download the installer for your operating 
system.

Follow the install instructions for your operating system
[Linux](https://golang.org/doc/install#tarball)
[Windows](https://golang.org/doc/install#windows)
[MacOS](https://golang.org/doc/install#macos)

Clone or download the ZIP of this repository, and place it in your go workspace.
`C:\Go\src` for Windows (by default)
`$HOME/go/src` for Linux/MacOS (by default)

### Run
If on Windows, use `\` instead of `/` for file paths.

To run use this command:
`go run ./cmd/lsb_encoder/` or `go run .\cmd\lsb_encoder\`
#### Flag Options
```sh
  - s # or --srcfile /path/to/input/source.png (.gif or .jpeg)
  - o # or --outdir /path/to/output/ Directory to save output.png (.gif or .jpeg)
  - t # or --text "The Secret Message to embed"
  - m # or --msgfile /path/to/secret_message.txt (can be anything)
  - i # or --stdin The secret message to embed comes from stdin (ex: pipe command)
  - d # or --decode Extract a message from an already embedded file
  - r13 # or --rot13 (Apply Rot13 pre encoding to the message before embedding)
  - b16 # or --base16 (Apply Base16 pre encoding to the message before embedding)
  - b32 # or --base32 (Apply Base32 pre encoding to the message before embedding)
  - b64 # or --base64 (Apply Base64 pre encoding to the message before embedding)
  - b85 # or --base85 (Apply Base85 pre encoding to the message before embedding)
  - c # or --complex (In the order they appear, apply pre encoding to the message before embedding)
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
```

# Why?
Maybe you want to send secret messages to a friend through images?

Maybe you just want to look at how much time someone can waste on something nobody really needs?

Or maybe you're building a Cyber Capture The Flag event for a hacker convention? (like me)

# Least Significant Bit(s) Steganography
Basically, take a secret message... "Hello!" for example and convert it from ASCII to binary:
```
    H        e        l        l        o        !
01001000 01100101 01101100 01101100 01101111 01000001
```
Break apart the message into an array of bits and hide them inside an image file's pixels using the **Least Significant Bit(s)** insertion.

| LSB      | Encoded Character | LSB      | Encoded Character | LSB      | Encoded Character |
|----------|:-----------------:|----------|:-----------------:|----------|:-----------------:|
|101101`01`|                   |010101`01`|                   |110101`01`|                   |
|101011`00`|       **H**       |101010`10`|       **e**       |101101`10`|       **l**       |
|101011`10`|                   |101101`01`|                   |110101`11`|                   |
|110101`00`|                   |010101`01`|                   |100101`00`|                   |
|||||||
|101001`01`|                   |101101`01`|                   |101011`01`|                   |
|110101`10`|       **l**       |100101`10`|       **o**       |101010`00`|       **!**       |
|101001`11`|                   |101101`11`|                   |101101`00`|                   |
|100101`00`|                   |101010`11`|                   |010101`01`|                   |
|||||||

Just do that until the message to encode runs out.
