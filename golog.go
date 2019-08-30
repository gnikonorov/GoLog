/*
	Log4* Like logger for the Go programming language

	Author: Gleb Nikonorov
*/

package golog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	outStreamStdErr = 10
	outStreamStdOut = 11
)

// NOTE: I may or may not be thread safe
var stringBuilder strings.Builder // Used to avoid costly string concatenation

// writeLog writes a formatted log line to the user specified outputs. If 'shouldPanic' is true,
// it will also raise a panic with the user provided log text
func (logger *Logger) writeLog(paintColor LoggingColor, resetColor LoggingColor, loggingLevel LoggingLevel, logText string, outputStream int, shouldPanic bool) {
	var logTime          = time.Now().String()
	var loggingLevelText = loggingLevel.String()
	var paintString      = paintColor.String()
	var resetString      = resetColor.String()

	if logger.loggingMode == ModeScreen || logger.loggingMode == ModeBoth {
		logStrings := []string{paintString, "[", logTime, "] ", loggingLevelText, ": ", logger.context, logText, resetString, "\n"}
		var logString = strings.Join(logStrings, "")
		if outputStream == outStreamStdErr {
			os.Stderr.WriteString(logString)
		} else {
			fmt.Printf(logString)
		}
	}

	if logger.loggingMode == ModeBoth || logger.loggingMode == ModeFile {
		stringBuilder.Reset()

		stringBuilder.WriteString(logger.loggingDirectory)
		stringBuilder.WriteString("/")
		stringBuilder.WriteString(logger.loggingFile)

		var fileName = stringBuilder.String()

		// append to the log file, creating if one does not exist. In case of any error, panic
		logHandle, err := logger.OsHandle.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		defer logHandle.Close()
		if err != nil {
			// can't open file
			stringBuilder.Reset()

			stringBuilder.WriteString("Unable to open log file '")
			stringBuilder.WriteString(fileName)
			stringBuilder.WriteString("' for writing because: ")
			stringBuilder.WriteString(err.Error())

			panic(stringBuilder.String())
		}

		stringBuilder.Reset()

		stringBuilder.WriteString(paintString)
		stringBuilder.WriteString("[")
		stringBuilder.WriteString(logTime)
		stringBuilder.WriteString("] ")
		stringBuilder.WriteString(loggingLevelText)
		stringBuilder.WriteString(": ")
		stringBuilder.WriteString(logger.context)
		stringBuilder.WriteString(logText)
		stringBuilder.WriteString(resetString)
		stringBuilder.WriteString("\n")

		var writeBytes = []byte(stringBuilder.String())
		_, err = logHandle.Write(writeBytes)
		if err != nil {
			stringBuilder.Reset()

			stringBuilder.WriteString("Unable to write to log file '")
			stringBuilder.WriteString(fileName)
			stringBuilder.WriteString("' because: ")
			stringBuilder.WriteString(err.Error())

			panic(stringBuilder.String())
		}
	}

	if shouldPanic {
		panic(logText)
	}
}

// Debug Outputs debug log information to the logging destination
func (logger *Logger) Debug(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorDebug
		resetColor = colorReset
	}

	logger.writeLog(paintColor, resetColor, levelDebug, logText, outStreamStdOut, false)
}

// Info Outputs info log information to the logging destination
func (logger *Logger) Info(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone

	logger.writeLog(paintColor, resetColor, levelInfo, logText, outStreamStdOut, false)
}

// Warning Outputs warning information to the logging destination
func (logger *Logger) Warning(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorWarn
		resetColor = colorReset
	}

	logger.writeLog(paintColor, resetColor, levelWarn, logText, outStreamStdErr, false)
}

// Err Outputs error information to the logging destination
func (logger *Logger) Err(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorErr
		resetColor = colorReset
	}

	logger.writeLog(paintColor, resetColor, levelErr, logText, outStreamStdErr, false)
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

	logger.writeLog(paintColor, resetColor, levelFatal, logText, outStreamStdErr, false)
}

// Panic Outputs fatal information to the logging desination and causes a panic
func (logger *Logger) Panic(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorPanic
		resetColor = colorReset
	}

	logger.writeLog(paintColor, resetColor, levelPanic, logText, outStreamStdErr, true)
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
