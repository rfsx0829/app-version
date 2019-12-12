package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	model "github.com/rfsx0829/little-tools/clouddisk/model"
)

// InsertFile insert a file info, won't use FID
func (d Database) InsertFile(file *model.File) error {
	stmt, err := d.db.Prepare(insertIntoFileTable)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(file.UID, file.Filename, file.Filepath, file.MD5Value)
	return err
}

// SelectFile select File infos
func (d Database) SelectFile(fields, extraRequirement string) ([]*model.File, error) {
	sqlStr := fmt.Sprintf(selectFromFileTable, fields) + extraRequirement

	rows, err := d.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	return rows2file(rows), nil
}

// DeleteFile delete a File info
func (d Database) DeleteFile(fid int) error {
	_, err := d.db.Exec(deleteFromFileTable + "where fid=" + strconv.Itoa(fid))
	return err
}

// UpdateFile update a File info
func (d Database) UpdateFile(fid int, changes string) error {
	_, err := d.db.Exec(updateFileTable + changes + " where fid=" + strconv.Itoa(fid))
	return err
}

func rows2file(rows *sql.Rows) []*model.File {
	slice := make([]*model.File, 0)
	for rows.Next() {
		var x model.File
		if err := rows.Scan(&x.FID, &x.UID, &x.Filename, &x.Filepath, &x.MD5Value); err != nil {
			log.Println(err)
		}
		slice = append(slice, &x)
	}
	return slice
}
