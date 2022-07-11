package model

import "github.com/jinzhu/gorm"

type Payment struct {
	gorm.Model
	Email  string `gorm:"size:100;unique"`
	Avatar string
	Sha    string
	Path   string
}
