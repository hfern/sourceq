package main

import (
	"flag"
	"github.com/jessevdk/go-flags"
)

type Context int

const (
	MASTERQUERY Context = iota
	SINGLESERVER
)

type MainOptions struct {
	Master MasterQueryOptions `command:"master"`
}

var ctx Context

var userSetContext *string = flag.String("query", "master", "One of master,server. ")

func main() {

	//flag.Parse()

	parser := flags.NewNamedParser("Source Query", flags.Default)
	parser.AddCommand("master", "Query Master Server",
		"Query the Master Server for a list of Source servers. "+
			"Display servers in row format.", &masterOptions)

	_, err := parser.Parse()

	if err != nil {
		return
	}

	switch parser.Active.Name {
	case "master":
		ctx = MASTERQUERY
		masterctx()
	case "server":
		ctx = SINGLESERVER
		serverctx()
	}
}
