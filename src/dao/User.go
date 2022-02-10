package dao

import (
	"database/sql"
)

type User struct {
	ID             string
	UserName       string
	UserPwd        string
	Email          string
	Phone          string
	EmailValidated string
	PhoneValidated string
	SignupAt       string
	LastActive     string
	Profile        sql.NullString
	Status         string
	UserToken      string
}
