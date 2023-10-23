package main

import (
	"diceDB/config"
	"diceDB/internal"
	"diceDB/internal/io"
	"diceDB/internal/server"
	"flag"
	"fmt"
)

func getConfigs() *config.Configs {
	var configs *config.Configs
	configs = config.InitConfigWithDefaultValues()
	flag.StringVar(&configs.Host, "host", "0.0.0.0", "Host for diceDB")
	flag.IntVar(&configs.Port, "port", 7379, "port for diceDB")
	flag.Parse()
	return configs
}

func main() {
	configs := getConfigs()
	fmt.Println("rolling the dice 🎲")
	respDecoder := io.NewRESPDecoder()
	respEncoder := io.NewRESPEncoder()
	syncTCPServer := server.NewSyncTCPServer(configs, respDecoder, internal.NewEvaluator(respEncoder))
	syncTCPServer.Run()

}
