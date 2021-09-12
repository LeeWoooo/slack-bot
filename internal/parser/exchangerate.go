package parser

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slack-bot/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/dustin/go-humanize"
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

	rate := &model.ExchangeRate{}

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

	rate.Date, rate.Bank = date, bank

	// KRW
	KRW := getKRW(doc)
	log.Println("KRW", KRW)
	rate.KRW = KRW

	// prev compareData
	compareData := getPrevDayCompreData(doc)
	log.Println("compareData", compareData)
	rate.DtD = compareData

	// transfer
	transferKWR := getTransferKWR(doc)
	log.Println("transferKWR", transferKWR)
	rate.TransferKWR = transferKWR

	URL := getGraphURL(doc)
	log.Println("URL", URL)
	rate.ImageURL = URL

	return rate, nil
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
	krwSelection := selection.Find("em > em")
	krw := strings.Replace(krwSelection.Text(), "\n", "", -1)

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

	var compareData string

	// get compare prev
	selection.Find("em").Each(func(i int, s *goquery.Selection) {
		compareData += s.Text()
	})

	// remove line break
	compareData = strings.Replace(compareData, "\n", "", -1)

	// remove white space
	compareData = strings.Replace(compareData, " ", "", -1)

	// trim and return
	compareData = strings.Trim(compareData, " ")

	// append "원"
	compareData = strings.Replace(compareData, "(", "원(", -1)

	// append sign
	if strings.Contains(compareData, "-") {
		return text + " -" + compareData
	}

	return text + " +" + compareData
}

func getTransferKWR(doc *goquery.Document) string {
	selection := doc.Find(".th_ex4")

	// get transferKWR
	KWR := selection.Next().Text()

	// get Preference
	Preference, _ := getPreference(KWR)

	// process string && return
	return Preference
}

func getPreference(KWR string) (string, error) {
	// separation essence, decimal
	arr := strings.Split(KWR, ".")
	essenceString := strings.Replace(arr[0], ",", "", -1)
	decimalString := arr[1]

	// get Preference
	essence, err := strconv.ParseInt(essenceString, 10, 64)
	if err != nil {
		logrus.Errorf("could not parse int essence error:%v", err)
		return "", nil
	}
	essence -= 6
	essenceString = humanize.Comma(essence)

	return essenceString + "." + decimalString + "원", nil
}

func getGraphURL(doc *goquery.Document) string {
	selection := doc.Find(".flash_area > img")

	//get graph URL
	URL, isEixst := selection.Attr("src")
	if !isEixst {
		logrus.Errorf("could not load imagePath")
		return ""
	}

	// // 57 == month string index
	// URL = strings.Replace(URL, "month3", "month", 57)

	//processing URL (Add query String)
	return URL + "?sidcode=" + string(time.Now().Format("20060102"))
}
