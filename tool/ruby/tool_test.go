package ruby

import (
	"testing"

	"github.com/broothie/ok/parameter"
	"github.com/stretchr/testify/assert"
)

func Test_parseType(t *testing.T) {
	assert.Equal(t, parameter.TypeString, parseType(`"hi"`))
	assert.Equal(t, parameter.TypeString, parseType(`'hi'`))

	assert.Equal(t, parameter.TypeBool, parseType("false"))
	assert.Equal(t, parameter.TypeBool, parseType("true"))

	assert.Equal(t, parameter.TypeInt, parseType("0"))
	assert.Equal(t, parameter.TypeInt, parseType("1"))
	assert.Equal(t, parameter.TypeInt, parseType("92"))

	assert.Equal(t, parameter.TypeFloat, parseType(".1"))
	assert.Equal(t, parameter.TypeFloat, parseType("2.1"))
	assert.Equal(t, parameter.TypeFloat, parseType("2."))
}
