package main

import (
	"github.com/deckarep/golang-set"
	"os"
	"strings"
)

var tagsToExclude = []interface{}{
	"podcasts",
	"punt counterpunt",
	"this week's obsession",
	"draftageddon 2015",
	"press conference transcripts",
	"press conference transcript",
	"2015 media day",
	"big ten media days",
	"guess the score",
}
var ExcludeTagSet = mapset.NewSetFromSlice(tagsToExclude)

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func removeSpecialCharacters(inputString, charsToStripOut string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(charsToStripOut, r) < 0 {
			return r
		}
		return -1
	}, inputString)
}

func SanitizeFilename(filename string) string {
	return strings.Replace(removeSpecialCharacters(filename, "';:!@#$%^&*/\\|"), " ", "-", -1)
}

var wordsToIgnore = []interface{}{
	"a", "about", "after", "all", "alright", "also", "an", "any", "and", "are", "around", "as", "at",
	"be", "because", "before", "both", "but", "by", "been",
	"can", "could",
	"did", "didn't", "do", "does", "don't", "down",
	"end", "either", "even",
	"for", "from",
	"get", "gets", "go", "got", "going",
	"had", "has", "have", "he", "here", "his", "him", "how",
	"i", "if", "in", "into", "is", "it",
	"just",
	"like", "lot",
	"maybe", "me", "more", "most", "much", "my",
	"neither", "no", "nor", "not", "now",
	"of", "off", "on", "only", "or", "other", "out", "over",
	"probably",
	"really", "right",
	"said", "say", "since", "she", "should", "so", "some", "still",
	"take", "than", "that", "the", "their", "them", "then", "there", "they", "this", "these", "those", "though", "through", "to",
	"up",
	"very",
	"was", "way", "we", "well", "were", "what", "where", "whether", "when", "which", "while", "will", "who", "with", "would",
	"yet", "you", "your",
	"", "$",
}
var ExcludeWordSet = mapset.NewSetFromSlice(wordsToIgnore)
