package utils

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type (
	dirFiles []fs.DirEntry
)

// Получаем все файлы из директории с нужным расширением
func GetFiles(extn string) []fs.DirEntry {
	files := make(dirFiles, 0)

	fl, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, f := range fl {
		if f.IsDir() {
			continue
		}
		if ext := filepath.Ext(f.Name()); ext == extn {
			files = append(files, f)
		}
	}

	return files
}

// Читаем файл и делим на строки
func ReadFileAndSplit(filename string) []string {
	// Читаем содержимое файла
	content, err := os.ReadFile(filename)

	if err != nil {
		log.Fatalf("Ошибка при чтении файла ошибок: %s\n", err)
	}

	// Разделяем содержимое на строки
	lines := strings.Split(string(content), "\n")
	return lines
}

// функция для изменения имени файла - оставляем только название и номер_сезона.номер_серии
func SplitFileNameByPattern(filename string) (string, error) {
	// 1. Yellowstone S03E01 WEB-DL 2160p.mkv			=> Yellowstone S03E01.720p.H265.mkv
	// 2. 01x00 Pilot [CBS Drama+OPT+Eng].mkv          	=> S01E00.Pilot.720p.H265.mkv
	// 3. 01. The One Where Monica Gets a Roommate.mkv 	=> E01.The One Where Monica Gets a Roommate.720p.H265.mkv

	const (
		desc = ".720p.H265"
	)

	// Паттерн 1: ([sS]\d\d[eE]\d\d-?\d?\d?)
	pattern1 := regexp.MustCompile(`([sS]\d\d[eE]\d\d-?\d?\d?)`)
	matches := pattern1.FindStringSubmatchIndex(filename)
	if len(matches) > 0 {
		ext := filepath.Ext(filename)
		return filename[:matches[0]] + filename[matches[2]:matches[3]] + desc + ext, nil
	}

	// Паттерн 2: (\d\d)x(\d\d)\s*(.*)\s*\[.*
	pattern2 := regexp.MustCompile(`(\d\d)x(\d\d)\s*(.*)\s*\[.*`)
	matches = pattern2.FindStringSubmatchIndex(filename)
	if len(matches) > 0 {
		ext := filepath.Ext(filename)
		return "S" + filename[matches[2]:matches[3]] + "E" + filename[matches[4]:matches[5]] + "." + strings.TrimSpace(filename[matches[6]:matches[7]]) + desc + ext, nil
	}

	// Паттерн 3: (\d\d)\.\s*(.*)\..*
	pattern3 := regexp.MustCompile(`(\d\d)\.\s*(.*)\..*`)
	matches = pattern3.FindStringSubmatchIndex(filename)
	if len(matches) > 0 {
		ext := filepath.Ext(filename)
		return "E" + filename[matches[2]:matches[3]] + "." + strings.TrimSpace(filename[matches[4]:matches[5]]) + desc + ext, nil
	}

	return "", fmt.Errorf("ни один из паттернов не найден в имени файла: %s", filename)
}

// функция для создания файла по имени
func CreateFile(name string) *os.File {
	outputFile, err := os.Create(name)
	if err != nil {
		log.Fatalf("Ошибка при создании файла %s: %s\n", name, err)
	}
	return outputFile
}

// функция для закрытия и удаления файла
func RemoveFile(name string, file *os.File) {
	file.Close()
	err := os.Remove(name)
	if err != nil {
		fmt.Printf("Ошибка при удалении файла %s: %v", name, err)
	}
}
