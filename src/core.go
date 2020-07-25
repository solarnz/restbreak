package restbreak

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/solarnz/restbreak/src/types"
	"github.com/solarnz/restbreak/src/x"
)

func Run() {
	restbreak := Parse()

	for _, action := range restbreak.Actions {
		action.LastInactive = time.Now()
		action.LastActive = time.Now()
	}

	// This only works on X, need to implement a listener for Wayland as well.
	// On wayland, the last active timer is just set to the time the application
	// starts
	xListener, err := x.NewXActivityListener()
	if err != nil {
		panic(err)
	}

	appTrayIcon := &tray{
		restbreak,
	}
	go appTrayIcon.showTray()

	stopChan := make(chan interface{})
	eventChan := make(chan types.ActivityEvent)

	// Assume that if we're just starting, we're currently active.
	lastActive := time.Now()

	go func() {
		for event := range eventChan {
			lastActive = event.LastActiveTimestamp
		}
	}()

	go func() {
		for {
			for _, action := range restbreak.Actions {
				action.LastActive = lastActive

				// Reset the time the user was last in-active if the user has
				// been inactive for longer than the specified time.
				// This handles the case where you leave your machine for a
				// different reason than the reminders.
				if action.LastActive.Before(time.Now().Add(-1 * action.InactiveResetDuration)) {
					action.LastInactive = time.Now()
					action.CurrentRefires = 0
				}

				// If the user has been active for longer than the allowed
				// duration, it's time to take action!
				if action.LastInactive.Before(time.Now().Add(-1 * action.ActiveLimitDuration)) {
					// We need to check how many times a notification has been fired.
					// If it's been fired too many times, the user might be too
					// busy and it's not worthwhile notifying them.
					if time.Now().Sub(action.LastFired) > action.RefireDuration && action.CurrentRefires <= action.RefireCount {
						action.LastFired = time.Now()
						action.CurrentRefires++

						// Run the command the user specified.
						go func() {
							cmd := exec.Command("sh", "-c", action.Command)
							err := cmd.Run()
							if err != nil {
								log.Printf("Unable to run command, %s\n", err)
							}
						}()

						if action.ResetOnLimit {
							action.LastInactive = time.Now()
							action.CurrentRefires = 0
						}
					}
				}

				if action.SystrayMenuItem != nil {
					durationUntilFire := action.LastInactive.Add(action.ActiveLimitDuration).Sub(time.Now())
					seconds := (durationUntilFire - durationUntilFire.Truncate(time.Minute)).Seconds()
					action.SystrayMenuItem.SetTitle(
						action.Name + "\n  Time Left: " +
							fmt.Sprintf("%02d:%02d", int(durationUntilFire.Minutes()), int(seconds)),
					)
				}
			}

			// Only check the status every second, as this application doesn't
			// need to be real-time.
			time.Sleep(time.Second)
		}
	}()

	xListener.Listen(stopChan, eventChan)
}
