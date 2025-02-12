package server

import (
	"math/rand"
)

var quotes = []string{
	"The only way to do great work is to love what you do. - Steve Jobs",
	"Innovation distinguishes between a leader and a follower. - Steve Jobs",
	"Stay hungry, stay foolish. - Steve Jobs",
	"The future belongs to those who believe in the beauty of their dreams. - Eleanor Roosevelt",
	"Success is not final, failure is not fatal: it is the courage to continue that counts. - Winston Churchill",
}

func GetRandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
