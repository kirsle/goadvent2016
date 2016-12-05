# Day 4

This program takes an input list of "room names", where each room consists of an
encrypted name (lowercase letters and hyphens) followed by a dash, a sector ID,
and a checksum in square brackets.

The room is valid if the checksum consists of the five most common letters in
the encrypted name, in order of frequency. When there are multiple letters with
the same frequency, they are listed in alphabetical order.

Examples:

- `aaaaa-bbb-z-y-x-123[abxyz]` is a real room because the most common letters
  are `a` (5), `b` (3), and then a tie between `x`, `y`, and `z`, which are
  listed alphabetically.
- `a-b-c-d-e-f-g-h-987[abcde]` is a real room because although the letters are
  all tied (1 of each), the first five are listed alphabetically.
- `not-a-real-room-404[oarel]` is a real room.
- `totally-real-room-200[decoy]` is not.

For the valid room names, the encryption on the name is a simple shift cipher,
where each letter is shifted the number of positions of its sector number.
Dashes become spaces. For example:

- `qzmt-zixmtkozy-ivhz-343` is `very encrypted name`
