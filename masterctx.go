package main

import (
	"errors"
	"fmt"
	"github.com/hfern/goseq"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MasterQueryOptions struct {
	Region string `long:"region" short:"r" default:"USW" description:"Region code to get results for. See --list-regions"`
	Serial bool   `long:"serial" short:"s" default:"false" description:"Force serial querying of the servers."`
	Fields string `long:"fields" default:"ip=21,name" description:"The fields to be included. Optionally includes the min-length space. See --list-fields" `
	// TODO(hunter): Add this
	ShowFields bool `long:"show-fields" default:"false" description:"Print details on each available field."`
	// TODO(hunter): Add this
	MasterIP string `long:"ip" default:"" description:"IP of the Master server to query."`
	Divider  string `long:"divider" default:" ¦ " description:"Characters used to seperate fields."`
	// TODO(hunter): Add this
	StartIP          string `long:"start" default:"" description:"Where to start reading IPs from. Defaults to start of list."`
	Limit            int    `long:"limit" short:"l" default:"0" description:"Limit the result set to n successful rows."`
	NoHeader         bool   `long:"no-header" default:"false" description:"Don't show header w/ column names."`
	ShowUnreachable  bool   `long:"unreachable" short:"U" default:"false" description:"Show unreachable servers (couldn't be connected to)."`
	ShowErrorSummary bool   `long:"errors" short:"E" default:"false" description:"Show error summary at end of list."`
	// TODO(hunter): Add this
	Filters map[string]string `long:"filter" short:"f" description:"Filters to use. See --list-filters"`
	// TODO(hunter): Add this
	ListFilters bool `long:"list-filters" default:"false" description:"List known filters." group:"Lists"`
	ListRegions bool `long:"list-regions" default:"false" description:"List Regions." group:"Lists"`
	ListFields  bool `long:"list-fields" default:"false" description:"List Server Fields." group:"Lists"`
}

var masterOptions MasterQueryOptions

var _fieldregexp = regexp.MustCompile(`\s*(([a-z]+)(=(\d+))?)\s*,?\s*`)

type infoPrinter func()

func masterctx() {
	log.SetFlags(0)

	if printInfo() {
		return
	}

	unreachable := 0
	errorsEncountererd := make([]error, 0)

	userRegionStr := strings.ToUpper(masterOptions.Region)

	region, found := regionData[userRegionStr]

	if !found {
		fmt.Errorf("Region '%s' does not exist.", masterOptions.Region)
		return
	}

	fields, err := parseFields(masterOptions.Fields, serverFieldProperties)

	if err != nil {
		log.Fatal(err)
	}

	master := goseq.NewMasterServer()
	master.SetRegion(region)
	master.SetAddr(goseq.MasterSourceServers[0])

	startIp := string(goseq.NoAddress)

	if masterOptions.StartIP != "" {
		startIp = masterOptions.StartIP
	}

	// Apply filters
	filt := master.GetFilter()
	for fname, fval := range masterOptions.Filters {
		filt.Set(fname, fval)
	}
	master.SetFilter(filt)

	servers, err := master.Query(startIp)
	numServers := len(servers)

	if err != nil {
		log.Fatal(err)
	}

	rec := make(chan SvResponse)
	printer := make(chan SvResponse)

	//tups := make([]SvResponse, 0, numServers)

	if !masterOptions.Serial {
		go AsyncQueryServers(rec, servers, 1*time.Second)
	} else {
		go serialQueryServers(rec, servers, 1*time.Second)
	}

	var writer Printer = &textWriter{}

	writer.Init(fields, printer)
	go writer.Run()

	if !masterOptions.NoHeader {
		printHeaderLine(fields, serverFieldProperties)
	}

	for i := 0; i < numServers; i++ {
		//tups = append(tups, <-rec)
		recd := <-rec

		if recd.err != nil {
			errorsEncountererd = append(errorsEncountererd, recd.err)
		}

		if recd.err != nil && !masterOptions.ShowUnreachable {
			unreachable++
			continue
		}

		if masterOptions.Limit > 0 && i >= masterOptions.Limit {
			continue
		}

		printer <- recd
	}

	writer.Done()

	close(rec)
	close(printer)

	log.Println()

	if !masterOptions.ShowUnreachable {
		log.Println(unreachable, "unreachable servers were hidden.")
	}

	if masterOptions.ShowErrorSummary {
		log.Printf("Errors Encountered (%dx):\n", len(errorsEncountererd))

		for _, detail := range errorsEncountererd {
			log.Println("\t", detail)
		}
	}
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

func printHeaderLine(fields []FieldSpec, props map[string]FieldProperty) {
	for i, field := range fields {

		if i > 0 {
			fmt.Print(masterOptions.Divider)
		}

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

type ErrorCount struct {
	err   error
	count int
}

func printInfo() (done bool) {
	done = false

	infos := []struct {
		f       infoPrinter
		enabled bool
	}{
		{printKnownFiltersInfo, masterOptions.ListFilters},
		{printRegionInfo, masterOptions.ListRegions},
		{printServerFieldProperties, masterOptions.ListFields},
	}

	for _, info := range infos {
		if info.enabled {
			info.f()
			done = true
		}
	}

	return
}
