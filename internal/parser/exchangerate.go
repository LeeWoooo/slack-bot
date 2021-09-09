package parser

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slack-bot/internal/model"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/sirupsen/logrus"
)

const (
	exchangeRateURL = "https://finance.naver.com/marketindex/exchangeDetail.nhn?marketindexCd=FX_USDKRW_SHB"
)

// ExchangerRate interface
type ExchangerRate interface {
	GetExchangerRate() (*model.ExchangeRate, error)
}

// ExchangeRateImpl implement ExchangeRate
type ExchangeRateImpl struct{}

// NewExchangeRate create new ExchangeRate instance
func NewExchangeRate() ExchangerRate {
	return &ExchangeRateImpl{}
}

// GetExchangerRate implement GetExchangerRate
func (e *ExchangeRateImpl) GetExchangerRate() (*model.ExchangeRate, error) {
	// get html
	resp, err := http.Get(exchangeRateURL)
	if err != nil {
		logrus.Errorf("could not get response error:%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// check http status
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("response status code error status:%s, statuscode:%d", resp.Status, resp.StatusCode)
		logrus.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}

	//convert euc-kr to utf-8
	utfBody, err := iconv.NewReader(resp.Body, "euc-kr", "utf-8")
	if err != nil {
		logrus.Errorf("could not convert euc.kr to utf-8 error:%v", err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		logrus.Errorf("could not get document error:%v", err)
		return nil, err
	}

	// 은행 정보 및 날짜
	date, bank := getBackandDate(doc)
	log.Println("date,bank", date, bank)

	// KRW
	KRW := getKRW(doc)
	log.Println("KRW", KRW)

	// prev compareData
	compareData := getPrevDayCompreData(doc)
	log.Println("compareData", compareData)

	return nil, nil
}

// getBackandDate get bank info and date
func getBackandDate(doc *goquery.Document) (string, string) {
	selection := doc.Find(".exchange_info")

	// find date
	date := selection.Find(".date").Text()

	//find bank
	bank := selection.Find(".standard").Text()

	return date, bank
}

func getKRW(doc *goquery.Document) string {
	selection := doc.Find(".no_today")

	// get KRW (숫자)
	// remove line break
	krwSelection := selection.Find(".no_up > .no_up")
	krw := strings.ReplaceAll(krwSelection.Text(), "\n", "")

	// trim
	krw = strings.Trim(krw, " ")

	// 단위
	unit := selection.Find(".txt_won").Text()

	return krw + unit
}

func getPrevDayCompreData(doc *goquery.Document) string {
	selection := doc.Find(".no_exday")

	// get text
	text := selection.Find(".txt_comparison").Text()
	log.Println("text", text)

	var compareData string

	// get compare prev
	selection.Find(".no_up").Each(func(i int, s *goquery.Selection) {
		compareData += s.Text()
	})

	// remove line break
	compareData = strings.ReplaceAll(compareData, "\n", "")

	// remove white space
	compareData = strings.ReplaceAll(compareData, " ", "")

	// trim and return
	return strings.Trim(compareData, " ")
}
