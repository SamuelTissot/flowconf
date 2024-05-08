package test

import "time"

type Config struct {
	MeaningOfLife   int
	Cats            []string
	Pi              float64
	Perfection      []int
	BackToTheFuture time.Time
	Secret          string
	Tag             string `json:"TagValue" toml:"TagValue"`
}
