package main

import (
	"fmt"
)

type bookInfo struct {
	id   int
	info string

	title     string
	author    string
	volume    string
	edition   string
	publisher string
	year      string
}

func (b *bookInfo) generateInfo() {
	var info string = ""
	info += b.title + ", "
	info += b.author + ", "
	info += "vol." + b.volume + ", "
	info += "ed." + b.edition + ", "
	info += b.publisher + ", "
	info += b.year
	b.info = info
}

func (b *bookInfo) getInfos() map[string]string {
	result := make(map[string]string)
	result["title"] = b.title
	result["author"] = b.author
	result["volume"] = b.volume
	result["edition"] = b.edition
	result["publisher"] = b.publisher
	result["year"] = b.year
	return result
}

type bookInfos []bookInfo

func newBookInfos() bookInfos {
	return make(bookInfos, 0)
}

func (b *bookInfos) bestMatch() *bookInfo {
	return &((*b)[0])
}

func powerSet(str string) []string {
	lenStr := len(str)

	var subs []string

	for start := 0; start < lenStr; start++ {
		for stop := start + 1; stop <= lenStr; stop++ {
			subStr := str[start:stop]
			subs = append(subs, subStr)
		}
	}
	return subs
}

func bookScore(subs []string, book string) int {
	var score int = 0
	var logStr string

	logStr = fmt.Sprintf("[bookScore] book: '%v'\n", book)
	for _, sub := range subs {
		lenStr := len(sub)
		lenBook := len(book)
		logStr += fmt.Sprintf("{sub: '%v', lenStr: %v, lenBook: %v}\n", sub, lenStr, lenBook)
		logStr += fmt.Sprintf(
			"{lenBook-lenStr: %v}\n",
			lenBook-lenStr,
		)

		for start := 0; start <= lenBook-lenStr; start++ {
			stop := start + lenStr
			logStr += fmt.Sprintf(
				"{start: %v, stop: %v, sub: %v, bookSub: %v}\n",
				start,
				stop,
				sub,
				book[start:stop],
			)
			if sub == book[start:stop] {
				score += lenStr + (lenBook - start)
				logStr += fmt.Sprintf("{score: %v}\n", score)
				break
			}
		}
	}

	// log.Println(logStr)

	return score
}
