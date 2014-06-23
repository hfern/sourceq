package main

import (
	"github.com/Ronny95/goseq"
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
