package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	model "github.com/rfsx0829/little-tools/clouddisk/model"
)

// InsertUser insert a user info, won't use UID
func (d Database) InsertUser(user *model.User) error {
	stmt, err := d.db.Prepare(insertIntoUserTable)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Uname, user.Password, user.Email)
	return err
}

// SelectUser select user infos
func (d Database) SelectUser(fields, extraRequirement string) ([]*model.User, error) {
	sqlStr := fmt.Sprintf(selectFromUserTable, fields) + extraRequirement

	rows, err := d.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	return rows2user(rows), nil
}

// DeleteUser delete a user info
func (d Database) DeleteUser(uid int) error {
	_, err := d.db.Exec(deleteFromUserTable + "where uid=" + strconv.Itoa(uid))
	return err
}

// UpdateUser update a user info
func (d Database) UpdateUser(uid int, changes string) error {
	_, err := d.db.Exec(updateUserTable + changes + " where uid=" + strconv.Itoa(uid))
	return err
}

func rows2user(rows *sql.Rows) []*model.User {
	slice := make([]*model.User, 0)
	for rows.Next() {
		var x model.User
		if err := rows.Scan(&x.UID, &x.Uname, &x.Password, &x.Email); err != nil {
			log.Println(err)
		}
		slice = append(slice, &x)
	}
	return slice
}
