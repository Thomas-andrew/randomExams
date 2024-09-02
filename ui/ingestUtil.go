package ui

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Twintat/randomExams/data"
	db "github.com/Twintat/randomExams/database"
)

func imageName() string {
	now := time.Now()

	_, week := now.ISOWeek()

	dayOfYear := now.YearDay()

	formattedTime := fmt.Sprintf("%02d%02d%04d-%02d%03d-%02d%02d%02d",
		now.Day(), now.Month(), now.Year(),
		week, dayOfYear,
		now.Hour(), now.Minute(), now.Second())
	return formattedTime + ".png"
}

func expandRanges(input string) ([]string, error) {
	var result []string

	input = strings.ReplaceAll(input, " ", "")

	ranges := strings.Split(input, ",")
	for _, r := range ranges {
		bounds := strings.Split(r, "-")
		start, err := strconv.Atoi(bounds[0])
		if err != nil {
			return nil, err
		}

		end := start
		if len(bounds) > 1 {
			end, err = strconv.Atoi(bounds[1])
			if err != nil {
				return nil, err
			}
		}

		for i := start; i <= end; i++ {
			result = append(result, strconv.Itoa(i))
		}

	}

	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] == result[j] {
				return nil, fmt.Errorf("overlaping ranges. Not permited")
			}
		}
	}

	return result, nil
}

type ExerciseColisions struct {
	colisions []string
}

func (e ExerciseColisions) Error() string {
	var msg string = "ERROR: exercises "
	for _, colision := range e.colisions {
		msg += colision + ", "
	}
	msg += "already exists in the database"
	return msg
}

func checkRanges(chapter data.Chapter, testRange []string) error {
	// get exercises
	exers, err := db.GetExercises(chapter.Id)
	if err != nil {
		if _, ok := err.(db.NoID); ok {
			// chapter is no id, probably new
			slog.Debug("[checkRanges] chapter has no id")
			return nil
		}
		return fmt.Errorf("[checkRanges] %w", err)
	}

	var test bool = true
	// colisions between exercises
	var colisions []string
	exersNumbers := exers.GetRange()
	for _, ex := range exersNumbers {
		for _, testEx := range testRange {
			if ex == testEx {
				test = false
				colisions = append(colisions, testEx)
			}
		}
	}
	if !test {
		return ExerciseColisions{colisions: colisions}
	}

	return nil
}
