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

+ `logMode (int)` - The mode the logger will operate in, as described in `Logging Modes`
+ `logFileStartupAction (int)` - Startup action taken by the logger as described in section `Startup Actions`
+ `logDirectory (string)` - If `logMode` is `Both` or `File`, this is the directory to which the logger will write logs.
+ `logFile (string)` - If the `logMode` is `Both` or `File`, this is the file to which the logger will write logs.
+ `shouldColorize (bool)` - If set to `true` logging will be colorized as described in section `Logging Modes`.

### Initialization From Struct

The logger may be initialized by creating a `LoggingConfig` struct and passing it to `SetupLoggerFromStruct`. 

The definition for the struct is shown below:

```
type LoggingConfig struct {
	Name                 string // The logger profile name
	LogMode              string // The logging mode
	LogFileStartupAction string // The action the logger will take on startup
	LogDirectory         string // The directory to which the logger writes
	LogFile              string // The name of the log file to write to
	ShouldColorize       bool   // Indicates if we should output information in color
}
```
A sample initialization would thus be as follows:

```
config := golog.LoggingConfig{LogMode: golog.Both, LogDirectory: "/home/gnikonorov/projects/go/src/tester", LogFile: "test.log", ShouldColorize: true}
logger := golog.SetupLoggerFromStruct(&config)
```
### Initialization From JSON File

The logger may be initialized by specifying the path to a `JSON` file that contains its definition to `SetupLoggerFromConfigFile` along with the configuration profile name `name`. The `name` is simply what you have chosen to name a current configuration.

A sample configuration file containing 3 profiles is shown below. Note how each entry has a configuration name `name`.

```
[{
	"name": "test1",
	"logMode": "FILE",
	"logFileStartupAction": "NONE",
	"logDirectory": "/home/gnikonorov/projects/go/src/tester",
	"logFile": "test1.log",
	"shouldColorize": true

}, {
	"name": "test2",
	"logMode": "SCREEN",
	"logFileStartupAction": "APPEND",
	"logDirectory": "/home/user/gnikonorov/logs",
	"logFile": "test2.log",
	"shouldColorize": true
}, {
	"name": "test3",
	"logMode": "BOTH",
	"logFileStartupAction": "COMPRESS",
	"logDirectory": "/home/gnikonorov/projects/go/src/tester/logs",
	"logFile": "test3.log",
	"shouldColorize": true
}]
```

If we wanted to initialize a logger to have the profile of `test1`, we would do the following:

```
var profile = "test1"
var configFile = "/home/gnikonorov/.golog/conf.json"
logger := golog.SetupLoggerFromConfigFile(configFile, profile)
```

### Initialization From Raw Parameters

It is also possible to initialize the logger by passing in raw parameters. For example:

```
logger := golog.SetupLoggerFromFields(golog.Both, golog.FileActionNone, "/home/gnikonorov/projects/go/src/tester", "test.log", true)
```
