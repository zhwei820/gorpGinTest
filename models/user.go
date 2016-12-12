package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

type User struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Email   string `db:"email" json:"mail"`
	Status  string `db:"status" json:"status"`
	Comment string `db:"comment, size:16384" json:"comment"`
	Pass    string `db:"pass" json:"pass"`
	Created string `db:"created" json:"created"` // or int64
	Updated string `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

// PreInsert set created an updated time before insert in db
func (a *User) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now().Format("2006-01-02 15:04:05") // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

// PreUpdate set updated time before insert in db
func (a *User) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now().Format("2006-01-02 15:04:05")
	return nil
}
