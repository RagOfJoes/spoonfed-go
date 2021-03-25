package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const (
	// length stores maximum slug length.
	// It's smart so it will cat slug after full word.
	// By default slugs aren't shortened.
	// If length is smaller than length of the first word, then returned
	// slug will contain only substring from the first word truncated
	// after length.
	length int = 64
	// SlugIDLength defines the length of ID to generate
	// for slug
	SlugIDLength int = 12
	// Defines valid letters to use to generate
	// id
	slugIDAlphabet string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Defines 6 bits to represent a letter index
	slugIDidxBits = 6
	// Defines all 1-bits, as many as slugIDidxBits
	slugIDidxMask = 1<<slugIDidxBits - 1
	// Defines # of letter indices fitting in 63 bits
	slugIDidxMax = 63 / slugIDidxBits
)

var (
	// CustomSub stores custom substitution map
	CustomSub map[string]string
	// CustomRuneSub stores custom rune substitution map
	CustomRuneSub map[rune]string

	// Lowercase defines if the resulting slug is transformed to lowercase.
	// Default is true.
	Lowercase = true

	regexpNonAuthorizedChars = regexp.MustCompile("[^a-zA-Z0-9-_]")
	regexpMultipleDashes     = regexp.MustCompile("-+")
)

var nanoTS = rand.NewSource(time.Now().UnixNano())

// Slug generates a unique slug given an initial string
func Slug(s string) (slug string) {
	slug = strings.TrimSpace(s)

	if Lowercase {
		slug = strings.ToLower(slug)
	}

	// Process all remaining symbols
	slug = regexpNonAuthorizedChars.ReplaceAllString(slug, "-")
	slug = regexpMultipleDashes.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-_")

	if length > 0 {
		slug = SmartTruncate(slug, length)
	}

	slug += fmt.Sprintf("-%s", SlugID(SlugIDLength))

	return slug
}

// SlugID generates secure URL-friendly unique ID.
// Refer: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func SlugID(length int) string {
	b := make([]byte, length)
	for i, cache, remain := length-1, nanoTS.Int63(), slugIDidxMax; i >= 0; {
		if remain == 0 {
			cache, remain = nanoTS.Int63(), slugIDidxMax
		}
		if idx := int(cache & slugIDidxMask); idx < len(slugIDAlphabet) {
			b[i] = slugIDAlphabet[idx]
			i--
		}
		cache >>= slugIDidxBits
		remain--
	}
	return string(b)
}

// SmartTruncate truncates a string down to given length
func SmartTruncate(text string, length int) string {
	l := length
	if length < 0 {
		l = length
	}
	if len(text) < l {
		return text
	}

	var truncated string
	words := strings.SplitAfter(text, "-")
	// If length is smaller than length of the first word return word
	// truncated after length.
	if len(words[0]) > l {
		return words[0][:l]
	}
	for _, word := range words {
		if len(truncated)+len(word)-1 <= l {
			truncated = truncated + word
		} else {
			break
		}
	}
	return strings.Trim(truncated, "-")
}
