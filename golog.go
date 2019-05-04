/*
	Log4* Like logger for the Go programming language

	Author: Gleb Nikonorov
*/

package golog

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	// Below are logging output modes

	// File indicates information will be outputted to a log file
	File = "FILE"

	// Screen indicates information will be outputted to the screen
	Screen = "SCREEN"

	// Both indicates information will be outted to both file and screen
	Both = "BOTH"

	// Below are init log file options

	// FileAppend instructs the logger to append onto an existing log if one exists
	FileAppend = "APPEND"

	// FileCompress instructs the logger to compress an existing log file if one exists
	FileCompress = "COMPRESS"

	// FileDelete instructs the logger to remove an existing log file if one exits
	FileDelete = "DELETE"

	// FileActionNone indicates no file actions ( e.g.: when user is writing to screen only
	FileActionNone = "NONE"

	outStreamStdErr = 10
	outStreamStdOut = 11
)

// NOTE: I may or may not be thread safe
var stringBuilder strings.Builder // Used to avoid costly string concatenation

// LoggingConfig holds a logging configuration for the logger and is used during logger initialization
type LoggingConfig struct {
	Name                 string // The logger profile name
	LogMode              string // The logging mode
	LogFileStartupAction string // The action the logger will take on startuip
	LogDirectory         string // The directory to which the logger writes
	LogFile              string // The name of the log file to write to
	ShouldColorize       bool   // Indicates if we should output information in color
}

// Logger is representative of the logger for use in other go programs
// Contains the following fields:
//	colorize: Boolean value indicating if logging out should be colorized
//	loggingMode: Must be either File, Screen, or Both
//	loggingDirectory: The directory to store log files in. Must be a valid directory if loggingMode is Screen or Both
//	loggingFile: The name of the log file to output to. Must be a valid file name of loggingMode is Screen or Both
//
// The following methods are exposed by this structure:
//	Debug(logText string): Log debug output to log destination
//	Info(logText string): Log info output to log destination
//	Warning(logText string): Log warning output to log destination
//	Err(logText string): Log error output to log destination
//	Fatal(logText string): Log fatal output to log destination
//	Is_Uninitialized: Returns true if this structure has not been allocated
type Logger struct {
	colorize         bool   // If true, print log output in color
	context          string // The context is the value prepended to each log line and set by the caller via 'SetContext'
	loggingDirectory string // The directory to store logs in
	loggingFile      string // The file to store logs in
	loggingMode      string // The mode of the logger (Should be FILE, SCREEN, or BOTH)
}

// Compress file validatedFilePath
func compressFile(filePath string) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		stringBuilder.Reset()

		stringBuilder.WriteString("Could not stat file '")
		stringBuilder.WriteString(filePath)
		stringBuilder.WriteString("' because: ")
		stringBuilder.WriteString(err.Error())

		panic(stringBuilder.String())
	}
	defer fileHandle.Close()

	fileInfo, err := fileHandle.Stat()
	if err != nil {
		stringBuilder.Reset()

		stringBuilder.WriteString("Could not get file info for '")
		stringBuilder.WriteString(filePath)
		stringBuilder.WriteString("' because: ")
		stringBuilder.WriteString(err.Error())

		panic(stringBuilder.String())
	}
	var fileSize = fileInfo.Size()

	// get file contents
	fileBytes := make([]byte, fileSize)
	fileReader := bufio.NewReader(fileHandle)
	_, err = fileReader.Read(fileBytes)
	if err != nil {
		stringBuilder.Reset()

		stringBuilder.WriteString("Could not get file bytes of '")
		stringBuilder.WriteString(filePath)
		stringBuilder.WriteString("' because: ")
		stringBuilder.WriteString(err.Error())

		panic(stringBuilder.String())
	}

	// write out .gz file, appending file last modified time to make a unique, identifiable name
	var byteBuffer bytes.Buffer
	zipWriter := gzip.NewWriter(&byteBuffer)
	zipWriter.Write(fileBytes)
	zipWriter.Close()

	var fileModTime = fileInfo.ModTime()
	var fileModTimeString = fileModTime.Format("20060102150405")

	stringBuilder.Reset()

	stringBuilder.WriteString(filePath)
	stringBuilder.WriteString(".")
	stringBuilder.WriteString(fileModTimeString)
	stringBuilder.WriteString(".gz")

	var gzipFileName = stringBuilder.String()
	err = ioutil.WriteFile(gzipFileName, byteBuffer.Bytes(), fileInfo.Mode())
	if err != nil {
		stringBuilder.Reset()

		stringBuilder.WriteString("Could not create zip of '")
		stringBuilder.WriteString(filePath)
		stringBuilder.WriteString("' because: ")
		stringBuilder.WriteString(err.Error())

		panic(stringBuilder.String())
	}

	// delete source file
	err = os.Remove(filePath)
	if err != nil {
		stringBuilder.Reset()

		stringBuilder.WriteString("Could not delete log file '")
		stringBuilder.WriteString(filePath)
		stringBuilder.WriteString("' because: ")
		stringBuilder.WriteString(err.Error())

		panic(stringBuilder.String())
	}

}

// Check to make sure that file fullPathToLogFile exists and return true/false
func doesLoggingFileExist(fullPathToLogFile string) bool {
	// check to see if a log file already exists. If it does, delete it
	fileInfo, err := os.Stat(fullPathToLogFile)
	if err != nil {
		// if the file does not exist there's nothing to do
		// if the error is anything else panic
		if !os.IsNotExist(err) {
			stringBuilder.Reset()

			stringBuilder.WriteString("Could not stat log file '")
			stringBuilder.WriteString(fullPathToLogFile)
			stringBuilder.WriteString("' because: ")
			stringBuilder.WriteString(err.Error())

			panic(stringBuilder.String())
		} else {
			return false
		}
	}

	if fileInfo.IsDir() {
		// this is unlikely, but possible
		panic("The log file you specified was a directory! Terminating.")
	}

	return true
}

// writeLog writes a formatted log line to the user specified outputs. If 'shouldPanic' is true,
// it will also raise a panic with the user provided log text
func (logger *Logger) writeLog(paintColor LoggingColor, resetColor LoggingColor, loggingLevel LoggingLevel, logText string, outputStream int, shouldPanic bool) {
	var logTime          = time.Now().String()
	var loggingLevelText = loggingLevel.String()
	var paintString      = paintColor.String()
	var resetString      = resetColor.String()

	if logger.loggingMode == Screen || logger.loggingMode == Both {
		logStrings := []string{paintString, "[", logTime, "] ", loggingLevelText, ": ", logger.context, logText, resetString, "\n"}
		var logString = strings.Join(logStrings, "")
		if outputStream == outStreamStdErr {
			os.Stderr.WriteString(logString)
		} else {
			fmt.Printf(logString)
		}
	}

	if logger.loggingMode == Both || logger.loggingMode == File {
		stringBuilder.Reset()

		stringBuilder.WriteString(logger.loggingDirectory)
		stringBuilder.WriteString("/")
		stringBuilder.WriteString(logger.loggingFile)

		var fileName = stringBuilder.String()

		// append to the log file, creating if one does not exist. In case of any error, panic
		logHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
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

// func validateLoggerConfig validate a loggers configuration as valid. If a configuration is invalid,
// an error is returned. Else, nil is returned
// TODO: Validate that passed log file action is a valid file action
func validateLoggerConfig(logMode string, logDirectory string, logFile string, logFileStartupAction string) error {
	if logMode != File && logMode != Screen && logMode != Both {
		stringBuilder.Reset()

		stringBuilder.WriteString("Log mode must either be 'File', 'Screen', or 'Both'. Not '")
		stringBuilder.WriteString(logMode)
		stringBuilder.WriteString("'.")

		return errors.New(stringBuilder.String())
	}

	if logMode == File || logMode == Both {
		// We're logging to a file, make sure that the directory given to us was valid.
		// Also try and create it for the user if it does not exist. MkdirAll is a no-op if the directory already
		// exists so no harm no foul
		mkdirErr := os.MkdirAll(logDirectory, os.ModePerm)
		if mkdirErr != nil {
			stringBuilder.Reset()

			stringBuilder.WriteString("Log directory '")
			stringBuilder.WriteString(logDirectory)
			stringBuilder.WriteString("' did not exist. Creation of directory failed because: '")
			stringBuilder.WriteString(mkdirErr.Error())
			stringBuilder.WriteString("'")

			return errors.New(stringBuilder.String())
		}

		// Errors here are extremely unlikely, but not harm in checking for them
		fileInfo, statErr := os.Stat(logDirectory)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				// The directory does not exist
				stringBuilder.Reset()

				stringBuilder.WriteString("Log directory '")
				stringBuilder.WriteString(logDirectory)
				stringBuilder.WriteString("' is not a valid log directory. Error.")

				return errors.New(stringBuilder.String())
			} else {
				stringBuilder.Reset()

				stringBuilder.WriteString("Could not stat log directory '")
				stringBuilder.WriteString(logDirectory)
				stringBuilder.WriteString("' because: '")
				stringBuilder.WriteString(statErr.Error())
				stringBuilder.WriteString("'")

				return errors.New(stringBuilder.String())
			}
		}

		// Check to make sure we actually gave a directory
		if !fileInfo.IsDir() {
			stringBuilder.Reset()

			stringBuilder.WriteString("'")
			stringBuilder.WriteString(logDirectory)
			stringBuilder.WriteString("' is a file not a directory.")

			return errors.New(stringBuilder.String())
		}

		// depending on the logFileStartupAction value perform the appropriate action on any existing log file
		// note that file append is default behavior
		stringBuilder.Reset()

		stringBuilder.WriteString(logDirectory)
		stringBuilder.WriteString("/")
		stringBuilder.WriteString(logFile)

		var fullPathToLogFile = stringBuilder.String()
		var fileExists = doesLoggingFileExist(fullPathToLogFile)
		if fileExists {
			if logFileStartupAction == FileCompress {
				// compress the file
				compressFile(fullPathToLogFile)
			} else if logFileStartupAction == FileDelete {
				// delete the file
				err := os.Remove(fullPathToLogFile)
				if err != nil {
					stringBuilder.Reset()

					stringBuilder.WriteString("Terminating. Could not delete log file '")
					stringBuilder.WriteString(fullPathToLogFile)
					stringBuilder.WriteString("' because: ")
					stringBuilder.WriteString(err.Error())

					return errors.New(stringBuilder.String())
				}
			}
		}
	}

	return nil
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
func (logger *Logger) IsUninitialized() bool {
	return logger.loggingMode == ""
}

// SetContext is called on the logger to the set its context. See 'Context' in the logging struct for more
// information
func (logger *Logger) SetContext(context string) {
	logger.context = context
}

// SetupLoggerFromConfigFile sets up and returns a logger instance as specified in fullFilePath for profile
func SetupLoggerFromConfigFile(fullFilePath string, profile string) (Logger, error) {
	var returnError error
	var logger      Logger

	// get bytes of file
	fileBytes, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		// don't tolerate any error while reading file
		stringBuilder.Reset()

		stringBuilder.WriteString("Could not read file '")
		stringBuilder.WriteString(fullFilePath)
		stringBuilder.WriteString("' because: ")
		stringBuilder.WriteString(err.Error())

		returnError = errors.New(stringBuilder.String())
		return logger, returnError
	}

	// parse out our json
	loggingConfigs := make([]LoggingConfig, 0)
	err = json.Unmarshal(fileBytes, &loggingConfigs)
	if err != nil {
		stringBuilder.Reset()

		stringBuilder.WriteString("Failed to decode config file because: '")
		stringBuilder.WriteString(err.Error())
		stringBuilder.WriteString("'. Ensure it is in a valid JSON format.")

		returnError = errors.New(stringBuilder.String())
		return logger, returnError
	}

	for _, config := range loggingConfigs {
		// TODO: Change to string builder
		fmt.Printf("The object is %+v\n", config)
		fmt.Printf("Configname is %s\n", config.Name)
		if config.Name == profile {
			returnError = validateLoggerConfig(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction)
			if returnError != nil {
				return logger, returnError
			}

			logger = Logger{loggingMode: config.LogMode, loggingDirectory: config.LogDirectory, loggingFile: config.LogFile, colorize: config.ShouldColorize}
			return logger, nil
		}
	}

	// if we get here we couldn't find any config for the profile
	stringBuilder.Reset()

	stringBuilder.WriteString("Configure profile '")
	stringBuilder.WriteString(profile)
	stringBuilder.WriteString("' not found in config file '")
	stringBuilder.WriteString(fullFilePath)
	stringBuilder.WriteString("'.")

	returnError = errors.New(stringBuilder.String())
	return logger, returnError
}

// SetupLoggerFromStruct sets up and returns a logger instance from a LoggingConfigStruct
func SetupLoggerFromStruct(config *LoggingConfig) (Logger, error) {
	var logger Logger

	returnError := validateLoggerConfig(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction)
	if returnError != nil {
		return logger, returnError
	}

	logger = Logger{loggingMode: config.LogMode, loggingDirectory: config.LogDirectory, loggingFile: config.LogFile, colorize: config.ShouldColorize}
	return logger, nil
}

// SetupLoggerFromFields Sets up and returns a logger instance from passed in individual fields
func SetupLoggerFromFields(logMode string, logFileStartupAction string, logDirectory string, logFile string, shouldColorize bool) (Logger, error) {
	var logger Logger

	returnError := validateLoggerConfig(logMode, logDirectory, logFile, logFileStartupAction)
	if returnError != nil {
		return logger, returnError
	}

	logger = Logger{loggingMode: logMode, loggingDirectory: logDirectory, loggingFile: logFile, colorize: shouldColorize}
	return logger, nil
}
