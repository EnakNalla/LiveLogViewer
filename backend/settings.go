package backend

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Settings struct {
	Theme             string `json:"theme"`
	TailThreshold     int64  `json:"tailThreshold"` // in mb
	TailLines         int    `json:"tailLines"`
	HighlightErrors   bool   `json:"highlightErrors"`
	HighlightWarnings bool   `json:"highlightWarnings"`
}

func (a *App) GetSettings() Response {
	return Success(a.settings)
}

func (a *App) getSettingsDir() {
	path := ""

	switch runtime.Environment(a.ctx).Platform {
	case "windows":
		path = os.Getenv("APPDATA") + "\\LiveLogViewer"
	case "darwin":
		path = os.Getenv("HOME") + "/Library/Application Support/LiveLogViewer/"
	case "linux":
		path = os.Getenv("HOME") + "/.config/LiveLogViewer"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.FileMode(0755)); err != nil {
			runtime.EventsEmit(a.ctx, "error", err.Error())
			runtime.LogErrorf(a.ctx, err.Error())
		}
	}

	a.settingsPath = path
}

func (a *App) readSettings() {
	a.getSettingsDir()

	a.settingsPath += "config.json"
	if _, err := os.Stat(a.settingsPath); os.IsNotExist(err) {
		a.settings = Settings{TailLines: 100, TailThreshold: 20, HighlightErrors: true, HighlightWarnings: false, Theme: "none"}
		a.WriteSettings(a.settings)
	} else {
		file, err := os.Open(a.settingsPath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "error", err.Error())
			runtime.LogErrorf(a.ctx, err.Error())
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		decoder.Decode(&a.settings)
	}
}

func (a *App) WriteSettings(settings Settings) Settings {
	data, err := json.Marshal(&settings)
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
		runtime.LogErrorf(a.ctx, err.Error())
	}

	ioutil.WriteFile(a.settingsPath, data, os.ModePerm)

	return a.settings
}
