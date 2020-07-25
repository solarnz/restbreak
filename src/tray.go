package restbreak

import (
	"github.com/solarnz/restbreak/src/icon"

	"github.com/getlantern/systray"
)

type tray struct {
	restbreak *RestBreak
}

func (t *tray) showTray() {
	systray.Run(t.onReady, t.onExit)
}

func (t *tray) onReady() {
	systray.SetTitle("RestBreak")
	systray.SetIcon(icon.Data)
	systray.AddMenuItem("RestBreak", "")
	systray.AddSeparator()

	for _, action := range t.restbreak.Actions {
		action.SystrayMenuItem = systray.AddMenuItem(action.Name, "")
	}
}

func (t *tray) onExit() {
}
