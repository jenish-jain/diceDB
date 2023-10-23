package internal

import (
	"diceDB/internal/io"
	"errors"
	"fmt"
	"log"
	"net"
)

type evalImpl struct {
	encoder io.Encoder
}

type Evaluator interface {
	Do(command *DiceCmd, con net.Conn) error
}

func InvalidArgumentCountError(commandName CommandName) error {
	return fmt.Errorf("ERR wrong number of arguments for `%s` command", commandName)
}

func (e evalImpl) ping(args []string, con net.Conn) error {
	var b []byte

	if len(args) > 2 {
		return InvalidArgumentCountError(PING)
	}

	if len(args) == 0 {
		b = e.encoder.Encode("PONG", true)
	} else {
		b = e.encoder.Encode(args[0], true)
	}

	_, err := con.Write(b)
	return err
}

func (e evalImpl) Do(command *DiceCmd, con net.Conn) error {
	log.Printf("evaluating command : %s \n", command)

	switch command.Cmd {
	case PING:
		return e.ping(command.Args, con)
	default:
		return errors.New("unrecognized command")
	}
}

func NewEvaluator(encoder io.Encoder) Evaluator {
	return &evalImpl{encoder: encoder}
}
