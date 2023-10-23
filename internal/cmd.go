package internal

type DiceCmd struct {
	Cmd  CommandName
	Args []string
}

type CommandName string

const PING CommandName = "PING"
