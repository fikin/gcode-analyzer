package gcode

import (
	"io/ioutil"
	"strings"
	"unicode"
)

// ParseGCodeFile is parsing given file into data structure
func ParseGCodeFile(filename string) (*File, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	l := parseLines(string(body))
	return &File{Lines: l}, nil
}

func parseLines(file string) []*Line {
	ret := []*Line{}
	for _, line := range strings.Split(file, "\n") {
		if line == "" {
			continue
		}
		line := parseLine(line)
		ret = append(ret, line)
	}
	return ret
}

func parseLine(line string) *Line {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	switch line[0] {
	case ';':
		return &Line{Comment: line}
	case '/':
		return &Line{DeleteFlg: true, Comment: line}
	default:
		return parseGCodeLine(line)
	}
}

func parseGCodeLine(line string) *Line {
	lineNbr, line := parseLineNumber(line)
	codes, line := parseLineCodes(line)
	return &Line{LineNbr: lineNbr, Comment: line, Codes: codes}
}

func parseLineNumber(line string) (lineNbr, remainder string) {
	if line[0] == 'N' || unicode.IsNumber(rune(line[0])) {
		return nextToken(line)
	}
	return "", line
}

func parseLineCodes(line string) (parsed []*GCode, remainder string) {
	parsed = []*GCode{}
	var gCode *GCode
	for true {
		gCode, line = parseGCode(line)
		if gCode == nil {
			break
		}
		parsed = append(parsed, gCode)
	}
	return parsed, line
}

func parseGCode(line string) (gCode *GCode, remainder string) {
	if line == "" || line[0] == ';' {
		return nil, line
	}
	token, line := nextToken(line)
	ret := &GCode{Letter: string(token[0]), Value: token[1:]}
	if peekIsCodeComment(line) {
		ret.Comment, line = nextToken(line)
	}
	return ret, line
}

func nextToken(line string) (token, remainder string) {
	if line == "" {
		return "", ""
	}
	breakSym := " "
	augment := 0
	if peekIsCodeComment(line) {
		breakSym = ")"
		augment = 1
	}
	i := strings.Index(line, breakSym) + augment
	if i > 0 {
		return line[:i], strings.TrimSpace(line[i:])
	}
	return line, ""
}

func peekIsCodeComment(line string) bool {
	if len(line) == 0 {
		return false
	}
	return line[0] == '('
}
