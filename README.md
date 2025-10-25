# gowatch

## about

**Gowatch** is a simple, multi-platform, terminal-based stopwatch that
is controlled by global hotkeys. It works with Linux (X11 only), MacOS, and
Windows (native, not WSL).

This project is an very simple alternative to
something like [LiveSplit](https://livesplit.org/), but much lighter-weight. I
built this as I couldn't find a good stopwatch tool to use in DBD 1v1s on Linux,
and other solutions never seemed to grab the hotkeys from my system while I was
in-game. So, rather than try and figure out why DWM wouldn't pass these keys
through, I spent an afternoon making this.

## requirements

You will need Go to build the project. You can find installation instructions
[here](https://go.dev/doc/install).

## installation

Clone the repo and `cd` into the directory:

```sh
git clone https://github.com/buxxket/gowatch.git
cd gowatch
```

### unix (Linux and MacOS)

Build the project, symlink the binary to your path, and copy the default config
file to `~/.config/gowatch/config.yaml`:

```sh
# Build the project
go build .

# Create a symbolic link to the binary in ~/bin
ln -s "$(pwd)/gowatch" ~/bin/

# Create the configuration directory
mkdir -p ~/.config/gowatch

# Copy the default configuration file
cp "$(pwd)/config.yaml.default" ~/.config/gowatch/config.yaml
```

## bad operating systems (Windows)

```powershell
# Build the project
go build .

# Copy the binary to a directory already in the PATH
copy .\gowatch.exe "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\gowatch.exe"

# Ensure the config directory exists in %APPDATA%
$APPDATA_PATH = "$env:APPDATA\gowatch"
mkdir $APPDATA_PATH

# Copy the default config to the config directory
copy .\config.yaml.default "$APPDATA_PATH\config.yaml"
```

## usage

Run `gowatch`.

Default keybinds are `Alt + W` to start/pause the timer (good for starting DBD
1v1s), and `Alt + R` to reset the timer.

Each time you pause or reset the timer, the elapsed time is written to
`/tmp/gowatch`, just in case you didn't catch the timestamp before you cleared
it.

## customisation

The supported *modifiers* are:
```
CMD, ALT, WIN, SUPER, CTRL, OPTION, SHIFT
```
The supported *keys* are:
```
0-1, A-Z
```

## todo
- [ ] create an AUR package
- [ ] add splits
