# Day 10

Example instructions:

```
value 5 goes to bot 2
bot 2 gives low to bot 1 and high to bot 0
value 3 goes to bot 1
bot 1 gives low to output 1 and high to bot 0
bot 0 gives low to output 2 and high to output 0
value 2 goes to bot 2
```

- Initially, bot 1 starts with a value-3 chip, and bot 2 starts with a value-2 chip and a value-5 chip.
- Because bot 2 has two microchips, it gives its lower one (2) to bot 1 and its higher one (5) to bot 0.
- Then, bot 1 has two microchips; it puts the value-2 chip in output 1 and gives the value-3 chip to bot 0.
- Finally, bot 0 has two microchips; it puts the 3 in output 2 and the 5 in output 0.

## Notes

My approach was to load all of the steps into a list, and repeatedly cycle
through the list on an event loop and test each step.

If a step is able to be fully carried out, it's marked as Done and not touched
again on future loops. For example, the "input bucket gives a bot a microchip"
step can fail if the receiving bot has no space left in their inventory to
receive the microchip.
