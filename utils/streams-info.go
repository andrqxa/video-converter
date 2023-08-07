package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

const (
	// Константы для регулярных выражений для поиска строк с информацией о потоках
	audioStartPtrn    = `^\s*Stream\s*#\d:(\d\d?)\(\w\w\w\): Audio:`
	subtitleStartPtrn = `^\s*Stream\s*#\d:(\d\d?)\(\w\w\w\): Subtitle:`

	// audioPatternRusPtrn    = `^\s*Stream\s*#0:(\d\d?)\(rus\): Audio:`
	// audioPatternEngPtrn    = `^\s*Stream\s*#0:(\d\d?)\(eng\): Audio:`
	// subtitlePatternRusPtrn = `^\s*Stream\s*#0:(\d\d?)\(rus\): Subtitle:`
	// subtitlePatternEngPtrn = `^\s*Stream\s*#0:(\d\d?)\(eng\): Subtitle:`

	videoPtrn = `^Stream\s*#\d:\d\d?\(\w+\):\s*Video:\s*.*,\s*(\d{3}\d?)x(\d{3}\d?)\s*.*$`
)

var (
	// Компилируем регулярные выражения для поиска строк с информацией о потоках
	audioStart    = regexp.MustCompile(audioStartPtrn)
	subtitleStart = regexp.MustCompile(subtitleStartPtrn)

	// audioPatternRus    = regexp.MustCompile(audioPatternRusPtrn)
	// audioPatternEng    = regexp.MustCompile(audioPatternEngPtrn)
	// subtitlePatternRus = regexp.MustCompile(subtitlePatternRusPtrn)
	// subtitlePatternEng = regexp.MustCompile(subtitlePatternEngPtrn)

	videoPattern = regexp.MustCompile(videoPtrn)
)

// Структура для хранения информации о видео
type VideoInfo struct {
	Width  int
	Height int
}

// Структура для хранения информации об аудиопотоке (аудио или субтитры)
type AudioInfo struct {
	Index int
	// Offset   int
	Title    string
	Language string
}

type Audios []AudioInfo

// Структура для хранения информации о субтитрах
type SubsInfo struct {
	Index int
	// Offset   int
	Title    string
	Language string
}

type Subs []SubsInfo

// Структура для хранения всей информации
type AllStreamInfo struct {
	v VideoInfo
	a Audios
	s Subs
}

func NewAllStreamInfo() *AllStreamInfo {
	v := VideoInfo{}
	// a := AudioInfo{Index: -1}
	// s := SubsInfo{Index: -1}
	a := make(Audios, 0)
	s := make(Subs, 0)
	return &AllStreamInfo{v: v, a: a, s: s}
}

func parseIndex(str string) int {
	var index int
	_, err := fmt.Sscanf(str, "%d", &index)
	if err != nil {
		return 0
	}
	return index
}

func GetRawInfo(file string, outputFile *os.File) []string {
	// Выполняем команду ffprobe и получаем вывод

	cmd := exec.Command("ffprobe", "-i", file)

	cmd.Stderr = outputFile

	// Выполняем команду
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка при выполнении команды: %s\n", err)
	}

	return ReadFileAndSplit(outputFile.Name())
}

// Функция для получения информации о потоках (аудио или субтитры) с помощью ffprobe
func GetStreamsInfo(file string) AllStreamInfo {
	// ffprobe -i Beforeigners.S01E01.1080p.HMAX.WEB-DL.DD5.1.H.264-BLS.mkv

	// Input #0, matroska,webm, from '.\Beforeigners.S01E01.1080p.HMAX.WEB-DL.DD5.1.H.264-BLS.mkv':
	//   Metadata:
	//     encoder         : libebml v1.4.2 + libmatroska v1.6.4
	//     creation_time   : 2022-04-07T06:47:45.000000Z
	//   Duration: 00:48:59.00, start: 0.000000, bitrate: 8984 kb/s
	//   Stream #0:0: Video: h264 (High), yuv420p(tv, bt709, progressive), 1920x1080 [SAR 1:1 DAR 16:9], 25 fps, 25 tbr, 1k tbn (default)
	//     Metadata:
	//       BPS             : 8214064
	//       DURATION        : 00:48:59.000000000
	//       NUMBER_OF_FRAMES: 73475
	//       NUMBER_OF_BYTES : 3017641915
	//       _STATISTICS_WRITING_APP: mkvmerge v65.0.0 ('Too Much') 64-bit
	//       _STATISTICS_WRITING_DATE_UTC: 2022-04-07 06:47:45
	//       _STATISTICS_TAGS: BPS DURATION NUMBER_OF_FRAMES NUMBER_OF_BYTES
	//   Stream #0:1(rus): Audio: ac3, 48000 Hz, 5.1(side), fltp, 384 kb/s (default)
	//     Metadata:
	//       title           : Кириллица
	//       BPS             : 384000
	//       DURATION        : 00:48:58.976000000
	//       NUMBER_OF_FRAMES: 91843
	//       NUMBER_OF_BYTES : 141070848
	//       _STATISTICS_WRITING_APP: mkvmerge v65.0.0 ('Too Much') 64-bit
	//       _STATISTICS_WRITING_DATE_UTC: 2022-04-07 06:47:45
	//       _STATISTICS_TAGS: BPS DURATION NUMBER_OF_FRAMES NUMBER_OF_BYTES
	//   Stream #0:2(nor): Audio: ac3, 48000 Hz, 5.1(side), fltp, 384 kb/s
	//     Metadata:
	//       BPS             : 384000
	//       DURATION        : 00:48:58.976000000
	//       NUMBER_OF_FRAMES: 91843
	//       NUMBER_OF_BYTES : 141070848
	//       _STATISTICS_WRITING_APP: mkvmerge v65.0.0 ('Too Much') 64-bit
	//       _STATISTICS_WRITING_DATE_UTC: 2022-04-07 06:47:45
	//       _STATISTICS_TAGS: BPS DURATION NUMBER_OF_FRAMES NUMBER_OF_BYTES
	//   Stream #0:3(rus): Subtitle: subrip
	//     Metadata:
	//       title           : Кириллица
	//       BPS             : 98
	//       DURATION        : 00:45:37.590000000
	//       NUMBER_OF_FRAMES: 547
	//       NUMBER_OF_BYTES : 33648
	//       _STATISTICS_WRITING_APP: mkvmerge v65.0.0 ('Too Much') 64-bit
	//       _STATISTICS_WRITING_DATE_UTC: 2022-04-07 06:47:45
	//       _STATISTICS_TAGS: BPS DURATION NUMBER_OF_FRAMES NUMBER_OF_BYTES
	//   Stream #0:4(nor): Subtitle: subrip
	//     Metadata:
	//       BPS             : 52
	//       DURATION        : 00:48:08.440000000
	//       NUMBER_OF_FRAMES: 401
	//       NUMBER_OF_BYTES : 18935
	//       _STATISTICS_WRITING_APP: mkvmerge v65.0.0 ('Too Much') 64-bit
	//       _STATISTICS_WRITING_DATE_UTC: 2022-04-07 06:47:45
	//       _STATISTICS_TAGS: BPS DURATION NUMBER_OF_FRAMES NUMBER_OF_BYTES
	//   Stream #0:5(eng): Subtitle: subrip
	//     Metadata:
	//       BPS             : 55
	//       DURATION        : 00:45:36.800000000
	//       NUMBER_OF_FRAMES: 416
	//       NUMBER_OF_BYTES : 18943
	//       _STATISTICS_WRITING_APP: mkvmerge v65.0.0 ('Too Much') 64-bit
	//       _STATISTICS_WRITING_DATE_UTC: 2022-04-07 06:47:45
	//       _STATISTICS_TAGS: BPS DURATION NUMBER_OF_FRAMES NUMBER_OF_BYTES

	// TODO: не брать форсированные субтитры. Пока же просто берутся ПОСЛЕДНИЕ в списке, потому что обычно первые это форсированные

	var (
		isRusA bool
		isEngA bool
		// isRusS  bool
		// isEngS  bool
		errFile = file + ".txt"
	)

	// Создаем файл куда будет записана информация о файле
	outputFile := CreateFile(errFile)
	// в defer делаем закрытие и удаление этого файла
	defer func() {
		RemoveFile(errFile, outputFile)
	}()

	res := NewAllStreamInfo()

	// Получаем информацию из файла вывода ffprobe
	lines := GetRawInfo(file, outputFile)

	leng := len(lines)
	for i := 0; i < leng; i++ {
		if i >= leng {
			break
		}
		line := lines[i]
		if match := videoPattern.FindStringSubmatch(line); match != nil {
			// обрабатывать пока не закончится видео раздел
			// handleVideo()
			// изменить i
			continue
		} else if match := audioStart.FindStringSubmatch(line); match != nil {
			// обрабатывать пока не закончится аудио раздел
			// handleAudio()
			// изменить i
			continue
		} else if match := subtitleStart.FindStringSubmatch(line); match != nil {
			// обрабатывать пока не закончится раздел субтитров
			// subtitleHandle()
			// изменить i
			continue
		}
	}

	for _, line := range lines {
		if match := videoPattern.FindStringSubmatch(line); match != nil {
			handleVideo()
		} else if match := audioStart.FindStringSubmatch(line); match != nil {
			handleAudio()
		} else if match := subtitleStart.FindStringSubmatch(line); match != nil {
			subtitleHandle()
		}

		// if match := audioPatternRus.FindStringSubmatch(line); match != nil && !isRusA {
		// 	isRusA = true
		// 	index := parseIndex(match[1]) - startAudio
		// 	res.UpdateIndex("rusAudio", index)
		// } else if match := audioPatternEng.FindStringSubmatch(line); match != nil && !isEngA {
		// 	isEngA = true
		// 	index := parseIndex(match[1]) - startAudio
		// 	res.UpdateIndex("engAudio", index)
		// } else if match := subtitlePatternRus.FindStringSubmatch(line); match != nil { //&& !isRusS {
		// 	// isRusS = true
		// 	index := parseIndex(match[1]) - startSub
		// 	res.UpdateIndex("rusSubs", index)
		// } else if match := subtitlePatternEng.FindStringSubmatch(line); match != nil { //&& !isEngS {
		// 	// isEngS = true
		// 	index := parseIndex(match[1]) - startSub
		// 	res.UpdateIndex("engSubs", index)
		// }
	}

	return *res
}
