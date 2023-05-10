package simulation

import (
	"bufio"
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volmedo/invasim/internal/aliens"
	"github.com/volmedo/invasim/internal/worldmap"
)

func Test_Run(t *testing.T) {
	world := worldmap.World{
		"Foo": worldmap.Roads{
			worldmap.Direction_North: "Bar",
			worldmap.Direction_West:  "Baz",
			worldmap.Direction_South: "Qu-ux",
		},
		"Bar": worldmap.Roads{
			worldmap.Direction_South: "Foo",
		},
		"Baz": worldmap.Roads{
			worldmap.Direction_East: "Foo",
		},
		"Qu-ux": worldmap.Roads{
			worldmap.Direction_North: "Foo",
		},
	}

	alienTracker := aliens.Tracker{
		"alien 0": "Foo",
		"alien 1": "Bar",
		"alien 2": "Baz",
		"alien 3": "Qu-ux",
	}

	maxIterations := 1
	out := &bytes.Buffer{}

	Run(world, alienTracker, maxIterations, out)

	assert.NotContains(t, world, "Foo")
	assert.NotContains(t, alienTracker, "alien 1")
	assert.NotContains(t, alienTracker, "alien 2")
	assert.NotContains(t, alienTracker, "alien 3")

	scanner := bufio.NewScanner(out)
	scanner.Scan()
	assert.Regexp(t, regexp.MustCompile(`Foo has been destroyed by alien \d, alien \d and alien \d!`), scanner.Text())
	scanner.Scan()
	assert.Equal(t, "Simulation finished!", scanner.Text())
	scanner.Scan()
	assert.Equal(t, "Max iterations reached, 1 alien(s) remaining", scanner.Text())
	scanner.Scan()
	assert.Equal(t, "This is what the world looks like after the invasion:", scanner.Text())
}
