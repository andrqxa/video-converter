package utils

import (
	"os"
	"testing"
)

func TestReadFileAndSplit(t *testing.T) {
	// Создаем временный файл и записываем в него тестовые данные
	testContent := "line 1\nline 2\nline 3"
	tmpFile := "testfile.txt"
	err := os.WriteFile(tmpFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %s", err)
	}
	defer os.Remove(tmpFile)

	// Вызываем тестируемую функцию
	lines := ReadFileAndSplit(tmpFile)

	// Проверяем результат
	expectedLines := []string{"line 1", "line 2", "line 3"}
	if !equalSlices(lines, expectedLines) {
		t.Errorf("Функция ReadFileAndSplit() вернула неверный результат. Ожидаемый: %v, Полученный: %v", expectedLines, lines)
	}
}

func equalSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func TestSplitFileNameByPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Yellowstone S03E01 WEB-DL 2160p.mkv",
			expected: "Yellowstone S03E01.720p.H265.mkv",
		},
		{
			input:    "01x00 Pilot [CBS Drama+OPT+Eng].mkv",
			expected: "S01E00.Pilot.720p.H265.mkv",
		},
		{
			input:    "01. The One Where Monica Gets a Roommate.mkv",
			expected: "E01.The One Where Monica Gets a Roommate.720p.H265.mkv",
		},
	}

	for _, test := range tests {
		result, err := SplitFileNameByPattern(test.input)
		if err != nil {
			t.Errorf("Ошибка при обработке файла: %v", err)
		}
		if result != test.expected {
			t.Errorf("Для входного значения %q, ожидался результат %q, получено %q", test.input, test.expected, result)
		}
	}
}
