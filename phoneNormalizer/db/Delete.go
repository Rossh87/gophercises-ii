package db

import (
	"context"
	"fmt"
)

func (db DB) Delete(id int) error {
	conn := db.conn

	if conn == nil {
		panic("accessing uninitialized database connection")
	}

	_, err := conn.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s WHERE id = %d", db.tablename, id))

	return err
}
