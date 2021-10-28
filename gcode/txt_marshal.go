package gcode

import "io"

// MarshallAsText is marshaling the content as gcode file
func MarshallAsText(data *File, w io.StringWriter) error {
	for _, line := range data.Lines {
		if err := txtMarshalLine(line, w); err != nil {
			return err
		}
	}
	return nil
}

func txtMarshalLine(line *Line, w io.StringWriter) (err error) {
	if line.DeleteFlg {
		return writeStrings(w, line.Comment, "\n")
	}
	if line.LineNbr != "" {
		if err = writeStrings(w, line.LineNbr, " "); err != nil {
			return
		}
	}
	for i, code := range line.Codes {
		if i > 0 {
			if _, err = w.WriteString(" "); err != nil {
				return
			}
		}
		if err = txtMarshalCode(code, w); err != nil {
			return
		}
	}
	if line.Comment != "" {
		if len(line.Codes) > 0 {
			if _, err = w.WriteString(" "); err != nil {
				return
			}
		}
		if _, err = w.WriteString(line.Comment); err != nil {
			return
		}
	}
	_, err = w.WriteString("\n")
	return
}

func txtMarshalCode(code *GCode, w io.StringWriter) (err error) {
	if err = writeStrings(w, code.Letter, code.Value); err != nil {
		return
	}
	if code.Comment != "" {
		err = writeStrings(w, " (", code.Comment, ")")
	}
	return
}

func writeStrings(w io.StringWriter, arr ...string) error {
	for _, str := range arr {
		if _, err := w.WriteString(str); err != nil {
			return err
		}
	}
	return nil
}
