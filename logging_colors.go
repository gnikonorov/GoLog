/*
	Terminal colors used by the logger
*/
package golog

// These are constants for terminal colors
type LoggingColor string

const (
	// Info is left in the native terminal color
	colorNone  LoggingColor = ""              // No color
	colorDebug LoggingColor = "\x1B[32m"      // This is green
	colorErr   LoggingColor = "\x1B[31m"      // This is red
	colorFatal LoggingColor = "\x1B[0;37;41m" // This is white text on a red background
	colorPanic LoggingColor = "\x1B[0;37;41m" // This is white text on a red background
	colorReset LoggingColor = "\x1B[0m"       // resets an applied terminal color
	colorWarn  LoggingColor = "\x1B[33m"      // This is yellow
)

func (color LoggingColor) String() string {
	return string(color)
}
