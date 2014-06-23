package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Ronny95/goseq"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var _region = flag.String("region", "USW", "Region code to get results for. \n"+
	"One of USE, USW, SA, EU, AS, AU, ME, AF, OTHER.")

var _asyncOkay = flag.Bool("async", true, "Allow async sub-querying of Source Servers to get info. \n"+
	"Extreme increase of performance with this enabled.")

var _printFields = flag.String("fields", "ip=21,name", "The fields to be included. Optionally includes the min-length space."+
	"See -showfields")

var _fieldDivider = flag.String("divider", " ¦ ", "Characters used to seperate fields.")

var _startIp = flag.String("start", goseq.NoAddress, "Where to start reading IPs from. Defaults to start.")

var _showHeader = flag.Bool("header", true, "Show header w/ column names.")

var _fieldregexp = regexp.MustCompile(`\s*(([a-z]+)(=(\d+))?)\s*,?\s*`)

type FieldSpec struct {
	name   string
	length int
}

type SvResponse struct {
	err    error
	server goseq.Server
	info   goseq.ServerInfo
}

func masterctx() {

	userRegionStr := strings.ToUpper(*_region)

	region, found := regionData[userRegionStr]

	if !found {
		fmt.Errorf("Region '%s' does not exist.", *_region)
		return
	}

	fields, err := parseFields(*_printFields, serverFieldProperties)

	if err != nil {
		log.Fatal(err)
	}

	master := goseq.NewMasterServer()
	master.SetRegion(region)
	master.SetAddr(goseq.MasterSourceServers[0])

	servers, err := master.Query(*_startIp)
	numServers := len(servers)

	if err != nil {
		log.Fatal(err)
	}

	rec := make(chan SvResponse)
	printer := make(chan SvResponse)

	//tups := make([]SvResponse, 0, numServers)

	//go serialQueryServers(rec, servers, 1*time.Second)
	go AsyncQueryServers(rec, servers, 1*time.Second)
	go printServerLine(fields, printer)

	if *_showHeader {
		printHeaderLine(fields, serverFieldProperties)
	}

	for i := 0; i < numServers; i++ {
		//tups = append(tups, <-rec)
		printer <- <-rec
	}

	close(rec)
	close(printer)
}

func serialQueryServers(send chan SvResponse, servers []goseq.Server, timeout time.Duration) {
	for _, server := range servers {
		info, err := server.Info(timeout)
		if err != nil {
			send <- SvResponse{err: err, server: server}
		}
		send <- SvResponse{err: err, server: server, info: info}
	}
}

func AsyncQueryServers(send chan SvResponse, servers []goseq.Server, timeout time.Duration) {
	for _, server := range servers {
		go serialQueryServers(send, []goseq.Server{server}, timeout)
	}
}

func parseFields(spec string, properties map[string]FieldProperty) ([]FieldSpec, error) {
	specs := make([]FieldSpec, 0, 1)
	found := _fieldregexp.FindAllStringSubmatch(spec, -1)

	if len(found) == 0 && len(strings.TrimSpace(spec)) != 0 {
		return nil, errors.New("Couldn't parse field list.")
	}

	for _, match := range found {

		if _, ok := serverMethodAccessors[match[2]]; !ok {
			if _, ok = serverProperties[match[2]]; !ok {
				return nil, errors.New("Attempted to use an unregistered field!")
			}
		}

		sp := FieldSpec{
			name:   match[2],
			length: -1,
		}

		if def, ok := properties[match[2]]; ok {
			sp.length = def.size
		}

		if match[4] != "" {
			val, err := strconv.Atoi(match[4])
			if err == nil {
				sp.length = val
			}
		}

		specs = append(specs, sp)
	}

	return specs, nil
}

func printServerLine(fields []FieldSpec, in <-chan SvResponse) {

	for sv := range in {
		for i, field := range fields {
			if i > 0 {
				fmt.Print(*_fieldDivider)
			}

			var val interface{}

			if handler, ok := serverMethodAccessors[field.name]; ok {
				val = handler(sv.info)
			}

			if handler, ok := serverProperties[field.name]; ok {
				val = handler(sv.server)
			}

			if transformer, ok := serverFieldTransformers[field.name]; ok {
				val = transformer(val)
			}

			written, _ := fmt.Print(val)

			for ; written < field.length; written++ {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func printHeaderLine(fields []FieldSpec, props map[string]FieldProperty) {
	for _, field := range fields {

		title := field.name

		if prop, ok := props[field.name]; ok {
			title = prop.name
		}

		sz := len(title)

		var padL int = 0
		var padR int = 0

		if sz < field.length {
			var rem int = field.length - sz
			padL = rem / 2
			padR = rem - padL
		}

		for j := 0; j < padL; j++ {
			fmt.Print(" ")
		}

		fmt.Print(title)

		for j := 0; j < padR; j++ {
			fmt.Print(" ")
		}
	}
	fmt.Print("\n")
}
