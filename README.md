# Elgato CLI

A simple CLI for controlling Elgato Key Light.
Tested with a single Elgato Key Light Mini.

## Usage
```
Control Elgato light

Usage:
  elgato [flags]
  elgato [command]

Available Commands:
  help        Help about any command
  off         Turn Elgato light OFF
  on          Turn Elgato light ON
  serve       Serve a webpage for controlling Elgato Light

Flags:
  -a, --addr ip           IP address of the light
  -b, --brightness int    brightness in percent; a value between 0 and 100
  -h, --help              help for elgato
  -t, --temperature int   temperature in Kelvins; a value between 2900 and 7000

Use "elgato [command] --help" for more information about a command.
```
