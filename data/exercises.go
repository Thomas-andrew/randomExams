package data

import "strconv"

type Exercise struct {
	Id     int
	Num    int
	Images map[int]string
}

func NewExercise(id, num int) Exercise {
	return Exercise{
		Id:     id,
		Num:    num,
		Images: make(map[int]string),
	}
}

type Exercises []Exercise

func NewExercises() Exercises {
	return make(Exercises, 0)
}

func (e Exercises) GetRange() []string {
	var result []string
	for _, ex := range e {
		result = append(result, strconv.Itoa(ex.Num))
	}
	return result
}
