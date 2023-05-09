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

	for name, testData := range testCases {
		t.Run(name, func(t *testing.T) {
			aliens, err := New(testData.numAliens, testData.world)

			if testData.expectsError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, testData.expectedAliensLength, len(aliens))
			}
		})
	}
}
