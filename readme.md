# Telegram Triggerbot

A Basic Telegram trigger bot.
Triggers can be configured as [go regular expressions](https://pkg.go.dev/regexp/syntax) in a configuration file.
Currently only text messages can be triggered.

[![Go](https://github.com/toke/go-triggerbot/actions/workflows/go.yml/badge.svg)](https://github.com/toke/go-triggerbot/actions/workflows/go.yml)

## Compiling

[Download](https://golang.org/) and install it.
In the Source code directory run: `go build -v .`
Done - proceed withe Section Usage.

## Usage

Use `config-example.yaml` as a template and place save it as `config.yml` in the same Directory as the binary.
A different name and location can be set via the command line parameter *-config*

Make sure to insert The Bot TOKEN into the config file. Use [@BotFather](https://telegram.me/BotFather) to create a bot
and a TOKEN. Make sure to disable privacy mode.

```yaml
telegram:
  token: "YOUR BOTTOKEN"
  timeout: 60

trigger:
  - # A normal Case sensitive Match
    match: "WTF?"
    text: "What the Fuck?"
  - # Case insensitive Match
    match: "(?i)guten morgen"
    text: "Nimm Dir nen Kaffee! ☕️"
  - # This will match "Die 5 Freunde" but not "Die fünf Freunde"
    match: "[Dd]ie \\d+ Freunde"
    text: "... und Timmy der Hund"
```
*NOTE:* The Configuration file uses YAML as format.
Indentation matters!
Also make sure not to mix Tabs and Spaces when doing indentation.