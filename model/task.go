package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"zhigui/pkg/errno"
)

type Task struct {
	gorm.Model
	Prefecture string `json:"prefecture_id" gorm:"column:prefecture"`
	Tag        string `json:"tag" gorm:"size:10"`
	Content    string `json:"content" gorm:"size:255"`
	Publisher  string `json:"publisher" gorm:"<-:create"`
	Accepter   string `json:"accepter" gorm:"<-:create"`
	Award      int    `json:"award"`
	Status     string `json:"status" gorm:"column:status;size:10"`
	Method     string `json:"method" gorm:"size:10"`
	WillingNum int    `json:"willing_num" gorm:"column:willing_num"`
	Willing    []Willing
	Confirm    Confirm
	Image      []TaskImage `gorm:"foreignKey:TaskID;reference:ID"`
}

type Confirm struct {
	gorm.Model
	TaskID    uint
	Accepter  string `json:"accepter"`
	Publisher string `json:"publisher"`
}

type Willing struct {
	gorm.Model
	TaskID uint
	Email  string
}

// 获取个人接受或发布的任务 0---发布，1---接受
func GetPersonalTasks(variety int, email string, offset int, limit int) ([]*Task, int, error) {
	if variety == 0 {
		item := make([]*Task, 0)
		d := DB.Table("tasks").
			Where("publisher = ?", email).
			Offset(offset).Limit(limit).
			Order("created_at desc").Find(&item)
		fmt.Println("item", item)
		if d.Error != nil {
			return nil, 0, d.Error
		}
		var num int
		for i, _ := range item {
			num = i + 1
		}
		return item, num, d.Error
	} else if variety == 1 {
		item := make([]*Task, 0)
		d := DB.Table("tasks").
			Where("accepter = ?", email).
			Offset(offset).Limit(limit).
			Order("created_at desc").Find(&item)
		if d.Error != nil {
			return nil, 0, d.Error
		}
		var num int
		tasks := make([]*Task, 0)
		d = DB.Table("tasks").
			Where("accepter = ?", email).
			Order("created_at desc").Find(&tasks)
		if d.Error != nil {
			return nil, 0, d.Error
		}
		for i, _ := range tasks {
			num = i + 1
		}
		return item, num, d.Error
	}
	return nil, 0, errno.ErrDatabase
}

// 获取不同分区的任务
func GetTasks(PrefectureID, offset, limit int) ([]*Task, int, error) {
	item := make([]*Task, 0)
	var prefecture string
	switch PrefectureID {
	case 1:
		prefecture = "数理专区"
	case 2:
		prefecture = "英语专区"
	case 3:
		prefecture = "专业课专区"
	case 4:
		prefecture = "竞赛专区"
	case 5:
		prefecture = "体育专区"
	case 6:
		prefecture = "游戏专区"
	case 7:
		prefecture = "赏乐专区"
	case 8:
		prefecture = "吃喝专区"

	}
	d := DB.Table("tasks").
		Where(" prefecture = ?", prefecture).
		Offset(offset).Limit(limit).
		Order("created_at desc").Find(&item)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	var num int
	for i, m := range item {
		var TaskImage []TaskImage
		DB.Where("task_id = ?", m.ID).Find(&TaskImage)
		item[i].Image = TaskImage

	}
	tasks := make([]*Task, 0)
	d = DB.Table("tasks").
		Where(" prefecture = ?", prefecture).
		Order("created_at desc").Find(&tasks)
	if d.Error != nil {
		return nil, 0, d.Error
	}
	for i, _ := range tasks {
		num = i + 1
	}
	return item, num, d.Error

}

// 获取指定任务的内容
func (task *Task) GetTask(id int, email string) (string, error) {
	var MyWilling string = "未发送意愿"
	if err := DB.Where("id = ?", id).First(&task).Error; err != nil {
		return "", err
	}
	var image []TaskImage
	if err := DB.Where("task_id = ?", id).Find(&image).Error; err != nil {
		return "", err
	}

	var confirm Confirm
	if err := DB.Where("task_id=?", id).First(&confirm).Error; err != nil {
		return "", err
	}
	var will []Willing

	if err := DB.Where("task_id=?", id).Find(&will).Error; err != nil {
		return "", err
	}
	for _, m := range will {
		if m.Email == email {
			MyWilling = "已发送意愿"
		}
	}

	task.Image = image
	task.Confirm = confirm
	task.Willing = will
	return MyWilling, nil
}

// 发布任务
func PublishTask(email, Prefecture, Tag, Method, Content string, Award int) (Task, error) {
	var publish Task
	var confirm Confirm
	if Method == "0" {
		publish.Method = "线上"
	}
	if Method == "1" {
		publish.Method = "线下"
	}
	switch Prefecture {
	case "1":
		publish.Prefecture = "数理专区"
		switch Tag {
		case "1":
			publish.Tag = "高等数学"
		case "2":
			publish.Tag = "线性代数"
		case "3":
			publish.Tag = "概率论"
		case "4":
			publish.Tag = "复变函数论"
		case "5":
			publish.Tag = "数学分析"
		default:
			publish.Tag = "其他"
		}
	case "2":
		publish.Prefecture = "英语专区"
		switch Tag {
		case "1":
			publish.Tag = "四六级"
		case "2":
			publish.Tag = "专四专八"
		case "3":
			publish.Tag = "考研"
		case "4":
			publish.Tag = "雅思托福"
		case "5":
			publish.Tag = "口语"
		default:
			publish.Tag = "其他"
		}
	case "3":
		publish.Prefecture = "专业课专区"
		switch Tag {
		case "1":
			publish.Tag = "学长学姐咨询"
		case "2":
			publish.Tag = "学霸分享"
		case "3":
			publish.Tag = "考试技巧"
		case "4":
			publish.Tag = "上岸经验"
		case "5":
			publish.Tag = "往年试题"
		default:
			publish.Tag = "其他"
		}
	case "4":
		publish.Prefecture = "竞赛专区"
		switch Tag {
		case "1":
			publish.Tag = "挑战杯"
		case "2":
			publish.Tag = "数学建模"
		case "3":
			publish.Tag = "创新创业"
		case "4":
			publish.Tag = "计算机设计"
		case "5":
			publish.Tag = "CCPC/ICPC"
		default:
			publish.Tag = "其他"
		}
	case "5":
		publish.Prefecture = "体育专区"
		switch Tag {
		case "1":
			publish.Tag = "田径"
		case "2":
			publish.Tag = "篮球"
		case "3":
			publish.Tag = "足球"
		case "4":
			publish.Tag = "羽毛球"
		case "5":
			publish.Tag = "乒乓球"
		default:
			publish.Tag = "其他"
		}
	case "6":
		publish.Prefecture = "游戏专区"
		switch Tag {
		case "1":
			publish.Tag = "陪玩"
		case "2":
			publish.Tag = "刷级"
		case "3":
			publish.Tag = "上分"
		case "4":
			publish.Tag = "组团开黑"
		case "5":
			publish.Tag = "代肝"
		default:
			publish.Tag = "其他"
		}
	case "7":
		publish.Prefecture = "赏乐专区"
		switch Tag {
		case "1":
			publish.Tag = "周末出游"
		case "2":
			publish.Tag = "寒暑假出游"
		case "3":
			publish.Tag = "组团出游"
		default:
			publish.Tag = "其他"
		}
	case "8":
		publish.Prefecture = "吃喝专区"
		switch Tag {
		case "1":
			publish.Tag = "外卖拼单"
		case "2":
			publish.Tag = "带饭"
		case "3":
			publish.Tag = "跑腿"
		case "4":
			publish.Tag = "代购"
		default:
			publish.Tag = "其他"
		}

	}
	publish.Publisher = email
	publish.Content = Content
	publish.Award = Award
	publish.WillingNum = 0
	publish.Status = "Unaccepted"

	tx := DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return publish, err
	}

	if err := tx.Create(&publish).Error; err != nil {
		tx.Rollback()
		return publish, err
	}
	confirm.Publisher = "未完成"
	confirm.Accepter = "未完成"
	confirm.TaskID = publish.Model.ID
	if err := tx.Create(&confirm).Error; err != nil {
		tx.Rollback()
		return publish, err
	}

	err := tx.Commit().Error
	publish.Confirm = confirm
	return publish, err

}

// 更新任务
func UpdateTask(id, PrefectureId, TagID, Method, content string, award int) (Task, error) {
	task := Task{}
	if Method == "0" {
		task.Method = "线上"
	}
	if Method == "1" {
		task.Method = "线下"
	}
	Id, _ := strconv.ParseUint(id, 10, 32)
	task.Model.ID = uint(Id)
	if PrefectureId != "" {

		switch PrefectureId {

		case "1":
			task.Prefecture = "数理专区"
			switch TagID {
			case "1":
				task.Tag = "高等数学"
			case "2":
				task.Tag = "线性代数"
			case "3":
				task.Tag = "概率论"
			case "4":
				task.Tag = "复变函数论"
			case "5":
				task.Tag = "数学分析"
			}
		case "2":
			task.Prefecture = "英语专区"

			switch TagID {
			case "1":
				task.Tag = "四六级"
			case "2":
				task.Tag = "专四专八"
			case "3":
				task.Tag = "考研"
			case "4":
				task.Tag = "雅思托福"
			case "5":
				task.Tag = "口语"

			}
		case "3":
			task.Prefecture = "专业课专区"

			switch TagID {
			case "1":
				task.Tag = "学长学姐咨询"
			case "2":
				task.Tag = "学霸分享"
			case "3":
				task.Tag = "考试技巧"
			case "4":
				task.Tag = "上岸经验"
			case "5":
				task.Tag = "往年试题"

			}
		case "4":
			task.Prefecture = "竞赛专区"

			switch TagID {
			case "1":
				task.Tag = "挑战杯"
			case "2":
				task.Tag = "数学建模"
			case "3":
				task.Tag = "创新创业"
			case "4":
				task.Tag = "计算机设计"
			case "5":
				task.Tag = "CCPC/ICPC"

			}
		case "5":
			task.Prefecture = "体育专区"

			switch TagID {
			case "1":
				task.Tag = "田径"
			case "2":
				task.Tag = "篮球"
			case "3":
				task.Tag = "足球"
			case "4":
				task.Tag = "羽毛球"
			case "5":
				task.Tag = "乒乓球"

			}
		case "6":
			task.Prefecture = "游戏专区"

			switch TagID {
			case "1":
				task.Tag = "陪玩"
			case "2":
				task.Tag = "刷级"
			case "3":
				task.Tag = "上分"
			case "4":
				task.Tag = "组团开黑"
			case "5":
				task.Tag = "代肝"
			default:
				task.Tag = "其他"

			}
		case "7":
			task.Prefecture = "赏乐专区"

			switch TagID {
			case "1":
				task.Tag = "周末出游"
			case "2":
				task.Tag = "寒暑假出游"
			case "3":
				task.Tag = "组团出游"
			}

		case "8":
			task.Prefecture = "吃喝专区"

			switch TagID {
			case "1":
				task.Tag = "外卖拼单"
			case "2":
				task.Tag = "带饭"
			case "3":
				task.Tag = "跑腿"
			case "4":
				task.Tag = "代购"

			}
		}
	}
	task.Content = content
	task.Award = award
	DB.Model(&Task{}).Where("id = ?", id).Updates(&task)
	TaskContent := Task{}
	return TaskContent, DB.Where("id = ?", id).First(&TaskContent).Error

}

// 删除任务
func DeleteTask(id string) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Where("id = ?", id).Delete(&Task{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func AccpetTask(email string, id string) error {
	var task Task
	if err := DB.Where("id = ?", id).First(&task).Error; err != nil {
		return err
	}
	var will1 Willing
	if err := DB.Where("task_id = ? and email = ?", id, email).First(&will1).Error; err == nil {
		return errors.New("无法重复接受")
	}

	task.WillingNum += 1
	ID, _ := strconv.Atoi(id)
	Id := uint(ID)
	var will = Willing{
		TaskID: Id,
		Email:  email,
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
	if err := tx.Model(Task{}).Where("id = ?", id).Updates(Task{WillingNum: task.WillingNum}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&will).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func ConfirmTask(id, user string) error {
	switch user {
	case "0":
		return DB.Model(&Confirm{}).Where("task_id=?", id).Updates(Confirm{Accepter: "已完成"}).Error
	case "1":
		return DB.Model(&Confirm{}).Where("task_id=?", id).Updates(Confirm{Publisher: "已完成"}).Error

	}
	return nil
}

func (task *Task) UpdateStatus(accepter string) error {
	var user User
	DB.Where("email = ?", accepter).First(&user)
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	user.Done += 1
	if err := tx.Model(&task).Where("id = ?", task.ID).
		Updates(Task{Status: "Unpaid"}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&User{}).Where("email = ?", accepter).
		Updates(User{Done: user.Done}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func GetWills(id string, email string) error {
	var will Willing
	err := DB.Where("task_id = ? and email = ?", id, email).First(&will).Error
	return err

}

func SelectAccepter(accepter string, id string) error {
	var user User
	DB.Where("email = ?", accepter).First(&user)
	user.Doing += 1
	var task Task
	DB.Where("id = ?", id).First(&task)
	if task.Accepter != "" {
		return errors.New("无法重复选择")
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
	if err := tx.Model(&User{}).Where("email = ?", accepter).Updates(User{Doing: user.Doing}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&Task{}).Where("id = ?", id).Updates(Task{Accepter: accepter, Status: "Accepted"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetAllTasks(offset int, limit int) ([]*Task, int, error) {
	var num int
	var tasks []Task
	DB.Table("tasks").Select("tasks.*").Scan(&tasks)
	for i, _ := range tasks {
		num = i + 1
	}
	item := make([]*Task, 0)
	d := DB.Table("tasks").
		Select("tasks.*").
		Offset(offset).Limit(limit).
		Order("created_at desc").Scan(&item)
	for i, m := range item {
		var confirm Confirm
		DB.Where("task_id=?", m.ID).First(&confirm)
		item[i].Confirm = confirm
		var image []TaskImage
		DB.Where("task_id=?", m.ID).Find(&image)
		item[i].Image = image
	}
	if d.Error != nil {
		return nil, num, d.Error
	}
	return item, num, nil
}

func ConfirmPaid(id string, email string) error {
	var task Task
	DB.Where("id = ?", id).First(&task)
	var user User
	DB.Where("email = ?", email).First(&user)
	if task.Status == "Finished" {
		return errors.New("交易已完成")
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
	if err := tx.Model(&Task{}).Where("id = ?", id).Updates(Task{Status: "Finished"}).Error; err != nil {
		tx.Rollback()
		return err
	}
	earn, _ := strconv.Atoi(user.Earning)
	award := task.Award + earn
	earning := strconv.Itoa(award)
	if err := tx.Model(&User{}).Where("email = ?", email).Update(User{Earning: earning}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
