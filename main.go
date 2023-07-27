package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func main() {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(
			"ffmpeg not found. Please make sure ffmpeg is installed and added to your system PATH.",
		)
	}

	// Определяем количество доступных ядер процессора
	numCores := runtime.NumCPU()
	fmt.Printf("Number of CPU cores: %d\n", numCores)

	// Создаем канал с буфером в размере, соответствующем количеству ядер
	semaphore := make(chan struct{}, numCores*3)

	// Используем wait group для ожидания завершения всех горутин
	var wg sync.WaitGroup

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if ext := filepath.Ext(file.Name()); ext == ".mkv" {
			wg.Add(1)
			semaphore <- struct{}{} // Захватываем слот в канале

			go func(inputFile string) {
				defer func() {
					<-semaphore // Освобождаем слот в канале
					wg.Done()
				}()

				outputFile, err := splitFileNameByPattern(inputFile)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("Processing file: %s\n", inputFile)

				// Получаем информацию о потоках аудио и субтитров с помощью ffprobe
				streams := getStreamsInfo(inputFile)

				// Переменные для хранения индексов русскоязычных и англоязычных потоков и субтитров
				var (
					russianAudioIndex    = ""
					englishAudioIndex    = ""
					russianSubtitleIndex = ""
					englishSubtitleIndex = ""
				)
				for _, srm := range streams {
					switch srm.Type {
					case "audio":
						switch srm.Language {
						case "rus":
							russianAudioIndex = strconv.Itoa(srm.Index)
						case "eng":
							englishAudioIndex = strconv.Itoa(srm.Index)
						}
					case "sub":
						switch srm.Language {
						case "rus":
							russianSubtitleIndex = strconv.Itoa(srm.Index)
						case "eng":
							englishSubtitleIndex = strconv.Itoa(srm.Index)
						}
					}
				}

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
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Printf("Error converting file %s: %v\n", inputFile, err)
					log.Printf("Output: %s\n", string(output))
					return
				}

				fmt.Printf("File %s successfully converted to %s\n", inputFile, outputFile)
			}(file.Name())
		}
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Conversion completed.")
}

func parseIndex(str string) int {
	var index int
	_, err := fmt.Sscanf(str, "%d", &index)
	if err != nil {
		return 0
	}
	return index
}

func ReadFileAndSplit(filename string) ([]string, error) {
	// Читаем содержимое файла
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Разделяем содержимое на строки
	lines := strings.Split(string(content), "\n")

	return lines, nil
}

// Функция для получения информации о потоках (аудио или субтитры) с помощью ffprobe
func getStreamsInfo(file string) [4]StreamInfo {
	var res [4]StreamInfo
	errFile := file + ".txt"
	outputFile, err := os.Create(errFile)
	if err != nil {
		log.Fatalf("Ошибка при создании файла %s: %s\n", errFile, err)
	}
	defer func() {
		outputFile.Close()
		err := os.Remove(errFile)
		if err != nil {
			fmt.Println("Ошибка при удалении файла:", err)
			return
		}
	}()

	// Выполняем команду ffprobe и получаем вывод
	cmd := exec.Command("ffprobe", "-i", file)

	cmd.Stderr = outputFile

	// Выполняем команду
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка при выполнении команды: %s\n", err)
	}

	fmt.Printf("Команда успешно выполнена. Вывод ошибок записан %s.\n", errFile)

	// Компилируем регулярные выражения для поиска строк с информацией о потоках
	audioStart := regexp.MustCompile(`^\s*Stream\s*#0:(\d)\(\w\w\w\): Audio:`)
	subtitleStart := regexp.MustCompile(`^\s*Stream\s*#0:(\d)\(\w\w\w\): Subtitle:`)

	audioPatternRus := regexp.MustCompile(`^\s*Stream\s*#0:(\d)\(rus\): Audio:`)
	audioPatternEng := regexp.MustCompile(`^\s*Stream\s*#0:(\d)\(eng\): Audio:`)
	subtitlePatternRus := regexp.MustCompile(`^\s*Stream\s*#0:(\d)\(rus\): Subtitle:`)
	subtitlePatternEng := regexp.MustCompile(`^\s*Stream\s*#0:(\d)\(eng\): Subtitle:`)

	// Обрабатываем вывод команды
	lines, err := ReadFileAndSplit(errFile)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла ошибок: %s\n", err)
	}
	var (
		isAudio bool
		isSubs  bool
		isRusA  bool
		isEngA  bool
		// isRusS     bool
		// isEngS     bool
		startAudio int
		startSub   int
	)
	// находим номер первого аудиопотока и первых сабов
	for _, ln := range lines {
		if match := audioStart.FindStringSubmatch(ln); match != nil && !isAudio {
			isAudio = true
			startAudio = parseIndex(match[1])
		} else if match := subtitleStart.FindStringSubmatch(ln); match != nil && !isSubs {
			isSubs = true
			startSub = parseIndex(match[1])
		}
	}

	// ffmpeg -i From.S02E01.WEB-DL.2160p.H265.SDR.mkv -c:v libx265 -crf 23 -vf "scale=-2:720" -c:a:0 copy -c:a:1 copy -c:s:0 copy -c:s:1 copy -map 0:v:0 -map 0:a:0 -map 0:a:4 -map 0:s:0 -map 0:s:2 From.S02E01.720p.H265.mkv
	// Stream #0:0: Video: hevc (Main), yuv420p(tv, bt709), 3840x2160 [SAR 1:1 DAR 16:9], 23.98 fps, 23.98 tbr, 1k tbn (default)
	// Stream #0:1(rus): Audio: ac3, 48000 Hz, 5.1(side), fltp, 448 kb/s (default)
	// Stream #0:2(rus): Audio: ac3, 48000 Hz, 5.1(side), fltp, 448 kb/s
	// Stream #0:3(rus): Audio: ac3, 48000 Hz, stereo, fltp, 384 kb/s
	// Stream #0:4(rus): Audio: ac3, 48000 Hz, stereo, fltp, 192 kb/s
	// Stream #0:5(eng): Audio: eac3, 48000 Hz, 5.1(side), fltp, 640 kb/s
	// Stream #0:6(rus): Subtitle: subrip (default) (forced)
	// Stream #0:7(rus): Subtitle: subrip
	// Stream #0:8(eng): Subtitle: subrip
	// Stream #0:9(eng): Subtitle: subrip (original) (hearing impaired)
	// Stream #0:10: Video: mjpeg (Baseline), yuvj420p(pc, bt470bg/unknown/unknown), 303x450 [SAR 120:120 DAR 101:150], 90k tbr, 90k tbn (attached pic)

	for _, line := range lines {
		if match := audioPatternRus.FindStringSubmatch(line); match != nil && !isRusA {
			isRusA = true
			index := parseIndex(match[1]) - startAudio
			res[0] = StreamInfo{"audio", index, "rus"}
		} else if match := audioPatternEng.FindStringSubmatch(line); match != nil && !isEngA {
			isEngA = true
			index := parseIndex(match[1]) - startAudio
			res[1] = StreamInfo{"audio", index, "eng"}
		} else if match := subtitlePatternRus.FindStringSubmatch(line); match != nil { //&& !isRusS {
			// isRusS = true
			index := parseIndex(match[1]) - startSub
			res[2] = StreamInfo{"sub", index, "rus"}
		} else if match := subtitlePatternEng.FindStringSubmatch(line); match != nil { //&& !isEngS {
			// isEngS = true
			index := parseIndex(match[1]) - startSub
			res[3] = StreamInfo{"sub", index, "eng"}
		}
	}

	return res
}

// функция для изменения имени файла
func splitFileNameByPattern(filename string) (string, error) {
	const desc = ".720p.H265"
	re := regexp.MustCompile(`([sS]\d\d[eE]\d\d)`)
	matches := re.FindStringSubmatchIndex(filename)
	if len(matches) > 0 {
		ext := filepath.Ext(filename)
		return filename[:matches[0]] + filename[matches[2]:matches[3]] + desc + ext, nil
	}
	return "", fmt.Errorf("Pattern not found\n")
}

// Структура для хранения информации о потоках (аудио или субтитры)
type StreamInfo struct {
	Type     string
	Index    int
	Language string
}
