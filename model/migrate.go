package model

import (
	"github.com/jinzhu/gorm"
)

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(&User{}, &Zone{}, &Prefecture{}, &Payment{}, &Task{}, &Willing{}, &UserImage{}, &TaskImage{}, &Confirm{}, &POST{}, &Comment{}, &PostImage{}, &LikeList{})
	var zone1 = Zone{
		Name: "学习区",
	}
	var zone2 = Zone{
		Name: "娱乐区",
	}
	DB.Create(&zone1)
	DB.Create(&zone2)
}
