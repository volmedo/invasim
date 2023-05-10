package worldmap

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// World is the graph of the map expressed as set of adjacency lists.
type World map[string]Roads

// Roads is a map of roads keyed by their directions where the value are the destination cities.
type Roads map[Direction]string

// Direction expresses the direction of a given road.
type Direction string

const (
	Direction_East  Direction = "east"
	Direction_North Direction = "north"
	Direction_South Direction = "south"
	Direction_West  Direction = "west"
)

// opposite returns the opposite direction of the direction given, as seen from the destination.
func (d Direction) opposite() (Direction, error) {
	switch d {
	case Direction_East:
		return Direction_West, nil
	case Direction_North:
		return Direction_South, nil
	case Direction_South:
		return Direction_North, nil
	case Direction_West:
		return Direction_East, nil
	default:
		return "", fmt.Errorf("invalid direction %s", d)
	}
}

// ReadFromFile reads a map file.
// The format for such files consists on a series of lines, where each line contains the declaration of a city along
// with the cities that can be reached from it taking roads in different directions. Each of these lines has the format
// '<city_name> [<road> [<road>]...]', where <city_name> is a string. <road> is a pair '<direction>=<destination_city_name>'.
// <direction> can only be one of "east", "north", "south" and "west".
//
// This format can be expressed in EBNF notation as:
//
//	map file = city line , { city line } ;
//	city line = city name , {" " , road} ;
//	city name = ( alpha | digit ) , { alpha | digit } ;
//	road = direction , "=" , city name ;
//	direction = "east" | "north" | "south" | "west" ;
func ReadFromFile(path string) (World, error) {
	file, err := os.Open(path)
	if err != nil {
		return World{}, err
	}
	defer file.Close()

	world := World{}
	lineNum := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := parseLine(world, scanner.Text(), lineNum); err != nil {
			return World{}, err
		}

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return World{}, err
	}

	return world, nil
}

// parseLine parses a single line from a map file, and adds the declared city and its roads to the passed World object.
func parseLine(world World, line string, lineNum int) error {
	if world == nil {
		return errors.New("world map must not be nil")
	}

	if line == "" {
		return nil
	}

	parts := strings.Split(line, " ")

	cityName := parts[0]
	world[cityName] = Roads{}
	for _, road := range parts[1:] {
		roadParts := strings.Split(road, "=")

		if len(roadParts) != 2 {
			return fmt.Errorf("malformed directions at line %d: %s", lineNum, line)
		}

		dir := Direction(roadParts[0])
		switch dir {
		case Direction_East, Direction_North, Direction_South, Direction_West:
		default:
			return fmt.Errorf("bad direction in directions at line %d: %s", lineNum, dir)
		}

		dest := roadParts[1]

		// add the road to the origin city and also to the destination one, carefully checking for conflicts
		if d, alreadyExists := world[cityName][dir]; alreadyExists {
			if d != dest {
				return fmt.Errorf(
					"conflict in road declaration at line %d: a road from %s direction %s to %s is declared, but there is already a road in that direction to %s",
					lineNum, cityName, dir, dest, d,
				)
			}
		} else {
			world[cityName][dir] = dest
		}

		// add an entry for the destination city to the map if it doesn't exist yet
		if _, ok := world[dest]; !ok {
			world[dest] = Roads{}
		}

		oppDir, _ := dir.opposite()
		if d, alreadyExists := world[dest][oppDir]; alreadyExists {
			if d != cityName {
				return fmt.Errorf(
					"conflict in road declaration at line %d: a road from %s direction %s to %s is declared, but the destination already has a road in the opposite direction to %s",
					lineNum, cityName, dir, dest, d,
				)
			}
		} else {
			world[dest][oppDir] = cityName
		}
	}

	return nil
}

// DestroyCity removes the given city from the World, along with the roads to other cities, leaving a big hole behind.
// The function will also take care to remove the road to the destroyed city from destination cities.
func (w World) DestroyCity(city string) {
	roads := w[city]
	for _, dest := range roads {
		destRoads := w[dest]
		for dir, origin := range destRoads {
			if origin == city {
				delete(w[dest], dir)
				break
			}
		}
	}

	delete(w, city)
}

// String fulfills the Stringer interface. It produces a representation of the given World instance in valid map
// file format.
func (w World) String() string {
	builder := strings.Builder{}
	for c, roads := range w {
		builder.WriteString(c)
		for dir, dest := range roads {
			builder.WriteString(fmt.Sprintf(" %s=%s", dir, dest))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
