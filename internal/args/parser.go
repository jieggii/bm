// Package args serves as a parser for command line arguments.
// It is written myself because I just wanted to do it myself.
// Also, I didn't find a command-line-parsing library that I liked enough to use it.
package args

import (
	"fmt"
)

// ProgramMeta represents meta information about a program.
type ProgramMeta struct {
	// Name of the program.
	Name string

	// Description of the program.
	Description string

	// Version of the program.
	Version string

	// Author of the program.
	Author string
}

type Command struct {
	// Name of the command.
	Name string

	// Usage of the command.
	Usage string

	// Description of the command.
	Description string

	// Expected number of arguments.
	ArgsNumber int

	// Function which will be called on command.
	Handler func(arguments []string) error
}

type Handler struct {
	// Basic information about the program.
	ProgramMeta ProgramMeta

	// Arguments to be parsed.
	Args []string

	// Slice of registered commands.
	Commands []*Command
}

func (h *Handler) getCommand(name string) (*Command, bool) {
	for _, command := range h.Commands {
		if name == command.Name {
			return command, true
		}
	}
	return nil, false
}

// PrintHelp prints help message.
func (h *Handler) PrintHelp() error {
	fmt.Println("Usage:")
	fmt.Printf("\t%v [--help/-h] <command> [--help/-h] [args...]\n\n", h.ProgramMeta.Name)

	fmt.Println("Description:")
	fmt.Printf("\t%v\n\n", h.ProgramMeta.Description)

	fmt.Println("Version:")
	fmt.Printf("\t%v\n\n", h.ProgramMeta.Version)

	fmt.Println("Commands:")
	maxUsageLen := 0
	for _, command := range h.Commands {
		usageLen := len(command.Usage)
		if usageLen > maxUsageLen {
			maxUsageLen = usageLen
		}
	}
	for _, command := range h.Commands {
		spacing := maxUsageLen - len(command.Usage) + 2
		fmt.Printf("\t%v %*s %v\n", command.Usage, spacing, "", command.Description)
	}
	fmt.Println()

	fmt.Println("Author:")
	fmt.Printf("\t%v\n", h.ProgramMeta.Author)

	return nil
}

func (h *Handler) Handle() error {
	if len(h.Args) < 2 { // if no arguments/commands were provided
		return h.PrintHelp()
	}
	if h.Args[1] == "-h" || h.Args[1] == "--help" { // if only help flag was provided
		return h.PrintHelp()
	}

	commandName := h.Args[1]
	command, found := h.getCommand(commandName)
	if !found {
		return fmt.Errorf("unknown command %v", commandName)
	}

	commandArgs := h.Args[2:]
	for _, arg := range commandArgs {
		if arg == "-h" || arg == "--help" { // if any of command flags is a help flag
			return h.PrintHelp()
		}
	}

	gotCommandArgs := len(commandArgs)
	if gotCommandArgs != command.ArgsNumber {
		return fmt.Errorf(
			"invalid number of arguments for command %v: expected %v, got %v",
			commandName, command.ArgsNumber, gotCommandArgs,
		)
	}

	return command.Handler(commandArgs)
}
