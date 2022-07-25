package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

type Gossip struct {
	Match      regexp.Regexp
	Text       string
	ParseMode  string
	Percentage int
}

type LimitBucket struct {
	Timeout   time.Duration
	Timestamp []time.Time
}

type Limits struct {
	Bucket []LimitBucket
}

type Config struct {
	Telegram struct {
		Token   string `yaml:"token"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"telegram"`
	Limit []struct {
		Bucket  string `yaml:"bucket"`
		Limit   int    `yaml:"limit"`
		BucketS time.Duration
	} `yaml:"limits"`
	Default struct {
		ShutupDisabled bool   `yaml:"shutup_disabled"`
		ShutupTime     string `yaml:"shutup_time"`
	} `yaml:"defaults"`
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

func (limits *Limits) update(cfg *Config) bool {
	ret := true
	for l, _ := range limits.Bucket {
		timeout := cfg.Limit[l].BucketS
		limit := cfg.Limit[l].Limit
		limits.Bucket[l].swipe(timeout)
		if limits.Bucket[l].enforce(limit) == false {
			ret = false
		}
		if ret == true {
			limits.Bucket[l].Timestamp = append(limits.Bucket[l].Timestamp, time.Now())
		}
	}
	return ret
}

func (bucket *LimitBucket) swipe(timeout time.Duration) {
	var tsa []time.Time
	for t, _ := range bucket.Timestamp {
		td := bucket.Timestamp[t].Add(bucket.Timeout)

		if time.Now().Before(td) {
			tsa = append(tsa, bucket.Timestamp[t])
		}
	}
	bucket.Timestamp = tsa
}

func (bucket *LimitBucket) enforce(limit int) bool {
	if len(bucket.Timestamp) >= limit {
		log.Printf("Too many requests in Bucket (%s)", bucket.Timeout.String())
		return false
	}
	return true
}

func main() {
	var gossip []Gossip
	var limit Limits
	var cfg Config
	var shutupEnd time.Time

	var fname string
	flag.StringVar(&fname, "config", "triggerbot.yaml", "Configuration Filename")
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
	for l, _ := range cfg.Limit {
		cfg.Limit[l].BucketS, _ = time.ParseDuration(cfg.Limit[l].Bucket)
		limit.Bucket = append(limit.Bucket, LimitBucket{
			Timeout: cfg.Limit[l].BucketS,
		})
		log.Printf("Bucket: %s", cfg.Limit[l].BucketS.String())
	}

	bot.Debug = *debugPtr

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Telegram.Timeout

	updates := bot.GetUpdatesChan(u)

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
			case "shutup":
				if cfg.Default.ShutupDisabled {
					break
				}
				var sut time.Duration
				args := update.Message.CommandArguments()
				if args == "off" {
					sut = 0
				} else if len(args) == 0 {
					sut, _ = time.ParseDuration(cfg.Default.ShutupTime)
				} else {
					sut, _ = time.ParseDuration(args)
				}
				shutupEnd = time.Now().Add(sut)
				log.Printf("Shut up for %s requested by %s [%d]", sut.String(), update.Message.From.String(), update.Message.From.ID)
			default:
				log.Printf("Unknown command: %s", update.Message.Command())
				//msg.Text = "I don't know that command"
			}
			if msg.Text != "" {
				bot.Send(msg)
			}
		} else {
			if shutupEnd.After(time.Now()) {
				continue
			}
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

							if limit.update(&cfg) {
								bot.Send(msg)
							}
						} else if v.Percentage == 0 {

							if limit.update(&cfg) {
								bot.Send(msg)
							}
						}
					}
				}
			}

		}
	}
}
