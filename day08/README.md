# Day 8

This program takes as input a list of instructions for how to manipulate
"pixels" on a tiny LCD display.

The screen is `50` pixels wide and `6` pixels tall.

The program inputs contain instructions in the following formats:

- `rect AxB` turns **on** all the pixels in a rectangle at the top-left corner
  of the screen which is `A` pixels wide and `B` pixels tall.
- `rotate row y=A by B` shifts all of the pixels in row `A` (0 is the top row)
  **right** by `B` pixels. Pixels that would fall off the right end appear at
  the left end of the row.
- `rotate column x=A by B` shifts all of the pixels in column `A` (0 is the
  left column) **down** by `B` pixels. Pixels that would fall off the bottom
  appear at the top of the column.

For example, here is a simple sequence on a smaller screen:

- `rect 3x2` creates a small rectangle in the top-left corner:

    ```
    ###....
    ###....
    .......
    ```

- `rotate column x=1 by 1` rotates the second column down by one pixel:

    ```
    #.#....
    ###....
    .#.....
    ```

- `rotate row y=0 by 4` rotates the top row right by four pixels:

    ```
    ....#.#
    ###....
    .#.....
    ```

- `rotate column x=1 by 1` rotates the second column down by one, causing the
  bottom pixel to wrap around to the top.

    ```
    .#..#.#
    #.#....
    .#.....
    ```

At the end, this program prints the number of lit pixels in the screen.

## Program Output

This will contain spoilers, but you can see what the output of my program looks
like by checking [Output.md](./Output.md)
