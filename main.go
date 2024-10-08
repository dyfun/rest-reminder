package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/ncruces/zenity"
	"strconv"
	"time"
)

var (
	appName         = "Rest Reminder"
	appTooltip      = "Stay focused, take breaks."
	defaultInterval = 25 * time.Minute
	timer           *time.Ticker
	stopCh          = make(chan bool)
	startTime       time.Time
	isStarted       = false
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle(appName)
	systray.SetTooltip(appTooltip)

	mStart := systray.AddMenuItem("Start", "Start the reminder")
	mStop := systray.AddMenuItem("Stop", "Stop the reminder")
	mSettings := systray.AddMenuItem("Settings", "Adjust settings")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit from the app")

	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				startTimer()
			case <-mStop.ClickedCh:
				stopTimer()
			case <-mSettings.ClickedCh:
				openSettings()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

	go func() {
		for {
			if isStarted {
				time.Sleep(1 * time.Second)
				elapsedSeconds := int(time.Since(startTime).Seconds())

				elapsedSeconds = elapsedSeconds % (24 * 3600)
				hour := elapsedSeconds / 3600
				elapsedSeconds %= 3600
				minute := elapsedSeconds / 60
				elapsedSeconds %= 60
				second := elapsedSeconds

				newAppName := fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
				systray.SetTitle(newAppName)
			}

		}
	}()
}

func startTimer() {
	message := fmt.Sprintf("Rest reminder started. You will be reminded every %s", defaultInterval)

	err := beeep.Notify("Rest Reminder Started", message, "assets/information.png")
	if err != nil {
		panic(err)
	}

	timer = time.NewTicker(defaultInterval)
	startTime = time.Now()
	isStarted = true

	go func() {
		for {
			select {
			case <-timer.C:
				err := beeep.Notify("Rest Reminder!", "Take a break", "assets/warning.png")
				if err != nil {
					panic(err)
				}

			case <-stopCh:
				return
			}
		}
	}()
}

func stopTimer() {
	fmt.Println("Stop clicked")
	timer.Stop()
	stopCh <- true

	isStarted = false
}

func openSettings() {
	newIntervalStr, err := zenity.Entry("Enter new interval in minutes:", zenity.Title("Set Interval"))
	if err != nil {
		fmt.Println("Settings dialog canceled or error:", err)
		return
	}

	newIntervalMinutes, err := strconv.Atoi(newIntervalStr)
	if err != nil || newIntervalMinutes <= 0 {
		zenity.Error("Please enter a valid positive integer.", zenity.Title("Invalid Input"))
		return
	}
	defaultInterval = time.Duration(newIntervalMinutes) * time.Second

	if isStarted {
		stopTimer()
		startTimer()
	}
}

func onExit() {
	// clean up here
}
