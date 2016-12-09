# Day 9

This program implements an experimental decompression algorithm.

The format composes a sequence of characters (whitespace is ignored). To
indicate that some sequence should be repeated, a marker is added to the file,
like `(10x2)`. To decompress the market, take the subsequent `10` characters
and repeat them `2` times. Then, continue reading the file **after** the
repeated data. The marker itself is not included in the decompressed output.

If parenthesis or other characters appear within the data referenced by a
marker, that's okay - treat it like normal data, not a marker, and then
resume looking for markers in the decompressed section.

Examples:

- `ADVENT` contains no markers and decompresses to itself (length `6`)
- `A(1x5)BC` repeats only the `B` a total of `5` times, becoming
  `ABBBBBC` (length `7`)
- `(3x3)XYZ` becomes `XYZXYZXYZ` (length `9`)
- `A(2x2)BCD(2x2)EFG` becomes `ABCBCDEFEFG`
- `(6x1)(1x3)A` simply becomes `(1x3)A` -- the `(1x3A)` looks like a marker,
  but it's within the `6` character data of another marker so it's not treated
  any differently than other characters. The decompressed length is `6`
- `X(8x2)(3x3)ABCY` becomes `X(3x3)ABC(3x3)ABCY` (length `18`), because the
  decompressed data from the `(8x2)` marker (the `(3x3)ABC`) is skipped and not
  processed further.

It turns out the input file actually uses **version two** of the compression
format: the only difference is that markers within compressed data **are**
decompressed.

For example:

- `(3x3)XYZ` still becomes `XYZXYZXYZ`
- `X(8x2)(3x3)ABCY` becomes `XABCABCABCABCABCABCY` because the decompressed data
  from the `(8x2)` marker is then further decompressed, thus triggering the
  `(3x3)` marker twice for a total of six `ABC` sequences.
- `(27x12)(20x12)(13x14)(7x10)(1x12)A` decompresses into a string of `A`
  repeated `241920` times.
- `(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN` becomes `445`
  characters long.
