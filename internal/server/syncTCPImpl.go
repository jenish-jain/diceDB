package server

import (
	"diceDB/config"
	"diceDB/internal"
	diceIO "diceDB/internal/io"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type syncTCPImpl struct {
	configs   *config.Configs
	decoder   diceIO.Decoder
	evaluator internal.Evaluator
}

func (s syncTCPImpl) decodeByteToTokens(data []byte) ([]string, error) {
	token, err := s.decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return []string{token.(string)}, nil
}

func (s syncTCPImpl) readCommand(c net.Conn) (*internal.DiceCmd, error) {
	// TODO: Max read in one shot is 512 bytes
	// To allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter
	var buf = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}
	tokens, err := s.decodeByteToTokens(buf[:n])
	if err != nil {
		return nil, err
	}
	return &internal.DiceCmd{
		Cmd:  internal.CommandName(strings.ToUpper(tokens[0])),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c net.Conn) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func (s syncTCPImpl) respond(cmd *internal.DiceCmd, c net.Conn) {
	err := s.evaluator.Do(cmd, c)
	if err != nil {
		respondError(err, c)
	}
}

func (s syncTCPImpl) Run() {
	log.Printf("starting a synchronous TCP server on %s:%d \n", s.configs.Host, s.configs.Port)

	var conClients int = 0

	//	listening to configured host:port
	listener, err := net.Listen("tcp", s.configs.Host+":"+strconv.Itoa(s.configs.Port))
	if err != nil {
		panic(err)
	}

	for {
		// blocking call waiting for new clients to connect

		client, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		// increment the number of concurrent clients
		conClients += 1
		log.Printf("client connected with address %+v , total concurrent client count %d \n", client.RemoteAddr(), conClients)

		for {
			// over the socket, continuously read the commands and print it out
			cmd, err := s.readCommand(client)
			if err != nil {
				err := client.Close()
				if err != nil {
					panic(err)
				}
				conClients -= 1
				log.Printf("client with remote address %+v disconnected !, remaining concurrent client count %d \n", client.RemoteAddr(), conClients)
				if err == io.EOF {
					break
				}
				log.Printf("error reading client command %+v \n", err)
			}

			s.respond(cmd, client)
		}
	}
}

func NewSyncTCPServer(configs *config.Configs, decoder diceIO.Decoder, evaluator internal.Evaluator) Server {
	return &syncTCPImpl{configs: configs, decoder: decoder, evaluator: evaluator}
}
