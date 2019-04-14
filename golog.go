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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	// Logging modes for the logger
	modeDebug = "DEBUG"   // This mode should be used debug information
	modeErr   = "ERROR"   // A recoverable error
	modeFatal = "FATAL"   // Non-recoverable error.
	modeInfo  = "INFO"    // Non severe log information. Should be used for things like user input
	modePanic = "PANIC"	  // Akin to an exception. Logs and throws a panic
	modeWarn  = "WARNING" // Indicator of potential problems

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

	// These are constants for terminal colors
	// Info is left in the native terminal color
	colorDebug = "\x1B[32m"      // This is green
	colorErr   = "\x1B[31m"      // This is red
	colorFatal = "\x1B[0;37;41m" // This is white text on a red background
	colorPanic = "\x1B[0;37;41m" // This is white text on a red background
	colorReset = "\x1B[0m"       // resets an applied terminal color
	colorWarn  = "\x1B[33m"      // This is yellow

	outStreamStdErr = 10
	outStreamStdOut = 11
)

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
	loggingDirectory string // The directory to store logs in
	loggingFile      string // The file to store logs in
	loggingMode      string // The mode of the logger (Should be FILE, SCREEN, or BOTH)
}

// Compress file validatedFilePath
func compressFile(filePath string) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		panic("Could not stat file " + filePath + " because " + err.Error())
	}
	defer fileHandle.Close()

	fileInfo, err := fileHandle.Stat()
	if err != nil {
		panic("Could not get file info for " + filePath + " because " + err.Error())
	}
	var fileSize = fileInfo.Size()

	// get file contents
	fileBytes := make([]byte, fileSize)
	fileReader := bufio.NewReader(fileHandle)
	_, err = fileReader.Read(fileBytes)
	if err != nil {
		panic("Could not get file bytes of " + filePath + " because " + err.Error())
	}

	// write out .gz file, appending file last modified time to make a unique, identifiable name
	var byteBuffer bytes.Buffer
	zipWriter := gzip.NewWriter(&byteBuffer)
	zipWriter.Write(fileBytes)
	zipWriter.Close()

	var fileModTime = fileInfo.ModTime()
	var fileModTimeString = fileModTime.Format("20060102150405")

	err = ioutil.WriteFile(filePath+"."+fileModTimeString+".gz", byteBuffer.Bytes(), fileInfo.Mode())
	if err != nil {
		panic("Could not create zip of " + filePath + " because " + err.Error())
	}

	// delete source file
	err = os.Remove(filePath)
	if err != nil {
		panic("Could not delete the log file because: " + err.Error() + ". Terminating.")
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
			panic("Could not stat the log file because: " + err.Error() + ". Terminating.")
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
func writeLog(colorString string, resetString string, loggingMode string, logText string, logger *Logger, outputStream int, shouldPanic bool) {
	var logTime = time.Now().String()

	if logger.loggingMode == Screen || logger.loggingMode == Both {
		logStrings := []string{colorString, "[", logTime, "] ", loggingMode, ": ", logText, resetString, "\n"}
		var logString = strings.Join(logStrings, "")
		if outputStream == outStreamStdErr {
			os.Stderr.WriteString(logString)
		} else {
			fmt.Printf(logString)
		}
	}

	if logger.loggingMode == Both || logger.loggingMode == File {
		var fileName = logger.loggingDirectory + "/" + logger.loggingFile

		// append to the log file, creating if one does not exist. In case of any error, panic
		logHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		defer logHandle.Close()
		if err != nil {
			// can't open file
			panic("Unable to open log file " + fileName + " for writing because " + err.Error())
		}

		var writeBytes = []byte(colorString + "[" + logTime + "] " + loggingMode + ": " + logText + resetString + "\n")
		_, err = logHandle.Write(writeBytes)
		if err != nil {
			panic("Unable to write to log file " + fileName + " because " + err.Error())
		}
	}

	if shouldPanic {
		panic(logText)
	}
}

// validate a loggers configuration as valid
func validateLoggerConfig(logMode string, logDirectory string, logFile string, logFileStartupAction string) {
	if logMode != File && logMode != Screen && logMode != Both {
		panic("Log mode must either be File, Screen, or Both. Goodbye")
	}

	if logMode == File || logMode == Both {
		// We're logging to a file, make sure that the directory given to us was valid
		fileInfo, err := os.Stat(logDirectory)
		if err != nil {
			if os.IsNotExist(err) {
				// The directory does not exist
				panic("Please provide a valid log directory. Goodbye.")
			} else {
				panic("Could not stat the log directory because: " + err.Error() + ". Terminating.")
			}
		}

		// Check to make sure we actually gave a directory
		if !fileInfo.IsDir() {
			panic("You must give a directory! Not a file!")
		}

		// depending on the logFileStartupAction value perform the appropriate action on any existing log file
		// note that file append is default behavior
		var fullPathToLogFile = logDirectory + "/" + logFile
		var fileExists = doesLoggingFileExist(fullPathToLogFile)
		if fileExists {
			if logFileStartupAction == FileCompress {
				// compress the file
				compressFile(fullPathToLogFile)
			} else if logFileStartupAction == FileDelete {
				// delete the file
				err := os.Remove(fullPathToLogFile)
				if err != nil {
					panic("Could not delete the log file because: " + err.Error() + ". Terminating.")
				}
			}
		}
	}
}

// Debug Outputs debug log information to the logging destination
func (logger *Logger) Debug(logText string) {
	var colorString = ""
	var resetString = ""
	if logger.colorize {
		colorString = colorDebug
		resetString = colorReset
	}

	writeLog(colorString, resetString, modeDebug, logText, logger, outStreamStdOut, false)
}

// Info Outputs info log information to the logging destination
func (logger *Logger) Info(logText string) {
	var colorString = ""
	var resetString = ""

	writeLog(colorString, resetString, modeInfo, logText, logger, outStreamStdOut, false)
}

// Warning Outputs warning information to the logging destination
func (logger *Logger) Warning(logText string) {
	var colorString = ""
	var resetString = ""
	if logger.colorize {
		colorString = colorWarn
		resetString = colorReset
	}

	writeLog(colorString, resetString, modeWarn, logText, logger, outStreamStdErr, false)
}

// Err Outputs error information to the logging destination
func (logger *Logger) Err(logText string) {
	var colorString = ""
	var resetString = ""
	if logger.colorize {
		colorString = colorErr
		resetString = colorReset
	}

	writeLog(colorString, resetString, modeErr, logText, logger, outStreamStdErr, false)
}

// Fatal Outputs fatal information to the logging desination but does not cause a panic,
// use 'Panic' instead.
func (logger *Logger) Fatal(logText string) {
	var colorString = ""
	var resetString = ""
	if logger.colorize {
		colorString = colorFatal
		resetString = colorReset
	}

	writeLog(colorString, resetString, modeFatal, logText, logger, outStreamStdErr, false)
}

// Panic Outputs fatal information to the logging desination and causes a panic
func (logger *Logger) Panic(logText string) {
	var colorString = ""
	var resetString = ""
	if logger.colorize {
		colorString = colorPanic
		resetString = colorReset
	}

	writeLog(colorString, resetString, modePanic, logText, logger, outStreamStdErr, true)
}

// IsUninitialized Returns true if this structure has not yet been allocated
func (logger *Logger) IsUninitialized() bool {
	return logger.loggingMode == ""
}

// SetupLoggerFromConfigFile sets up and returns a logger instance as specified in fullFilePath for profile
func SetupLoggerFromConfigFile(fullFilePath string, profile string) Logger {
	// get bytes of file
	fileBytes, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		// don't tolerate any error while reading file
		panic("Could not read file because: " + err.Error() + "!")
	}

	// parse out our json
	loggingConfigs := make([]LoggingConfig, 0)
	err = json.Unmarshal(fileBytes, &loggingConfigs)
	if err != nil {
		panic("Failed to decode config file because: " + err.Error() + ". Ensure it is in JSON format.")
	}

	for _, config := range loggingConfigs {
		fmt.Printf("The object is %+v\n", config)
		fmt.Printf("Configname is " + config.Name + "\n")
		if config.Name == profile {
			validateLoggerConfig(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction)

			logger := Logger{loggingMode: config.LogMode, loggingDirectory: config.LogDirectory, loggingFile: config.LogFile, colorize: config.ShouldColorize}
			return logger
		}
	}

	// if we get here we couldn't find any config for the profile
	panic("Configuration profile " + profile + " not found in config file " + fullFilePath + "!")
}

// SetupLoggerFromStruct sets up and returns a logger instance from a LoggingConfigStruct
func SetupLoggerFromStruct(config *LoggingConfig) Logger {
	validateLoggerConfig(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction)

	logger := Logger{loggingMode: config.LogMode, loggingDirectory: config.LogDirectory, loggingFile: config.LogFile, colorize: config.ShouldColorize}
	return logger
}

// SetupLoggerFromFields Sets up and returns a logger instance from passed in individual fields
func SetupLoggerFromFields(logMode string, logFileStartupAction string, logDirectory string, logFile string, shouldColorize bool) Logger {
	validateLoggerConfig(logMode, logDirectory, logFile, logFileStartupAction)

	logger := Logger{loggingMode: logMode, loggingDirectory: logDirectory, loggingFile: logFile, colorize: shouldColorize}
	return logger
}
