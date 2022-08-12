package db

import (
	"context"
	"fmt"
)

func (db DB) Update(id int, value string) error {
	conn := db.conn

	if conn == nil {
		panic("accessing uninitialized database connection")
	}

	_, err := conn.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET number = %s WHERE id = %d", db.tablename, value, id))

	return err
}
