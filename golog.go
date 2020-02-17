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
func (logger *Logger) writeLog(loggingMessage logMessage) {
	if logger.loggingMode == ModeScreen || logger.loggingMode == ModeBoth {
		logStrings := []string{loggingMessage.paintColor, "[", loggingMessage.logTime, "] ", loggingMessage.loggingLevel, ": ", logger.context, loggingMessage.logText, loggingMessage.resetColor, "\n"}
		var logString = strings.Join(logStrings, "")
		if loggingMessage.outputStream == outStreamStdErr {
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

		stringBuilder.WriteString(loggingMessage.paintColor)
		stringBuilder.WriteString("[")
		stringBuilder.WriteString(loggingMessage.logTime)
		stringBuilder.WriteString("] ")
		stringBuilder.WriteString(loggingMessage.loggingLevel)
		stringBuilder.WriteString(": ")
		stringBuilder.WriteString(logger.context)
		stringBuilder.WriteString(loggingMessage.logText)
		stringBuilder.WriteString(loggingMessage.resetColor)
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

	if loggingMessage.shouldPanic {
		panic(loggingMessage.logText)
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

	loggingMessage := logMessage{time.Now().String(), levelDebug.String(), paintColor.String(), resetColor.String(), logText, outStreamStdOut, false}
	logger.writeLog(loggingMessage)
}

// Info Outputs info log information to the logging destination
func (logger *Logger) Info(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone

	loggingMessage := logMessage{time.Now().String(), levelInfo.String(), paintColor.String(), resetColor.String(), logText, outStreamStdOut, false}
	logger.writeLog(loggingMessage)
}

// Warning Outputs warning information to the logging destination
func (logger *Logger) Warning(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorWarn
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelWarn.String(), paintColor.String(), resetColor.String(), logText, outStreamStdOut, false}
	logger.writeLog(loggingMessage)
}

// Err Outputs error information to the logging destination
func (logger *Logger) Err(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorErr
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelErr.String(), paintColor.String(), resetColor.String(), logText, outStreamStdErr, false}
	logger.writeLog(loggingMessage)
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

	loggingMessage := logMessage{time.Now().String(), levelFatal.String(), paintColor.String(), resetColor.String(), logText, outStreamStdErr, false}
	logger.writeLog(loggingMessage)
}

// Panic Outputs fatal information to the logging desination and causes a panic
func (logger *Logger) Panic(logText string) {
	var paintColor = colorNone
	var resetColor = colorNone
	if logger.colorize {
		paintColor = colorPanic
		resetColor = colorReset
	}

	loggingMessage := logMessage{time.Now().String(), levelPanic.String(), paintColor.String(), resetColor.String(), logText, outStreamStdErr, true}
	logger.writeLog(loggingMessage)
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
