package main

import (
	"github.com/hfern/goseq"
	"log"
	"time"
)

type ServerQueryOptions struct {
	NoPlayers    bool `long:"no-players" short:"P" default:"false" description:"Don't list players."`
	NoInfo       bool `long:"no-info" short:"I" default:"false" description:"Don't list general info."`
	NoRules      bool `long:"no-rules" short:"R" default:"false" description:"Don't list server rules."`
	Serial       bool `long:"serial" short:"s" default:"false" description:"Force serial querying of server attributes."`
	Json         bool `long:"json" default:"false" description:"Output as JSON to StdOut"`
	Timeout      uint `long:"timeout" short:"t" default:"2" description:"Timeout for attribute queries in seconds."`
	OnlyKeywords bool `long:"only-keywords" short:"K" default:"false" description:"Only list the keywords of the servers one per line."`
}

type DoneChannel chan int

const DONE int = 0

type MaybePlayers struct {
	Error   error
	Players []goseq.Player
}

type MaybeInfo struct {
	Error error
	Info  goseq.ServerInfo
}

type MaybeRules struct {
	Error error
	Rules goseq.RuleMap
}

type ServerQueryableAttributes struct {
	Address string
	Players MaybePlayers
	Info    MaybeInfo
	Rules   MaybeRules
}

type ServerAttrPair struct {
	Server goseq.Server
	Attrs  ServerQueryableAttributes
}

type doIf struct {
	cond bool
	call func(DoneChannel)
}

var serverSingleOptions ServerQueryOptions

func serverctx(serverAddresses []string) {
	options := &serverSingleOptions

	if !assertLogicalServerFlags(options) {
		return
	}

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
		servers[i].Attrs.Address = serverAddr
		server := goseq.NewServer()
		err := server.SetAddress(serverAddr)
		if err != nil {
			log.Fatal(err)
		}
		servers[i].Server = server
	}

	for i, _ := range servers {
		server := &servers[i]
		go queryServer(options, &server.Server, &server.Attrs, timeout, done)
		if options.Serial {
			<-done
		}
	}

	if !options.Serial {
		for _ = range servers {
			<-done
		}
	}

	if options.Json {
		viewServerJSON(options, servers)
	} else {
		viewServerText(options, servers)
	}
}

func assertLogicalServerFlags(options *ServerQueryOptions) bool {
	if options.OnlyKeywords {
		if options.Json {
			log.Fatal("--only-keywords cannot be used with --json")
			return false
		}
		options.NoRules = true
		options.NoPlayers = true
		options.NoInfo = false
	}
	return true
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
				getServerInfo(server, &attrs.Info, timeout, donner)
			},
		},
		{
			cond: !options.NoRules,
			call: func(donner DoneChannel) {
				getServerRules(server, &attrs.Rules, timeout, donner)
			},
		},
		{
			cond: !options.NoPlayers,
			call: func(donner DoneChannel) {
				getServerPlayers(server, &attrs.Players, timeout, donner)
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
	info.Info, info.Error = (*server).Info(timeout)
}

func getServerRules(server *goseq.Server, rules *MaybeRules, timeout time.Duration, donner DoneChannel) {
	defer func() { donner <- DONE }()
	rules.Rules, rules.Error = (*server).Rules(timeout)
}

func getServerPlayers(server *goseq.Server, plys *MaybePlayers, timeout time.Duration, donner DoneChannel) {
	defer func() { donner <- DONE }()
	plys.Players, plys.Error = (*server).Players(timeout)
}
