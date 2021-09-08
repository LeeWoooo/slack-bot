package parser

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slack-bot/internal/model"

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
	info := doc.Find(".exchange_info")
	date, bank := getBackandDate(info)

	log.Println("date,bank", date, bank)
	return nil, nil
}

func getBackandDate(selection *goquery.Selection) (string, string) {
	// find date
	date := selection.Find(".date").Text()

	//find bank
	bank := selection.Find(".standard").Text()

	return date, bank
}
