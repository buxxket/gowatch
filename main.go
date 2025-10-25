package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/spf13/viper"
	"golang.design/x/hotkey"
)

type HotkeyConfig struct {
	StartPause string `mapstructure:"startpause"`
	Reset      string `mapstructure:"reset"`
}

type AppConfig struct {
	Hotkeys HotkeyConfig `mapstructure:"hotkeys"`
}

func printTime(d time.Duration) {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	fmt.Printf("\r%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func writeElapsedToFile(elapsed time.Duration) {
	var path string
	if runtime.GOOS == "windows" {
		path = os.TempDir() + "\\gowatch"
	} else {
		path = "/tmp/gowatch"
	}
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%02d:%02d.%03d\n", int(elapsed.Minutes()), int(elapsed.Seconds())%60, int(elapsed.Milliseconds())%1000)
}

func main() {
	configPath, err := ConfigFilePath()
	if err != nil {
		log.Fatalf("Error determining config file path: %v", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		fmt.Println("You can copy the default config file with:")
		fmt.Println("cp /usr/share/gowatch/config.yaml $HOME/.config/gowatch/config.yaml")
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	startPauseMods, startPauseKey, err := ParseHotkeyString(config.Hotkeys.StartPause)
	if err != nil {
		log.Fatalf("Invalid StartPause hotkey: %v", err)
	}
	resetMods, resetKey, err := ParseHotkeyString(config.Hotkeys.Reset)
	if err != nil {
		log.Fatalf("Invalid Reset hotkey: %v", err)
	}

	startPause := hotkey.New(startPauseMods, startPauseKey)
	reset := hotkey.New(resetMods, resetKey)

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
					writeElapsedToFile(elapsed)
				}
			case <-reset.Keydown():
				running = false
				writeElapsedToFile(elapsed)
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
