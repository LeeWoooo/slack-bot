package parser

import (
	"net/http"
	"testing"

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
		//given = doc

		//when
		_, bank := getBackandDate(doc)

		//then
		assert.Equal("신한은행", bank)

		//TODO: 얻어온 Date를 format을 이용하여 parsing 후 compare
	})
}
