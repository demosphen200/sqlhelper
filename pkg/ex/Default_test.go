package ex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
	defInt         int
	defString      string
	assignedInt    int
	assignedString string
}

func TestDefault(t *testing.T) {
	var test = testStruct{
		assignedInt:    25,
		assignedString: "assigned",
	}

	assert.Equal(t, 33, Default(test.defInt, 33))
	assert.Equal(t, 25, Default(test.assignedInt, 33))
	assert.Equal(t, "default", Default(test.defString, "default"))
	assert.Equal(t, "assigned", Default(test.assignedString, "default"))
}
