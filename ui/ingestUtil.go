package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	ndb "github.com/Twintat/randomExams/db"
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

// recive "1-3,5-6"
// return []int{1,2,3,5,6}
func expandRanges(input string) ([]int, error) {
	var result []int

	input = strings.ReplaceAll(input, " ", "")

	// split on commas
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
			result = append(result, i)
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
	colisions []int
}

func (e ExerciseColisions) Error() string {
	var msg string = "ERROR: exercises "
	for _, colision := range e.colisions {
		msg += fmt.Sprintln(colision) + ", "
	}
	msg += "already exists in the database"
	return msg
}

func checkRanges(chapter ndb.Chapter, testRange []int) error {
	// open db
	dbSource, err := ndb.OpenDB()
	if err != nil {
		return fmt.Errorf("[checkRanges] %w", err)
	}
	defer dbSource.Close()

	qdb := ndb.New(dbSource)
	ctx := context.Background()
	// get exercises
	exs, err := qdb.GetExeRange(ctx, chapter.ID)
	if err != nil {
		return fmt.Errorf("[checkRanges] %w", err)
	}

	var test bool = true
	// colisions between exercises
	var colisions []int
	for _, ex := range exs {
		for _, testEx := range testRange {
			if int(ex) == testEx {
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
