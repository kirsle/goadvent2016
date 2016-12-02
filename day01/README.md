# Day 1

This program follows a series of steps to navigate through a city, where each
step includes a direction to turn (Left or Right) and the number of city blocks
to travel in the direction you are now facing.

After traversing the steps this program will tell you the shortest distance to
get to the final destination.

You start at an arbitrary place (coordinate `(0,0)`) and are facing North.
Here are some example inputs and outputs:

* Following `R2, L3` leaves you `2` blocks East and `3` blocks North, or `5`
  blocks away from where you started.
* `R2, R2, R2` leaves you `2` blocks due South of your starting position, which
  is `2` blocks away.
* `R5, L5, R5, R3` leaves you `12` blocks away.

The instructions continue on the back and say that the destination you're
looking for is actually at the first location that you visit twice.

For example if your instructions are `R8, R4, R4, R8`, the first location you
visit twice is `4` blocks away, due East.
