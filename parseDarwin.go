//go:build darwin

package main

import (
	"fmt"
	"golang.design/x/hotkey"
	"strings"
)

func ParseHotkeyString(s string) ([]hotkey.Modifier, hotkey.Key, error) {
	parts := strings.Split(s, "+")
	var mods []hotkey.Modifier
	var key hotkey.Key

	for _, part := range parts {
		switch strings.ToUpper(part) {
		case "CTRL":
			mods = append(mods, hotkey.ModCtrl)
		case "SHIFT":
			mods = append(mods, hotkey.ModShift)
		case "OPTION", "ALT":
			mods = append(mods, hotkey.ModOption)
		case "CMD":
			mods = append(mods, hotkey.ModCmd)
		default:
			keyMap := map[string]hotkey.Key{
				"A": hotkey.KeyA, "B": hotkey.KeyB, "C": hotkey.KeyC, "D": hotkey.KeyD,
				"E": hotkey.KeyE, "F": hotkey.KeyF, "G": hotkey.KeyG, "H": hotkey.KeyH,
				"I": hotkey.KeyI, "J": hotkey.KeyJ, "K": hotkey.KeyK, "L": hotkey.KeyL,
				"M": hotkey.KeyM, "N": hotkey.KeyN, "O": hotkey.KeyO, "P": hotkey.KeyP,
				"Q": hotkey.KeyQ, "R": hotkey.KeyR, "S": hotkey.KeyS, "T": hotkey.KeyT,
				"U": hotkey.KeyU, "V": hotkey.KeyV, "W": hotkey.KeyW, "X": hotkey.KeyX,
				"Y": hotkey.KeyY, "Z": hotkey.KeyZ, "SPACE": hotkey.KeySpace,
			}
			up := strings.ToUpper(part)
			if val, ok := keyMap[up]; ok {
				key = val
			} else {
				return nil, 0, fmt.Errorf("unknown key: %s", part)
			}
		}
	}
	return mods, key, nil
}
