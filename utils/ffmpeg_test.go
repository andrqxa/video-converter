package utils

import (
	// "reflect"
	"fmt"
	"reflect"
	"testing"
)

func TestSetArguments(t *testing.T) {
	// Test case 1: Both audio indexes are undefined
	expectedError := fmt.Sprintf("Can't convert because both indexes undefined:\n\trussianAudioIndex = %s\n\tenglishAudioIndex = %s\n\trussianSubtitleIndex = %s\n\tenglishSubtitleIndex = %s\n\t", "-1", "-1", "1", "1")
	_, err := setArguments("-1", "-1", "1", "1", "input.mkv", "output.mkv")
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %s, but got: %v", expectedError, err)
	}

	// Test case 2: Both subtitle indexes are undefined
	expectedError = fmt.Sprintf("Can't convert because both indexes undefined:\n\trussianAudioIndex = %s\n\tenglishAudioIndex = %s\n\trussianSubtitleIndex = %s\n\tenglishSubtitleIndex = %s\n\t", "1", "1", "-1", "-1")
	_, err = setArguments("1", "1", "-1", "-1", "input.mkv", "output.mkv")
	if err == nil || err.Error() != expectedError {
		t.Errorf("Expected error: %s, but got: %v", expectedError, err)
	}

	// "-i", "input.mkv", "-c:v", "libx265", "-crf", "23", "-vf", "scale=-2:720",
	// "-c:a:0", "copy", "-c:a:1", "copy", "-c:s:0", "copy", "-c:s:1", "copy", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:a:1", "-map", "0:s:0", "-map", "0:s:1", "output.mkv"

	// Test case 2: Only Russian audio index is defined
	expectedArgs := []string{"-i", "input.mkv", "-c:v", "libx265", "-crf", "23", "-vf", "scale=-2:720",
		"-c:a:0", "copy", "-c:s:0", "copy", "-c:s:1", "copy", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:s:0", "-map", "0:s:1", "output.mkv"}
	actualArgs, _ := setArguments("0", "-1", "0", "1", "input.mkv", "output.mkv")
	if !reflect.DeepEqual(expectedArgs, actualArgs) {
		t.Errorf("Expected: %v, but got: %v", expectedArgs, actualArgs)
	}

	// Test case 3: Only English audio index is defined
	expectedArgs = []string{"-i", "input.mkv", "-c:v", "libx265", "-crf", "23", "-vf", "scale=-2:720",
		"-c:a:0", "copy", "-c:s:0", "copy", "-c:s:1", "copy", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:s:0", "-map", "0:s:1", "output.mkv"}
	actualArgs, _ = setArguments("-1", "0", "0", "1", "input.mkv", "output.mkv")
	if !reflect.DeepEqual(expectedArgs, actualArgs) {
		t.Errorf("Expected: %v, but got: %v", expectedArgs, actualArgs)
	}

	// Test case 4: Both audio indexes are defined
	expectedArgs = []string{"-i", "input.mkv", "-c:v", "libx265", "-crf", "23", "-vf", "scale=-2:720",
		"-c:a:0", "copy", "-c:a:1", "copy", "-c:s:0", "copy", "-c:s:1", "copy", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:a:1", "-map", "0:s:0", "-map", "0:s:1", "output.mkv"}
	actualArgs, _ = setArguments("0", "1", "0", "1", "input.mkv", "output.mkv")
	if !reflect.DeepEqual(expectedArgs, actualArgs) {
		t.Errorf("Expected: %v, but got: %v", expectedArgs, actualArgs)
	}

	// Test case 5: Only Russian subtitle index is defined
	expectedArgs = []string{"-i", "input.mkv", "-c:v", "libx265", "-crf", "23", "-vf", "scale=-2:720",
		"-c:a:0", "copy", "-c:a:1", "copy", "-c:s:0", "copy", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:a:1", "-map", "0:s:0", "output.mkv"}
	actualArgs, _ = setArguments("0", "1", "0", "-1", "input.mkv", "output.mkv")
	if !reflect.DeepEqual(expectedArgs, actualArgs) {
		t.Errorf("Expected: %v, but got: %v", expectedArgs, actualArgs)
	}

	// Test case 6: Only English subtitle index is defined
	expectedArgs = []string{"-i", "input.mkv", "-c:v", "libx265", "-crf", "23", "-vf", "scale=-2:720",
		"-c:a:0", "copy", "-c:a:1", "copy", "-c:s:0", "copy", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:a:1", "-map", "0:s:0", "output.mkv"}
	actualArgs, _ = setArguments("0", "1", "-1", "0", "input.mkv", "output.mkv")
	if !reflect.DeepEqual(expectedArgs, actualArgs) {
		t.Errorf("Expected: %v, but got: %v", expectedArgs, actualArgs)
	}

	// Test case 7: Both subtitle indexes are defined
	expectedArgs = []string{"-i", "input.mkv", "-c:v", "libx265", "-crf", "23", "-vf", "scale=-2:720",
		"-c:a:0", "copy", "-c:a:1", "copy", "-c:s:0", "copy", "-c:s:1", "copy", "-map", "0:v:0", "-map", "0:a:0", "-map", "0:a:1", "-map", "0:s:0", "-map", "0:s:1", "output.mkv"}
	actualArgs, _ = setArguments("0", "1", "0", "1", "input.mkv", "output.mkv")
	if !reflect.DeepEqual(expectedArgs, actualArgs) {
		t.Errorf("Expected: %v, but got: %v", expectedArgs, actualArgs)
	}
}
