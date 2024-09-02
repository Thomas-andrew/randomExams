package data

import (
	"strconv"
)

type Chapter struct {
	Id     int
	BookID int
	Num    int
	Name   string
	Info   string
}

// generate info
func (c *Chapter) GenerateInfo() {
	c.Info = strconv.Itoa(c.Num) + " - " + c.Name
}

type Chapters []Chapter

func (c *Chapters) BestMatch() Chapter {
	return (*c)[0]
}

func (c *Chapters) IsEqual(t Chapter) bool {
	for _, currChap := range *c {
		if currChap.Id == t.Id {
			return true
		}
	}
	return false
}
