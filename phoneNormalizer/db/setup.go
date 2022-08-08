package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

const (
	host      = "localhost"
	port      = "5432"
	user      = "postgres"
	dbname    = "gophercises_phonenumber"
	tablename = "phone_numbers"
)

func main() {
	rawdata := [][]interface{}{
		{"123 456 7891"},
		{"(123) 456 7892"},
		{"(123) 456-7893"},
		{"123-456-7894"},
		{"123-456-7890"},
		{"1234567892"},
		{"(123)456-7892"},
		{"1234567890"},
	}

	conStr := fmt.Sprintf("postgres://%s@%s:%s/%s", user, host, port, dbname)

	conn, err := pgx.Connect(context.Background(), conStr)

	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS phone_numbers(id SERIAL NOT NULL PRIMARY KEY, number TEXT NOT NULL);")

	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
}

func setupTable(conn *pgx.Conn, data [][]interface{}) {
	_, err := conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS phone_numbers(id SERIAL NOT NULL PRIMARY KEY, number TEXT NOT NULL);")

	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	// column name is number, table name is phone_numbers
	conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"phone_numbers"},
		[]string{"number"},
		pgx.CopyFromRows(data),
	)
}

func reset(conn *pgx.Conn) {
	_, err := conn.Exec(context.Background(), "DROP TABLE $1;1", tablename)

}
