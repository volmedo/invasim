package alienmap

import (
	"fmt"

	"github.com/volmedo/invasim/internal/worldmap"
)

// Aliens keeps track of the city each alien is currently at
type Aliens map[string]string

// New creates a new Aliens map with numAliens aliens placed randomly in one of the cities of world.
// Since there can only be an alien in a city, numAliens cannot be greater than the number of cities in world.
func New(numAliens int, world worldmap.World) (Aliens, error) {
	if numAliens > len(world) {
		return Aliens{}, fmt.Errorf("not enough cities (%d) to place %d aliens", len(world), numAliens)
	}

	aliens := Aliens{}

	// to place the aliens randomly we can take advantage of the inherently non-deterministic ordering of elements
	// observed when ranging over a map in Go
	alien := 0
	for city := range world {
		if alien >= numAliens {
			break
		}

		name := fmt.Sprintf("alien %d", alien)
		aliens[name] = city
		alien++
	}

	return aliens, nil
}
