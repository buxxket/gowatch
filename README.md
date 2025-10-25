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

Build the project, symlink the binary to your path, and copy the default config
file to `~/.config/gowatch/config.yaml`:

```sh
go build .
ln -s "$(pwd)/gowatch" ~/bin/
mkdir -p ~/.config/gowatch
cp "$(pwd)/config.yaml.default" ~/.config/gowatch/config.yaml
```

## usage

Run `gowatch`.

Default keybinds are `Alt + W` to start/pause the timer (good for starting DBD
1v1s), and `Alt + R` to reset the timer.

Each time you pause or reset the timer, the elapsed time is written to
`/tmp/gowatch`, just in case you didn't catch the timestamp before you cleared
it.

## todo
- [ ] create an AUR package
- [ ] add splits
