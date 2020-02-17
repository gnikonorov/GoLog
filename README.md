[![Build Status](https://travis-ci.org/gnikonorov/golog.svg?branch=master)](https://travis-ci.org/gnikonorov/golog)

# golog

`golog` is a logger for the Go programming language.

## Logging Levels

The logging levels shown below are provided by this logger, and are colorized if the logger is initialized with 
`colorize = true`.

+ `Info`    ( Presented in the terminal's native colors )
+ `Debug`   ( green text on the terminal's native background )
+ `Warning` ( yellow text on the terminal's native background )
+ `Error`   ( red text on the terminal's native background )
+ `Fatal`   ( white text on a red background )
+ `Panic`   ( white text on a red background )

## Logging Modes

The logger may be set up to run in the three modes listed below. These modes are defined in `logging_output_modes.go`:

+ `ModeScreen` - Output logs to `STDOUT`
+ `ModeFile`   - Output logs to a user provided file
+ `ModeBoth`   - Output logs to both `STDOUT` and a user provided file

## Startup Actions

Upon initialization of the logger, the user may specify what to do with an existing log file if the user has specified `ModeFile` or `ModeBoth` as their logging mode.

The following actions may be taken. These actions are defined in `logging_file_actions.go`:

+ `FileActionAppend`   - If the log file to which the user is writing already exists, add new logs onto it
+ `FileActionCompress` - If the log file to which the user is writing already exists, compress it into a `.gz` file
+ `FileActionDelete`   - If the log file to which the user is writing already exists, delete it

The user may also specify `FileActionNone`. For now, this is logically equivalent to passing in `FileAppend`. 

## Initialization

The logger is a structure that may be initialized by the functions defined below

### Raw Initializaion

The following arguments are passed in during initialization by calling `SetupLoggerFromFields`:

+ `logMode (LoggingOutputMode)` - The mode the logger will operate in, as described in `Logging Modes`
+ `logFileStartupAction (LoggingFileAction)` - Startup action taken by the logger as described in section `Startup Actions`
+ `logDirectory (string)` - If `logMode` is `Both` or `File`, this is the directory to which the logger will write logs.
+ `logFile (string)` - If the `logMode` is `Both` or `File`, this is the file to which the logger will write logs.
+ `shouldColorize (bool)` - If set to `true` logging will be colorized as described in section `Logging Modes`.
+ `isMock (bool)` - If set to `true` all filesystem ops will be stubbed by a fake filesystem as defined by [afero](https://github.com/spf13/afero)
+ `isAsynch (bool)` - If set to `true` perform all logging asynchly in a seperate goroutine

For example:
```
logger := golog.SetupLoggerFromFields(golog.ModeBoth, golog.FileActionNone, "/path/to/log/file", "file.log", true, false, true)
```

### Initialization From Struct

The logger may be initialized by creating a `LoggingConfig` struct and passing it to `SetupLoggerFromStruct`. 

The definition for the struct is shown below:

```
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
```
A sample initialization would thus be as follows:

```
config := golog.LoggingConfig{LogMode: golog.Both, LogDirectory: "/dir/to/my/go/project", LogFile: "important.log", ShouldColorize: true, IsMock: false, IsAsynch: true}
logger := golog.SetupLoggerFromStruct(&config)
```
### Initialization From JSON File

The logger may be initialized by specifying the path to a `JSON` file that contains its definition to `SetupLoggerFromConfigFile` along with the configuration profile name `name`. The `name` is what you have chosen to name a current configuration.

A sample configuration file containing 3 profiles is shown below. Note how each entry has a configuration name `name`.

```
[{
	"name": "test1",
	"logMode": 1,                            // ModeFile
	"logFileStartupAction": 1,               // FileActionAppend
	"logDirectory": "/dir/to/my/go/project",
	"logFile": "test1.log",
	"shouldColorize": true,
	"isMock": true,
	"isAsynch": false
}, {
	"name": "test2",
	"logMode": 2,                            // ModeScreen 
	"logFileStartupAction": 2,               // FileActionCompress
	"logDirectory": "/dir/to/my/go/project",
	"logFile": "test2.log",
	"shouldColorize": true,
	"isMock": false,
	"isAsynch": false
}, {
	"name": "test3",
	"logMode": 3,                            // ModeBoth
	"logFileStartupAction": 3,               // FileActionDelete
	"logDirectory": "/dir/to/my/go/project",
	"logFile": "test3.log",
	"shouldColorize": true,
	"isMock": true,
	"isAsynch": true
}]
```

If we wanted to initialize a logger to have the profile of `test1`, we would do the following:

```
var profile = "test1"
var configFile = "/path/to/my/config/config.json"
logger := golog.SetupLoggerFromConfigFile(configFile, profile)
```
