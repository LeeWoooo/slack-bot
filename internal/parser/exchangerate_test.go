package parser

import (
	"net/http"
	"testing"

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
