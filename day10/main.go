package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kirsle/goadvent2016/advent"
)

func main() {
	// Get the input instructions.
	input, err := advent.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	// Parse the instruction set so we know who all our bots and inputs are.
	steps := ParseInstructions(input)
	var (
		outputs = NewOutputs()
		bots    = NewBots()
	)

	// Cycle through the steps until we can do no more of them.
	var tick = 0
	for {
		time.Sleep(1000000000)
		advent.Debug("### TICK %d ###\n", tick)
		tick++

		done := DoOneLoop(steps, bots, outputs)
		if done {
			break
		}
	}

	// Pretty print things by making these booleans true.
	if true {
		fmt.Println("Summary of the Bots")
		fmt.Println("===================")
		for _, bot := range bots.bots {
			fmt.Printf("## Bot %s\n", bot.ID)
			fmt.Printf("   Inventory: %v\n", bot.Inventory)
			fmt.Printf("   Comparison History:\n")
			for _, h := range bot.History {
				fmt.Printf("      %d <> %d\n", h.A, h.B)
			}
		}
	}
	if true {
		fmt.Println("\nSummary of the Outputs")
		fmt.Println("======================")
		for _, out := range outputs.outputs {
			fmt.Printf("Output: %s\n", out.ID)
			fmt.Printf("Inventory: %v\n", out.Inventory)
		}
	}

	// Find the bot that had to compare 17 and 61.
	for _, bot := range bots.bots {
		for _, h := range bot.History {
			if h.A == 17 && h.B == 61 {
				fmt.Printf("Found bot %s (comped 17 <> 61)\n", bot.ID)
			}
		}
	}

	// Multiplying the values of outputs 0, 1 and 2.
	outA := outputs.Find("0")
	outB := outputs.Find("1")
	outC := outputs.Find("2")
	fmt.Printf("Final output: %d\n",
		outA.Inventory[0]*
			outB.Inventory[0]*
			outC.Inventory[0],
	)
}

// DoOneLoop processes a loop of A.I. instructions until it runs out of things
// to do.
func DoOneLoop(steps *Steps, bots *Bots, outputs *Outputs) bool {
	// Number of steps executed.
	var executed int

	// Loop through the steps and see if any can be done yet.
	for _, step := range steps.steps {
		// Skip spent steps.
		if step.Done {
			continue
		}

		// What sort of step?
		if step.Action == InputAction {
			// The input gives a token to a bot.
			bot := bots.Find(step.BotID)
			if bot.Give(step.Value) {
				advent.Debug("[ OK ] Input gave %d to Bot %s\n", step.Value, step.BotID)
				step.Done = true
				executed++
			} else {
				advent.Debug("[FAIL] Input can't give %d to Bot %s: no room\n", step.Value, step.BotID)
			}
		} else if step.Action == GiveAction {
			// Give away a token if we have 2.
			bot := bots.Find(step.BotID)
			if !bot.Full() {
				// Bots only give away chips when their hands are full.
				continue
			}

			// Find out the lower and higher value chips. The inventory is
			// already sorted.
			lower, higher := bot.Inventory[0], bot.Inventory[1]
			bot.History = append(bot.History, History{
				A: lower,
				B: higher,
			})

			// Who is the low one for?
			var to string
			var ok bool
			if step.LowTo == ToBot {
				recip := bots.Find(step.LowID)
				to = "Bot"
				ok = bot.Transfer(recip, lower)
			} else if step.LowTo == ToOutput {
				recip := outputs.Find(step.LowID)
				to = "Output"
				ok = bot.Deposit(recip, lower)
			}

			if ok {
				advent.Debug("[ OK ] Bot %s gave L%d to %s %s\n", bot.ID, lower, to, step.LowID)
				executed++
			}

			// Who is the high one for?
			if step.HighTo == ToBot {
				recip := bots.Find(step.HighID)
				to = "Bot"
				ok = bot.Transfer(recip, higher)
			} else if step.HighTo == ToOutput {
				recip := bots.Find(step.HighID)
				to = "Output"
				ok = bot.Transfer(recip, higher)
			}

			if ok {
				advent.Debug("[ OK ] Bot %s gave H%d to %s %s\n", bot.ID, lower, to, step.LowID)
				executed++
			}
		}
	}

	advent.Debug("Executed %d steps\n", executed)
	return executed == 0
}

// ParseInstructions turns the input lines into Step objects.
func ParseInstructions(input []string) *Steps {
	result := &Steps{
		steps: []*Step{},
	}

	for _, line := range input {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Test the regexps.
		var match []string

		// The Input giving a microchip to a bot.
		match = RE_Input.FindStringSubmatch(line)
		if len(match) > 0 {
			value, _ := strconv.Atoi(match[1])
			result.steps = append(result.steps, &Step{
				Action: InputAction,
				BotID:  match[2],
				Value:  Microchip(value),
			})
			continue
		}

		// The give conditions.
		match = RE_Give.FindStringSubmatch(line)
		if len(match) > 0 {
			// Parse the integers out.
			result.steps = append(result.steps, &Step{
				Action: GiveAction,
				BotID:  match[1],
				LowTo:  WhoTo(match[2]),
				LowID:  match[3],
				HighTo: WhoTo(match[4]),
				HighID: match[5],
			})
		}
	}

	return result
}

// WhoTo identifies who the recipient of a microchip is.
func WhoTo(name string) bool {
	return name == "bot"
}

// Regular expression to match the A.I. steps from the input file.
var (
	RE_Input = regexp.MustCompile(`^value (\d+) goes to bot (\d+)$`)
	RE_Give  = regexp.MustCompile(`^bot (\d+) gives low to (bot|output) (\d+) and high to (bot|output) (\d+)$`)
)

type ID int

type Steps struct {
	steps []*Step
}

// Type Step represents a step of A.I. from the input file to try.
type Step struct {
	Action int  // What sort of action this step is
	Done   bool // When the step was able to be finished

	// Parsed regexp values.
	BotID  string    // The actor
	Value  Microchip // Value to be given to the actor
	LowTo  bool      // What sort of thing the low is given to (bot or output)
	LowID  string    // Recipient of the value given by the actor
	HighTo bool
	HighID string
}

// Recipients.
const (
	// Valid types for a step.
	GiveAction = iota
	InputAction

	// Values for the LowTo and HighTo values of a Step.
	ToOutput = false
	ToBot    = true
)

// Type Bots is a collection of bots.
type Bots struct {
	bots map[string]*Bot
}

func NewBots() *Bots {
	return &Bots{
		bots: map[string]*Bot{},
	}
}

type Outputs struct {
	outputs map[string]*Output
}

func NewOutputs() *Outputs {
	return &Outputs{
		outputs: map[string]*Output{},
	}
}

// Find finds a bot or creates it if it doesn't exist.
func (s *Bots) Find(id string) *Bot {
	if bot, ok := s.bots[id]; ok {
		return bot
	}
	s.bots[id] = &Bot{
		ID:        id,
		Inventory: []Microchip{},
		History:   []History{},
	}
	return s.bots[id]
}

// Find finds an output or creates it if it doesn't exist.
func (s *Outputs) Find(id string) *Output {
	if output, ok := s.outputs[id]; ok {
		return output
	}
	s.outputs[id] = &Output{
		ID:        id,
		Inventory: []Microchip{},
	}
	return s.outputs[id]
}

// Type Bot represents a robot.
type Bot struct {
	ID        string
	Inventory []Microchip // What the bot is carrying in its hands (max 2)
	History   []History   // Things this bot has had to do.
}

// Full checks if a bot's inventory is full.
func (b *Bot) Full() bool {
	return len(b.Inventory) >= 2
}

// Give sees if a bots inventory has room to take a microchip and gives it to
// the bot if so. Also ensures that the chips are sorted by lower to higher.
func (b *Bot) Give(chip Microchip) bool {
	if len(b.Inventory) < 2 {
		b.Inventory = append(b.Inventory, chip)

		// Keep them sorted.
		if len(b.Inventory) >= 2 && b.Inventory[0] > b.Inventory[1] {
			b.Inventory[0], b.Inventory[1] = b.Inventory[1], b.Inventory[0]
		}

		return true
	}
	return false
}

// Transfer gives a token to another bot.
func (b *Bot) Transfer(to *Bot, chip Microchip) bool {
	// If successful, prepare our new inventory.
	inventory := []Microchip{}
	for _, current := range b.Inventory {
		if current == chip {
			continue
		}
		inventory = append(inventory, current)
	}

	// We don't have this chip?
	if len(inventory) == len(b.Inventory) {
		return false
	}

	// Can the recipient take it?
	if to.Give(chip) {
		b.Inventory = inventory
		return true
	}

	return false
}

// Deposit puts a bot's chip into an output container.
func (b *Bot) Deposit(to *Output, chip Microchip) bool {
	// If successful, prepare our new inventory.
	inventory := []Microchip{}
	for _, current := range b.Inventory {
		if current == chip {
			continue
		}
		inventory = append(inventory, current)
	}

	// We don't have this chip?
	if len(inventory) == len(b.Inventory) {
		return false
	}

	to.Inventory = append(to.Inventory, chip)
	b.Inventory = inventory
	return true
}

// Type History is the things a bot has had to do.
type History struct {
	A Microchip // Two values the bot has had to compare.
	B Microchip
}

// Type Output is an output bin a bot can give chips to.
type Output struct {
	ID        string
	Inventory []Microchip
}

// Type Microchip is a chip that a bot is carrying. Each microchip contains
// a single number.
type Microchip int
