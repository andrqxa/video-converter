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
	audioStartPtrn    = `^\s*Stream\s*#0:(\d\d?)\(\w\w\w\): Audio:`
	subtitleStartPtrn = `^\s*Stream\s*#0:(\d\d?)\(\w\w\w\): Subtitle:`

	audioPatternRusPtrn    = `^\s*Stream\s*#0:(\d\d?)\(rus\): Audio:`
	audioPatternEngPtrn    = `^\s*Stream\s*#0:(\d\d?)\(eng\): Audio:`
	subtitlePatternRusPtrn = `^\s*Stream\s*#0:(\d\d?)\(rus\): Subtitle:`
	subtitlePatternEngPtrn = `^\s*Stream\s*#0:(\d\d?)\(eng\): Subtitle:`
)

var (
	// Компилируем регулярные выражения для поиска строк с информацией о потоках
	audioStart    = regexp.MustCompile(audioStartPtrn)
	subtitleStart = regexp.MustCompile(subtitleStartPtrn)

	audioPatternRus    = regexp.MustCompile(audioPatternRusPtrn)
	audioPatternEng    = regexp.MustCompile(audioPatternEngPtrn)
	subtitlePatternRus = regexp.MustCompile(subtitlePatternRusPtrn)
	subtitlePatternEng = regexp.MustCompile(subtitlePatternEngPtrn)
)

// Структура для хранения информации о потоках (аудио или субтитры)
type StreamInfo struct {
	Type     string
	Index    int
	Language string
}

type AllStreamInfo map[string]StreamInfo

func NewAllStreamInfo() *AllStreamInfo {
	res := make(AllStreamInfo)
	return &res
}

func (a *AllStreamInfo) Put(idx string, s StreamInfo) {
	(*a)[idx] = s
}

func (a AllStreamInfo) Get(idx string) StreamInfo {
	return a[idx]
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

func GetFirstIndexes(lines []string) (int, int) {
	var (
		isAudio    bool
		isSubs     bool
		startAudio int
		startSub   int
	)

	// fmt.Printf("Команда успешно выполнена. Вывод ошибок записан %s.\n", file)

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

	return startAudio, startSub
}

// Функция для получения информации о потоках (аудио или субтитры) с помощью ffprobe
func GetStreamsInfo(file string) AllStreamInfo {
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

	// находим номер первого аудиопотока и первых сабов
	startAudio, startSub := GetFirstIndexes(lines)

	for _, line := range lines {
		if match := audioPatternRus.FindStringSubmatch(line); match != nil && !isRusA {
			isRusA = true
			index := parseIndex(match[1]) - startAudio
			res.Put("rusAudio", StreamInfo{"audio", index, "rus"})
		} else if match := audioPatternEng.FindStringSubmatch(line); match != nil && !isEngA {
			isEngA = true
			index := parseIndex(match[1]) - startAudio
			res.Put("engAudio", StreamInfo{"audio", index, "eng"})
		} else if match := subtitlePatternRus.FindStringSubmatch(line); match != nil { //&& !isRusS {
			// isRusS = true
			index := parseIndex(match[1]) - startSub
			res.Put("rusSubs", StreamInfo{"sub", index, "rus"})
		} else if match := subtitlePatternEng.FindStringSubmatch(line); match != nil { //&& !isEngS {
			// isEngS = true
			index := parseIndex(match[1]) - startSub
			res.Put("engSubs", StreamInfo{"sub", index, "eng"})
		}
	}

	return *res
}
