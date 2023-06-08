# miau

A simple CLI for controlling Elgato Key Light.
Tested with a single Elgato Key Light Mini.

## Installation
Run `make` to compile `miau` and generate completions.
Run `make install` to install `miau` and completions to `/usr/local/`.
Clean up with `make clean` and `make uninstall`, respectively.

To override `/usr/local/` PREFIX variable use `make -e PREFIX=/foo/bar/baz/`.

## Usage
```
miau - Control Elgato Key Light from CLI

Usage:
  miau [flags]
  miau [command]

Available Commands:
  help        Help about any command
  off         Turn the light OFF
  on          Turn the light ON
  serve       Serve a webpage for controlling Elgato Key Light
  settings    Update light's settings

Flags:
  -a, --addr ip            IP address of the light
  -b, --brightness int     brightness in percent; a value between 0 and 100
  -h, --help               help for miau
  -t, --temperature int    temperature in Kelvins; a value between 2900 and 7000
      --timeout duration   light API timeout duration (default 15s)
  -v, --version            version for miau

Use "miau [command] --help" for more information about a command.
```

### Serve
This CLI comes with a command to serve a simple webpage.
The webpage lists current settings and allows turning the light on or off.
```
Serve a webpage for controlling Elgato Key Light

Usage:
  miau serve [flags]

Flags:
  -h, --help       help for serve
  -p, --port int   port to listen on (default 9123)

Global Flags:
  -a, --addr ip            IP address of the light
  -b, --brightness int     brightness in percent; a value between 0 and 100
  -t, --temperature int    temperature in Kelvins; a value between 2900 and 7000
      --timeout duration   light API timeout duration (default 15s)
```

### Settings
Currently only updating battery bypass setting is supported.

```
Running settings command without any flags will fetch and print the current settings.

Usage:
  miau settings [flags]

Flags:
      --bypass int   configure the light's battery bypass (0 or 1)
  -h, --help         help for settings

Global Flags:
  -a, --addr ip            IP address of the light
  -b, --brightness int     brightness in percent; a value between 0 and 100
  -t, --temperature int    temperature in Kelvins; a value between 2900 and 7000
      --timeout duration   light API timeout duration (default 5s)
```

## Name
Elgato is the company that produces the equipment.
"el gato" means "cat" in Spanish.
Cat goes "miau" in Polish.
