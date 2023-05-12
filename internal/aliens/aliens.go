package aliens

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/volmedo/invasim/internal/worldmap"
)

// Tracker keeps track of the city each alien is currently at
type Tracker map[string]string

// NewTracker creates a new alien Tracker with numAliens aliens placed randomly in one of the cities of world.
// Since there can only be an alien in a city, numAliens cannot be greater than the number of cities in world.
func NewTracker(numAliens int, world worldmap.World) (Tracker, error) {
	if numAliens > len(world) {
		return Tracker{}, fmt.Errorf("not enough cities (%d) to place %d aliens", len(world), numAliens)
	}

	tracker := Tracker{}

	randomCities := randomizeCities(world)
	randomCities = randomCities[:numAliens]

	for _, city := range randomCities {
		name := randomAlienName()
		tracker[name] = city
	}

	return tracker, nil
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

var vowels = []string{"a", "e", "i", "o", "u"}
var consonants = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "y", "z"}
var alphabet = append(consonants, vowels...)

// randomAlienName creates a random alien name between 4 and 8 characters long, with a 30% chance of having a hyphen
// for extra alienness. Thanks ChatGPT.
func randomAlienName() string {
	length := rand.Intn(4) + 4
	name := []string{strings.ToUpper(vowels[rand.Intn(len(vowels))])}
	for i := 0; i < length-2; i++ {
		name = append(name, alphabet[rand.Intn(len(alphabet))])
	}
	name = append(name, vowels[rand.Intn(len(vowels))])
	nameStr := ""
	for _, c := range name {
		nameStr += c
	}
	if rand.Float64() < 0.3 {
		index := rand.Intn(length-2) + 1
		nameStr = nameStr[:index] + "-" + nameStr[index:]
	}
	return nameStr
}

// VisitedCities offers the opposite view than what Aliens provides. It maps each city being visited to a list of
// aliens currently placed in that location. Cities with no alien presence will not appear in this map. It is used
// as an auxiliary data structure to enable quick checking of cities and aliens that should be destroyed during fights.
type VisitedCities map[string][]string

// MoveRandomly moves all the aliens in the Tracker randomly through one the roads available from the city each of them
// is currently at. Once an available road is chosen, the alien's position is updated to the destination.
// As the function moves aliens around, it also collects visited cities to make checking which cities have more than
// one alien more convenient.
func (t Tracker) MoveRandomly(world worldmap.World) VisitedCities {
	visited := VisitedCities{}
	for a, currCity := range t {
		roads := world[currCity]
		if len(roads) == 0 {
			// TODO: consider the possibility of removing the alien from the tracker, as it won't be able to move any further
			continue
		}

		destCity := pickRandomDestination(roads)

		t[a] = destCity

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

// DestroyAliens removes the passed aliens from the Tracker as they were horribly destroyed by their enemies.
func (t Tracker) DestroyAliens(aliens []string) {
	for _, a := range aliens {
		delete(t, a)
	}
}
