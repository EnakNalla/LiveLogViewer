package backend

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LogFile struct {
	Path  string   `json:"path"`
	Lines []string `json:"lines"`
}

func (a *App) SelectLog() Response {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Log Files (*.log;*.txt)",
		Filters: []runtime.FileFilter{{DisplayName: "Logs (*.txt;*.log)", Pattern: "*.txt;*.log"}},
	})
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return Failure(err.Error())
	}
	if selection == "" {
		return Failure("No selection was made.")
	}

	lines, err := a.readLines(a.settings.TailLines, selection)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return Failure(err.Error())
	}

	a.logs[selection] = 1
	a.watcher.Add(selection)

	return Success(LogFile{selection, lines})
}

func (a *App) readLines(nLines int, path string) ([]string, error) {
	lines := make([]string, nLines)

	fileHandle, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer fileHandle.Close()

	stat, _ := fileHandle.Stat()
	filesize := stat.Size()

	mb := filesize / (1 << 20)
	if mb < a.settings.TailThreshold {
		fileScanner := bufio.NewScanner(fileHandle)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			lines = append(lines, fileScanner.Text())
		}

		return lines, nil
	}

	line := ""
	lineCounter := 0

	var cursor int64 = 0
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			lines[lineCounter] = line
			lineCounter++
			line = ""

			if lineCounter == nLines {
				break
			}
		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -filesize {
			break
		}
	}

	return lines, nil
}

func (a *App) readNextLine(path string) (string, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -filesize {
			break
		}
	}

	return line, nil
}

func (a *App) RemoveLog(path string) {
	delete(a.logs, path)
	a.watcher.Remove(path)
}
