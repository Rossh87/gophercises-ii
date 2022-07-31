package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gophercises/taskCLI/db"

	"github.com/urfave/cli"
)

// type stubDB struct {
// }

// func (db stubDB) Add(task string) error {
// 	return nil
// }

// func (db stubDB) List() ([]string, error) {

// 	v := []string{
// 		"Fuck bitches",
// 		"Get money",
// 		"Do crimes",
// 	}

// 	return v, nil
// }

// func (db stubDB) Do(task string) error {
// 	return nil
// }

func initApp(db db.DB) *cli.App {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add a task to list",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if c.NArg() > 0 {
					task := strings.Join(args, ` `)
					fmt.Printf("Adding task \"%s\"\n", task)

					err := db.Add(task)

					if err != nil {
						fmt.Printf("error adding task \"%s\"\n", task)
						fmt.Printf("reason: %+v\n", err)
					} else {
						fmt.Printf("successfully added \"%s\"\n", task)
					}

					return err
				}

				return errors.New("length of task to add must be at least 1")
			}},
		{
			Name:  "list",
			Usage: "show all tasks",
			Action: func(c *cli.Context) error {
				list, err := db.List()

				if err != nil {
					fmt.Println("Failed listing tasks")
					fmt.Printf("reason: %+v\n", err)
					return err
				}

				for _, task := range list {
					fmt.Printf("%d. %s\n", task.Key, task.Value)
				}

				return nil
			}}, {
			Name:  "do",
			Usage: "complete a task",
			Action: func(c *cli.Context) error {
				if c.NArg() == 1 {
					taskNumber := c.Args().Get(0)

					id, err := strconv.ParseInt(taskNumber, 10, 64)

					if err != nil {
						fmt.Printf("Failed to convert argument \"%s\" to integer", taskNumber)
						return err
					}

					err = db.Do(id)

					if err != nil {
						fmt.Printf("error completing task \"%s\"\n", taskNumber)
						fmt.Printf("reason: %+v\n", err)
					} else {
						fmt.Printf("successfully completed \"%s\"\n", taskNumber)
					}

					return err
				}

				return errors.New("must specify exactly one task to complete")
			}},
	}

	return app
}

func main() {
	db := db.New("./data.db")

	defer db.Close()

	app := initApp(db)

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("Unable to run app\n%+v\n", err)
	}
}
