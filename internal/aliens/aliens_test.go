package aliens

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volmedo/invasim/internal/worldmap"
)

func Test_New(t *testing.T) {
	testCases := map[string]struct {
		numAliens         int
		world             worldmap.World
		expectedNumAliens int
		expectsError      bool
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
			expectedNumAliens: 2,
			expectsError:      false,
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
			expectedNumAliens: 0,
			expectsError:      true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tracker, err := NewTracker(tc.numAliens, tc.world)

			if tc.expectsError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedNumAliens, len(tracker))
			}
		})
	}
}

func Test_randomizeCities(t *testing.T) {
	world := worldmap.World{
		"Foo": worldmap.Roads{},
		"Bar": worldmap.Roads{},
		"Baz": worldmap.Roads{},
	}

	resultCounts := map[string]int{
		"FooBarBaz": 0,
		"FooBazBar": 0,
		"BarFooBaz": 0,
		"BarBazFoo": 0,
		"BazFooBar": 0,
		"BazBarFoo": 0,
	}

	// since the results from the function are random, we'll call it a given number of times and collect results.
	// We will then check those results for statistical randomness
	numIterations := 2000
	for i := 0; i < numIterations; i++ {
		randCitites := randomizeCities(world)
		resultCounts[strings.Join(randCitites, "")]++
	}

	// assert results are selected uniformly
	assert.Condition(t, func() bool {
		// allow 15% deviation
		even := numIterations / len(resultCounts)
		delta := even * 15 / 100
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

func Test_MoveRandomly(t *testing.T) {
	tracker := Tracker{
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

	visitedCities := tracker.MoveRandomly(world)

	// alien 0 can go to Bar or Baz, while aliens 1 and 2 can only go to Foo
	assert.Condition(t, func() bool {
		alien0Dest := tracker["alien 0"]
		return alien0Dest == "Bar" || alien0Dest == "Baz"
	})
	assert.Equal(t, "Foo", tracker["alien 1"])
	assert.Equal(t, "Foo", tracker["alien 2"])

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
	assert.Contains(t, visitedCities["Foo"], "alien 1")
	assert.Contains(t, visitedCities["Foo"], "alien 2")
	assert.NotContains(t, visitedCities["Foo"], "alien 0")
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

	// since the results from the function are random, we will call it a given number of times and collect results.
	// We will then check those results for statistical randomness
	numIterations := 2000
	for i := 0; i < numIterations; i++ {
		dest := pickRandomDestination(roads)
		resultCounts[dest]++
	}

	// assert that roads are selected uniformly (around 2000/4 or 500 times each)
	assert.Condition(t, func() bool {
		// allow 15% deviation
		even := numIterations / len(resultCounts)
		delta := even * 15 / 100
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
		tracker         Tracker
		aliensToDestroy []string
		expectedTracker Tracker
	}{
		"happy path": {
			tracker: Tracker{
				"alien 0": "Foo",
				"alien 1": "Bar",
				"alien 2": "Baz",
				"alien 3": "Qu-ux",
			},
			aliensToDestroy: []string{"alien 1", "alien 3"},
			expectedTracker: Tracker{
				"alien 0": "Foo",
				"alien 2": "Baz",
			},
		},
		"destroying a non-existant alien is a no-op": {
			tracker: Tracker{
				"alien 0": "Foo",
				"alien 1": "Bar",
			},
			aliensToDestroy: []string{"alien 2", "alien 3"},
			expectedTracker: Tracker{
				"alien 0": "Foo",
				"alien 1": "Bar",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.tracker.DestroyAliens(tc.aliensToDestroy)

			assert.Equal(t, tc.expectedTracker, tc.tracker)
		})
	}
}
