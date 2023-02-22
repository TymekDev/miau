# miau

A simple CLI for controlling Elgato Key Light.
Tested with a single Elgato Key Light Mini.

## Usage
```
Control Elgato Key Light

Usage:
  miau [flags]
  miau [command]

Available Commands:
  help        Help about any command
  off         Turn the light OFF
  on          Turn the light ON
  serve       Serve a webpage for controlling Elgato Key Light

Flags:
  -a, --addr ip           IP address of the light
  -b, --brightness int    brightness in percent; a value between 0 and 100
  -h, --help              help for miau
  -t, --temperature int   temperature in Kelvins; a value between 2900 and 7000

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
  -a, --addr ip           IP address of the light
  -b, --brightness int    brightness in percent; a value between 0 and 100
  -t, --temperature int   temperature in Kelvins; a value between 2900 and 7000
```

## Name
Elgato is the company that produces the equipment.
"el gato" means "cat" in Spanish.
Cat goes "miau" in Polish.
