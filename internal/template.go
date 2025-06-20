package internal

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RandomString(length int) string {
	// Seeds the random number generator with the current time in nanoseconds to ensure a different result each time.
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	str := make([]byte, length)
	for index := range str {
		str[index] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func FormatAsTitle(title string) string {
	title = cases.Title(language.Und, cases.NoLower).String(title)
	formatted := strings.ReplaceAll(title, "_", " ")
	return formatted
}

func FormatAsDate(createdAt time.Time) string {
	return createdAt.Format(time.RFC822)
}
