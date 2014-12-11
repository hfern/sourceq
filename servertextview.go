package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
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
	}
}

func textFormatServer(server ServerAttrPair, ident Ident, options *ServerQueryOptions) {

	if options.OnlyKeywords {
		listKeywords(server.Attrs.Info)
		return
	}

	ident.Println("Server: ", server.Attrs.Address)
	ident.level++

	textFormatServerInfo(server.Attrs.Info, ident, options)
	textFormatPlayers(server.Attrs.Players, ident, options)

	ident.level--
	ident.Println("")
}

func listKeywords(info MaybeInfo) {
	if info.Error != nil {
		log.Println("Error fetching AS_INFO keywords: ", info.Error.Error())
	}
	keywords := strings.Split(info.Info.GetKeywords(), ",")

	for _, kw := range keywords {
		fmt.Println(strings.TrimSpace(kw))
	}
}

func textFormatPlayers(players MaybePlayers, ident Ident, options *ServerQueryOptions) {
	if options.NoPlayers {
		return
	}

	ident.Println("Players:")
	ident.level++

	if players.Error != nil {
		ident.Println("Error fetching player list: ", players.Error.Error())
		return
	}

	type PlayerRow struct {
		RowID    string
		Name     string
		Index    string
		Score    string
		Duration string
	}

	plrRows := make([]PlayerRow, 0, len(players.Players)+1)
	plrRows = append(plrRows, PlayerRow{
		RowID:    "  ",
		Name:     "Name",
		Index:    "Id",
		Score:    "Scr",
		Duration: "Time",
	})

	for i, player := range players.Players {
		col := PlayerRow{
			RowID: strconv.Itoa(i + 1),
			Name:  player.Name(),
			Index: strconv.Itoa(player.Index()),
			Score: strconv.Itoa(player.Score()),
			// round to 1s
			Duration: (player.Duration() - (player.Duration() % time.Second)).String(),
		}
		plrRows = append(plrRows, col)
	}

	largest := func(things []PlayerRow, stringer func(PlayerRow) string) int {
		max := 0
		for _, thing := range things {
			sz := len(stringer(thing))
			if sz > max {
				max = sz
			}
		}
		return max
	}

	maxColumnSizes := struct {
		RowID    int
		Name     int
		Index    int
		Score    int
		Duration int
	}{
		RowID:    largest(plrRows, func(r PlayerRow) string { return r.RowID }),
		Name:     largest(plrRows, func(r PlayerRow) string { return r.Name }),
		Index:    largest(plrRows, func(r PlayerRow) string { return r.Index }),
		Score:    largest(plrRows, func(r PlayerRow) string { return r.Score }),
		Duration: largest(plrRows, func(r PlayerRow) string { return r.Duration }),
	}

	for i, row := range plrRows {
		ident.Printf(
			" %v | %v | %v | %v | %v \n",
			paddedR(row.RowID, maxColumnSizes.RowID),
			padded(row.Name, maxColumnSizes.Name),
			padded(row.Index, maxColumnSizes.Index),
			paddedR(row.Score, maxColumnSizes.Score),
			padded(row.Duration, maxColumnSizes.Duration),
		)

		if i == 0 {
			lengths := []int{
				maxColumnSizes.RowID,
				maxColumnSizes.Name,
				maxColumnSizes.Index,
				maxColumnSizes.Score,
				maxColumnSizes.Duration,
			}
			divider := ""

			for j, length := range lengths {
				if j != 0 {
					divider = divider + "+"
				}
				divider = divider + strings.Repeat("-", length+2)
			}

			ident.Printf("%s\n", divider)
		}
	}
	ident.Println("")
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

	filters := map[string]func(interface{}) interface{}{
		"SteamID": func(interface{}) interface{} {
			return info.Info.GetSteamID()
		},
	}

	keys := getKeys(mapped)
	sort.Strings(keys)

	ident.prompt = "| "
	ident.promptActive = true

	for _, key := range keys {
		val := mapped[key]

		if _, ok := filters[key]; ok {
			val = filters[key](val)
		}

		ident.Printf("%s:%s %v\n", key, pad(maxKeySize-len(key)), val)
	}

	ident.promptActive = false
	ident.Println("")
}

func (ident *Ident) Println(args ...interface{}) {
	fmt.Print(ident.GetPrefix())
	fmt.Println(args...)
}

func (ident *Ident) Printf(format string, args ...interface{}) (int, error) {
	fmt.Print(ident.GetPrefix())
	return fmt.Printf(format, args...)
}

func (ident *Ident) GetPrefix() string {
	cursor := ""
	if ident.promptActive {
		cursor = ident.prompt
	}
	return strings.Repeat(ident.padding, ident.level) + cursor
}

func getKeys(amap map[string]interface{}) []string {
	keys := make([]string, 0, len(amap))
	for key, _ := range amap {
		keys = append(keys, key)
	}
	return keys
}

func pad(n int) string {
	return strings.Repeat(" ", n)
}

func padded(text string, padLength int) string {
	return text + pad(padLength-len(text))
}

func paddedR(text string, padLength int) string {
	return pad(padLength-len(text)) + text
}
