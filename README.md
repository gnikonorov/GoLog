# GoLog

GoLog is a logger for the Go programming language, written in Go. It aims to be easy to use, and provides basic
logging functionality for the Go programmer.

## Logging Levels

The logging levels shown below are provided by this logger, and are colorized if the logger is initialized with 
`colorize = true`.

+ `Info`    ( Presented in the terminal's native colors )
+ `Debug`   ( green text on the terminal's native background )
+ `Warning` ( yellow text on the terminal's native background )
+ `Error`   ( red text on the terminal's native background )
+ `Fatal`   ( white text on a red background )

## Logging Modes

The logger may be set up to run in the three modes listed below:

+ `Screen` - Output logs to `STDOUT`
+ `File`   - Output logs to a user provided file
+ `Both`   - Output logs to both `STDOUT` and a user provided file

## Startup Actions

Upon initialization of the logger, the user may specify what to do with an existing log file if the user has specified 
`Screen` or `Both` as their logging mode.

The following actions may be taken:

+ `FileAppend`   - If the log file to which the user is writing already exists, add new logs onto it
+ `FileCompress` - If the log file to which the user is writing already exists, compress it into a `.gz` file
+ `FileDelete`   - If the log file to which the user is writing already exists, delete it

The user may also specify `FileActionNone`. For now, this is logically equivalent to passing in `FileAppend`. 

## Initialization

The logger is a struct initialized by calling `SetupLoggerWithFilename`.
The following arguments are passed in during initialization:

+ `logMode (int)` - This is the mode the logger will operate in, as described in `Logging Modes`
+ `logFileStartupAction (int)` - The startup action taken by the logger as described in section `Startup Actions`
+ `logDirectory (string)` - If `logMode` is `Both` or `File`, this is the directory to which the logger will write logs.
+ `logFile (string)` - If the `logMode` is `Both` or `File`, this is the file to which the logger will write logs.
+ `shouldColorize (bool)` - If set to `true` logging will be colorized as described in section `Logging Modes`.
