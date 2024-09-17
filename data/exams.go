package data

import (
	"context"

	"github.com/Twintat/randomExams/db"
)

type Exam struct {
	Gui      *GUI
	Pull     []db.Exercise
	Num      int
	Duration string
}

func NewExam(g *GUI) *Exam {
	return &Exam{
		Gui:  g,
		Pull: []db.Exercise{},
	}
}

type SetTable struct {
	BookIDs    []int
	ChapterIDs []int
}

func NewSetTable() SetTable {
	return SetTable{
		BookIDs:    []int{},
		ChapterIDs: []int{},
	}
}

func (e *SetTable) AddBookID(id int) {
	e.BookIDs = append(e.BookIDs, id)
}

func (e *SetTable) AddChapterID(id int) {
	e.ChapterIDs = append(e.ChapterIDs, id)
}

func (e *SetTable) GenExerPull() ([]db.Exercise, error) {
	fail := func(err error) ([]db.Exercise, error) {
		return nil, err
	}

	// open db
	dbSrc, err := db.OpenDB()
	if err != nil {
		return fail(err)
	}
	defer dbSrc.Close()

	qdb := db.New(dbSrc)
	ctx := context.Background()

	exercises := []db.Exercise{}
	// get exercises from bookID
	for _, bookID := range e.BookIDs {
		chapterIDs, err := qdb.GetChapterIDs(ctx, int64(bookID))
		if err != nil {
			return fail(err)
		}
		for _, chapterID := range chapterIDs {
			pulled, err := qdb.GetExercises(ctx, chapterID)
			if err != nil {
				return fail(err)
			}
			exercises = append(exercises, pulled...)
		}
	}
	// get exercises from chapterID
	for _, chapterID := range e.ChapterIDs {
		pulled, err := qdb.GetExercises(ctx, int64(chapterID))
		if err != nil {
			return fail(err)
		}
		exercises = append(exercises, pulled...)
	}
	return exercises, nil
}
