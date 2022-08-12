package db

import (
	"context"
	"fmt"
	"gophercises/phoneNormalizer/phonenumber"
)

func (db DB) GetRecords() ([]phonenumber.PhoneRecord, error) {
	conn := db.conn

	if conn == nil {
		panic("accessing uninitialized database connection")
	}

	rows, err := conn.Query(context.Background(), fmt.Sprintf("SELECT * FROM %s", db.tablename))

	if err != nil {
		return nil, err
	}

	ret := []phonenumber.PhoneRecord{}

	for rows.Next() {
		var id int
		var pn string

		if err := rows.Scan(&id, &pn); err != nil {
			return nil, err
		}

		ret = append(ret, phonenumber.PhoneRecord{Id: id, Pn: pn})
	}

	rows.Close()

	return ret, nil
}
