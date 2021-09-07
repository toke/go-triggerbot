package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

type Gossip struct {
	match     regexp.Regexp
	text      string
	parseMode string
}

type Config struct {
	Telegram struct {
		Token   string `yaml:"token"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"telegram"`
	Trigger []struct {
		Match     string `yaml:"match"`
		Text      string `yaml:"text"`
		ParseMode string `yaml:"parseMode"`
	} `yaml:"trigger"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func main() {
	var gossip []Gossip
	var cfg Config
	readFile(&cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	for k, _ := range cfg.Trigger {
		log.Printf("Compile Trigger: \"%s\" \t \"%s\"", cfg.Trigger[k].Match, cfg.Trigger[k].Text)
		gossip = append(gossip, Gossip{match: *regexp.MustCompile(cfg.Trigger[k].Match), text: cfg.Trigger[k].Text})
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Telegram.Timeout

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if bot.Debug == true {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "No Help for this Bot"
			default:
				log.Printf("Unknown command: %s", update.Message.Command())
				//msg.Text = "I don't know that command"
			}
			if msg.Text != "" {
				bot.Send(msg)
			}
		} else {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			for k, v := range gossip {
				if bot.Debug == true {
					log.Printf("%s: %s", k, v)
				}
				if gossip[k].match.MatchString(update.Message.Text) {
					log.Printf("%s", v.text)
					msg.Text = v.text
					if v.parseMode == "" {
						msg.ParseMode = "markdown"
					} else {
						msg.ParseMode = v.parseMode
					}
					if msg.Text != "" {
						bot.Send(msg)
					}
				}
			}

		}
	}
}
