package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Profile struct {
	gorm.Model
	IdUser         uint      `json:"id_user"`
	Avatar         string    `json:"avatar" gorm:"type:longtext"` //longtext no BD (mysql-MariaDB)
	DataNascimento time.Time `json:"datanascimento" `             //maximo 8 digitos
}
type ProfileInput struct {
	Avatar         string    `json:"avatar" gorm:"type:longtext"` //longtext no BD (mysql-MariaDB)
	DataNascimento time.Time `json:"datanascimento" `             //maximo 8 digitos
}

/*
func (p *Profile) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":

	case "login":

	default:

	}
	return nil
}
*/