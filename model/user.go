package model

import (
	"encoding/base64"
	"errors"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email      string `json:"email" gorm:"size:100;not null;unique"`
	PhoneNum   string `json:"phone_num"gorm:"size:11;null"`
	NickName   string `json:"nick_name"gorm:"size:20;not null;unique"`
	Password   string `json:"password" gorm:"required;not null"`
	Gender     string `json:"gender"gorm:"null"`
	Degree     string `json:"degree"  gorm:"size:10; null"`
	Publisheds []Task `gorm:"foreignKey:Pubilsher;reference:Email"`
	Accpeteds  []Task `gorm:"foreignKey:Accepter;reference:Email"`
	Doing      int
	Done       int
	School     string `gorm:"null"`
	StudentId  string `gorm:"null"`
	Payment    Payment
	Image      UserImage `gorm:"foreignKey:Owner;reference:Email"`
	Earning    string    `gorm:"size:10"`
}

// 查看用户是否已被创建
func IfExist(email, name string) (err error) {
	var user User
	err1 := DB.Where("email=?", email).First(&user).Error
	err2 := DB.Where("nick_name=?", name).First(&user).Error
	if err1 == nil && err2 == nil {
		return errors.New("邮箱以及用户名均被注册")
	}
	if err1 == nil && err2 != nil {
		return errors.New("邮箱已被注册")
	}
	if err1 != nil && err1 == nil {
		return errors.New("用户名已被注册")
	}
	if err1 != nil && err2 != nil {
		return nil
	}
	return nil
}

func CreateAccount(email, nickname, password string) error {
	var user User
	user.Email = email
	user.NickName = nickname
	user.Degree = ""
	user.Gender = "未知"
	user.Earning = "0"
	var image UserImage
	image.Owner = email
	image.Avatar = "https://cdn.jsdelivr.net/gh/Wishforpeace/ZhiguiImage@master/Users/user.png"
	var payment Payment
	payment.Email = email
	user.Password = base64.StdEncoding.EncodeToString([]byte(password))
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&image).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
	return nil
}

func GetInfo(email string) (User, error) {
	var user User
	var image UserImage
	user.Image.Avatar = image.Avatar
	DB.Where("owner = ?", email).First(&image)
	err := DB.Where("email = ?", email).First(&user).Error
	user.Image = image
	return user, err
}

func UpdateUserInfo(email, phone, nickname, degree, password string) error {
	PassWord := base64.StdEncoding.EncodeToString([]byte(password))
	return DB.Model(&User{}).Where("email = ?", email).Update(User{
		NickName: nickname,
		PhoneNum: phone,
		Degree:   degree,
		Password: PassWord,
	}).Error
}
