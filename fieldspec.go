package main

import (
	"github.com/hfern/goseq"
)

type Any interface{}

type ServerMethod func(sv goseq.Server) interface{}
type ServerInfoMethod func(sv goseq.ServerInfo) interface{}

var serverMethodAccessors = map[string]ServerInfoMethod{
	"bots":          func(sv goseq.ServerInfo) interface{} { return sv.GetBots() },
	"duration":      func(sv goseq.ServerInfo) interface{} { return sv.GetDuration() },
	"environment":   func(sv goseq.ServerInfo) interface{} { return sv.GetEnvironment() },
	"folder":        func(sv goseq.ServerInfo) interface{} { return sv.GetFolder() },
	"game":          func(sv goseq.ServerInfo) interface{} { return sv.GetGame() },
	"gameid":        func(sv goseq.ServerInfo) interface{} { return sv.GetGameID() },
	"id":            func(sv goseq.ServerInfo) interface{} { return sv.GetID() },
	"keywords":      func(sv goseq.ServerInfo) interface{} { return sv.GetKeywords() },
	"map":           func(sv goseq.ServerInfo) interface{} { return sv.GetMap() },
	"maxplayers":    func(sv goseq.ServerInfo) interface{} { return sv.GetMaxPlayers() },
	"mode":          func(sv goseq.ServerInfo) interface{} { return sv.GetMode() },
	"name":          func(sv goseq.ServerInfo) interface{} { return sv.GetName() },
	"players":       func(sv goseq.ServerInfo) interface{} { return sv.GetPlayers() },
	"port":          func(sv goseq.ServerInfo) interface{} { return sv.GetPort() },
	"servertype":    func(sv goseq.ServerInfo) interface{} { return sv.GetServertype() },
	"spectatorname": func(sv goseq.ServerInfo) interface{} { return sv.GetSpectatorName() },
	"spectatorport": func(sv goseq.ServerInfo) interface{} { return sv.GetSpectatorPort() },
	"steamid":       func(sv goseq.ServerInfo) interface{} { return sv.GetSteamID() },
	"vac":           func(sv goseq.ServerInfo) interface{} { return sv.GetVAC() },
	"version":       func(sv goseq.ServerInfo) interface{} { return sv.GetVersion() },
	"visibility":    func(sv goseq.ServerInfo) interface{} { return sv.GetVisibility() },
	"witnesses":     func(sv goseq.ServerInfo) interface{} { return sv.GetWitnesses() },
}

var serverProperties = map[string]ServerMethod{
	"ip": func(sv goseq.Server) interface{} { return sv.Address() },
}

type FieldProperty struct {
	name string
	size int
}

var serverFieldProperties = map[string]FieldProperty{
	"bots":          FieldProperty{name: "Bots", size: 5},
	"duration":      FieldProperty{name: "Arrest In", size: 7},
	"environment":   FieldProperty{name: "Env", size: 3},
	"folder":        FieldProperty{name: "Folder", size: 10},
	"game":          FieldProperty{name: "Game", size: 5},
	"gameid":        FieldProperty{name: "GameID", size: 6},
	"id":            FieldProperty{name: "ID", size: 5},
	"keywords":      FieldProperty{name: "Keywords", size: 9},
	"map":           FieldProperty{name: "Map", size: 10},
	"maxplayers":    FieldProperty{name: "Max", size: 3},
	"mode":          FieldProperty{name: "Mode", size: 4},
	"name":          FieldProperty{name: "Name", size: 15},
	"players":       FieldProperty{name: "Ply", size: 3},
	"port":          FieldProperty{name: "Port", size: 5},
	"servertype":    FieldProperty{name: "Type", size: 5},
	"spectatorname": FieldProperty{name: "Spectator", size: 9},
	"spectatorport": FieldProperty{name: "SpPort", size: 7},
	"steamid":       FieldProperty{name: "SteamID", size: 10},
	"vac":           FieldProperty{name: "VAC", size: 3},
	"version":       FieldProperty{name: "Version", size: 5},
	"visibility":    FieldProperty{name: "Pw.", size: 3},
	"witnesses":     FieldProperty{name: "Witnesses", size: 10},
	"ip":            FieldProperty{name: "IP Addr", size: 21},
}

type FieldTransformer func(Any) Any

var serverFieldTransformers = map[string]FieldTransformer{
	"environment": transformEnvironment,
}

func transformEnvironment(in Any) Any {
	switch t := in.(type) {
	case goseq.ServerEnvironment:
		switch t {
		case goseq.Linux:
			return "Lnx"
		case goseq.Windows:
			return "Win"
		default:
			return "?"
		}
	default:
		return "_"
	}
}
