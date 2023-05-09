package simulation

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volmedo/invasim/internal/alienmap"
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

	aliens := alienmap.Aliens{
		"alien 0": "Foo",
		"alien 1": "Bar",
		"alien 2": "Baz",
		"alien 3": "Qu-ux",
	}

	maxIterations := 1
	out := &bytes.Buffer{}

	Run(world, aliens, maxIterations, out)

	assert.NotContains(t, world, "Foo")
	assert.NotContains(t, aliens, "alien 1")
	assert.NotContains(t, aliens, "alien 2")
	assert.NotContains(t, aliens, "alien 3")

	scanner := bufio.NewScanner(out)
	scanner.Scan()
	assert.Equal(t, "Foo has been destroyed by alien 1, alien 2 and alien 3!", scanner.Text())
}
