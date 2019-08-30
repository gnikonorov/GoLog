package golog

import "testing"

func TestIsValidFileActionAcceptsAllValidFileActions(t *testing.T) {
	if !FileActionAppend.IsValidFileAction() {
		t.Errorf("Expected 'FileActionAppend' to be valid file action but was not.")
	}

	if !FileActionCompress.IsValidFileAction() {
		t.Errorf("Expected 'FileActionCompress' to be valid file action but was not.")
	}

	if !FileActionDelete.IsValidFileAction() {
		t.Errorf("Expected 'FileActionDelete' to be valid file action but was not.")
	}

	if !FileActionNone.IsValidFileAction() {
		t.Errorf("Expected 'FileActionNone' to be valid file action but was not.")
	}
}

func TestIsValidFileActionRejectsFileActionsThatAreInvalid(t *testing.T) {
	// NOTE: Any action outside range of 1 -> 4 is invalid, and we tested validity above.
	var badLoggingFileAction LoggingFileAction

	badLoggingFileAction = 0
	if badLoggingFileAction.IsValidFileAction() {
		t.Errorf("Expected invalid logging file action '%d' to be invalid but it was valid.", badLoggingFileAction)
	}

	badLoggingFileAction = 5
	if badLoggingFileAction.IsValidFileAction() {
		t.Errorf("Expected invalid logging file action '%d' to be invalid but it was valid.", badLoggingFileAction)
	}

	badLoggingFileAction = 23531616
	if badLoggingFileAction.IsValidFileAction() {
		t.Errorf("Expected invalid logging file action '%d' to be invalid but it was valid.", badLoggingFileAction)
	}
}
