package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Quotes struct {
	ID     int
	Author string `json:"author"`
	Text   string `json:"quote"`
}

type Quote struct {
	Author string `json:"author"`
	Text   string `json:"quote"`
}
