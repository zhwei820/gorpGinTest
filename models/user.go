package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

type User struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Email   string `db:"email" json:"email"`
	Status  string `db:"status" json:"status"`
	Comment string `db:"comment, size:16384" json:"comment"`
	Pass    string `db:"pass" json:"pass"`
	Created string `db:"created" json:"created"`
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

type ZUser struct {
	Uid        int64  `db:"uid" json:"uid"`
	Pnum       int64  `db:"pnum" json:"pnum"`
	PnumMd5    string `db:"pnum_md5" json:"pnummd5"`
	Password   string `db:"password" json:"password"`
	Status     int64  `db:"status" json:"status"`
	DeviceId   string `db:"device_id" json:"deviceid"`
	Imsi       string `db:"imsi" json:"imsi"`
	OsType     string `db:"os_type" json:"ostype"`
	Ctime      string `db:"ctime" json:"ctime"`
	RegisterIp string `db:"register_ip" json:"registerip"`
	InviteCode int64  `db:"invite_code" json:"invitecode"`
	Channel    string `db:"channel" json:"channel"`
	Ulevel     int64  `db:"ulevel" json:"ulevel"`
	FromApp    int64  `db:"from_app" json:"fromapp"`
	UpdateTime string `db:"update_time" json:"updatetime"`
}
