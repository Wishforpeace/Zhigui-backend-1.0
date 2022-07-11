package model

import "github.com/jinzhu/gorm"

type Prefecture struct {
	gorm.Model
	Name string `json:"name"gorm:"size:20;unique"`
}
