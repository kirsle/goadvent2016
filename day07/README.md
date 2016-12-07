# Day 7

This takes an input list of "IPv7 Addresses" and counts the number of addresses
that support "TLS" (Transport-Layer Snooping)

An IP supports TLS if it has an "ABBA" sequence of characters anywhere in its
address, for example a sequence like `xyyx` or `abba`. However, if it contains
an ABBA inside the square bracket parts of the address, then it **DOES NOT**
support TLS.

For examples:

- `abba[mnop]qrst` supports TLS (`abba` outside square brackets)
- `abcd[bddb]xyyx` does **not** (`bddb` is inside square brackets, even though
  `xyyx` is outside brackets)
- `aaaa[qwer]tyui` does **not** (`aaaa` is invalid; the interior characters
  must be different)
- `ioxxoj[asdfgh]zxcvbn` supports TLS (`oxxo` is outside square brackets,
  even though it's within a larger string)

It also counts the IP addresses that support "SSL" (Super-Secret Listening).
An IP supports SSL if it has an "ABA" sequence anywhere in the non-bracketed
parts, and additionally a "BAB" inside the bracketed parts.

An ABA is any three-character sequence that consists of the same character
twice with a different character between them, such as `xyx` or `aba`. A
corresponding BAB is the same characters but in reversed positions: `yxy` and
`bab`, respectively.

Examples:

- `aba[bab]xyz` supports SSL (`aba` outside sqaure brackets and `bab` inside)
- `xyx[xyx]xyx` does NOT support SSL (`xyx` but no corresponding `yxy`)
- `aaa[kek]eke` supports SSL (`eke` outside the brackets and `kek` inside)
- `zazbz[bzb]cdb` supports SSL (`zaz` has no corresponding `aza`, but `zbz`
  has a corresponding `bzb`, even though `zaz` and `zbz` overlap)

## Implementation Details

Something I particularly like in my code is my `SlidingWindow` function, which
handles the logic of moving a sliding window of arbitrary size across an input
string, calls a comparator function for each slice to determine whether it's
what the caller wanted, and lets the caller abort the loop and return the last
result checked.

For example, to find the first ABBA sequence and stop:

```go
abba := SlidingWindow(input, 4, func(slice string) bool {
    return slice[0] == slice[3] && slice[1] == slice[2] && slice[0] != slice[1]
})
```

Or to populate a map of *all* valid ABA/BAB sequences (the loop doesn't
terminate early like in the above example):

```go
SlidingWindow(input, 3, func(slice string) bool {
    if slice[0] == slice[2] && slice[0] != slice[1] {
        result[slice] = true
    }
    return false
})
```
