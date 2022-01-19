package main

import (
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

func (a *Account) SetPasswordHash(plainPass string) {
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
