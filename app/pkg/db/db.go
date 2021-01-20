package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func AccessData() *sql.DB {
	dataBase, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Println("db err:", err)
	}
	return dataBase
}

func Initiate(dataBase *sql.DB) {
	statement, err := dataBase.Prepare(
		"CREATE TABLE IF NOT EXISTS settings (id INTEGER PRIMARY KEY, settingName TEXT, settingValue INTEGER)")
	if err != nil {
		log.Println("db err:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println("db err:", err)
	}
}

func AddSetting(dataBase *sql.DB, settingName string, settingValue int) {
	statement, err := dataBase.Prepare("INSERT INTO settings (settingName, settingValue) VALUES (?, ?)")
	if err != nil {
		log.Println("db err:", err)
	}

	_, err = statement.Exec(settingName, settingValue)
	if err != nil {
		log.Println("db err:", err)
	}
}

func EditSetting(dataBase *sql.DB, id int, settingValue int) {
	statement, err := dataBase.Prepare("UPDATE settings SET settingValue=? WHERE id=?")
	if err != nil {
		log.Println("db err:", err)
	}
	_, err = statement.Exec(settingValue, id)
	if err != nil {
		log.Println("db err:", err)
	}
}

func DeleteSetting(dataBase *sql.DB, id int) {
	statement, err := dataBase.Prepare("DELETE FROM settings WHERE id=?")
	if err != nil {
		log.Println("db err:", err)
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Println("db err:", err)
	}
}

func ReadSettings(dataBase *sql.DB) []int {
	var result []int

	rows, err := dataBase.Query("SELECT * from settings")
	if err != nil {
		log.Println("db err:", err)
	}

	for rows.Next() {
		var (
			tmpId int
			tmpName string
			tmpVal int
		)
		err = rows.Scan(&tmpId, &tmpName, &tmpVal)
		if err != nil {
			log.Println(err)
		}

		result = append(result, tmpVal)
	}
	return result
}