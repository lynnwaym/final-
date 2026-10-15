package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema string = `CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date char(8) NOT NULL DEFAULT '',
	title VARCHAR(255) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
	);
	
CREATE INDEX scheduler_date ON scheduler(date)`

func Init(dbFile string) error {

	install := false
	_, err := os.Stat(dbFile)

	if os.IsNotExist(err) {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)

	if err != nil {
		return err
	}

	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			DB.Close()
			return err
		}
	}
	return nil
}
