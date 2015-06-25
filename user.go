package main

import (
	"database/sql"
	"path"
	"strconv"
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

func (user User) Dir(base, file string) string {
	dir := "staff"
	if user.staffYear == 0 && user.year != 0 {
		dir = strconv.Itoa(user.year)
	}
	return path.Join(base, dir, user.name, file)
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

func (s PostgresStore) ListUsers(year int, order bool) ([]User, error) {
	query := "select name, year, year_no, staff_year from users"
	params := []interface{}{}
	if year != 0 {
		query += " where year = $1 or staff_year = $1"
		params = append(params, year)
	}

	if order {
		if year == 0 {
			query += " order by year_no, staff_no, name"
		} else {
			query += " order by staff_no, year_no, name"
		}
	} else {
		query += " order by name"
	}

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
