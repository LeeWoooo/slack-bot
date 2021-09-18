package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"slack-bot/internal/bot"
	"slack-bot/internal/parser"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// create exchanger
	exchanger := parser.NewExchangeRate()

	// create slackbot
	bot := bot.NewSlackBot(exchanger)

	// create cron
	l, _ := time.LoadLocation("Asia/Seoul")
	c := cron.New(cron.WithLocation(l))

	// every year every month Mon~Fri AM 9 : 15
	c.AddFunc("15 9 * * MON-FRI", func() {
		err := bot.SendTransfer()
		if err != nil {
			log.Fatal(err)
		}
	})

	// every 10m request
	c.AddFunc("@every 10m", func() {
		bot.PreventSleeping()
	})

	c.Start()

	//for heroku
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

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
