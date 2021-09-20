package parser

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/stretchr/testify/assert"
)

func TestGetHtml(t *testing.T) {
	assert := assert.New(t)

	//givne
	//when
	resp, err := http.Get(exchangeRateURL)

	//then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestGetExchangerRate(t *testing.T) {
	assert := assert.New(t)

	//givne
	//when
	_, err := NewExchangeRate().GetExchangerRate()

	//then
	assert.NoError(err)
}

func TestIntegrated(t *testing.T) {
	assert := assert.New(t)

	//prepared
	resp, err := http.Get(exchangeRateURL)
	assert.NoError(err)
	defer resp.Body.Close()
	assert.Equal(http.StatusOK, resp.StatusCode)

	//convert euc-kr to utf-8
	utfBody, err := iconv.NewReader(resp.Body, "euc-kr", "utf-8")
	assert.NoError(err)

	doc, err := goquery.NewDocumentFromReader(utfBody)
	assert.NoError(err)

	t.Run("Get Bank and Data", func(t *testing.T) {
		// given = doc

		// when
		date, bank := getBackandDate(doc)

		// then (bank)
		assert.Equal("신한은행", bank)

		// then (date)
		format := "2006.01.02 15:04"

		// 넘어오는 데이터 format ex) 2021.09.17 18:00
		// 토요일 일요일일 경우 날짜를 금요일 기준으로 변경해야 함
		now := time.Now()

		switch now.Weekday() {
		case time.Saturday:
			now = now.Add(time.Hour * 24 * -1)
		case time.Sunday:
			now = now.Add(time.Hour * 24 * -2)
		}
		ft, err := time.Parse(format, date)
		assert.NoError(err)

		// year
		assert.Equal(now.Year(), ft.Year())

		// month
		assert.Equal(now.Month(), ft.Month())

		// day
		assert.Equal(now.Day(), ft.Day())

		// weekday
		assert.Equal(now.Weekday(), ft.Weekday())
	})

	t.Run("Get KRW", func(t *testing.T) {
		//givne = doc

		//when
		KRW := getKRW(doc)

		//then
		arr := strings.Split(KRW, ".")
		assert.Equal(2, len(arr))

		integer := arr[0]
		decimal := arr[1]

		// 정수 부분의 길이가 5일 경우 (x,xxx)
		if len(integer) == 5 {
			assert.Equal(",", string(integer[1]))
		}
		assert.Equal(2, len(decimal))
	})

	t.Run("Get Prev", func(t *testing.T) {
		// given = doc

		// when
		prev := getPrevDayCompreData(doc)

		t.Log(prev)
	})
}
