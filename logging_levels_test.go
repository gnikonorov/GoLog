package golog

import "testing"

func TestStringProperlyConvertsAllLoggingLevelsToAString(t *testing.T) {
	wantStringForLevelDebug := "DEBUG"
	if levelDebug.String() != wantStringForLevelDebug {
		t.Errorf("Expected levelDebug to be '%s' but got '%s'", wantStringForLevelDebug, levelDebug.String())
	}

	wantStringForLevelErr := "ERROR"
	if levelErr.String() != wantStringForLevelErr {
		t.Errorf("Expected levelErr to be '%s' but got '%s'", wantStringForLevelErr, levelErr.String())
	}

	wantStringForLevelFatal := "FATAL"
	if levelFatal.String() != wantStringForLevelFatal {
		t.Errorf("Expected levelFatal to be '%s' but got '%s'", wantStringForLevelFatal, levelFatal.String())
	}

	wantStringForLevelInfo := "INFO"
	if levelInfo.String() != wantStringForLevelInfo {
		t.Errorf("Expected levelInfo to be '%s' but got '%s'", wantStringForLevelInfo, levelInfo.String())
	}

	wantStringForLevelPanic := "PANIC"
	if levelPanic.String() != wantStringForLevelPanic {
		t.Errorf("Expected levelPanic to be '%s' but got '%s'", wantStringForLevelPanic, levelPanic.String())
	}

	wantStringForLevelWarn := "WARNING"
	if levelWarn.String() != wantStringForLevelWarn {
		t.Errorf("Expected levelWarn to be '%s' but got '%s'", wantStringForLevelWarn, levelWarn.String())
	}
}
