/*
	Logging output modes
*/
package golog

import "testing"

func TestIsValidModeRecognizesAllValidModes(t *testing.T) {
	var loggingOutputMode LoggingOutputMode

	loggingOutputMode = ModeFile
	if !loggingOutputMode.IsValidMode() {
		t.Errorf("Expected output mode 'ModeFile' to be a valid logging output mode but it was not.")
	}

	loggingOutputMode = ModeScreen
	if !loggingOutputMode.IsValidMode() {
		t.Errorf("Expected output mode 'ModeScreen' to be a valid logging output mode but it was not.")
	}

	loggingOutputMode = ModeBoth
	if !loggingOutputMode.IsValidMode() {
		t.Errorf("Expected output mode 'ModeBoth' to be a valid logging output mode but it was not.")
	}
}

func TestIsValidModeRejectsAllInvalidModes(t *testing.T) {
	// NOTE: Any mode outside range of 1 -> 3 is invalid, and we tested validity above.
	var badLoggingOutputMode LoggingOutputMode

	badLoggingOutputMode = 0
	if badLoggingOutputMode.IsValidMode() {
		t.Errorf("Expected output mode '%d' to be invalid but it was valid.", badLoggingOutputMode)
	}

	badLoggingOutputMode = 4
	if badLoggingOutputMode.IsValidMode() {
		t.Errorf("Expected output mode '%d' to be invalid but it was valid.", badLoggingOutputMode)
	}

	badLoggingOutputMode = 123
	if badLoggingOutputMode.IsValidMode() {
		t.Errorf("Expected output mode '%d' to be invalid but it was valid.", badLoggingOutputMode)
	}
}

func TestIntConvertsTheLoggingOutputModeToAnInteger(t *testing.T) {
	var loggingOutputMode LoggingOutputMode
	var expectedLoggingOutputModeAsInt int
	var actualLoggingOutputModeAsInt int

	loggingOutputMode = ModeFile
	expectedLoggingOutputModeAsInt = 1
	actualLoggingOutputModeAsInt = loggingOutputMode.Int()
	if actualLoggingOutputModeAsInt != expectedLoggingOutputModeAsInt {
		t.Errorf("Expected output mode 'ModeFile' to be converted to int value '%d' but it was converted to '%d'.", expectedLoggingOutputModeAsInt, actualLoggingOutputModeAsInt)
	}

	loggingOutputMode = ModeScreen
	expectedLoggingOutputModeAsInt = 2
	actualLoggingOutputModeAsInt = loggingOutputMode.Int()
	if actualLoggingOutputModeAsInt != expectedLoggingOutputModeAsInt {
		t.Errorf("Expected output mode 'ModeScreen' to be converted to int value '%d' but it was converted to '%d'.", expectedLoggingOutputModeAsInt, actualLoggingOutputModeAsInt)
	}

	loggingOutputMode = ModeBoth
	expectedLoggingOutputModeAsInt = 3
	actualLoggingOutputModeAsInt = loggingOutputMode.Int()
	if actualLoggingOutputModeAsInt != expectedLoggingOutputModeAsInt {
		t.Errorf("Expected output mode 'ModeBoth' to be converted to int value '%d' but it was converted to '%d'.", expectedLoggingOutputModeAsInt, actualLoggingOutputModeAsInt)
	}
}

func TestStringConvertsTheLoggingOutputModeToAStringRepresentingTheIntValue(t *testing.T) {
	var loggingOutputMode LoggingOutputMode
	var expectedLoggingOutputModeAsString string
	var actualLoggingOutputModeAsString string

	loggingOutputMode = ModeFile
	expectedLoggingOutputModeAsString = "1"
	actualLoggingOutputModeAsString = loggingOutputMode.String()
	if actualLoggingOutputModeAsString != expectedLoggingOutputModeAsString {
		t.Errorf("Expected output mode 'ModeFile' to be converted to string value '%s' but it was converted to '%s'.", expectedLoggingOutputModeAsString, actualLoggingOutputModeAsString)
	}

	loggingOutputMode = ModeScreen
	expectedLoggingOutputModeAsString = "2"
	actualLoggingOutputModeAsString = loggingOutputMode.String()
	if actualLoggingOutputModeAsString != expectedLoggingOutputModeAsString {
		t.Errorf("Expected output mode 'ModeScreen' to be converted to string value '%s' but it was converted to '%s'.", expectedLoggingOutputModeAsString, actualLoggingOutputModeAsString)
	}

	loggingOutputMode = ModeBoth
	expectedLoggingOutputModeAsString = "3"
	actualLoggingOutputModeAsString = loggingOutputMode.String()
	if actualLoggingOutputModeAsString != expectedLoggingOutputModeAsString {
		t.Errorf("Expected output mode 'ModeBoth' to be converted to string value '%s' but it was converted to '%s'.", expectedLoggingOutputModeAsString, actualLoggingOutputModeAsString)
	}
}
