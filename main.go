package main

import (
	"flag"
	"fmt"
)

type Context int

const (
	MASTERQUERY Context = iota
	SINGLESERVER
)

var ctx Context

var userSetContext *string = flag.String("query", "master", "One of master,server. ")

func main() {

	flag.Parse()

	switch *userSetContext {
	case "master":
		ctx = MASTERQUERY
		masterctx()
	case "server":
		ctx = SINGLESERVER
		serverctx()
	default:
		fmt.Errorf("-query must be either \"master\" or \"server\"!\n")
	}
}
