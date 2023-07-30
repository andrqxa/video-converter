package utils

import (
	"fmt"
	"io/fs"
	"io/ioutil"
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
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Ошибка при чтении файла ошибок: %s\n", err)
	}

	// Разделяем содержимое на строки
	lines := strings.Split(string(content), "\n")
	return lines
}

// функция для изменения имени файла - оставляем только название и номер_сезона.номер_серии
func SplitFileNameByPattern(filename string) string {
	const (
		desc = ".720p.H265"
		se   = `([sS]\d\d[eE]\d\d)`
	)
	re := regexp.MustCompile(se)
	matches := re.FindStringSubmatchIndex(filename)
	if len(matches) > 0 {
		ext := filepath.Ext(filename)
		return filename[:matches[0]] + filename[matches[2]:matches[3]] + desc + ext
	}
	log.Fatalf("Pattern %s not found\n", se)
	return ""
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
