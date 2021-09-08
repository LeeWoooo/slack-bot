package model

// ExchangeRate for slack ExchangeRate Message
type ExchangeRate struct {
	Date      string //slack message를 보내는 시점의 날짜.
	DtD       string // 전일 대비 환율 증가 감소 percent
	KRW       string // slack message를 보내는 시점의 환율(원화)
	MonthHigh string // slack message를 보내는 시점의 한달의 최고
	MonthLow  string // slack message를 보내는 시점의 한달의 최저
}
