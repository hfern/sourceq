package main

import (
	"fmt"
	"github.com/hfern/goseq"
	"log"
	"time"
)

type ServerQueryOptions struct {
	NoPlayers bool `long:"no-players" short:"P" default:"false" description:"Don't list players."`
	NoInfo    bool `long:"no-info" short:"I" default:"false" description:"Don't list general info."`
	NoRules   bool `long:"no-rules" short:"R" default:"false" description:"Don't list server rules."`
	Serial    bool `long:"serial" short:"s" default:"false" description:"Force serial querying of server attributes."`
	Json      bool `long:"json" default:"false" description:"Output as JSON to StdOut"`
	Timeout   uint `long:"timeout" short:"t" default="1" description:"Timeout for attribute queries in seconds."`
}

type DoneChannel chan int

const DONE int = 0

type MaybePlayers struct {
	err     error
	players []goseq.Player
}

type MaybeInfo struct {
	err  error
	info goseq.ServerInfo
}

type MaybeRules struct {
	err   error
	rules goseq.RuleMap
}

type ServerQueryableAttributes struct {
	address string
	players MaybePlayers
	info    MaybeInfo
	rules   MaybeRules
}

type ServerAttrPair struct {
	server goseq.Server
	attrs  ServerQueryableAttributes
}

type doIf struct {
	cond bool
	call func(DoneChannel)
}

var serverSingleOptions ServerQueryOptions

func serverctx(serverAddresses []string) {

	options := &serverSingleOptions

	if len(serverAddresses) == 0 {
		log.Fatal(
			"The first argument to sourceq server must be the address of the server. \n" +
				"e.g. sourceq server google.com\n" +
				"or sourceq server google.com:8080\n" +
				"or sourceq server 192.168.1.1:6060\n" +
				"or, for multiple servers: \n\tsourceq server google.com:80 example.org 123.156.178")
		return
	}

	timeout := time.Duration(options.Timeout) * time.Second
	done := make(DoneChannel)

	servers := make([]ServerAttrPair, len(serverAddresses))

	for i, serverAddr := range serverAddresses {
		servers[i].attrs.address = serverAddr
		server := goseq.NewServer()
		err := server.SetAddress(serverAddr)
		if err != nil {
			log.Fatal(err)
		}
		servers[i].server = server
	}

	for i, _ := range servers {
		server := &servers[i]
		go queryServer(options, &server.server, &server.attrs, timeout, done)
		if options.Serial {
			<-done
		}
	}

	if !options.Serial {
		for _ = range servers {
			<-done
		}
	}

	fmt.Println(servers)
}

func queryServer(
	options *ServerQueryOptions,
	server *goseq.Server,
	attrs *ServerQueryableAttributes,
	timeout time.Duration,
	serverIsDone DoneChannel) {

	defer func() { serverIsDone <- DONE }()

	done := make(DoneChannel)

	callbacks := []doIf{
		{
			cond: !options.NoInfo,
			call: func(donner DoneChannel) {
				getServerInfo(server, &attrs.info, timeout, donner)
			},
		},
		{
			cond: !options.NoRules,
			call: func(donner DoneChannel) {
				getServerRules(server, &attrs.rules, timeout, donner)
			},
		},
		{
			cond: !options.NoPlayers,
			call: func(donner DoneChannel) {
				getServerPlayers(server, &attrs.players, timeout, donner)
			},
		},
	}

	for _, cb := range callbacks {
		if cb.cond {
			go cb.call(done)
			if options.Serial {
				<-done
			}
		}
	}

	if !options.Serial {
		for _ = range callbacks {
			<-done
		}
	}

}

func getServerInfo(server *goseq.Server, info *MaybeInfo, timeout time.Duration, donner DoneChannel) {
	defer func() { donner <- DONE }()
	info.info, info.err = (*server).Info(timeout)
}

func getServerRules(server *goseq.Server, rules *MaybeRules, timeout time.Duration, donner DoneChannel) {
	defer func() { donner <- DONE }()
	rules.rules, rules.err = (*server).Rules(timeout)
}

func getServerPlayers(server *goseq.Server, plys *MaybePlayers, timeout time.Duration, donner DoneChannel) {
	defer func() { donner <- DONE }()
	plys.players, plys.err = (*server).Players(timeout)
}
