/*
	Logging level constants
*/
package golog

// Logging levels for the logger
type LoggingLevel string

const (
	levelDebug LoggingLevel = "DEBUG"   // This level should be used debug information
	levelErr   LoggingLevel = "ERROR"   // A recoverable error
	levelFatal LoggingLevel = "FATAL"   // Non-recoverable error.
	levelInfo  LoggingLevel = "INFO"    // Non severe log information. Should be used for things like user input
	levelPanic LoggingLevel = "PANIC"   // Akin to an exception. Logs and throws a panic
	levelWarn  LoggingLevel = "WARNING" // Indicator of potential problems
)

func (level LoggingLevel) String() string {
	return string(level)
}
