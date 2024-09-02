package database

import (
	"fmt"

	"github.com/Twintat/randomExams/data"
)

func GenExerPull(set *data.SetTable) (data.Exercises, error) {
	// get chapters from books
	chapters := data.Chapters{}
	for _, bookID := range set.BookIDs {
		bookChaps, err := ListChapters(bookID)
		if err != nil {
			return nil, fmt.Errorf("[GenExerPull] %v", err)
		}
		chapters = append(chapters, bookChaps...)
	}
	for _, chapID := range set.ChapterIDs {
		chap, err := GetChapter(chapID)
		if err != nil {
			return nil, fmt.Errorf("[GenExerPull] %v", err)
		}
		if !chapters.IsEqual(chap) {
			chapters = append(chapters, chap)
		}
	}
	// get exercises from chapters
	exs := data.NewExercises()
	for _, chap := range chapters {
		ex, err := GetExercises(chap.Id)
		if err != nil {
			return nil, fmt.Errorf("[GenExerPull] %v", err)
		}
		exs = append(exs, ex...)
	}
	return exs, nil
}
