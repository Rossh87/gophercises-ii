package main

import (
	"fmt"
	"os"

	"github.com/rossh87/gophercises/passwordReader/cmd"
)

func main() {
	fmt.Println("Enter PW:")
	res, err := cmd.Run(cmd.StdinPasswordReader{})
	if err != nil {
		fmt.Printf("problem!: %v", err)
		os.Exit(1)
	}
	fmt.Printf("your PW is: %s", res)
}
