package data

import (
	"strconv"
)

type Chapter struct {
	Id   int
	Num  int
	Name string
	Info string
}

// generate info
func (c *Chapter) GenerateInfo() {
	c.Info = strconv.Itoa(c.Num) + " - " + c.Name
}

type Chapters []Chapter

func (c *Chapters) BestMatch() Chapter {
	return (*c)[0]
}
