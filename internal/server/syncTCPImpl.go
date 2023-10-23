package server

import (
	"diceDB/config"
	"io"
	"log"
	"net"
	"strconv"
)

type syncTCPImpl struct {
	configs *config.Configs
}

func readCommand(c net.Conn) (string, error) {
	// TODO: Max read in one shot is 512 bytes
	// To allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter
	var buf = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
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
			cmd, err := readCommand(client)
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

			log.Printf("received command %s", cmd)
			if err = respond(cmd, client); err != nil {
				log.Println("write error : ", err)
			}
		}
	}
}

func NewSyncTCPServer(configs *config.Configs) Server {
	return &syncTCPImpl{configs: configs}
}
