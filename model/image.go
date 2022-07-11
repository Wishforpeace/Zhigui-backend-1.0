package model

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type UserImage struct {
	gorm.Model
	Owner  string `gorm:"<-:create;type:varchar(30);unique"`
	Avatar string
	Sha    string
	Path   string
}
type TaskImage struct {
	gorm.Model
	TaskID uint
	Url    string
	Sha    string
	Path   string
}

type PostImage struct {
	gorm.Model
	PostID uint
	Url    string
	Sha    string
	Path   string
}

func UploadAvatar(email, picUrl, picSha, picPath string) (UserImage, error) {
	var Image UserImage
	return Image, DB.Model(&Image).Where("owner = ?", email).Updates(UserImage{Avatar: picUrl, Sha: picSha, Path: picPath}).Error
}

func GetImage(email string) (UserImage, error) {
	var image UserImage
	return image, DB.Model(&image).Where("owner = ?", email).Find(&image).Error
}

func GetTaskPicture(id string) (TaskImage, error) {
	var TaskImg TaskImage
	err := DB.Where("id = ?", id).First(&TaskImg).Error
	return TaskImg, err
}
func GetPayment(email string) (Payment, error) {
	var image Payment
	return image, DB.Model(&image).Where("email = ?", email).Find(&image).Error
}

func UploadPayment(email, picUrl, picSha, picPath string) (Payment, error) {
	var Image Payment
	return Image, DB.Model(&Image).Where("email = ?", email).Updates(Payment{Avatar: picUrl, Sha: picSha, Path: picPath}).Error
}

func (task *Task) UploadTask(ID, Url, Sha, Path string) error {
	u64, _ := strconv.ParseUint(ID, 10, 32)
	id := uint(u64)
	var image TaskImage = TaskImage{
		TaskID: id,
		Url:    Url,
		Sha:    Sha,
		Path:   Path,
	}
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&image).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := DB.Where("id = ?", ID).Find(task).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	var IMG []TaskImage
	if err := DB.Where("task_id = ?", ID).Find(&IMG).Error; err != nil {
		tx.Rollback()
		return err
	}
	task.Image = IMG
	return nil
}

func DeletePicture(id string) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&TaskImage{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

}

func GetPicTask(id string) (uint, error) {
	var image TaskImage
	err := DB.Where("id = ?", id).First(&image).Error
	return image.TaskID, err
}

func (img *PostImage) UploadPostPicture() error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(img).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
