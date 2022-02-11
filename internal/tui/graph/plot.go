package graph

import (
	"math/rand"
	"time"

	"github.com/guptarohit/asciigraph"
)

var (
	vals                       = 100
	fps                float64 = 24
	realTimeDataBuffer         = 63
	height                     = 2
)

func Plot() string {
	rand.Seed(time.Now().Local().Unix())

	graph := func(len int) string {
		randVals := make([]float64, len)

		for i := range randVals {
			randVals[i] = rand.Float64()
		}
		return asciigraph.Plot(randVals, asciigraph.Height(height), asciigraph.Caption("Packets per millisecond"))
	}

	return graph(100)
}
