# lsb_encoder
A Least Significant Bit(s) Steganography encoder/decoder. Takes a secret message and embeds it into an image file using LSB insertion.

Supported Formats:
* PNG
* JPEG - outputs as output_jpeg.png to avoid compression destroying the secret
* BMP
* GIF - Currently only supports smaller messages (~200 characters) due to GIF being a crappy format with predefined color tables.

### Install Dependencies
Visit the [Golang](https://golang.org/dl/) downloads page and download the installer for your operating 
system.

Follow the install instructions for your operating system
* [Linux](https://golang.org/doc/install#tarball)
* [Windows](https://golang.org/doc/install#windows)
* [MacOS](https://golang.org/doc/install#macos)

Clone or download the ZIP of this repository, and place it in your go workspace.
`C:\Go\src` for Windows (by default)
`$HOME/go/src` for Linux/MacOS (by default)

### Run
If on Windows CMD, use `\` instead of `/` for file paths.

To run, use this command:
`lsb_encoder <options>` or `lsb_encoder.exe <options>`
#### Flag Options
```sh
  Flag Options:
  - s # or --srcfile /path/to/input/source.png (.gif, .bmp, or .jpeg)
  - o # or --outdir /path/to/output/ Directory to save output.png (.gif, .bmp, or .jpeg)
  - t # or --text "The Secret Message to embed"
  - m # or --msgfile /path/to/secret_message.txt (can be anything, just has to fit in srcfile)
  - i # or --stdin The secret message to embed comes from stdin (ex: pipe command)
  - b # or --bitOpt Option for which LSBs to embed: 1 = Last Bit Only, 2 = 2-3-3/R-G-B method (default)
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
  lsb_encoder \
    --srcfile ~/Desktop/Pics/kitty_cat.jpeg \
    --outdir ~/Desktop/Pics -base64 \
    --text "Kitty Cat"

  lsb_encoder --decode --b64 \
    -s ~/Desktop/Pics/output_jpeg.png \
    -o ~/Desktop/Pics \

  # Fancy
  # embed a message in a small image file, like my_avatar.png
  lsb_encoder \
    -s ~/Desktop/Pics/my_avatar.png \
    -o ~/Desktop/Pics/Output \
    --text "Shhhh, don't tell anyone this is hidden in my avatar."
  
  # embed the output from above in a wallpaper
  lsb_encoder \
    -s ~/Desktop/Pics/really_cool_wallpaper.jpeg \
    -o ~/Desktop/Pics/Output \
    --msgfile ~/Desktop/Pics/Output/output.png
```

# Least Significant Bit(s) Steganography
Basically, take a secret message... "Hello!" for example and convert it from ASCII to binary:
```
    H        e        l        l        o        !
01001000 01100101 01101100 01101100 01101111 01000001
```
Break apart the message into an array of bits and hide them inside an image file's pixels using the **Least Significant Bit(s)** insertion. The emphasized text in the table below is the ASCII character split into 2-3-3 and embedded in an R-G-B pixel.

| RGB      | Encoded Character | RGB      | Encoded Character | RGB      | Encoded Character |
|----------|:-----------------:|----------|:-----------------:|----------|:-----------------:|
|101101`01`|                   |010101`01`|                   |110101`01`|                   |
|10101`001`|       **H**       |10101`100`|       **e**       |10110`101`|       **l**       |
|10101`000`|                   |10110`101`|                   |11010`100`|                   |
|||||||
|101001`01`|                   |101101`01`|                   |101011`01`|                   |
|11010`101`|       **l**       |10010`101`|       **o**       |10101`000`|       **!**       |
|10100`100`|                   |10110`111`|                   |10110`001`|                   |
|||||||

