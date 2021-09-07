# Telegram Triggerbot

A Basic Telegram trigger bot.
Triggers can be configured as [go regular expressions](https://pkg.go.dev/regexp/syntax) in a configuration file.
Currently only text messages can be triggered.

## Usage

Use 'config-example.yaml' as a template and place save it as 'config.yml' in the same Directory as the binary.
A different name and location can be set via the command line parameter *-config*

Make sure to insert The Bot TOKEN into the config file. Use [@BotFather](https://telegram.me/BotFather) to create a bot
and a TOKEN.

'''
telegram:
  token: "YOURBOTTOKEN"
  timeout: 60

# Match syntax https://pkg.go.dev/regexp/syntax
trigger:
  - match: "WTF?"
    text: "What the Fuck?"
  - match: "(?i)guten morgen"
    text: "Lass mich weiterschlafen!"
  - match: "[Dd]ie \d+ Freunde"
    text: "... und Timmy der Hund"
'''

