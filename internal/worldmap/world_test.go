package worldmap

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadFromFile(t *testing.T) {
	testCases := map[string]struct {
		mapFileContents string
		expectedWorld   World
		expectsError    bool
	}{
		"happy path": {
			mapFileContents: "Foo north=Bar west=Baz south=Qu-ux\nBar south=Foo west=Bee",
			expectedWorld: World{
				"Foo": Roads{
					Direction_North: "Bar",
					Direction_West:  "Baz",
					Direction_South: "Qu-ux",
				},
				"Bar": Roads{
					Direction_South: "Foo",
					Direction_West:  "Bee",
				},
				"Baz": Roads{
					Direction_East: "Foo",
				},
				"Qu-ux": Roads{
					Direction_North: "Foo",
				},
				"Bee": Roads{
					Direction_East: "Bar",
				},
			},
			expectsError: false,
		},
		"malformed line": {
			mapFileContents: "Foo Bar west=Baz south=Qu-ux",
			expectedWorld:   nil,
			expectsError:    true,
		},
		"unsupported direction": {
			mapFileContents: "Foo southeast=Bar west=Baz south=Qu-ux",
			expectedWorld:   nil,
			expectsError:    true,
		},
		"conflicting road declaration": {
			mapFileContents: "Foo north=Bar west=Baz south=Qu-ux\nBar south=Qu-ux west=Bee",
			expectedWorld:   nil,
			expectsError:    true,
		},
		"non-conflicting road re-declarations work": {
			mapFileContents: "Foo north=Bar west=Baz north=Bar",
			expectedWorld: World{
				"Foo": Roads{
					Direction_North: "Bar",
					Direction_West:  "Baz",
				},
				"Bar": Roads{
					Direction_South: "Foo",
				},
				"Baz": Roads{
					Direction_East: "Foo",
				},
			},
			expectsError: false,
		},
	}

	tmpDir := t.TempDir()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mapFilePath, err := writeTestFile(tmpDir, name, tc.mapFileContents)
			if err != nil {
				t.Fatal("Error writing test file")
			}

			w, err := ReadFromFile(mapFilePath)
			if tc.expectsError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedWorld, w)
			}
		})
	}
}

// writeTestFile writes a temporary file with the given contents and returns its path.
func writeTestFile(tmpDir, testName, contents string) (string, error) {
	f, err := os.CreateTemp(tmpDir, testName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	f.Write([]byte(contents))
	return f.Name(), nil
}

func Test_DestroyCity(t *testing.T) {
	testCases := map[string]struct {
		world         World
		cityToDestroy string
		expectedWorld World
	}{
		"happy path": {
			world: World{
				"Foo": Roads{
					Direction_North: "Bar",
					Direction_West:  "Baz",
				},
				"Bar": Roads{
					Direction_South: "Foo",
				},
				"Baz": Roads{
					Direction_East: "Foo",
				},
			},
			cityToDestroy: "Foo",
			expectedWorld: World{
				"Bar": Roads{},
				"Baz": Roads{},
			},
		},
		"passing a non-existent city is a no-op": {
			world: World{
				"Foo": Roads{
					Direction_North: "Bar",
					Direction_West:  "Baz",
				},
				"Bar": Roads{
					Direction_South: "Foo",
				},
				"Baz": Roads{
					Direction_East: "Foo",
				},
			},
			cityToDestroy: "Qu-ux",
			expectedWorld: World{
				"Foo": Roads{
					Direction_North: "Bar",
					Direction_West:  "Baz",
				},
				"Bar": Roads{
					Direction_South: "Foo",
				},
				"Baz": Roads{
					Direction_East: "Foo",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.world.DestroyCity(tc.cityToDestroy)

			assert.Equal(t, tc.expectedWorld, tc.world)
		})
	}
}
