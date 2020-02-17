/*
	Log4* Like logger for the Go programming language

	Author: Gleb Nikonorov
*/

package golog

import (
	"time"
)

const (
	outStreamStdErr = 10
	outStreamStdOut = 11
)

// Debug Outputs debug log information to the logging destination
func (logger *Logger) Debug(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorDebug
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelDebug.String(), paintColor.String(), resetColor.String(), logText, outStreamStdOut, false, logger}
	writeLog(loggingMessage)
}

// Info Outputs info log information to the logging destination
func (logger *Logger) Info(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone

	loggingMessage := logMessage{time.Now().String(), levelInfo.String(), paintColor.String(), resetColor.String(), logText, outStreamStdOut, false, logger}
	writeLog(loggingMessage)
}

// Warning Outputs warning information to the logging destination
func (logger *Logger) Warning(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorWarn
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelWarn.String(), paintColor.String(), resetColor.String(), logText, outStreamStdOut, false, logger}
	writeLog(loggingMessage)
}

// Err Outputs error information to the logging destination
func (logger *Logger) Err(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorErr
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelErr.String(), paintColor.String(), resetColor.String(), logText, outStreamStdErr, false, logger}
	writeLog(loggingMessage)
}

// Fatal Outputs fatal information to the logging desination but does not cause a panic,
// use 'Panic' instead.
func (logger *Logger) Fatal(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorFatal
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelFatal.String(), paintColor.String(), resetColor.String(), logText, outStreamStdErr, false, logger}
	writeLog(loggingMessage)
}

// Panic Outputs fatal information to the logging desination and causes a panic
func (logger *Logger) Panic(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorPanic
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelPanic.String(), paintColor.String(), resetColor.String(), logText, outStreamStdErr, true, logger}
	writeLog(loggingMessage)
}

// IsUninitialized Returns true if this structure has not yet been allocated
// since logging mode is private to golog, package users can never set 'logging mode' without
// using a logger setup method
func (logger *Logger) IsUninitialized() bool {
	return logger.loggingMode == 0
}

// SetContext is called on the logger to the set its context. See 'Context' in the logging struct for more
// information
func (logger *Logger) SetContext(context string) {
	logger.context = context
}
