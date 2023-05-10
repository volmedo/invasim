package simulation

import (
	"fmt"
	"io"
	"strings"

	"github.com/volmedo/invasim/internal/aliens"
	"github.com/volmedo/invasim/internal/worldmap"
)

// Run runs a new simulation with the given parameters.
//
// The simulation is implemented as a loop. In each iteration, aliens move randomly to any of the cities that are
// reachable from the city they are currently in, one city at a time. When aliens end up in the same city, they
// unleash their futuristic weapons and destroy each other, along with the city itself and any roads leading into or
// out of it.
// The simulation ends when there are no more aliens alive or maxIterations iterations have been executed, whatever
// happens first.
// The function accepts an io.Writer where city destruction messages will be printed to make testing for correct output
// easier.
func Run(world worldmap.World, alienTracker aliens.Tracker, maxIterations int, out io.Writer) {
	for i := 0; i < maxIterations && len(alienTracker) > 0; i++ {
		// move aliens
		// at this point no city should have more than 1 alien (it would've already been destroyed otherwise)
		visitedCities := alienTracker.MoveRandomly(world)

		// check if aliens are in the same place using the visited cities view
		for city, aliens := range visitedCities {
			if len(aliens) > 1 {
				world.DestroyCity(city)
				alienTracker.DestroyAliens(aliens)

				fmt.Fprintf(
					out,
					"%s has been destroyed by %s and %s!\n",
					city, strings.Join(aliens[:len(aliens)-1], ", "), aliens[len(aliens)-1],
				)
			}
		}
	}

	// check final conditions: either all aliens were destroyed or we reached maxIterations
	fmt.Fprintf(out, "Simulation finished!\n")
	if len(alienTracker) == 0 {
		fmt.Fprintf(out, "All aliens were destroyed!\n")
	} else {
		fmt.Fprintf(out, "Max iterations reached, %d alien(s) remaining\n", len(alienTracker))
	}

	fmt.Fprintln(out, "This is what the world looks like after the invasion:")
	fmt.Fprintln(out, world)
}
