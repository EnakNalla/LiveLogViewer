package backend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

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
		return Success(nil)
	}

	_, exists := a.logs[selection]
	if exists {
		return Failure("Log already loaded")
	}

	lines, err := a.readLines(a.settings.TailLines, selection)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return Failure(err.Error())
	}

	a.logs[selection] = 1
	if a.settings.PollingEnabled {
		a.pollFile(selection)
	} else {
		err := a.watcher.Add(selection)
		if err != nil {
			delete(a.logs, selection)
			return Failure(err.Error())
		}
	}

	return Success(LogFile{selection, lines})
}

func (a *App) readLines(nLines int, path string) ([]string, error) {
	lines := make([]string, nLines)

	fileHandle, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return lines, err
	}
	defer func(fileHandle *os.File) {
		_ = fileHandle.Close()
	}(fileHandle)

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
		_, err2 := fileHandle.Seek(cursor, io.SeekEnd)
		if err2 != nil {
			return nil, err2
		}

		char := make([]byte, 1)
		_, err3 := fileHandle.Read(char)
		if err3 != nil {
			return nil, err3
		}

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
	fileHandle, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer func(fileHandle *os.File) {
		_ = fileHandle.Close()
	}(fileHandle)

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		_, err2 := fileHandle.Seek(cursor, io.SeekEnd)
		if err2 != nil {
			return "", err2
		}

		char := make([]byte, 1)
		_, err3 := fileHandle.Read(char)
		if err3 != nil {
			return "", err3
		}

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

func (a *App) pollFile(path string) {
	fileHandle, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
	}

	stat, err := fileHandle.Stat()
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
	}
	lastReadSize := stat.Size()

	err = fileHandle.Close()
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
		return
	}

	go func() {
		for {
			time.Sleep(time.Duration(a.settings.PollInterval) * time.Millisecond)
			// TODO: better implement polling exit
			if _, exists := a.logs[path]; !exists {
				// runtime.EventsEmit(a.ctx, "error", "File not found!")
				break
			}

			fileHandle, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				runtime.EventsEmit(a.ctx, "error", err.Error())
			}

			stat, _ := fileHandle.Stat()
			filesize := stat.Size()

			if filesize > lastReadSize {
				_, err := fileHandle.Seek(lastReadSize, io.SeekStart)
				if err != nil {
					delete(a.logs, path)
					runtime.EventsEmit(a.ctx, "error", err.Error())
					break
				}
				fileScanner := bufio.NewScanner(fileHandle)
				fileScanner.Split(bufio.ScanLines)

				for fileScanner.Scan() {
					runtime.EventsEmit(a.ctx, path, fileScanner.Text())
				}

				lastReadSize = filesize
				err = fileHandle.Close()
				if err != nil {
					delete(a.logs, path)
					runtime.EventsEmit(a.ctx, "error", err.Error())
					break
				}
			}
		}
	}()
}

func (a *App) RemoveLog(path string) {
	delete(a.logs, path)

	if a.watcher == nil {
		return
	}

	_ = a.watcher.Remove(path)
}
