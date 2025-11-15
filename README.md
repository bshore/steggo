# steggo

A Least Significant Bits Steganography (**LSB**) embedder/extractor. Takes a message and embeds it into an image file using LSB insertion.

## Supported Input Formats

- PNG
- JPEG - outputs as `<input_name>_jpeg_output.png` (Outputs as PNG because JPEG is [lossy](https://youtu.be/jmaUIyvy8E8?si=uj2WBSBmbSfRlAT3) which destroys the message)
- BMP - outputs as `<input_name>_bmp_output.png` (Outputs as PNG because BMP is hard-capped at 256 colors))
- GIF

## Run Embed

```bash
steggo embed --help

Embeds --input {message} into --target {file} placing output in --dest {path}

Usage:
  steggo embed [flags]

Flags:
  -d, --dest string            The destination path to output the target file after embedding (default ".")
  -h, --help                   help for embed
  -i, --input string           The input path or message to embed into the target file
  -p, --pre-encoding strings   (Optional) A comma separated list of pre-encoders to apply before embedding, 5 max: r13, b16, b32, b64, b85, gzip.
                               Each encoder is applied in the order they are specified.

                               NOTE: The gzip option compresses the message and may not be used with other encoders.

  -t, --target string          The path to the image file being targeted for embedding
```

## Run Extract

```bash
steggo extract --help

Extracts the message from --target {file} and outputs it to --dest {path}

Usage:
  steggo extract [flags]

Flags:
  -d, --dest string            The destination path to output the extracted message (default ".")
  -h, --help                   help for extract
  -t, --target string          The path to the image file being targeted for extraction
```

## What is it? How?

Take the example string input "Hello!" and convert it from ASCII to an array of it's binary representation.

```js
[
  01001000, // H
  01100101, // e
  01100100, // l
  01100100, // l
  01100101, // o
  01000000, // !
]
```

Break apart each binary representation further into groups of 3 that will fit into one R-G-B color.

```js
[
  _01,_001,_000, // H
  _01,_100,_101, // e
  _01,_100,_100, // l
  _01,_100,_100, // l
  _01,_100,_101, // o
  _01,_000,_000, // !
]
```

Insert the values of each group of 3 into the **Least Significant Bits** of an RGB color.

The emphasized text in the table below shows the 2-3-3 insertion into existing RGB color values.

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

## What does it look like?

<table>
  <thead>
    <tr>
      <td>Before</td>
      <td>After</td>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><img src="./img/bsotp.jpg" width="500" /></td>
      <td><img src="./img/bsotp_jpeg_output.png" width="500" /></td>
    </tr>
  </tbody>
</table>
