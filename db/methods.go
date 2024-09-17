package db

import "fmt"

func (b *Book) Info() string {
	var vol string
	var ed string
	var pub string
	if b.Volume.Valid {
		vol = " vol " + b.Volume.String
	}
	if b.Edition.Valid {
		ed = " ed " + b.Edition.String
	}
	if b.Publisher.Valid {
		pub = " pub " + b.Publisher.String
	}
	return fmt.Sprintf(
		"%v, %v, %v, %v%v%v",
		b.Title,
		b.Author,
		b.Year,
		vol, ed, pub,
	)
}

func (c *Chapter) Info() string {
	return fmt.Sprintf("%v - %v", c.Number, c.Name)
}
