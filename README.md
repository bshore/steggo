# steggo

A Least Significant Bit(s) Steganography (**LSB**) embedder/extracter. Takes a secret message and embeds it into an image file using LSB insertion.

Supported Input Formats:

- PNG
- JPEG (outputs as `<input_name>_jpeg_output.png`)
  - The `jpeg` format has it's own built-in compression algorithm that will alter image pixels _on save_ and the embedded message's integrity can't be guaranteed.
- BMP (outputs as `<input_name>_bmp_output.png`)
  - The `bmp` format only supports 256 colors. tl;dr - The LSB process modifies enough of these pixel colors that having a hard cap of 256 means there is no guarantee that the message can be extracted back out.

### Run

Example Commands:

```bash
steggo embed --target path/picture.png --dest path/outputs/ --input "Words go here"

steggo extract --target path/outputs/picture_output.png

steggo embed --target path/another.png --dest path/outputs/ --input somefile.txt --pre-encoding r13,b64 # apply rot13 & base64 encoding

steggo extract --target path/outputs/another_output.png --dest path/extracted/
cat path/extracted/message.txt
```

# Least Significant Bit(s) Steganography

Basically, take a secret message... "Hello!" for example and convert it from ASCII to binary:

```
    H        e        l        l        o        !
01001000 01100101 01101100 01101100 01101111 01000001
```

Break apart the message into an array of bits and hide them inside an image file's pixels using the **Least Significant Bit(s)** insertion. The emphasized text in the table below is the ASCII character split into 2-3-3 and embedded in an R-G-B pixel.

| RGB        | Encoded Character | RGB        | Encoded Character | RGB        | Encoded Character |
| ---------- | :---------------: | ---------- | :---------------: | ---------- | :---------------: |
| 101101`01` |                   | 010101`01` |                   | 110101`01` |                   |
| 10101`001` |       **H**       | 10101`100` |       **e**       | 10110`101` |       **l**       |
| 10101`000` |                   | 10110`101` |                   | 11010`100` |                   |
|            |                   |            |                   |            |                   |
| 101001`01` |                   | 101101`01` |                   | 101011`01` |                   |
| 11010`101` |       **l**       | 10010`101` |       **o**       | 10101`000` |       **!**       |
| 10100`100` |                   | 10110`111` |                   | 10110`001` |                   |
|            |                   |            |                   |            |                   |

## TODOs

- Need to try a different approach to `.gif` embedding, seeing as it's a larger file type there's more opportunity for embedding larger messages (entire books, perhaps?) inside them.
- Frontend & REST server so there's a more user-friendly way of using `steggo`
