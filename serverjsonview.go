package main

import (
	"encoding/json"
	"fmt"
)

type jsonResponseWrapper struct {
	Servers []interface{}
}

func viewServerJSON(options *ServerQueryOptions, servers []ServerAttrPair) {
	formattedServers := make([]interface{}, len(servers))
	for i, server := range servers {
		formattedServers[i] = jsonFormatServer(server, options)
	}

	wrapper := jsonResponseWrapper{Servers: formattedServers}

	encoded, err := json.Marshal(wrapper)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(encoded))
}

func jsonFormatServer(server ServerAttrPair, opts *ServerQueryOptions) interface{} {
	return map[string]interface{}{
		"Address": server.Attrs.Address,
		"Players": jsonFormatPlayers(server.Attrs.Players, opts),
		"Info":    jsonFormatInfo(server.Attrs.Info, opts),
		"Rules":   jsonFormatRules(server.Attrs.Rules, opts),
	}
}

func jsonFormatRules(rules MaybeRules, opts *ServerQueryOptions) interface{} {
	if opts.NoRules {
		return nil
	}

	ret := struct {
		Error interface{}
		Rules interface{}
	}{}

	if rules.Error != nil {
		ret.Error = rules.Error.Error()
		ret.Rules = nil
		return ret
	}

	ret.Rules = make(map[string]interface{})

	return ret
}

func jsonFormatInfo(info MaybeInfo, opts *ServerQueryOptions) interface{} {
	if opts.NoInfo {
		return nil
	}
	ret := struct {
		Error interface{}
		Info  interface{}
	}{
		Error: nil,
		Info:  nil,
	}

	if info.Error == nil {
		ret.Info = info.Info
	} else {
		ret.Error = info.Error.Error()
	}

	return ret
}

func jsonFormatPlayers(mbplys MaybePlayers, opts *ServerQueryOptions) interface{} {
	players := mbplys.Players
	if opts.NoPlayers {
		return nil
	}
	fmtPlayers := make([]map[string]interface{}, len(players))

	for i, player := range players {
		ply := make(map[string]interface{})
		ply["Index"] = player.Index()
		ply["Name"] = player.Name()
		ply["Duration"] = player.Duration()
		ply["Score"] = player.Score()
		fmtPlayers[i] = ply
	}

	type ReturnStruct struct {
		Error   interface{}
		Players interface{}
	}

	ret := ReturnStruct{
		Error:   nil,
		Players: fmtPlayers,
	}

	if mbplys.Error != nil {
		ret.Error = mbplys.Error.Error()
	}

	return ret
}
