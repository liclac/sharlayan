package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONNode(t *testing.T) {
	assert.Equal(t,
		fsDumpNode{Filename: "test.txt", Mode: 0666, Data: "{\n  \"key\": \"value\"\n}\n"},
		renderDump(t, JSON(Named("test.txt"), map[string]string{"key": "value"})))
}
