package entity

import (
	"time"
)

type Account struct {
	Id        int64     `json:"id"`
	Balance   float64   `json:"balance"`
	CreatedDt time.Time `json:"created_dt"`
}
