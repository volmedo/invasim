package alienmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volmedo/invasim/internal/worldmap"
)

func Test_New(t *testing.T) {
	testCases := map[string]struct {
		numAliens            int
		world                worldmap.World
		expectedAliensLength int
		expectsError         bool
	}{
		"happy path": {
			numAliens: 2,
			world: worldmap.World{
				"Foo": worldmap.Roads{
					worldmap.Direction_South: "Qu-ux",
				},
				"Qu-ux": worldmap.Roads{
					worldmap.Direction_North: "Foo",
				},
			},
			expectedAliensLength: 2,
			expectsError:         false,
		},
		"too many aliens": {
			numAliens: 3,
			world: worldmap.World{
				"Foo": worldmap.Roads{
					worldmap.Direction_South: "Qu-ux",
				},
				"Qu-ux": worldmap.Roads{
					worldmap.Direction_North: "Foo",
				},
			},
			expectedAliensLength: 0,
			expectsError:         true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			aliens, err := New(tc.numAliens, tc.world)

			if tc.expectsError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedAliensLength, len(aliens))
			}
		})
	}
}

func Test_MoveRandomly(t *testing.T) {
	aliens := Aliens{
		"alien 0": "Foo",
		"alien 1": "Bar",
		"alien 2": "Baz",
	}

	world := worldmap.World{
		"Foo": worldmap.Roads{
			worldmap.Direction_North: "Bar",
			worldmap.Direction_West:  "Baz",
		},
		"Bar": worldmap.Roads{
			worldmap.Direction_South: "Foo",
		},
		"Baz": worldmap.Roads{
			worldmap.Direction_East: "Foo",
		},
	}

	visitedCities := aliens.MoveRandomly(world)

	// alien 0 can go to Bar or Baz, while aliens 1 and 2 can only go to Foo
	assert.Condition(t, func() bool {
		alien0Dest := aliens["alien 0"]
		return alien0Dest == "Bar" || alien0Dest == "Baz"
	})
	assert.Equal(t, "Foo", aliens["alien 1"])
	assert.Equal(t, "Foo", aliens["alien 2"])

	_, okBar := visitedCities["Bar"]
	_, okBaz := visitedCities["Baz"]
	assert.Condition(t, func() bool {
		return okBar || okBaz
	})

	if okBar {
		assert.Equal(t, []string{"alien 0"}, visitedCities["Bar"])
	}
	if okBaz {
		assert.Equal(t, []string{"alien 0"}, visitedCities["Baz"])
	}

	assert.Contains(t, visitedCities, "Foo")
	assert.Equal(t, []string{"alien 1", "alien 2"}, visitedCities["Foo"])
}

func Test_pickRandomDestination(t *testing.T) {
	roads := worldmap.Roads{
		worldmap.Direction_East:  "Foo",
		worldmap.Direction_North: "Bar",
		worldmap.Direction_South: "Baz",
		worldmap.Direction_West:  "Qu-ux",
	}

	resultCounts := map[string]int{
		"Foo":   0,
		"Bar":   0,
		"Baz":   0,
		"Qu-ux": 0,
	}

	numIterations := 1000
	for i := 0; i < numIterations; i++ {
		dest := pickRandomDestination(roads)
		resultCounts[dest]++
	}

	t.Log(resultCounts)

	// assert that roads are selected uniformly (around 1000/4 or 250 times each)
	assert.Condition(t, func() bool {
		// allow 10% deviation
		even := numIterations / len(resultCounts)
		delta := even * 10 / 100
		lowThreshold := even - delta
		highThreshold := even + delta

		for _, count := range resultCounts {
			if count < lowThreshold || count > highThreshold {
				return false
			}
		}

		return true
	})
}

func Test_DestroyAliens(t *testing.T) {
	testCases := map[string]struct {
		aliens          Aliens
		aliensToDestroy []string
		expectedAliens  Aliens
	}{
		"happy path": {
			aliens: Aliens{
				"alien 0": "Foo",
				"alien 1": "Bar",
				"alien 2": "Baz",
				"alien 3": "Qu-ux",
			},
			aliensToDestroy: []string{"alien 1", "alien 3"},
			expectedAliens: Aliens{
				"alien 0": "Foo",
				"alien 2": "Baz",
			},
		},
		"destroying a non-existant alien is a no-op": {
			aliens: Aliens{
				"alien 0": "Foo",
				"alien 1": "Bar",
			},
			aliensToDestroy: []string{"alien 2", "alien 3"},
			expectedAliens: Aliens{
				"alien 0": "Foo",
				"alien 1": "Bar",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.aliens.DestroyAliens(tc.aliensToDestroy)

			assert.Equal(t, tc.expectedAliens, tc.aliens)
		})
	}
}
