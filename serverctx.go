package main

import (
	"fmt"
	"github.com/hfern/goseq"
	"log"
)

type ServerQueryOptions struct {
	NoPlayers bool `long:"no-players" short:"P" default:"false" description:"Don't list players."`
	NoInfo    bool `long:"no-info" short:"I" default:"false" description:"Don't list general info."`
	NoRules   bool `long:"no-rules" short:"R" default:"false" description:"Don't list server rules."`
}

var serverSingleOptions ServerQueryOptions

func serverctx(args []string) {

	if len(args) == 0 {
		log.Fatal("The first argument to sourceq server must be the address of the server. \n" +
			"e.g. sourceq server google.com\n" +
			"or sourceq server google.com:8080\n" +
			"or sourceq server 192.168.1.1:6060")
		return
	}

	fmt.Println(args)

	server := goseq.NewServer()
	server.SetAddress("sd")
}
