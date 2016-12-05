# Day 5

This program produces an eight-character password, found one character at a
time, by finding the MD5 hash of some Door ID (the puzzle input) and an
increasing integer index starting with 0.

A hash indicates the **next character** in the password if its hexadecimal
representation starts with **five zeroes**. If the hash begins with five zeroes,
then the **sixth** character represents the position (0-7) and the **seventh**
character is the symbol to put in that position.

A hash result of `000001f` means that `f` is the second character in the
password. Use only the **first result** for each position, and ignore invalid
positions.

For example, if the Door ID is `abc`:

- The first interesting hash is from `abc3231929` which produces
  `0000015...`; so `5` goes in position `1`: `_5______`
- The second interesting hash is at index `5357525`, which produces
  `000004e...`; so `e` goes in position `4`: `_5__e___`
- The final password ends up being `05ace8e3`

## Run It

```bash
go run main.go reyedfim
```
