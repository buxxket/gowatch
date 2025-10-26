package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.design/x/hotkey"
)

type HotkeyConfig struct {
	StartPause string `mapstructure:"startpause"`
	Reset      string `mapstructure:"reset"`
	Split      string `mapstructure:"split"`
}

type AppConfig struct {
	Hotkeys    HotkeyConfig `mapstructure:"hotkeys"`
	OutputPath string       `mapstructure:"outputpath"`
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "$HOME") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[len("$HOME/"):])
		}
	}
	return path
}

func printTime(d time.Duration) {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	fmt.Printf("\033[2K\r")
	fmt.Printf("\r%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func printSplitToConsole(splitCount int, d time.Duration) {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	fmt.Printf("\033[2K\r")
	fmt.Printf("\r"+splitCounterStringFormatter(splitCount)+":\t%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func writeElapsedToFile(outputFile *os.File, passedText string,
	elapsedTime time.Duration) {
	fmt.Fprintf(outputFile, "%v\t%02d:%02d.%03d\n",
		passedText, int(elapsedTime.Minutes()),
		int(elapsedTime.Seconds())%60,
		int(elapsedTime.Milliseconds())%1000)
}

func writeStartTimeToFile(outputFile *os.File, startTime time.Time) error {
	_, err := fmt.Fprintf(outputFile, "\n\n[STARTTIME]:\t%02d/%02d/%04d - %02d:%02d:%02d\n",
		startTime.Day(),
		startTime.Month(),
		startTime.Year(),
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second())
	return err
}

func splitCounterStringFormatter(count int) string {
	countString := fmt.Sprintf("[SPLIT%v]", count)
	return countString
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
		log.Fatalf("Error reading config file, %s\n", err)
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
	splitMods, splitKey, err := ParseHotkeyString(config.Hotkeys.Split)
	if err != nil {
		log.Fatalf("Invalid Split hotkey: %v", err)
	}

	outputPath := config.OutputPath
	if outputPath == "" {
		if runtime.GOOS == "windows" {
			outputPath = filepath.Join(os.TempDir(), "gowatch")
		} else {
			outputPath = "/tmp/gowatch"
		}
	} else {
		outputPath = expandPath(outputPath)
	}
	outputFile, err := os.OpenFile(outputPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer outputFile.Close()

	startPause := hotkey.New(startPauseMods, startPauseKey)
	reset := hotkey.New(resetMods, resetKey)
	split := hotkey.New(splitMods, splitKey)

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
	if err := split.Register(); err != nil {
		fmt.Println("Failed to register split hotkey:", err)
		return
	}
	defer split.Unregister()

	isTimerRunning := false
	isRunActive := false
	var startTime time.Time
	var elapsed time.Duration
	var splitElapsed time.Duration
	var lastSplit time.Duration = -1 // so first split always writes
	var splitCount int = 0
	var splitDifference time.Duration

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if isTimerRunning {
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
				if !isTimerRunning {
					isTimerRunning = true
					startTime = time.Now()
					if !isRunActive {
						writeStartTimeToFile(outputFile, startTime)
						isRunActive = true
					}
				} else {
					isTimerRunning = false
					elapsed += time.Since(startTime)
					writeElapsedToFile(outputFile, "[PAUSED]:", elapsed)
				}
			case <-reset.Keydown():
				if isRunActive {
					if isTimerRunning {
						elapsed += time.Since(startTime)
					}
					writeElapsedToFile(outputFile, "[FINAL]:", elapsed)
					isTimerRunning = false
					elapsed = 0
					isRunActive = false
				}
			case <-split.Keydown():
				if isTimerRunning {
					splitElapsed = elapsed + time.Since(startTime)
				} else {
					splitElapsed = elapsed
				}
				if lastSplit != splitElapsed {
					splitDifference = splitElapsed - lastSplit
					splitCount++
					writeElapsedToFile(outputFile,
						splitCounterStringFormatter(splitCount)+":", splitElapsed)
					fmt.Println()
					printSplitToConsole(splitCount, splitDifference)
					printTime(splitElapsed)
					lastSplit = splitElapsed
				}
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
