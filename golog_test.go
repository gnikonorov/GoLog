package golog

import "testing"

func makeLoggerInstance() (*Logger, error) {
	logDirectory := ""
	logFile := ""
	logConfig := LoggingConfig{ LogMode: ModeScreen, LogFileStartupAction: FileActionAppend, LogDirectory: logDirectory, LogFile: logFile, ShouldColorize: true, IsMock: true }

	logger, err := SetupLoggerFromStruct(&logConfig)
	if err != nil {
		return nil, err
	}

	return &logger, nil
}

func TestSetContextProperlySetsLogContext(t *testing.T) {
	logger, err := makeLoggerInstance()
	if err != nil {
		t.Errorf("Failed to set up logger because: '%s'", err.Error())
		return
	}

	testContext := "CONTEXT"
	logger.SetContext(testContext)
	if logger.context != testContext {
		t.Errorf("Expected logger context to be '%s' but it was '%s'", testContext, logger.context)
	}

	testContext = ""
	logger.SetContext(testContext)
	if logger.context != testContext {
		t.Errorf("Expected logger context to be '%s' but it was '%s'", testContext, logger.context)
	}
}

func TestIsUninitializedReturnsFalseWhenLoggerIsInstantiated(t *testing.T) {
	logger, err := makeLoggerInstance()
	if err != nil {
		t.Errorf("Failed to set up logger because: '%s'", err.Error())
		return
	}

	if logger.IsUninitialized() {
		t.Errorf("Expected logger to be initialized, but it was uninitialized")
	}
}

func TestIsUninitializedReturnsTrueWhenLoggingModeIsUnset(t *testing.T) {
	// NOTE: This test works since all logger setup methods specify a logging mode
	//       and thus it is impossible to have an unset logging mode in the logger
	//       if it is being used as returned from a setup method.
	logger, err := makeLoggerInstance()
	if err != nil {
		t.Errorf("Failed to set up logger because: '%s'", err.Error())
		return
	}

	logger.loggingMode = 0
	if !logger.IsUninitialized() {
		t.Errorf("Expected logger to be uninitialized, but it was initialized")
	}
}
