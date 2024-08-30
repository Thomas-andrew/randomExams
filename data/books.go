package data

import (
	"fmt"
)

type BookInfo struct {
	Id   int
	Info string

	Title     string
	Author    string
	Volume    string
	Edition   string
	Publisher string
	Year      string
}

func (b *BookInfo) GenerateInfo() {
	var info string = ""
	info += b.Title + ", "
	info += b.Author + ", "
	info += "vol." + b.Volume + ", "
	info += "ed." + b.Edition + ", "
	info += b.Publisher + ", "
	info += b.Year
	b.Info = info
}

func (b *BookInfo) GetInfos() map[string]string {
	result := make(map[string]string)
	result["title"] = b.Title
	result["author"] = b.Author
	result["volume"] = b.Volume
	result["edition"] = b.Edition
	result["publisher"] = b.Publisher
	result["year"] = b.Year
	return result
}

type BookInfos []BookInfo

func NewBookInfos() BookInfos {
	return make(BookInfos, 0)
}

func (b *BookInfos) BestMatch() *BookInfo {
	return &((*b)[0])
}

func PowerSet(str string) []string {
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

func BookScore(subs []string, book string) int {
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
