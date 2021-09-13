package bot

import (
	"fmt"
	"log"
	"slack-bot/internal/parser"
	"testing"

	"github.com/slack-go/slack"
)

const format = "%s %s 기준 환율 보고 드립니다.\n1$당 KWR(원화)는 %s이며 %s 입니다.\n 해외 송금 기준 %s 입니다.(우대 환율 적용)\n"

func TestSendMessage(t *testing.T) {
	//
	api := slack.New("")

	td, _ := parser.NewExchangeRate().GetExchangerRate()

	text := fmt.Sprintf(format, td.Date, td.Bank, td.KRW, td.DtD, td.TransferKWR)

	att := slack.Attachment{
		Text:     text,
		ImageURL: td.ImageURL,
	}

	channelID, timeStamp, err := api.PostMessage(
		"",
		slack.MsgOptionText("", false),
		slack.MsgOptionAttachments(att),
		slack.MsgOptionAsUser(false),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timeStamp)
}
