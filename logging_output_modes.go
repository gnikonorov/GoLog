/*
	Logging output modes
*/
package golog

import (
	"strconv"
)

type LoggingOutputMode int

const (
	ModeFile   LoggingOutputMode = iota + 1 // ModeFile indicates information will be outputted to a log file
	ModeScreen                              // ModeScreen indicates information will be outputted to the screen
	ModeBoth                                // ModeBoth indicates information will be outted to both file and screen
)

func (mode LoggingOutputMode) Int() int {
	return int(mode)
}

func (mode LoggingOutputMode) String() string {
	return strconv.Itoa(int(mode))
}
