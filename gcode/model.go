// Package gcode is providing parsing and marshalling functions to work with gcode data in text form.
package gcode

// File represents *.gcode content (gcode application)
type File struct {
	Lines []*Line
}

// Line represents one line (gcode block) from *.gcode file
type Line struct {
	DeleteFlg bool
	LineNbr   string
	Codes     []*GCode
	Comment   string
}

// GCode represents one command
type GCode struct {
	Letter  string
	Value   string
	Comment string
}
