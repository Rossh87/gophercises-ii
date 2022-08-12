package main

import (
	"gophercises/phoneNormalizer/db"
)

func main() {
	d := db.DB{}

	d.Init()

	d.Setup()

	d.Close()
}
