/*
	Actions logger will take regarding physical log files on disk
*/
package golog

type LoggingFileAction int

const (
	FileActionAppend   LoggingFileAction = iota + 1 // Instructs the logger to append onto an existing log if one exists
	FileActionCompress                              // Instructs the logger to compress an existing log file if one exists
	FileActionDelete                                // Instructs the logger to remove an existing log file if one exits
	FileActionNone                                  // Indicates no file actions ( e.g.: when user is writing to screen only )
)

func (fileAction LoggingFileAction) IsValidFileAction() bool {
	return (fileAction == FileActionAppend   ||
	        fileAction == FileActionCompress ||
	        fileAction == FileActionDelete   ||
	        fileAction == FileActionNone)
}
