/*
	Self contained representation of a log message
*/

package golog

// Log message is a self contained representation of a golog log message
type logMessage struct {
	logTime       string // The time the intent to log occurred
	loggingLevel  string // The string representation of the coresponding LoggingLevel
	paintColor    string // The color to print the log as a string, corresponding to a LoggingColor
	resetColor    string // The color to reset the output streams rendering to, corresponding to a LoggingColor
	logText       string // The text to log
	outputStream  int    // The output stream to write to
	shouldPanic   bool   // If true, raise a panic while logging
	logger       *Logger // The logger that will be used to write the message
}
