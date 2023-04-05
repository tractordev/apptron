package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"tractor.dev/apptron/bridge/api/window"
	"tractor.dev/apptron/bridge/event"
	"tractor.dev/apptron/bridge/misc"
)

type WindowSettings struct {
	Key      string
	Position misc.Position
	Size     misc.Size
}

func SaveWindowSettings(win *window.Window, identifier string, key string) bool {
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Println("[WindowSettings] Failed to get user cache dir")
		return false
	}

	dirpath := filepath.Join(dir, identifier)

	// create directory if not exists
	if _, err = os.Stat(dirpath); os.IsNotExist(err) {
		err = os.Mkdir(dirpath, os.ModePerm)
		if err != nil {
			log.Println("[WindowSettings] Failed to create save directory:", dirpath, err)
			return false
		}
	}

	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		log.Println("[WindowSettings] Directory doesn't exist:", dirpath, err)
		return false
	}

	settings := WindowSettings{Key: key, Position: win.GetOuterPosition(), Size: win.GetInnerSize()}

	data, _ := json.MarshalIndent(settings, "", " ")

	fname := "window_" + key + ".json"
	path := filepath.Join(dirpath, fname)

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Println("[WindowSettings] Did not write file:", path, err)
		return false
	}

	return true
}

func RestoreWindowSettings(win *window.Window, identifier string, key string) bool {
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Println("[WindowSettings] Failed to get user cache dir")
		return false
	}

	dirpath := filepath.Join(dir, identifier)

	fname := "window_" + key + ".json"
	path := filepath.Join(dirpath, fname)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("[WindowSettings] Failed to read file:", path, err)
		return false
	}

	settings := WindowSettings{}
	err = json.Unmarshal(content, &settings)
	if err != nil {
		log.Println("[WindowSettings] Failed to parse JSON:", path, err)
		return false
	}

	win.SetPosition(settings.Position)
	win.SetSize(settings.Size)

	return true
}

func setupWindowRestoreListener(identifier string) {
	event.Listen("__APPTRON_Platform_listener__WindowRestore__", func(e event.Event) error {
		if e.Type == event.Created {
			win, _ := window.Get(e.Window)
			if win != nil && len(win.ID) > 0 {
				RestoreWindowSettings(win, identifier, win.ID)
			}
		}

		// TODO: event.Close is not fired on MacOS
		if e.Type == event.Close || e.Type == event.Destroyed {
			win, _ := window.Get(e.Window)
			if win != nil && len(win.ID) > 0 {
				SaveWindowSettings(win, identifier, win.ID)
			}
		}

		return nil
	})
}
