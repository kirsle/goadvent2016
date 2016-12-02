# Day 2

This program figures out a passcode using spatial instructions. Imagine a
pass code entry form that looks like this:

```
      1
    2 3 4
  5 6 7 8 9
    A B C
      D
```

You start on "5" and follow inputs like these, where each line points to a
number relative to where you previously were:

```
  ULL
  RRDDD
  LURDL
  UUUUD
```

- You start at "5" and don't move (up and left are edges)
- Continuing from "5" you move right twice and down three times and land
  on "D"
- Then from "D" you move five more times and end at "B"
- Finally after five more moves, you end at "3"

So in this example the code is: `5DB3`
