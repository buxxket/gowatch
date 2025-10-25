package main

import (
	"fmt"
	"strings"

	"golang.design/x/hotkey"
)

func ParseHotkeyString(hotkeyString string) ([]hotkey.Modifier, hotkey.Key, error) {
	parts := strings.Split(hotkeyString, "+")
	var mods []hotkey.Modifier
	var key hotkey.Key
	for _, part := range parts {
		switch strings.ToUpper(part) {
		case "MOD1":
			mods = append(mods, hotkey.Mod1)
		case "MOD2":
			mods = append(mods, hotkey.Mod2)
		case "MOD3":
			mods = append(mods, hotkey.Mod3)
		case "MOD4":
			mods = append(mods, hotkey.Mod4)
		case "MOD5":
			mods = append(mods, hotkey.Mod5)
		case "CTRL":
			mods = append(mods, hotkey.ModCtrl)
		case "ALT":
			mods = append(mods, hotkey.Mod1)
		case "SHIFT":
			mods = append(mods, hotkey.ModShift)
		case "WIN":
			mods = append(mods, hotkey.Mod4)
		default:
			k := "Key" + strings.ToUpper(part)
			keyMap := map[string]hotkey.Key{
				"KeyA": hotkey.KeyA,
				"KeyB": hotkey.KeyB,
				"KeyC": hotkey.KeyC,
				"KeyD": hotkey.KeyD,
				"KeyE": hotkey.KeyE,
				"KeyF": hotkey.KeyF,
				"KeyG": hotkey.KeyG,
				"KeyH": hotkey.KeyH,
				"KeyI": hotkey.KeyI,
				"KeyJ": hotkey.KeyJ,
				"KeyK": hotkey.KeyK,
				"KeyL": hotkey.KeyL,
				"KeyM": hotkey.KeyM,
				"KeyN": hotkey.KeyN,
				"KeyO": hotkey.KeyO,
				"KeyP": hotkey.KeyP,
				"KeyQ": hotkey.KeyQ,
				"KeyR": hotkey.KeyR,
				"KeyS": hotkey.KeyS,
				"KeyT": hotkey.KeyT,
				"KeyU": hotkey.KeyU,
				"KeyV": hotkey.KeyV,
				"KeyW": hotkey.KeyW,
				"KeyX": hotkey.KeyX,
				"KeyY": hotkey.KeyY,
				"KeyZ": hotkey.KeyZ,
				"Key1": hotkey.Key1,
				"Key2": hotkey.Key2,
				"Key3": hotkey.Key3,
				"Key4": hotkey.Key4,
				"Key5": hotkey.Key5,
				"Key6": hotkey.Key6,
				"Key7": hotkey.Key7,
				"Key8": hotkey.Key8,
				"Key9": hotkey.Key9,
				"Key0": hotkey.Key0,
			}
			if val, ok := keyMap[k]; ok {
				key = val
			} else {
				return nil, 0, fmt.Errorf("unknown key: %s", part)
			}
		}
	}
	return mods, key, nil
}
