package backend

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	settings     Settings
	settingsPath string
	watcher      *fsnotify.Watcher
	logs         map[string]uint8
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.readSettings()
	a.logs = make(map[string]uint8)

	var err error = nil
	a.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err.Error())
	}
	defer a.watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-a.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					_, exists := a.logs[event.Name]
					if exists {
						line, err := a.readNextLine(event.Name)
						if err != nil {
							runtime.LogErrorf(ctx, err.Error())
						} else {
							runtime.EventsEmit(ctx, event.Name, line)
						}
					}
				}
			}
		}
	}()

	<-make(chan struct{})
}
