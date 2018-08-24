package model

type Admin struct {
	ID       int    `storm:"unique,increment=10000"`
	Email    string `storm:"unique"`
	Password string
}
