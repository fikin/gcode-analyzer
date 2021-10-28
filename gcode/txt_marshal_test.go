package gcode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalCode(t *testing.T) {
	buf := bytes.NewBufferString("")
	assert.NoError(t, txtMarshalCode(&GCode{Letter: "G", Value: "10", Comment: "; comment"}, buf))
	assert.Equal(t, "G10 (; comment)", buf.String())
}

func TestMarshalDeletedLine(t *testing.T) {
	buf := bytes.NewBufferString("")
	assert.NoError(t, txtMarshalLine(&Line{
		DeleteFlg: true,
		Comment:   "/ this is deleted",
	}, buf))
	assert.Equal(t, "/ this is deleted\n", buf.String())
}

func TestMarshalCommentLine(t *testing.T) {
	buf := bytes.NewBufferString("")
	assert.NoError(t, txtMarshalLine(&Line{
		Comment: "; this is comment",
	}, buf))
	assert.Equal(t, "; this is comment\n", buf.String())
}

func TestMarshalGCodesLine(t *testing.T) {
	buf := bytes.NewBufferString("")
	assert.NoError(t, txtMarshalLine(&Line{
		Comment: "; this is line comment",
		Codes: []*GCode{
			{Letter: "G", Value: "10"},
			{Letter: "X", Value: "+2", Comment: "x comment"},
		},
	}, buf))
	assert.Equal(t, "G10 X+2 (x comment) ; this is line comment\n", buf.String())
}
