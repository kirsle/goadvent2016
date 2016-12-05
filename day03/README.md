# Day 3

This program takes as input three sides of a triangle and it identifies
triangles that are impossible.

In a valid triangle the sum of any two sides must be greater than the remaining
side. In the example `5 10 25`, the "triangle" is invalid because `5 + 10` is
not greater than `25`

The format of the input file is that the triangle coordinates are listed
vertically, not horizontally; so in this example each number with the same
hundreds digit is part of the same triangle:

```
101 301 501
102 302 502
103 303 503
201 401 601
202 402 602
203 403 604
```
