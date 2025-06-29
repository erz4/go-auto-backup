package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type Config struct {
	SourceDirs []string  `json:"source_dirs"`
	TargetDir  string    `json:"target_dir"`
	SyncTime   string    `json:"sync_time"` // e.g. "02:00" or interval like "6h"
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".go-dir-backup-config.json")
}

func LoadOrInitConfig() (*Config, error) {
	path := configPath()
	if _, err := os.Stat(path); err == nil {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		var cfg Config
		if err := json.NewDecoder(f).Decode(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return InitConfig(path)
}

func InitConfig(path string) (*Config, error) {
	a := app.New()
	w := a.NewWindow("Backup Configuration")

	var cfg *Config
	var submitErr error
	done := make(chan struct{})

	sourceEntry := widget.NewMultiLineEntry()
	sourceEntry.SetPlaceHolder("/path/to/source1\n/path/to/source2\n...")

	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("/path/to/target")
	targetBtn := widget.NewButton("Choose Target Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				targetEntry.SetText(uri.Path())
			}
		}, w)
	})

	syncEntry := widget.NewEntry()
	syncEntry.SetPlaceHolder("02:00 or 6h")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Source Directories (one per line)", Widget: sourceEntry},
			{Text: "Target Directory", Widget: container.NewHBox(targetEntry, targetBtn)},
			{Text: "Sync Time", Widget: syncEntry},
		},
		OnSubmit: func() {
			sourceDirs := []string{}
			for _, line := range strings.Split(sourceEntry.Text, "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					sourceDirs = append(sourceDirs, line)
				}
			}
			targetDir := strings.TrimSpace(targetEntry.Text)
			syncTime := strings.TrimSpace(syncEntry.Text)
			cfg = &Config{SourceDirs: sourceDirs, TargetDir: targetDir, SyncTime: syncTime}
			f, err := os.Create(path)
			if err != nil {
				submitErr = err
				dialog.ShowError(err, w)
				return
			}
			defer f.Close()
			json.NewEncoder(f).Encode(cfg)
			close(done)
			w.Close()
		},
	}

	w.SetContent(container.NewVBox(form))
	w.Resize(fyne.NewSize(500, 400))

	go func() {
		<-done
		a.Quit()
	}()

	w.ShowAndRun()

	return cfg, submitErr
}

func splitLines(s string) []string {
	lines := []string{}
	for _, l := range strings.Split(s, "\n") {
		lines = append(lines, strings.TrimSpace(l))
	}
	return lines
}
