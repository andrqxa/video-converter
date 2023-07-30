package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Функция возвращает путь к ffmpeg
func Ffmpeg() string {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(
			"ffmpeg not found. Please make sure ffmpeg is installed and added to your system PATH.",
		)
	}
	return ffmpegPath
}

func ConvertFile(
	ffmpegPath string,
	inputFile string,
	outputFile string,
	russianAudioIndex string,
	englishAudioIndex string,
	russianSubtitleIndex string,
	englishSubtitleIndex string,
) error {
	// Формируем команду ffmpeg для сохранения выбранных потоков и субтитров
	cmd := exec.Command(
		ffmpegPath,
		"-i",
		inputFile,
		"-c:v",
		"libx265",
		"-crf",
		"23",
		"-vf",
		"scale=-2:720",
		"-c:a:0",
		"copy",
		"-c:a:1",
		"copy",
		"-c:s:0",
		"copy",
		"-c:s:1",
		"copy",
		"-map",
		"0:v:0",
		"-map",
		"0:a:"+russianAudioIndex,
		"-map",
		"0:a:"+englishAudioIndex,
		"-map",
		"0:s:"+russianSubtitleIndex,
		"-map",
		"0:s:"+englishSubtitleIndex,
		outputFile)

	// Запускаем команду и выводим результат
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error converting file %s: %v\n", inputFile, err)
		log.Printf("File %s is removing\n", outputFile)
		err2 := os.Remove(outputFile)
		if err2 != nil {
			return fmt.Errorf("1. Ошибка при конвертировании файла %w\n2.Ошибка при удалении файла: %w", err, err2)
		}
		return fmt.Errorf("1. Ошибка при конвертировании файла %w\n", err)
	}
	return nil
}
