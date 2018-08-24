package model

import "encoding/gob"

func init() {
	gob.Register(Admin{})
}
