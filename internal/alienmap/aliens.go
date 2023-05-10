package alienmap

import (
	"fmt"
	"math/rand"

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

	randomCities := randomizeCities(world)
	randomCities = randomCities[:numAliens]

	// to place the aliens randomly we can take advantage of the inherently non-deterministic ordering of elements
	// observed when ranging over a map in Go
	alien := 0
	for _, city := range randomCities {
		name := fmt.Sprintf("alien %d", alien)
		aliens[name] = city
		alien++
	}

	return aliens, nil
}

// randomizeCities returns a slice with the names of the cities in world in a random order.
func randomizeCities(world worldmap.World) []string {
	randomCities := make([]string, 0, len(world))
	for city := range world {
		randomCities = append(randomCities, city)
	}

	rand.Shuffle(len(randomCities), func(i, j int) {
		randomCities[i], randomCities[j] = randomCities[j], randomCities[i]
	})

	return randomCities
}

// VisitedCities offers the opposite view than what Aliens provides. It maps each city being visited to a list of
// aliens currently placed in that location. Cities with no alien presence will not appear in this map. It is used
// as an auxiliary data structure to enable quick checking of cities and aliens that should be destroyed during fights.
type VisitedCities map[string][]string

// MoveRandomly moves all the aliens in the Aliens tracker randomly through one the roads available from the city
// each of them is currently at. Once an available road is chosen, the alien's position is updated to the destination.
// As the function moves aliens around, it also collects visited cities to make checking which cities have more than
// one alien more convenient.
func (as Aliens) MoveRandomly(world worldmap.World) VisitedCities {
	visited := VisitedCities{}
	for a, currCity := range as {
		roads := world[currCity]
		if len(roads) == 0 {
			// TODO: consider the possibility of removing the alien from the tracker, as it won't be able to move any further
			continue
		}

		destCity := pickRandomDestination(roads)

		as[a] = destCity

		visited[destCity] = append(visited[destCity], a)
	}

	return visited
}

// pickRandomDestination picks a random road from the set of roads being passed and return the city it leads to.
// It does so by choosing a random index and enumerating the available roads until the chosen index is found.
func pickRandomDestination(roads worldmap.Roads) string {
	randIdx := rand.Intn(len(roads))

	var destCity string
	i := 0
	for _, destCity = range roads {
		if i == randIdx {
			break
		}

		i++
	}

	return destCity
}

// DestroyAliens removes the passed aliens from the tracker as they were horribly destroyed by their enemies.
func (as Aliens) DestroyAliens(aliens []string) {
	for _, a := range aliens {
		delete(as, a)
	}
}
