package main

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"
	u "video-converter/utils"
)

const (
	fileExt = ".mkv"
)

func main() {

	// Получаем все файлы в текущем каталоге с расширением .mkv
	files := u.GetFiles(fileExt)

	// Получаем путь к ffmpeg
	ffmpegPath := u.Ffmpeg()
	fmt.Printf("FFMPEG = %s\n", ffmpegPath)

	// Определяем количество доступных ядер процессора
	numCores := runtime.NumCPU()
	fmt.Printf("Number of CPU cores: %d\n", numCores)

	// Создаем канал с буфером в размере, соответствующем количеству ядер
	semaphore := make(chan struct{}, numCores*3)

	// Используем wait group для ожидания завершения всех горутин
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		semaphore <- struct{}{} // Захватываем слот в канале

		go func(inputFile string) {
			defer func() {
				<-semaphore // Освобождаем слот в канале
				wg.Done()
			}()

			// получаем новое имя для перeкодированного файла
			outputFile, err := u.SplitFileNameByPattern(inputFile)
			if err != nil {
				log.Printf("ERROR: patern for file %s hasn't found", inputFile)
				return
			}

			fmt.Printf("Processing file: %s\n", inputFile)

			// Получаем информацию о потоках аудио и субтитров с помощью ffprobe
			streams := u.GetStreamsInfo(inputFile)

			// Получаем индексы для рус/англ аудиопотока и субтитров
			// TODO: убрать магические строки
			russianAudioIndex := strconv.Itoa(streams.Get("rusAudio").Index)
			englishAudioIndex := strconv.Itoa(streams.Get("engAudio").Index)
			russianSubtitleIndex := strconv.Itoa(streams.Get("rusSubs").Index)
			englishSubtitleIndex := strconv.Itoa(streams.Get("engSubs").Index)

			// fmt.Printf("russianAudioIndex = %s\n", russianAudioIndex)
			// fmt.Printf("englishAudioIndex = %s\n", englishAudioIndex)
			// fmt.Printf("russianSubtitleIndex = %s\n", russianSubtitleIndex)
			// fmt.Printf("englishSubtitleIndex = %s\n", englishSubtitleIndex)

			// время начала конвертации
			start := time.Now()
			// Выполняем конвертацию
			err = u.ConvertFile(
				ffmpegPath,
				inputFile,
				outputFile,
				russianAudioIndex,
				englishAudioIndex,
				russianSubtitleIndex,
				englishSubtitleIndex,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Calculate the elapsed time
			elapsed := time.Since(start)
			// Format and print the elapsed time
			hours := int(elapsed.Hours())
			minutes := int(elapsed.Minutes()) % 60
			seconds := int(elapsed.Seconds()) % 60
			fmt.Printf("File %s was successfully converted to %s in %02d:%02d:%02d\n", inputFile, outputFile, hours, minutes, seconds)
		}(file.Name())
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Conversion completed.")
}
