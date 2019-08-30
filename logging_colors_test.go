package golog

import "testing"

func TestStringProperlyConvertsAllLoggingColorsToAString(t *testing.T) {
	wantColorForColorNone := ""
	if colorNone.String() != wantColorForColorNone {
		t.Errorf("Expected colorNone to be %q but got %q", wantColorForColorNone, colorNone.String())
	}

	wantColorForColorDebug := "\x1B[32m"
	if colorDebug.String() != wantColorForColorDebug {
		t.Errorf("Expected colorDebug to be %q but got %q", wantColorForColorDebug, colorDebug.String())
	}

	wantColorForColorErr := "\x1B[31m"
	if colorErr.String() != wantColorForColorErr {
		t.Errorf("Expected colorErr to be %q but got %q", wantColorForColorErr, colorErr.String())
	}

	wantColorForColorFatal := "\x1B[0;37;41m"
	if colorFatal.String() != wantColorForColorFatal {
		t.Errorf("Expected colorFatal to be %q but got %q", wantColorForColorFatal, colorFatal.String())
	}

	wantColorForColorPanic := "\x1B[0;37;41m"
	if colorPanic.String() != wantColorForColorPanic {
		t.Errorf("Expected colorPanic to be %q but got %q", wantColorForColorPanic, colorPanic.String())
	}

	wantColorForColorReset := "\x1B[0m"
	if colorReset.String() != wantColorForColorReset {
		t.Errorf("Expected colorReset to be %q but got %q", wantColorForColorReset, colorReset.String())
	}

	wantColorForColorWarn := "\x1B[33m"
	if colorWarn.String() != wantColorForColorWarn {
		t.Errorf("Expected colorWarn to be %q but got %q", wantColorForColorWarn, colorWarn.String())
	}
}
