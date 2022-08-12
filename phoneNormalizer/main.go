package main

import (
	"gophercises/phoneNormalizer/db"
	"gophercises/phoneNormalizer/phonenumber"
)

func main() {
	d := db.DB{}

	defer d.Close()

	d.Init()

	phonenumber.Format(d)
}
