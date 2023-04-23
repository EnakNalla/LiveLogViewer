package backend

import (
	"encoding/json"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Settings struct {
	Theme           string `json:"theme"`
	TailThreshold   int64  `json:"tailThreshold"` // in mb
	TailLines       int    `json:"tailLines"`
	HighlightLevels bool   `json:"highlightLevels"`
	PollingEnabled  bool   `json:"pollingEnabled"`
	PollInterval    int    `json:"pollInterval"` // ms
	IgnoreCase      bool   `json:"ignoreCase"`
	TextWrap        bool   `json:"textWrap"`
	LineNumbers     bool   `json:"lineNumbers"`
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
		a.settings = Settings{TailLines: 100, TailThreshold: 5, HighlightLevels: false, Theme: "none", PollingEnabled: false, PollInterval: 1000, IgnoreCase: false, TextWrap: false, LineNumbers: true}
		a.WriteSettings(a.settings)
	} else {
		file, err := os.Open(a.settingsPath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "error", err.Error())
			runtime.LogErrorf(a.ctx, err.Error())
			return
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&a.settings)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (a *App) WriteSettings(settings Settings) Settings {
	data, err := json.Marshal(&settings)
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
		runtime.LogErrorf(a.ctx, err.Error())
	}

	err = os.WriteFile(a.settingsPath, data, os.ModePerm)
	if err != nil {
		runtime.EventsEmit(a.ctx, "error", err.Error())
		runtime.LogErrorf(a.ctx, err.Error())
	}

	return a.settings
}
