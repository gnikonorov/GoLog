/*

	Log4Go Like logger for the go programming language

	Structure initialization:
		Appender_Mode: Either DEBUG, SCREEN, or BOTH
		Appender_File: 

	Author: Gleb Nikonorov
*/

package golog;

import "io/ioutil"
import "fmt"
import "os"
import "time"

const FATAL = "FATAL"			// Non-recoverable error.
const ERROR = "ERROR"			// A recoverable error
const WARNING = "WARNING"		// Indicator of potential problems
const INFO = "INFO"			// Non severe log information. Should be used for things like user input
const DEBUG = "DEBUG"			// This mode should be used debug information 

// Logging modes for the logger
const FILE = 1			// Outputs information to a log file
const SCREEN = 2		// Outputs information to the screen
const BOTH = 3			// Outputs information to both file and screen

// Struct containing the logger object used by the Go Logger
// Contains the following fields:
//	logging_mode: Must be either FILE, SCREEN, or BOTH
//	logging_directory: The directory to store log files in. Must be a valid directory if logging_mode is SCREEN or BOTH
//	logging_file: The name of the log file to output to. Must be a valid file name of logging_mode is SCREEN or BOTH
//
// The following methods are exposed by this structure:
//	Debug(log_text string): Log debug output to log destination
//	Info(log_text string): Log info output to log destination
//	Warning(log_text string): Log warning output to log destination
//	Err(log_text string): Log error output to log destination
//	Fatal(log_text string): Log fatal output to log destination
//	Setup_logger_with_filename(log_mode int, log_directory string, log_file string)
//	Is_Uninitialized: Returns true if this structure has not been allocated
type Logger struct {
	logging_mode int	 // The mode of the logger (Should be FILE, SCREEN, or BOTH)
	logging_directory string // The directory to store logs in
	logging_file string	 // The file to store logs in
}

// func Debug Outputs DEBUG log information to the logging destination
func (logger *Logger) Debug(log_text string) {
	if(logger.logging_mode == SCREEN || logger.logging_mode == BOTH) {
		fmt.Printf("[%s] %s: %s\n", time.Now().String(), DEBUG, log_text)
	} 

	if(logger.logging_mode == BOTH || logger.logging_mode == FILE) {
		var file_name = logger.logging_directory + "/" + logger.logging_file;
		var write_bytes = []byte("[" + time.Now().String() + "] " + DEBUG + ":" + log_text)
		ioutil.WriteFile(file_name, write_bytes, 0644)
	}
}

// func Info Outputs INFO log information to the logging destination
func (logger *Logger) Info(log_text string) {
	if(logger.logging_mode == SCREEN || logger.logging_mode == BOTH) {
		fmt.Printf("[%s] %s: %s\n", time.Now().String(), INFO, log_text)
	} 

	if(logger.logging_mode == BOTH || logger.logging_mode == FILE) {
		var file_name = logger.logging_directory + "/" + logger.logging_file;
		var write_bytes = []byte("[" + time.Now().String() + "] " + DEBUG + ":" + log_text)
		ioutil.WriteFile(file_name, write_bytes, 0644)
	}
}

// func Warning Outputs WARNING information to the logging destination
func (logger *Logger) Warning(log_text string) {
	if(logger.logging_mode == SCREEN || logger.logging_mode == BOTH) {
		fmt.Printf("[%s] %s: %s\n", time.Now().String(), WARNING, log_text)
	} 

	if(logger.logging_mode == BOTH || logger.logging_mode == FILE) {
		var file_name = logger.logging_directory + "/" + logger.logging_file;
		var write_bytes = []byte("[" + time.Now().String() + "] " + DEBUG + ":" + log_text)
		ioutil.WriteFile(file_name, write_bytes, 0644)
	}
}

// func Err Outputs ERROR information to the logging destination
func (logger *Logger) Err(log_text string) {
	if(logger.logging_mode == SCREEN || logger.logging_mode == BOTH) {
		fmt.Printf("[%s] %s: %s\n", time.Now().String(), ERROR, log_text)
	} 

	if(logger.logging_mode == BOTH || logger.logging_mode == FILE) {
		var file_name = logger.logging_directory + "/" + logger.logging_file;
		var write_bytes = []byte("[" + time.Now().String() + "] " + DEBUG + ":" + log_text)
		ioutil.WriteFile(file_name, write_bytes, 0644)
	}
}

// func Fatal Outputs FATAL information to the logging desination
func (logger *Logger) Fatal(log_text string) {
	if(logger.logging_mode == SCREEN || logger.logging_mode == BOTH) {
		fmt.Printf("[%s] %s: %s\n", time.Now().String(), FATAL, log_text)
	} 

	if(logger.logging_mode == BOTH || logger.logging_mode == FILE) {
		var file_name = logger.logging_directory + "/" + logger.logging_file;
		var write_bytes = []byte("[" + time.Now().String() + "] " + DEBUG + ":" + log_text)
		ioutil.WriteFile(file_name, write_bytes, 0644)
	}
}

// func Is_Uninitialized: Returns true if this structure has not yet been allocated
func (logger *Logger) Is_Uninitialized() bool {
	return logger.logging_mode == 0
}

// func Setup_logger_with_filename Sets up and returns a logger instance.
func Setup_logger_with_filename (log_mode int, log_directory string, log_file string) Logger{
	// Validate our parmaters
	if(log_mode != FILE && log_mode != SCREEN && log_mode != BOTH) {
		panic("Log mode must either be FILE, SCREEN, or BOTH. Goodbye")
	}

	if(log_mode == FILE || log_mode == BOTH) {
		// We're logging to a file, make sure that the directory given to us was valid
		fileInfo, err := os.Stat(log_directory)
		if(err != nil) {
			if (os.IsNotExist(err)) {
				// The file does not exist
				panic("Please provide a valid log directory. Goodbye.")
			}
		}

		// Check to make sure we actually gave a directory
		if(!fileInfo.IsDir()) {
			panic("You must give a directory! Not a file!")
		}
	}

	logger := Logger{logging_mode:log_mode, logging_directory:log_directory, logging_file:log_file}
	return logger
}
