package model

// ExchangeRate for slack ExchangeRate Message
type ExchangeRate struct {
	Date        string //slack message를 보내는 시점의 날짜.
	Bank        string //  환율의 정보를 제공해주는 은행명
	KRW         string // slack message를 보내는 시점의 환율(원화)
	DtD         string // 전일 대비 환율 증가 감소 DATA
	TransferKWR string // 송금보낼 때의 가격
	Preference  string // Preference 우대 환율 KWR
	ImageURL    string // 한달 환율 변동 그래프 URL
}
