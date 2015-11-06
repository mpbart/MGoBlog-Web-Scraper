package main

import (
	"fmt"
	"sort"
	"strings"
)

type WordCounter map[string]int
type WordFrequency struct {
	word      string
	frequency int
}
type SortableFrequency []WordFrequency

var finalResults = WordCounter{}

func isSpecialCharacter(word string) bool {
	if len(word) == 0 || isNumber(rune(word[0])) || isPunctuation(rune(word[0])) {
		return true
	}
	return false
}

func isNumber(c rune) bool {
	return strings.ContainsRune("1234567890", c)
}

func isPunctuation(c rune) bool {
	return strings.ContainsRune("!@#$%^&*?.|~,/\\=+", c)
}

func CountWords(article Article) WordCounter {
	articleCounter := WordCounter{}
	for _, str := range article.Content {
		words := strings.Split(str, " ")
		for _, word := range words {
			if !isSpecialCharacter(word) {
				baseWord := removeExtras(word)
				if _, exists := articleCounter[baseWord]; exists == false {
					articleCounter[baseWord] = 1
				} else {
					articleCounter[baseWord] += 1
				}
			}
		}
	}
	return articleCounter
}

func removeExtras(word string) string {
	word = strings.ToLower(word)
	if len(word) > 2 && word[len(word)-2:] == "'s" {
		word = word[:len(word)-2]
	}
	if strings.ContainsRune("'\"([{", rune(word[0])) {
		word = word[1:]
	}
	for {
		if len(word) == 0 {
			break
		}
		if strings.ContainsRune("\"?.,!)]}*':;", rune(word[len(word)-1])) {
			word = word[:len(word)-1]
		} else {
			break
		}
	}
	return word
}

func (arr SortableFrequency) Len() int {
	return len(arr)
}

func (arr SortableFrequency) Less(i, j int) bool {
	if arr[i].frequency > arr[j].frequency {
		return true
	}
	return false
}

func (arr SortableFrequency) Swap(i, j int) {
	a, b := arr[i], arr[j]
	arr[i], arr[j] = b, a
}

func aggregateResults(words WordCounter) {
	for word, count := range words {
		if _, exists := finalResults[word]; exists == false {
			finalResults[word] = count
		} else {
			finalResults[word] += count
		}
	}
}
func printTopNResults(numberToPrint int) {
	frequencies := make(SortableFrequency, len(finalResults))
	idx := 0
	for key, value := range finalResults {
		if ExcludeWordSet.Contains(key) {
			continue
		}
		frequencies[idx] = WordFrequency{word: key, frequency: value}
		idx += 1
	}
	sort.Sort(frequencies)
	for _, i := range frequencies[:numberToPrint] {
		fmt.Println(i)
	}
}
