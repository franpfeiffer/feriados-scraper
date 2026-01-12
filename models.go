package main

import "time"

type Feriado struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}
