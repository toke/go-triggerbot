package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

type Gossip struct {
	Match      regexp.Regexp
	Text       string
	ParseMode  string
	Percentage int
}

type Config struct {
	Telegram struct {
		Token   string `yaml:"token"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"telegram"`
	Trigger []struct {
		Match      string `yaml:"match"`
		Text       string `yaml:"text"`
		ParseMode  string `yaml:"parseMode"`
		Percentage int    `yaml:"percentage"`
	} `yaml:"trigger"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config, filename string) {
	f, err := os.Open(filename)
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

	var fname string
	flag.StringVar(&fname, "config", "config.yaml", "Configuration Filename")
	debugPtr := flag.Bool("debug", false, "Debug Output")
	flag.Parse()
	fmt.Println("Loading config:", fname)

	readFile(&cfg, fname)

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	for k, _ := range cfg.Trigger {
		if *debugPtr == true {
			log.Printf("Compile Trigger: \"%s\" \t \"%s\"", cfg.Trigger[k].Match, cfg.Trigger[k].Text)
		}
		gossip = append(gossip, Gossip{
			Match:      *regexp.MustCompile(cfg.Trigger[k].Match),
			Text:       cfg.Trigger[k].Text,
			Percentage: cfg.Trigger[k].Percentage})
	}

	bot.Debug = *debugPtr

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Telegram.Timeout

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if *debugPtr == true {
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

				if gossip[k].Match.MatchString(update.Message.Text) {
					if *debugPtr == true {
						log.Printf("%s", v.Text)
					}
					msg.Text = v.Text
					if v.ParseMode == "" {
						msg.ParseMode = "markdown"
					} else {
						msg.ParseMode = v.ParseMode
					}
					if msg.Text != "" {
						if v.Percentage > 0 && rand.Intn(100) < v.Percentage {
							bot.Send(msg)
						}
					}
				}
			}

		}
	}
}
