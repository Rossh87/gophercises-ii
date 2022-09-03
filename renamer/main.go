package main

import (
	"os"
	"renamer/fileManager"
	"renamer/parser"
)

func main() {
	args := os.Args[1:]

	fm := fileManager.FileManager{}

	parser.Rename(args[0], &fm)
}
