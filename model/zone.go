package model

import "github.com/jinzhu/gorm"

type Zone struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}
