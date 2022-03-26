package utils

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
)

//Errors to return
var ErrNoArgument = errors.New("No user argument found!")
var ErrWrongCommand = errors.New("Unknown command type was used!")
var ErrWrongSearchCommand = errors.New("One argument was given for search command!")
var ErrWrongBuyCommand = errors.New("Two positive number arguments are required for buy command!")
var ErrWrongDeleteCommand = errors.New("A positive number argument are required for delete command!")

var checkYourCommandMessage = "Please check your command."

//Check args to check if there is a mistake with the given arguments
func CheckUserArguments(args []string) (string, error) {
	projectName := path.Base(args[0])
	if err := checkSize(args); err != nil {
		printNoArgument()
		printOptions(projectName)
		return "", err
	} else {
		command, err := checkCommand(args)
		if err != nil {
			switch err {
			case ErrWrongCommand:
				printWrongCommandUsed()
			case ErrWrongSearchCommand:
				PrintLineMessage(ErrWrongBuyCommand.Error() + checkYourCommandMessage)
			case ErrWrongBuyCommand:
				PrintLineMessage(ErrWrongBuyCommand.Error() + checkYourCommandMessage)
			case ErrWrongDeleteCommand:
				PrintLineMessage(ErrWrongDeleteCommand.Error() + checkYourCommandMessage)
			}
			printOptions(projectName)
		}
		return command, err
	}
}

//Check argument size
//Should be at least 2.
func checkSize(args []string) error {
	if len(args) == 1 {
		return ErrNoArgument
	}
	return nil
}

//Check argument size if search command is needed.
//Should be at least 3
func checkSearchCommandArgumentSize(args []string) (string, error) {
	if len(args) <= 2 {
		return SearchCommand, ErrWrongSearchCommand
	}
	return SearchCommand, nil
}

//Check argument size if buy command is needed.
//Should be at least 4
//Two positive integers needed for this command, give error otherwise.
func checkBuyCommandArguments(args []string) (string, error) {
	if len(args) <= 3 {
		return BuyCommand, ErrWrongBuyCommand
	}
	number, err := strconv.Atoi(args[2])
	if number <= 0 || err != nil {
		return BuyCommand, ErrWrongBuyCommand
	}
	number, err = strconv.Atoi(args[3])
	if number <= 0 || err != nil {
		return BuyCommand, ErrWrongBuyCommand
	}
	return BuyCommand, nil
}

//Check argument size if search command is needed.
//Should be at least 3
//One positive integer needed for this command, give error otherwise
func checkDeleteCommandArguments(args []string) (string, error) {
	if len(args) <= 2 {
		return DeleteCommand, ErrWrongDeleteCommand
	}
	number, err := strconv.Atoi(args[2])
	if number <= 0 || err != nil {
		return DeleteCommand, ErrWrongDeleteCommand
	}
	return DeleteCommand, nil
}

//Check if second program argument contains one of list,
//search, buy, or delete commands
func checkCommand(args []string) (string, error) {
	command := strings.ToLower(args[1])
	if strings.EqualFold(command, ListCommand) {
		return ListCommand, nil
	} else if strings.EqualFold(command, SearchCommand) {
		return checkSearchCommandArgumentSize(args)
	} else if strings.EqualFold(command, BuyCommand) {
		return checkBuyCommandArguments(args)
	} else if strings.EqualFold(command, DeleteCommand) {
		return checkDeleteCommandArguments(args)
	} else {
		return "", ErrWrongCommand
	}
}

func printNoArgument() {
	PrintLineMessage("No command found. Please check your command!")
}

func printWrongCommandUsed() {
	PrintLineMessage("This command is unknown to the program. Please check your command!")
}

func printOptions(projectName string) {
	fmt.Printf("%s Application Usage:\n", projectName)
	PrintLineMessage("\tgo run main.go <command> [arguments]")
	PrintLineMessage("The commands are:")
	PrintLineMessage("\tlist")
	PrintLineMessage("\tsearch [book-name or ISBN or stock-code]")
	PrintLineMessage("\tbuy [ID count]")
	PrintLineMessage("\tdelete [ID]")
}
