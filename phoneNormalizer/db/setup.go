package db

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
	dbname    = "phone_normalizer"
	tablename = "phone_numbers"
)

type DB struct {
	tablename string
	conn      *pgx.Conn
}

func (db *DB) Init() {
	conStr := fmt.Sprintf("postgres://%s@%s:%s/%s", user, host, port, dbname)

	conn, err := pgx.Connect(context.Background(), conStr)

	handleErr(err)

	db.conn = conn

	db.tablename = tablename
}

func (db *DB) Close() {
	db.conn.Close(context.Background())
}

func (db *DB) Setup() {
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

	handleErr(err)

	defer conn.Close(context.Background())

	err = reset(conn, tablename)

	handleErr(err)

	err = table(conn)

	handleErr(err)

	err = populate(conn, tablename, rawdata)

	handleErr(err)
}

func table(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS phone_numbers(id SERIAL NOT NULL PRIMARY KEY, number TEXT NOT NULL);")

	return err
}

func populate(conn *pgx.Conn, tablename string, data [][]interface{}) error {
	_, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{tablename},
		[]string{"number"},
		pgx.CopyFromRows(data),
	)

	return err
}

func reset(conn *pgx.Conn, tablename string) error {
	cmd := fmt.Sprintf("DROP TABLE IF EXISTS %s", tablename)

	_, err := conn.Exec(context.Background(), cmd)

	return err
}

func handleErr(e error) {
	if e != nil {
		log.Fatalf("%v", e)
		os.Exit(1)
	}
}
