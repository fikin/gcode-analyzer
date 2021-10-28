package gcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertGCode(t *testing.T, letter, value, comment string, gCode *GCode) {
	assert.Equal(t, letter, gCode.Letter)
	assert.Equal(t, value, gCode.Value)
	assert.Equal(t, comment, gCode.Comment)
}

func TestPeekIsComment(t *testing.T) {
	assert.True(t, peekIsCodeComment("("))
	assert.False(t, peekIsCodeComment(" "))
}

func TestNextToken1(t *testing.T) {
	var token string
	var line = "M10 F11 (comment  F11)   Z+12  (comment Z)"
	token, line = nextToken(line)
	assert.Equal(t, "M10", token)
	token, line = nextToken(line)
	assert.Equal(t, "F11", token)
	token, line = nextToken(line)
	assert.Equal(t, "(comment  F11)", token)
	token, line = nextToken(line)
	assert.Equal(t, "Z+12", token)
	token, line = nextToken(line)
	assert.Equal(t, "(comment Z)", token)
	assert.Equal(t, "", line)
}

func TestParseGCode(t *testing.T) {
	var gCode *GCode
	var line = "M10 F11 (comment  F11)   Z+12  (comment Z)"
	gCode, line = parseGCode(line)
	assertGCode(t, "M", "10", "", gCode)
	gCode, line = parseGCode(line)
	assertGCode(t, "F", "11", "(comment  F11)", gCode)
	gCode, line = parseGCode(line)
	assertGCode(t, "Z", "+12", "(comment Z)", gCode)
	gCode, line = parseGCode(line)
	assert.Nil(t, gCode)
	assert.Equal(t, "", line)
}

func TestParseLineNumber(t *testing.T) {
	var token, line string
	token, line = parseLineNumber("N123 G10 ; Comment")
	assert.Equal(t, "N123", token)
	assert.Equal(t, "G10 ; Comment", line)
	token, line = parseLineNumber("123 G10 ; Comment")
	assert.Equal(t, "123", token)
	assert.Equal(t, "G10 ; Comment", line)
	token, line = parseLineNumber("G10 ; Comment")
	assert.Equal(t, "", token)
	assert.Equal(t, "G10 ; Comment", line)
}

func TestParseGCodeLine(t *testing.T) {
	line := parseGCodeLine("; this is comment")
	assert.Equal(t, "; this is comment", line.Comment)
	assert.Empty(t, line.Codes)
	assert.Empty(t, line.LineNbr)
	assert.False(t, line.DeleteFlg)

	line = parseGCodeLine("G10")
	assert.Empty(t, line.Comment)
	assert.Equal(t, 1, len(line.Codes))
	assert.Equal(t, GCode{Letter: "G", Value: "10"}, *line.Codes[0])
	assert.Empty(t, line.LineNbr)
	assert.False(t, line.DeleteFlg)

	line = parseGCodeLine("G10 (; this code comment) ; this is line comment")
	assert.Equal(t, "; this is line comment", line.Comment)
	assert.Equal(t, 1, len(line.Codes))
	assert.Equal(t, GCode{Letter: "G", Value: "10", Comment: "(; this code comment)"}, *line.Codes[0])
	assert.Empty(t, line.LineNbr)
	assert.False(t, line.DeleteFlg)

	line = parseGCodeLine("010")
	assert.Empty(t, line.Comment)
	assert.Empty(t, line.Codes)
	assert.Equal(t, "010", line.LineNbr)
	assert.False(t, line.DeleteFlg)

	line = parseGCodeLine("010")
	assert.Empty(t, line.Comment)
	assert.Empty(t, line.Codes)
	assert.Equal(t, "010", line.LineNbr)
	assert.False(t, line.DeleteFlg)

	line = parseGCodeLine("N010 G11 X1 Y-2 (cmnt Y) Z+3 ; (cmnt 2 )")
	assert.Equal(t, "; (cmnt 2 )", line.Comment)
	assert.Equal(t, 4, len(line.Codes))
	assert.Equal(t, GCode{Letter: "G", Value: "11"}, *line.Codes[0])
	assert.Equal(t, GCode{Letter: "X", Value: "1"}, *line.Codes[1])
	assert.Equal(t, GCode{Letter: "Y", Value: "-2", Comment: "(cmnt Y)"}, *line.Codes[2])
	assert.Equal(t, GCode{Letter: "Z", Value: "+3"}, *line.Codes[3])
	assert.Equal(t, "N010", line.LineNbr)
	assert.False(t, line.DeleteFlg)
}

func TestParseLine(t *testing.T) {
	line := parseLine("/ deleted G10")
	assert.True(t, line.DeleteFlg)
	assert.Equal(t, "/ deleted G10", line.Comment)
	assert.Empty(t, line.LineNbr)

	line = parseLine("; G10 line is the comment")
	assert.Equal(t, "; G10 line is the comment", line.Comment)
	assert.Empty(t, line.LineNbr)
	assert.Empty(t, line.Codes)
	assert.False(t, line.DeleteFlg)

	line = parseLine("11 G10 ; comment")
	assert.Equal(t, "; comment", line.Comment)
	assert.Equal(t, "11", line.LineNbr)
	assert.Equal(t, 1, len(line.Codes))
	assert.Equal(t, GCode{Letter: "G", Value: "10"}, *line.Codes[0])
	assert.False(t, line.DeleteFlg)
}

func TestParseLines(t *testing.T) {
	lines := parseLines(`
G10

/ deleted`)
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, 1, len(lines[0].Codes))
	assert.Equal(t, GCode{Letter: "G", Value: "10"}, *(lines[0].Codes[0]))
	assert.Empty(t, lines[0].Comment)
	assert.False(t, lines[0].DeleteFlg)
	assert.Equal(t, Line{DeleteFlg: true, Comment: "/ deleted"}, *lines[1])
}
