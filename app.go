package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LogFile struct {
	Path  string   `json:"path"`
	Lines []string `json:"lines"`
}

type Result struct {
	Success bool    `json:"success"`
	Data    LogFile `json:"data"`
	Error   string  `json:"error"`
}

type App struct {
	ctx context.Context
}

var watcher, err = fsnotify.NewWatcher()
var logs = map[string]LogFile{}

func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					file, exists := logs[event.Name]
					if exists {
						line, err := readLastLine(file.Path)
						if err != nil {
							runtime.LogErrorf(ctx, err.Error())
						} else {
							runtime.EventsEmit(ctx, file.Path, line)
						}
					}
				}
			}
		}
	}()

	<-make(chan struct{})
}

func (a *App) SelectFile() Result {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select log",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Logs (*.txt;*.log)",
				Pattern:     "*.txt;*.log",
			},
		},
	})
	if err != nil || selection == "" {
		if err != nil {
			runtime.LogErrorf(a.ctx, err.Error())
		} else {
			runtime.LogErrorf(a.ctx, "No file selection made")
		}
		return Result{Data: LogFile{}, Success: false, Error: "Failed to open file."}
	}

	_, exists := logs[selection]
	if exists {
		return Result{Data: LogFile{}, Success: false, Error: "File with path exists."}
	}

	logs[selection] = LogFile{Path: selection}

	return addWatcher(selection, a)
}

func addWatcher(path string, a *App) Result {
	logs[path] = LogFile{Path: path}

	err = watcher.Add(path)
	if err != nil {
		delete(logs, path)

		runtime.LogErrorf(a.ctx, err.Error())
		return Result{Data: LogFile{}, Success: false, Error: err.Error()}
	}

	result, err := readLines(logs[path])
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
	}

	return result
}

func readLines(logFile LogFile) (Result, error) {
	file, err := os.Open(logFile.Path)
	if err != nil {
		delete(logs, logFile.Path)

		return Result{Data: logFile, Success: false, Error: "Failed to read file"}, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		logFile.Lines = append(logFile.Lines, fileScanner.Text())
	}

	return Result{Data: logFile, Success: true, Error: ""}, nil
}

func readLastLine(file string) (string, error) {
	fileHandle, err := os.Open(file)
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

func (a *App) RemoveWatcher(path string) Result {
	delete(logs, path)

	err = watcher.Remove(path)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return Result{Data: LogFile{}, Success: false, Error: err.Error()}
	}

	return Result{Data: LogFile{}, Success: true, Error: ""}
}
