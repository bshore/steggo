# lsb_encoder
A Least Significant Bit(s) Stenography encoder

# Why?
Maybe you want to send secret messages to a friend through images?
Maybe you just want to look at how much time someone can waste on something nobody really needs?
Or maybe you're building a Cyber Capture The Flag event for a hacker convention? (like me)

# Least Significant Bit(s) Stenography
Basically, take a secret message... "Hello Secret World!" for example:
```
01001000 01100101 01101100 01101100 01101111 // Hello
00100000 // Space
01010011 01100101 01100011 01110010 01100101 01110100 // Secret
00100000 // Space
01010111 01101111 01110010 01101100 01100100 00100001 // World!
```
>Break apart the message into n bits (usually 1 or 2) and hide them inside an image/video/audio file using the **Least Significant Bit(s)**

| LSB      | Encoded Character | LSB      | Encoded Character |
|----------|:-----------------:|----------|:-----------------:|
|101101`01`|                   |010101`01`|                   | 
|101011`00`|         H         |101010`10`|         e         |
|101011`10`|                   |101101`01`|                   |
|110101`00`|                   |010101`01`|                   |
|||||
|110101`01`|                   |101001`01`|                   | 
|101101`10`|         l         |110101`10`|         l         |
|110101`11`|                   |101001`11`|                   |
|100101`00`|                   |100101`00`|                   |
|||||
|101101`01`|                   |101011`01`|                   | 
|100101`10`|         o         |101010`00`|      (space)      |
|101101`11`|                   |101101`00`|                   |
|101010`11`|                   |010101`00`|                   |
|||||