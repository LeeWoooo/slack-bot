package main

import (
	"log"
	"os"
	"os/signal"
	"slack-bot/internal/bot"
	"slack-bot/internal/parser"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// create exchanger
	exchanger := parser.NewExchangeRate()

	// create slackbot
	bot := bot.NewSlackBot(exchanger)

	// create cron
	c := cron.New()
	// 매년 매월 월요일~금요일 아침 9시 15분
	c.AddFunc("59 0 * * MON-FRI", func() {
		err := bot.SendTransfer()
		if err != nil {
			log.Fatal(err)
		}
	})
	c.Start()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	logrus.Info("starting slack bot...")
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down slack bot...")
}
