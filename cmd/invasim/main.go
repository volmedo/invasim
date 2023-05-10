package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/volmedo/invasim/internal/aliens"
	"github.com/volmedo/invasim/internal/simulation"
	"github.com/volmedo/invasim/internal/worldmap"
)

const MAX_ITERATIONS = 10_000

func main() {
	var mapFilePath string
	flag.StringVar(&mapFilePath, "map", "", "path to a file to read the world map from")

	var numAliens int
	flag.IntVar(&numAliens, "aliens", 0, "number of aliens to unleash. It must not be greater than the number of cities in the map")

	flag.Parse()

	if mapFilePath == "" {
		fmt.Println("-map: a path to a map file is required and cannot be blank")
		flag.Usage()
		os.Exit(42)
	}

	if numAliens == 0 {
		fmt.Println("-aliens: a number of aliens greater than 0 is required")
		flag.Usage()
		os.Exit(42)
	}

	world, err := worldmap.ReadFromFile(mapFilePath)
	if err != nil {
		fatalf("Error reading map file: %v", err)
	}

	alienTracker, err := aliens.NewTracker(numAliens, world)
	if err != nil {
		fatalf("Error placing aliens on their starting positions: %v", err)
	}

	simulation.Run(world, alienTracker, MAX_ITERATIONS, os.Stdout)
}

func fatalf(format string, v ...any) {
	fmt.Printf(format+"\n", v...)
	os.Exit(42)
}
