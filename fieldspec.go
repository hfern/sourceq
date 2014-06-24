package main

import (
	"fmt"
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
	full string
}

var serverFieldProperties = map[string]FieldProperty{
	"bots":          FieldProperty{name: "Bots", size: 5, full: "Number of Bots"},
	"duration":      FieldProperty{name: "Arrest In", size: 7, full: "Will arrest in (The Ship)"},
	"environment":   FieldProperty{name: "Env", size: 3, full: "Environment OS"},
	"folder":        FieldProperty{name: "Folder", size: 10},
	"game":          FieldProperty{name: "Game", size: 5},
	"gameid":        FieldProperty{name: "GameID", size: 6},
	"id":            FieldProperty{name: "ID", size: 5},
	"keywords":      FieldProperty{name: "Keywords", size: 9},
	"map":           FieldProperty{name: "Map", size: 10},
	"maxplayers":    FieldProperty{name: "Max", size: 3, full: "Maximum number of players allowed"},
	"mode":          FieldProperty{name: "Mode", size: 4},
	"name":          FieldProperty{name: "Name", size: 15, full: "Name of Server"},
	"players":       FieldProperty{name: "Ply", size: 3, full: "Number Players"},
	"port":          FieldProperty{name: "Port", size: 5},
	"servertype":    FieldProperty{name: "Type", size: 5, full: "Hosting Type (eg dedicated)"},
	"spectatorname": FieldProperty{name: "Spectator", size: 9},
	"spectatorport": FieldProperty{name: "SpPort", size: 7},
	"steamid":       FieldProperty{name: "SteamID", size: 10},
	"vac":           FieldProperty{name: "VAC", size: 3, full: "Is the server VAC protected?"},
	"version":       FieldProperty{name: "Version", size: 5},
	"visibility":    FieldProperty{name: "Pw.", size: 3, full: "Is a password required to join?"},
	"witnesses":     FieldProperty{name: "Witnesses", size: 10, full: "# Witnesses for The Ship."},
	"ip":            FieldProperty{name: "IP Addr", size: 21, full: "IP Address of the Server"},
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

func printServerFieldProperties() {
	fmt.Println("Server Fields:")

	longest := 0
	for field, _ := range serverFieldProperties {
		if len(field) > longest {
			longest = len(field)
		}
	}

	for field, prop := range serverFieldProperties {

		fmt.Print("    ", field)

		for i := len(field); i < longest; i++ {
			fmt.Print(" ")
		}

		desc := prop.name

		if prop.full != "" {
			desc = prop.full
		}

		fmt.Printf("    %s (size %d)\n", desc, prop.size)
	}

	fmt.Println()
}
