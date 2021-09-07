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
and a TOKEN.

```
telegram:
  token: "YOUR BOTTOKEN"
  timeout: 60

trigger:
  - match: "WTF?"
    text: "What the Fuck?"
  - match: "(?i)guten morgen"
    text: "Lass mich weiterschlafen!"
  - match: "[Dd]ie \d+ Freunde"
    text: "... und Timmy der Hund"
```
