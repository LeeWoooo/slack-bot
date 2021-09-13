package bot

import (
	"fmt"
	"log"
	"slack-bot/internal/parser"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// SlackBot interface
type SlackBot interface {
	SendTransfer() error
}

// SlackBotImpl implement SlackBot
type SlackBotImpl struct {
	exchanger parser.ExchangerRate
}

// NewSlackBot create New SlackBot instance
func NewSlackBot(exchanger parser.ExchangerRate) SlackBot {
	return &SlackBotImpl{exchanger: exchanger}
}

// SendTransfer send slack message what daily transferRate
func (s *SlackBotImpl) SendTransfer() error {
	td, err := s.exchanger.GetExchangerRate()
	if err != nil {
		return err
	}

	// make text with format
	format := "%s %s 기준 환율 보고 드립니다.\n1$당 KWR(원화)는 %s이며 %s 입니다.\n 해외 송금 기준 %s 입니다.(우대 환율 적용)\n"
	text := fmt.Sprintf(format, td.Date, td.Bank, td.KRW, td.DtD, td.TransferKWR)

	attachment := slack.Attachment{
		Title:    "Daily TransferRate",
		Text:     text,
		ImageURL: td.ImageURL,
	}

	// get slack api
	api := slack.New("xoxb-2249906725653-2465451160295-cXK8yCQ3zWz3BL0pVK1S3Obn")

	channelID, timeStamp, err := api.PostMessage(
		"C028ZJJJ2TY",
		slack.MsgOptionText("", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	logrus.Infof("Message successfully sent to channel %s at %s", channelID, timeStamp)
	return nil
}
