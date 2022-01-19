package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBConnection interface {
	Initialize()
	Update(query string) bool
	Query(query string) *sql.Rows
}

type DatabaseConnection struct {
	DB *sql.DB
}

func NewDatabase() *DatabaseConnection {
	return &DatabaseConnection{
		DB: DBConnect(),
	}
}

func DBConnect() *sql.DB {
	db, err := sql.Open("mysql", "mudhost:B@ckstab69@tcp(127.0.0.1:3306)/gomud")
	if err != nil {
		panic(err.Error())
	}

	// TODO(jon): handle during shutdown
	//defer db.Close()
	return db
}

func (db *DatabaseConnection) Initialize() {
	db.Update("CREATE TABLE IF NOT EXISTS Account(" +
		"username VARCHAR(20) PRIMARY KEY," +
		"password CHAR(60) NOT NULL," +
		"max_chars TINYINT UNSIGNED NOT NULL DEFAULT 3," +
		"last_ip VARCHAR(15)," +
		"last_date TIMESTAMP," +
		"email VARCHAR(319) NOT NULL" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS CharClass(" +
		"name VARCHAR(20) PRIMARY KEY," +
		"realm TINYINT NOT NULL DEFAULT 0," +
		"enabled BOOLEAN NOT NULL DEFAULT 0" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Race(" +
		"name VARCHAR(20) PRIMARY KEY," +
		"enabled BOOLEAN NOT NULL DEFAULT 0" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Area(" +
		"name VARCHAR(50) PRIMARY KEY," +
		"realm TINYINT UNSIGNED NOT NULL DEFAULT 0" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Room(" +
		"room_id VARCHAR(5) NOT NULL," +
		"area VARCHAR(50) NOT NULL," +
		"description VARCHAR(512)," +
		"links VARCHAR(1024)," +
		"PRIMARY KEY (area, room_id)," +
		"FOREIGN KEY (area) REFERENCES Area(name) ON UPDATE CASCADE ON DELETE RESTRICT" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Portal(" +
		"name VARCHAR(50) NOT NULL," +
		"room VARCHAR(60) NOT NULL," +
		"dest_room VARCHAR(60)," +
		"PRIMARY KEY (name, room)" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Item(" +
		"article VARCHAR(3)," +
		"name VARCHAR(50) NOT NULL," +
		"keyword VARCHAR(50) NOT NULL," +
		"owner VARCHAR(60) NOT NULL," +
		"value INT UNSIGNED NOT NULL," +
		"item_type TINYINT UNSIGNED NOT NULL," +
		"PRIMARY KEY (article, name)" +
		")")

	db.Update("CREATE TABLE IF NOT EXISTS Player(" +
		"name VARCHAR(20) PRIMARY KEY," +
		"account VARCHAR(20)," +
		"class VARCHAR(20)," +
		"race VARCHAR(20)," +
		"room VARCHAR(60)," +
		"coins BIGINT UNSIGNED NOT NULL," +
		"stre TINYINT UNSIGNED NOT NULL," +
		"cons TINYINT UNSIGNED NOT NULL," +
		"agil TINYINT UNSIGNED NOT NULL," +
		"dext TINYINT UNSIGNED NOT NULL," +
		"inte TINYINT UNSIGNED NOT NULL," +
		"wisd TINYINT UNSIGNED NOT NULL," +
		"con_loss TINYINT UNSIGNED NOT NULL," +
		"level TINYINT UNSIGNED NOT NULL," +
		"exp BIGINT UNSIGNED NOT NULL," +
		"rp INT UNSIGNED NOT NULL," +
		"hits SMALLINT UNSIGNED NOT NULL," +
		"fat SMALLINT UNSIGNED NOT NULL," +
		"power SMALLINT UNSIGNED NOT NULL," +
		"trains SMALLINT UNSIGNED NOT NULL DEFAULT 0," +
		"guild VARCHAR(30)," + // FK
		"guild_rank TINYINT UNSIGNED NOT NULL DEFAULT 0," +
		"last_date TIMESTAMP," +
		"FOREIGN KEY (account) REFERENCES Account(username) ON UPDATE CASCADE ON DELETE SET NULL," +
		"FOREIGN KEY (class) REFERENCES CharClass(name) ON UPDATE CASCADE ON DELETE SET NULL," +
		"FOREIGN KEY (race) REFERENCES Race(name) ON UPDATE CASCADE ON DELETE SET NULL" +
		")")
}

func (db *DatabaseConnection) Update(query string) bool {
	insert, err := db.DB.Query(query)
	if err != nil {
		log.Panic(err.Error())
		panic(err.Error())
		return false
	}
	defer insert.Close()
	return true
}

func (db *DatabaseConnection) Query(query string) *sql.Rows {
	results, err := db.DB.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	return results
}
