package entity

import (
	"time"
)

type Transaction struct {
	Id        int64     `json:"id"`
	TransDt   time.Time `json:"trans_dt"`
	AccountId int64     `json:"account_id"`
	DocNum    int64     `json:"doc_num"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
}
