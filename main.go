package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
	"golang.design/x/hotkey"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type HotkeyConfig struct {
	StartPause string `mapstructure:"startpause"`
	Reset      string `mapstructure:"reset"`
}

type AppConfig struct {
	Hotkeys HotkeyConfig `mapstructure:"hotkeys"`
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("%s/.config/gowatch", os.Getenv("HOME")))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	startPauseStr := config.Hotkeys.StartPause
	resetStr := config.Hotkeys.Reset

	startPauseMods, startPauseKey, err := ParseHotkeyString(startPauseStr)
	if err != nil {
		log.Fatalf("Invalid StartPause hotkey: %v", err)
	}
	resetMods, resetKey, err := ParseHotkeyString(resetStr)
	if err != nil {
		log.Fatalf("Invalid Reset hotkey: %v", err)
	}

	startPause := hotkey.New(startPauseMods, startPauseKey)
	reset := hotkey.New(resetMods, resetKey)

	file, err := os.Create("/tmp/gowatch")
	check(err)

	if err := startPause.Register(); err != nil {
		fmt.Println("Failed to register start/pause hotkey:", err)
		return
	}
	defer startPause.Unregister()

	if err := reset.Register(); err != nil {
		fmt.Println("Failed to register reset hotkey:", err)
		return
	}
	defer reset.Unregister()

	running := false
	var startTime time.Time
	var elapsed time.Duration

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if running {
					now := time.Now()
					printTime(elapsed + now.Sub(startTime))
				} else {
					printTime(elapsed)
				}
			case <-done:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-startPause.Keydown():
				if !running {
					running = true
					startTime = time.Now()
				} else {
					running = false
					elapsed += time.Since(startTime)
					file.WriteString(writeTimeToFile(elapsed))
				}
			case <-reset.Keydown():
				running = false
				file.WriteString("\n")
				elapsed = 0
			case <-done:
				return
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
	fmt.Println("\nExiting.")
	os.Exit(0)
}

func printTime(d time.Duration) {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	fmt.Printf("\r%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func writeTimeToFile(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	time := fmt.Sprintf("\r%02d:%02d.%03d\n", minutes, seconds, milliseconds)
	return time
}
