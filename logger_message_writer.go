/*
	File holding functions related to writing log messages to output channels
*/

package golog

import (
	"fmt"
	"os"
	"strings"
)

// writeLog writes a formatted log line to the user specified outputs. If 'shouldPanic' is true,
// it will also raise a panic with the user provided log text
func writeLog(loggingMessage logMessage) {
	if loggingMessage.logger.loggingMode == ModeScreen || loggingMessage.logger.loggingMode == ModeBoth {
		logStrings := []string{loggingMessage.paintColor, "[", loggingMessage.logTime, "] ", loggingMessage.loggingLevel, ": ", loggingMessage.logger.context, loggingMessage.logText, loggingMessage.resetColor, "\n"}
		var logString = strings.Join(logStrings, "")
		if loggingMessage.outputStream == outStreamStdErr {
			os.Stderr.WriteString(logString)
		} else {
			fmt.Printf(logString)
		}
	}

	var stringBuilder strings.Builder

	if loggingMessage.logger.loggingMode == ModeBoth || loggingMessage.logger.loggingMode == ModeFile {
		stringBuilder.Reset()

		stringBuilder.WriteString(loggingMessage.logger.loggingDirectory)
		stringBuilder.WriteString("/")
		stringBuilder.WriteString(loggingMessage.logger.loggingFile)

		var fileName = stringBuilder.String()

		// append to the log file, creating if one does not exist. In case of any error, panic
		logHandle, err := loggingMessage.logger.osHandle.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
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

		stringBuilder.WriteString("[")
		stringBuilder.WriteString(loggingMessage.logTime)
		stringBuilder.WriteString("] ")
		stringBuilder.WriteString(loggingMessage.loggingLevel)
		stringBuilder.WriteString(": ")
		stringBuilder.WriteString(loggingMessage.logger.context)
		stringBuilder.WriteString(loggingMessage.logText)
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
