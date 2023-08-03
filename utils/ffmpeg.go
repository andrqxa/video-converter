package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// сколько threads для одной команды ffmpeg
const numThreads = "4"

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

func setArguments(russianAudioIndex string, englishAudioIndex string, russianSubtitleIndex string, englishSubtitleIndex string, inputFile string, outputFile string) ([]string, error) {
	if russianAudioIndex == "-1" && englishAudioIndex == "-1" {
		return nil, fmt.Errorf("Can't convert because both audio indexes in %s are undefined:\n\trussianAudioIndex = %s\n\tenglishAudioIndex = %s\n\trussianSubtitleIndex = %s\n\tenglishSubtitleIndex = %s\n\t", inputFile, russianAudioIndex, englishAudioIndex, russianSubtitleIndex, englishSubtitleIndex)

		// TODO: очистить ресурсы?
	}

	res := make([]string, 0)
	res = append(res, "-i")
	res = append(res, inputFile)
	res = append(res, "-c:v")
	res = append(res, "libx265")
	// res = append(res, "-threads")
	// res = append(res, "numThreads")
	res = append(res, "-crf")
	res = append(res, "23")
	res = append(res, "-vf")
	res = append(res, "scale=-2:720")

	// предполагаем что хоть одна аудиодорожка есть
	if russianAudioIndex == "-1" || englishAudioIndex == "-1" {
		res = append(res, "-c:a:0")
		res = append(res, "copy")
	} else {
		res = append(res, "-c:a:0")
		res = append(res, "copy")
		res = append(res, "-c:a:1")
		res = append(res, "copy")
	}

	switch {
	case russianSubtitleIndex == "-1" && englishSubtitleIndex != "-1", russianSubtitleIndex != "-1" && englishSubtitleIndex == "-1":
		{
			res = append(res, "-c:s:0")
			res = append(res, "copy")

		}
	case russianSubtitleIndex == "-1" && englishSubtitleIndex == "-1":

	default:
		{
			res = append(res, "-c:s:0")
			res = append(res, "copy")
			res = append(res, "-c:s:1")
			res = append(res, "copy")
		}
	}

	res = append(res, "-map")
	res = append(res, "0:v:0")

	switch {
	case russianAudioIndex != "-1" && englishAudioIndex == "-1":
		{
			res = append(res, "-map")
			res = append(res, "0:a:"+russianAudioIndex)
		}
	case russianAudioIndex == "-1" && englishAudioIndex != "-1":
		{
			res = append(res, "-map")
			res = append(res, "0:a:"+englishAudioIndex)
		}
	case russianAudioIndex != "-1" && englishAudioIndex != "-1":
		{
			res = append(res, "-map")
			res = append(res, "0:a:"+russianAudioIndex)
			res = append(res, "-map")
			res = append(res, "0:a:"+englishAudioIndex)
		}
	}
	switch {
	case russianSubtitleIndex != "-1" && englishSubtitleIndex == "-1":
		{
			res = append(res, "-map")
			res = append(res, "0:s:"+russianSubtitleIndex)
		}
	case russianSubtitleIndex == "-1" && englishSubtitleIndex != "-1":
		{
			res = append(res, "-map")
			res = append(res, "0:s:"+englishSubtitleIndex)
		}
	case russianSubtitleIndex == "-1" && englishSubtitleIndex == "-1":
	default:
		{
			res = append(res, "-map")
			res = append(res, "0:s:"+russianSubtitleIndex)
			res = append(res, "-map")
			res = append(res, "0:s:"+englishSubtitleIndex)
		}
	}

	res = append(res, outputFile)

	return res, nil
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
	args, err := setArguments(russianAudioIndex, englishAudioIndex, russianSubtitleIndex, englishSubtitleIndex, inputFile, outputFile)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Args for ffmeg = %v\n", args)
	cmd := exec.Command(ffmpegPath, args...)

	// Запускаем команду и выводим результат
	_, err = cmd.CombinedOutput()
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
