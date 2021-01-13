package fantasy

import (
	"fmt"
)

// Commands.
const (
	CommandGET    = "GET"
	CommandSET    = "SET"
	CommandDELETE = "DELETE"
	CommandPURGE  = "PURGE"
	CommandLEN    = "LEN"
)

func requireMinimum(require, has int) {
	if has < require {
		panic("Not enough arguments passed. Please check documentation")
	}
}

func parseCommand(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}
		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}
		if c == '\\' {
			escapeNext = true
			continue
		}
		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}
		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}
		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}
	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}
	if current != "" {
		args = append(args, current)
	}
	return args, nil
}
