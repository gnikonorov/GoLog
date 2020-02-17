/*
	File holding functions and structures related to logger setup
*/

package golog

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/afero"
)

// Logger is representative of the logger for use in other go programs
//
// The following methods are exposed by this structure ( defined in golog.go ):
//	Debug(logText string): Log debug output to log destination
//	Info(logText string): Log info output to log destination
//	Warning(logText string): Log warning output to log destination
//	Err(logText string): Log error output to log destination
//	Fatal(logText string): Log fatal output to log destination
//	Is_Uninitialized: Returns true if this structure has not been allocated
type Logger struct {
	colorize         bool              // If true, print log output in color
	context          string            // The context is the value prepended to each log line and set by the caller via 'SetContext'
	loggingDirectory string            // The directory to store logs in
	loggingFile      string            // The file to store logs in
	loggingMode      LoggingOutputMode // The mode of the logger ( see 'logging_output_modes.go' )
	osHandle         afero.Fs          // We are using afero to enable mocking and stubbing the native FS during tests.
	isAsynch         bool              // If true, Asynchly handle log requests
	queueMgr         queueManager      // The asynch message handler, populated only if 'isAsynch' is true
}

// LoggingConfig holds a logging configuration for the logger and is used during logger initialization
type LoggingConfig struct {
	Name                 string            // The logger profile name
	LogMode              LoggingOutputMode // The logging mode
	LogFileStartupAction LoggingFileAction // The action the logger will take on startup
	LogDirectory         string            // The directory to which the logger writes
	LogFile              string            // The name of the log file to write to
	ShouldColorize       bool              // Indicates if we should output information in color
	IsMock               bool              // If true, mock the filesystem via 'afero'
	IsAsynch             bool              // If true, Asynchly handle log requests
}

// func compressFile compresses the file pointed to by 'filePath'
func compressFile(filePath string, osPtr afero.Fs) {
	var stringBuilder strings.Builder

	fileHandle, err := osPtr.Open(filePath)
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
	err = osPtr.Remove(filePath)
	if err != nil {
		stringBuilder.Reset()

		stringBuilder.WriteString("Could not delete log file '")
		stringBuilder.WriteString(filePath)
		stringBuilder.WriteString("' because: ")
		stringBuilder.WriteString(err.Error())

		panic(stringBuilder.String())
	}

}

// func doesLoggingFileExist checks to make sure that file 'fullPathToLogFile' exists and returns assertion of its existance
func doesLoggingFileExist(fullPathToLogFile string, osPtr afero.Fs) bool {
	var stringBuilder strings.Builder

	// check to see if a log file already exists. If it does, delete it
	fileInfo, err := osPtr.Stat(fullPathToLogFile)
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

// func getOSPtr returns a pointer to the os file system we are using. Choices are native filesystem or an in memory map
// based on the value of 'isMock'
func getOSPtr(isMock bool) afero.Fs {
	osPtr := afero.NewOsFs()
	if isMock {
		osPtr = afero.NewMemMapFs()
	}

	return osPtr
}

// func handleOldLogFile performs any necessary setup work on existing log files, if we are logging to a file based off the logging
// output mode
func handleOldLogFile(logMode LoggingOutputMode, logDirectory string, logFile string, logFileStartupAction LoggingFileAction, osPtr afero.Fs) error {
	var stringBuilder strings.Builder

	if logMode == ModeFile || logMode == ModeBoth {
		// depending on the logFileStartupAction value perform the appropriate action on any existing log file
		// note that file append is default behavior
		stringBuilder.Reset()

		stringBuilder.WriteString(logDirectory)
		stringBuilder.WriteString("/")
		stringBuilder.WriteString(logFile)

		var fullPathToLogFile = stringBuilder.String()
		var fileExists = doesLoggingFileExist(fullPathToLogFile, osPtr)
		if fileExists {
			if logFileStartupAction == FileActionCompress {
				// compress the file
				compressFile(fullPathToLogFile, osPtr)
			} else if logFileStartupAction == FileActionDelete {
				// delete the file
				err := osPtr.Remove(fullPathToLogFile)
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

// func validateLogDirectory validates the provided log directory. If it does not exist, it is created.
// an error is returned if the log directory is invalid. Else, nil is returned
func validateLogDirectory(logDirectory string, logMode LoggingOutputMode, osPtr afero.Fs) error {
	var stringBuilder strings.Builder

	// the log directory is not used if we're not logging to a file
	if !(logMode == ModeFile || logMode == ModeBoth) {
		return nil
	}

	mkdirErr := osPtr.MkdirAll(logDirectory, os.ModePerm)
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
	fileInfo, statErr := osPtr.Stat(logDirectory)
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

	return nil
}

// func validateLoggerConfig validate a loggers configuration as valid. If a configuration is invalid,
// an error is returned. Else, nil is returned
func validateLoggerConfig(logMode LoggingOutputMode, logDirectory string, logFile string, logFileStartupAction LoggingFileAction, osPtr afero.Fs) error {
	if !logMode.IsValidMode() {
		return errors.New("Invalid log mode provided. See log modes in 'logging_output_modes.go'")
	}

	if !logFileStartupAction.IsValidFileAction() {
		return errors.New("Invalid log file startup action provided. See actions in 'logging_file_actions.go'")
	}

	err := validateLogDirectory(logDirectory, logMode, osPtr)
	if err != nil {
		return err
	}

	return nil
}

// func SetupLoggerFromConfigFile sets up and returns a logger instance as specified in 'fullFilePath' for 'profile'
func SetupLoggerFromConfigFile(fullFilePath string, profile string) (Logger, error) {
	var returnError   error
	var logger        Logger
	var stringBuilder strings.Builder

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
	// NOTE: This should all be retested after refactors and type renames
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
		if config.Name == profile {
			osPtr := getOSPtr(config.IsMock)

			returnError = validateLoggerConfig(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction, osPtr)
			if returnError != nil {
				return logger, returnError
			}

			returnError = handleOldLogFile(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction, osPtr)
			if returnError != nil {
				return logger, returnError
			}

			var queueMgr queueManager
			if config.IsAsynch {
				queueMgr = createQueueMgr()
				queueMgr.start()
			}

			logger = Logger{loggingMode: config.LogMode, loggingDirectory: config.LogDirectory, loggingFile: config.LogFile, colorize: config.ShouldColorize, osHandle: osPtr, isAsynch: config.IsAsynch, queueMgr: queueMgr}
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

// func SetupLoggerFromFields sets up and returns a logger instance from passed in individual fields
func SetupLoggerFromFields(logMode LoggingOutputMode, logFileStartupAction LoggingFileAction, logDirectory string, logFile string, shouldColorize bool, isMock bool, isAsynch bool) (Logger, error) {
	var logger Logger

	osPtr := getOSPtr(isMock)

	returnError := validateLoggerConfig(logMode, logDirectory, logFile, logFileStartupAction, osPtr)
	if returnError != nil {
		return logger, returnError
	}

	returnError = handleOldLogFile(logMode, logDirectory, logFile, logFileStartupAction, osPtr)
	if returnError != nil {
		return logger, returnError
	}

	var queueMgr queueManager
	if isAsynch {
		queueMgr = createQueueMgr()
		queueMgr.start()
	}

	logger = Logger{loggingMode: logMode, loggingDirectory: logDirectory, loggingFile: logFile, colorize: shouldColorize, osHandle: osPtr, isAsynch: isAsynch, queueMgr: queueMgr}
	return logger, nil
}

// func SetupLoggerFromStruct sets up and returns a logger instance from a LoggingConfigStruct
func SetupLoggerFromStruct(config *LoggingConfig) (Logger, error) {
	var logger Logger

	osPtr := getOSPtr(config.IsMock)

	returnError := validateLoggerConfig(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction, osPtr)
	if returnError != nil {
		return logger, returnError
	}

	returnError = handleOldLogFile(config.LogMode, config.LogDirectory, config.LogFile, config.LogFileStartupAction, osPtr)
	if returnError != nil {
		return logger, returnError
	}

	var queueMgr queueManager
	if config.IsAsynch {
		queueMgr = createQueueMgr()
		queueMgr.start()
	}

	logger = Logger{loggingMode: config.LogMode, loggingDirectory: config.LogDirectory, loggingFile: config.LogFile, colorize: config.ShouldColorize, osHandle: osPtr, isAsynch: config.IsAsynch, queueMgr: queueMgr}
	return logger, nil
}
