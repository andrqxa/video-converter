package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
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
			outputFile := u.SplitFileNameByPattern(inputFile)
			fmt.Printf("Processing file: %s\n", inputFile)

			// Получаем информацию о потоках аудио и субтитров с помощью ffprobe
			streams := u.GetStreamsInfo(inputFile)

			// Получаем индексы для рус/англ аудиопотока и субтитров
			// TODO: убрать магические строки
			russianAudioIndex := strconv.Itoa(streams.Get("rusAudio").Index)
			englishAudioIndex := strconv.Itoa(streams.Get("engAudio").Index)
			russianSubtitleIndex := strconv.Itoa(streams.Get("rusSubs").Index)
			englishSubtitleIndex := strconv.Itoa(streams.Get("engSubs").Index)

			// Выполняем конвертацию
			err := u.ConvertFile(
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
			fmt.Printf("File %s successfully converted to %s\n", inputFile, outputFile)
		}(file.Name())
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Conversion completed.")
}
