# Differences

Comparatively, this library offers combined functionality not found within most if not all other COBS libraries on the market as of writing this. To name a few:

 - **Choose whether the Delimiter (Special byte) is added.** Libraries tend to choose this for you.
 - **Choose what the Special Byte is.** Libraries tend to provide the NULL (0x00) byte only.
 - **Use COBS as a Layer of Integrity.** By ensuring that the special byte does not occur (expect with a delimiter) and by ensuring that the flags lead to the end of the data, COBS can provide a small layer of integrity.
 - **Use a TONS of different extensions (types).** Found in [usage](https://github.com/justincpresley/go-cobs/blob/master/docs/USAGE.md), rare types are included and even ones only included in this library.
 - **Additional API Commands.** Not only is there functions to calculate the min/max overhead for COBS but also flag-related functions.

## Research

There is a lot of research potential out of this library. For example, It would be an interesting to find what the optimal amount of flags is to provide the best integrity checking for some given data.

The license on **go-cobs** was chosen to help facilitate research. If you do conduct research, you are welcomed to open an [Issue](https://github.com/justincpresley/go-cobs/issues) to document your findings/papers/etc. This way, **go-cobs** can take in your research (given the approach attributions) and grow as a library.
