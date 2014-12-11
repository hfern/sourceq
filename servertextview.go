package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Ident struct {
	level        int
	padding      string
	prompt       string
	promptActive bool
}

var defaultIdent Ident = Ident{
	level:        0,
	padding:      "  ",
	prompt:       "* ",
	promptActive: false,
}

func viewServerText(options *ServerQueryOptions, servers []ServerAttrPair) {
	for _, server := range servers {
		textFormatServer(server, defaultIdent, options)
		fmt.Println("\n")
	}
}

func textFormatServer(server ServerAttrPair, ident Ident, options *ServerQueryOptions) {
	ident.Println("Server: ", server.Attrs.Address)
	ident.level++

	textFormatServerInfo(server.Attrs.Info, ident, options)
}

func textFormatServerInfo(info MaybeInfo, ident Ident, options *ServerQueryOptions) {
	if options.NoInfo {
		return
	}

	ident.Println("Info:")
	ident.level++

	if info.Error != nil {
		ident.Println("Error fetching server info: ", info.Error.Error())
		return
	}

	// TODO: Use fields by hand
	// Let package JSON do the hard reflection work for us.

	jsobj, err := json.Marshal(info.Info)
	if err != nil {
		ident.Println("Error marshalling json: ", err.Error())
	}

	mapped := make(map[string]interface{})

	err = json.Unmarshal(jsobj, &mapped)
	if err != nil {
		ident.Println("Error unmarshalling json: ", err.Error())
	}

	maxKeySize := 0

	for key, _ := range mapped {
		keySize := len(key)
		if keySize > maxKeySize {
			maxKeySize = keySize
		}
	}

	pad := func(n int) string {
		return strings.Repeat(" ", n)
	}

	filters := map[string]func(interface{}) interface{}{
		"SteamID": func(interface{}) interface{} {
			return info.Info.GetSteamID()
		},
	}

	keys := getKeys(mapped)
	sort.Strings(keys)

	for _, key := range keys {
		val := mapped[key]

		if _, ok := filters[key]; ok {
			val = filters[key](val)
		}

		ident.Printf("| %s:%s %v\n", key, pad(maxKeySize-len(key)), val)
	}
}

func (ident *Ident) Println(args ...interface{}) {
	fmt.Print(ident.GetPrefix())
	fmt.Println(args...)
}

func (ident *Ident) Printf(format string, args ...interface{}) {
	fmt.Print(ident.GetPrefix())
	fmt.Printf(format, args...)
}

func (ident *Ident) GetPrefix() string {
	return strings.Repeat(ident.padding, ident.level)
}

func getKeys(amap map[string]interface{}) []string {
	keys := make([]string, 0, len(amap))
	for key, _ := range amap {
		keys = append(keys, key)
	}
	return keys
}
