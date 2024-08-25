package main

import "strconv"

type exercise struct {
	id     int
	num    int
	images map[int]string
}

func newExercise(id, num int) exercise {
	return exercise{
		id:     id,
		num:    num,
		images: make(map[int]string),
	}
}

type exercises []exercise

func newExercises() exercises {
	return make(exercises, 0)
}

func (e exercises) getRange() []string {
	var result []string
	for _, ex := range e {
		result = append(result, strconv.Itoa(ex.num))
	}
	return result
}
