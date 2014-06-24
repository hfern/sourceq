package main

import (
	"fmt"
	"github.com/hfern/goseq"
)

var regionData = map[string]goseq.Region{
	"USE":   goseq.USEast,
	"USW":   goseq.USWest,
	"SA":    goseq.SouthAmerica,
	"EU":    goseq.Europe,
	"AS":    goseq.Asia,
	"AU":    goseq.Australia,
	"ME":    goseq.MiddleEast,
	"AF":    goseq.Africa,
	"OTHER": goseq.RestOfWorld,
}

var knownFilters = map[string]string{
	"type":       "Servers running (d)edicated, (l)isten, or (p) SourceTV.",
	"secure":     "(1) Servers using anti-cheat technology (VAC, but potentially others as well).",
	"gamedir":    "Servers running the specified modification (ex. cstrike)",
	"map":        "Servers running the specified map (ex. cs_italy)",
	"linux":      "Servers running on a Linux (1) platform",
	"empty":      "Servers that are not empty (1)",
	"full":       "Servers that are not full (1)",
	"proxy":      "Servers that are spectator proxies (1)",
	"napp":       "Servers that are NOT running game ([appid])",
	"noplayers":  "Servers that are empty (1)",
	"white":      "Servers that are whitelisted (1)",
	"gametype":   "Servers with all of the given tag(s) in sv_tags (tag1,tag2,...)",
	"gamedata":   "Servers with all of the given tag(s) in their 'hidden' tags (L4D2) (tag1,tag2,...)",
	"gamedataor": "Servers with any of the given tag(s) in their 'hidden' tags (L4D2) (tag1,tag2,...)",
}

func printKnownFiltersInfo() {
	fmt.Println("Known Filters")
	fmt.Println("(See https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol#Filter)")
	fmt.Println()

	for filter, detail := range knownFilters {
		fmt.Println("\t", filter, "\t", detail)
	}

	fmt.Println()
}
