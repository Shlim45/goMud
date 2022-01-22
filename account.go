package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type Account struct {
	username   string
	password   string
	maxChars   uint8
	lastIp     string
	lastDate   time.Time
	email      string
	characters map[string]*MOB
}

func NewAccount() *Account {
	return &Account{
		characters: make(map[string]*MOB),
	}
}

func (a *Account) UserName() string {
	return a.username
}

func (a *Account) SetUserName(newName string) {
	a.username = newName
}

func (a *Account) PasswordHash() string {
	return a.password
}

func (a *Account) SetPasswordHash(hashedPass string) {
	a.password = hashedPass
}

func (a *Account) HashAndSetPassword(plainPass string) {
	passHash, err := HashPassword(plainPass)
	if err != nil {
		log.Panic(err.Error())
	}
	a.password = passHash
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (a *Account) CheckPasswordHash(toCheck string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.password), []byte(toCheck))
	return err == nil
}

func (a *Account) MaxChars() uint8 {
	return a.maxChars
}

func (a *Account) SetMaxChars(newMax uint8) {
	a.maxChars = newMax
}

func (a *Account) LastIP() string {
	return a.lastIp
}

func (a *Account) SetLastIP(newIP string) {
	a.lastIp = newIP
}

func (a *Account) LastDate() time.Time {
	return a.lastDate
}

func (a *Account) SetLastDate(newDate time.Time) {
	a.lastDate = newDate
}

func (a *Account) UpdateLastDate() {
	a.lastDate = time.Now()
}

func (a *Account) Email() string {
	return a.email
}

func (a *Account) SetEmail(newEmail string) {
	a.email = newEmail
}

func (a *Account) SaveAccountToDBQuery() string {
	lastDate := TimeString(a.LastDate())
	return fmt.Sprintf("INSERT INTO Account VALUES ('%s', '%s', %d, '%s', '%s', '%s') AS new ON DUPLICATE KEY UPDATE "+
		"username=new.username, password=new.password, max_chars=new.max_chars, last_ip=new.last_ip, last_date=new.last_date, email=new.email",
		a.UserName(), a.PasswordHash(), a.MaxChars(), a.LastIP(), lastDate, a.Email())
}

func CreateAccountTableDBQuery() string {
	return "CREATE TABLE IF NOT EXISTS Account(" +
		"username VARCHAR(20) PRIMARY KEY," +
		"password CHAR(60) NOT NULL," +
		"max_chars TINYINT UNSIGNED NOT NULL DEFAULT 3," +
		"last_ip VARCHAR(15)," +
		"last_date TIMESTAMP," +
		"email VARCHAR(319) NOT NULL" +
		")"
}
