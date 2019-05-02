package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var trickyByteValues = [][]byte{
	[]byte(":"),
	[]byte(`\`),
	[]byte("::"),
	[]byte(`\\`),
	[]byte(`\:`),
	[]byte(`:\`),
	[]byte(`\\:`),
	[]byte(`::\`),
	[]byte(`\:\:`),
	[]byte(`:\:\`),
	[]byte(`:\\`),
	[]byte(`\::`),
	[]byte(`::\\`),
	[]byte(`\\::`),
}

func TestEscapeUnescape(t *testing.T) {
	for _, expected := range trickyByteValues {
		actual := unescape(escape(expected))
		assert.Equal(t, expected, actual)
	}
}

func TestFindWithValueWithEscape(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	ageIndex := col.AddIndex("age", func(m Model) []byte {
		// Note: We add the ':' to the index value to try and trip up the escaping
		// algorithm.
		return []byte(fmt.Sprintf(":%d:", m.(*testModel).Age))
	})
	models := make([]*testModel, len(trickyByteValues))
	// Use the trickyByteValues as the names for each model.
	for i, name := range trickyByteValues {
		models[i] = &testModel{
			Name: string(name),
			Age:  i,
		}
	}
	for i, expected := range models {
		require.NoError(t, col.Insert(expected), "testModel %d", i)
		actual := []*testModel{}
		err := col.FindWithValue(ageIndex, []byte(fmt.Sprintf(":%d:", expected.Age)), &actual)
		require.NoError(t, err, "testModel %d", i)
		require.Len(t, actual, 1, "testModel %d", i)
		assert.Equal(t, expected, actual[0])
	}
}