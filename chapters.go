package main

import (
	"fmt"
	"strconv"
)

type chapter struct {
	id   int
	num  int
	name string
	info string
}

// generate info
func (c *chapter) generateInfo() {
	c.info = strconv.Itoa(c.num) + " - " + c.name
}

type NoID struct {
	objectName string
}

func (e NoID) Error() string {
	return fmt.Sprintf("ERROR: '%v' has no id yet! Maybe new or not in the db\n", e.objectName)
}

func (c *chapter) getExercises() (exercises, error) {
	if c.id == 0 {
		return nil, NoID{objectName: "chapter"}
	}
	db, err := openDB()
	if err != nil {
		return nil, fmt.Errorf("[getExercises] %w", err)
	}
	defer db.Close()

	// get exercise IDs from this chapter
	query := "SELECT id, exNum FROM exerciseId WHERE chapterID = ?"
	rowsExerIDs, err := db.Query(query, c.id)
	if err != nil {
		return nil, fmt.Errorf("[getExercises] %w", err)
	}
	defer rowsExerIDs.Close()

	exers := newExercises()
	// from exercise IDs get images
	for rowsExerIDs.Next() {
		// get exercise id and number
		var exID, exNum int
		err := rowsExerIDs.Scan(&exID, &exNum)
		if err != nil {
			return nil, fmt.Errorf("[getExercises] %w", err)
		}
		// from id and number get exercise order and image name
		query = "SELECT imageName, imageOrder FROM exerciseData WHERE exID = ?"
		rowsExerData, err := db.Query(query, exID)
		if err != nil {
			return nil, fmt.Errorf("[getExercises] %w", err)
		}
		ex := newExercise(exID, exNum)
		for rowsExerData.Next() {
			var imageName string
			var imageOrder int
			err := rowsExerData.Scan(&imageName, &imageOrder)
			if err != nil {
				return nil, fmt.Errorf("[getExercises] %w", err)
			}
			ex.images[imageOrder] = imageName
		}
		exers = append(exers, ex)
	}

	return exers, nil
}

type chapters []chapter

func (c *chapters) bestMatch() chapter {
	return (*c)[0]
}

// retrive from db
// listChapters
func listChapters(bookID int) (chapters, error) {
	if bookID == 0 {
		return nil, nil
	}
	db, err := openDB()
	if err != nil {
		return nil, fmt.Errorf("[listChapters] %w", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT chapterID, number, name FROM chapters WHERE bookID = ?", bookID)
	if err != nil {
		return nil, fmt.Errorf("[listChapters] error retriving chapters of a book: %w", err)
	}
	defer rows.Close()

	chapters := make(chapters, 0)

	for rows.Next() {
		var chapterID int
		var chapterNum int
		var chapterName string

		err = rows.Scan(&chapterID, &chapterNum, &chapterName)
		if err != nil {
			return nil, fmt.Errorf("[listChapters] error scanning chapterNum, chapterName: %w", err)
		}
		chapter := chapter{id: chapterID, num: chapterNum, name: chapterName}
		chapter.generateInfo()
		chapters = append(chapters, chapter)
	}
	return chapters, nil
}
