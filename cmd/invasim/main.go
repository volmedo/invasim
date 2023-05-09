package main

import (
	"flag"
	"log"
)

const MAX_ITERATIONS = 10_000

func main() {
	// parse command-line arguments
	var mapFilePath string
	flag.StringVar(&mapFilePath, "map", "", "path to a file to read the world map from")

	var numAliens int
	flag.IntVar(&numAliens, "aliens", 0, "number of aliens to unleash. It must not be greater than the number of cities in the map")

	flag.Parse()

	// build the world map

	// validate arguments

	// random initial placement for aliens

	// start simulation loop

	// check final conditions: either we reached MAX_ITERATIONS or all aliens were destroyed
	log.Printf("Simulation finished!\n")
	if numAliens == 0 {
		log.Printf("All aliens were destroyed!\n")
	} else {
		log.Printf("Max iterations reached, %d aliens remaining\n", numAliens)
	}

	// print how the world looks like after the invasion
}
