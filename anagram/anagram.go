package anagram

import (
	"strings"
	"unicode"
)

// const maxChars = 36
const maxChars = 16
const keepText = false

// Anagram data
type Anagram struct {
	Length    int
	Mask      uint32
	Counts    [maxChars]uint8
	Text      string
	positions [256]int
}

// Word represents an analyzed word
type Word struct {
	Used         [maxChars]uint8
	Remaining    [maxChars]uint8
	RestLength   int
	Length       int
	UsedMask     uint32
	DoNotUseMask uint32
	text         string
}

// Part of the word
type Part struct {
	Remaining    [maxChars]uint8
	RestLength   int
	DoNotUseMask uint32
	text         string
}

func isValidChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char)
}

// InitizalizeAnagram transforms the text to data representation
func InitizalizeAnagram(text string) Anagram {
	anagram := Anagram{
		Text: text,
	}

	positionCount := 0
	for i := range anagram.positions {
		anagram.positions[i] = -1
	}

	for _, char := range strings.ToUpper(text) {
		if isValidChar(char) {
			index := int(char)
			position := anagram.positions[index]
			if position < 0 {
				position = positionCount
				positionCount++
				anagram.positions[index] = position
			}

			anagram.Counts[position]++
			anagram.Length++
			anagram.Mask |= 1 << position
		}
	}

	return anagram
}

// IsComplete checks if the word combined with a rest is a valid anagram
func (w *Word) IsComplete(rest *Part) bool {
	for position, used := range w.Used {
		if used != rest.Remaining[position] {
			return false
		}
	}

	return true
}

// Combine creates writes the combination of word and rest to a new rest.
func (w *Word) Combine(entry *Part, target *Part) bool {
	target.Remaining = entry.Remaining
	target.DoNotUseMask = entry.DoNotUseMask
	target.RestLength = entry.RestLength - w.Length

	for position, used := range w.Used {
		if used > target.Remaining[position] {
			return false
		}

		target.Remaining[position] -= used
		if target.Remaining[position] == 0 {
			target.DoNotUseMask |= 1 << position
		}
	}

	if keepText {
		target.text = entry.text + " " + w.text
	}

	return true
}

// ToPart converts the word to a rest
func (w *Word) ToPart(target *Part) {
	target.Remaining = w.Remaining
	target.DoNotUseMask = w.DoNotUseMask
	target.RestLength = w.RestLength

	if keepText {
		target.text = w.text
	}
}

// Combine builds a word from an anagram and a text
func (a *Anagram) Combine(text string) *Word {
	w := Word{}
	w.Remaining = a.Counts
	w.UsedMask = 0
	w.DoNotUseMask = 0

	for _, char := range strings.ToUpper(text) {
		if isValidChar(char) {
			position := a.positions[int(char)]
			if position < 0 {
				return nil
			}

			if position < 0 || w.Remaining[position] == 0 {
				return nil
			}

			w.Remaining[position]--
			w.Used[position]++
			w.Length++

			w.UsedMask |= 1 << position
			if w.Remaining[position] == 0 {
				w.DoNotUseMask |= 1 << position
			}
		}
	}

	w.RestLength = a.Length - w.Length
	w.text = text
	return &w
}
