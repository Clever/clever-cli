package clevertable

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/azylman/optimus.v1"
	"testing"
)

func TestFlattenRow(t *testing.T) {
	out, err := FlattenRow(optimus.Row{
		"key1": map[string]interface{}{
			"key2": map[string]interface{}{
				"key3": "val1",
				"key4": "val2",
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, out, optimus.Row{
		"key1.key2.key3": "val1",
		"key1.key2.key4": "val2",
	})
}

func TestStringifyArrayVals(t *testing.T) {
	out, err := StringifyArrayVals(optimus.Row{"key": []interface{}{"val1", "val2"}})
	assert.Nil(t, err)
	assert.Equal(t, out, optimus.Row{"key": `["val1","val2"]`})
}
