// +build postgres

package main

import (
	"database/sql"
)

type User struct {
	name      string
	year      int
	yearNo    int
	staffYear int
}

func NewUser(name string, year, yearNo, staffYear sql.NullInt64) *User {
	return &User{name, int(year.Int64), int(yearNo.Int64), int(staffYear.Int64)}
}

func (s PostgresStore) GetUser(name string) (*User, error) {
	var year, yearNo, staffYear sql.NullInt64
	sql := "select year, year_no, staff_year from users where name = $1"
	err := s.db.QueryRow(sql, name).Scan(&year, &yearNo, &staffYear)
	if err != nil {
		return nil, err
	}

	return NewUser(name, year, yearNo, staffYear), nil
}

func (s PostgresStore) ListUsers(year int) ([]User, error) {
	query := "select name, year, year_no, staff_year from users"
	params := []interface{}{}
	if year != 0 {
		query += " where year = $1 or staff_year = $1"
		params = append(params, year)
	}
	query += " order by year_no, name"
	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var name string
		var year, yearNo, staffYear sql.NullInt64
		err = rows.Scan(&name, &year, &yearNo, &staffYear)
		if err != nil {
			return nil, err
		}
		users = append(users, *NewUser(name, year, yearNo, staffYear))
	}

	return users, nil
}
