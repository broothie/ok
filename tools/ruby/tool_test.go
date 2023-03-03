package ruby

import (
	"testing"

	"github.com/broothie/ok/task"
	"github.com/stretchr/testify/assert"
)

func Test_parseType(t *testing.T) {
	assert.Equal(t, task.TypeString, parseType(`"hi"`))
	assert.Equal(t, task.TypeString, parseType(`'hi'`))

	assert.Equal(t, task.TypeBool, parseType("false"))
	assert.Equal(t, task.TypeBool, parseType("true"))

	assert.Equal(t, task.TypeInt, parseType("0"))
	assert.Equal(t, task.TypeInt, parseType("1"))
	assert.Equal(t, task.TypeInt, parseType("92"))

	assert.Equal(t, task.TypeFloat, parseType(".1"))
	assert.Equal(t, task.TypeFloat, parseType("2.1"))
	assert.Equal(t, task.TypeFloat, parseType("2."))
}
