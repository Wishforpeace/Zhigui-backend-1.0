package tasks

import (
	"errors"
	"log"
	"strconv"
	"zhigui/model"
)

type PublishedTask struct {
	ID         uint
	Prefecture string
	Tag        string
	Content    string
	Award      int
	Publisher  string
	Method     string
}

func Publish(email, PrefectureID, Tag, Method, Content string, Award int) (PublishedTask, error) {
	var published PublishedTask

	Task, err := model.PublishTask(email, PrefectureID, Tag, Method, Content, Award)
	log.Println("Task", Task)
	log.Println("err", err)
	if err != nil {
		return published, err
	} else {
		published.ID = Task.ID
		published.Publisher = Task.Publisher
		published.Prefecture = Task.Prefecture
		published.Tag = Task.Tag
		published.Content = Task.Content
		published.Award = Task.Award
		published.Method = Task.Method
		log.Println(published.ID)
		return published, err
	}
}

func UpdateTask(ID, PrefectureID, TagID, Metohd, Content string, Award int) (PublishedTask, error) {
	var publish PublishedTask
	task, err := model.UpdateTask(ID, PrefectureID, TagID, Metohd, Content, Award)
	if err != nil {
		return publish, err
	} else {
		publish.ID = task.Model.ID
		publish.Prefecture = task.Prefecture
		publish.Tag = task.Tag
		publish.Content = task.Content
		publish.Award = task.Award
		publish.Publisher = task.Publisher

		return publish, nil
	}

}

func AccpetTask(email string, ID string) error {
	var task model.Task
	id, _ := strconv.Atoi(ID)
	task.GetTask(id, "")
	if task.Publisher == email {
		return errors.New("不能接受自己发布的任务")
	}
	if err := model.AccpetTask(email, ID); err != nil {
		return err
	}
	return nil

}

func ConfirmFinish(email string, id string) error {
	var task model.Task
	ID, _ := strconv.Atoi(id)
	task.GetTask(ID, "")
	var err error
	if task.Accepter == email {
		user := "0"
		err = model.ConfirmTask(id, user)
	}
	if task.Publisher == email {
		user := "1"
		err = model.ConfirmTask(id, user)

	}
	if task.Accepter != email && task.Publisher != email {
		return errors.New("没有资格")
	}
	task.GetTask(ID, email)
	if task.Confirm.Accepter == "已完成" && task.Confirm.Publisher == "已完成" {
		err = task.UpdateStatus(task.Accepter)

	}
	return err
}

func Payment(id string, email string) (model.Payment, error) {
	var task model.Task
	ID, _ := strconv.Atoi(id)
	task.GetTask(ID, "")
	if task.Publisher != email {
		return model.Payment{}, errors.New("没有资格")
	}
	if task.Status != "Finished" {
		return model.Payment{}, errors.New("尚未完成")
	}

	payment, err := model.GetPayment(task.Accepter)
	if err != nil {
		return model.Payment{}, err

	}
	return payment, nil

}

func SelectAccepter(email string, accepter string, id string) error {
	var task model.Task
	ID, _ := strconv.Atoi(id)
	task.GetTask(ID, email)
	if task.Publisher != email {
		return errors.New("您不是改任务的创建者")
	}

	err := model.GetWills(id, accepter)
	if err != nil {
		return errors.New("该用户未加入意愿")
	}
	if err := model.SelectAccepter(accepter, id); err != nil {
		return errors.New("选择失败")
	}
	return nil
}

func ConfirmPaid(id string, email string) error {
	var task model.Task
	ID, _ := strconv.Atoi(id)
	task.GetTask(ID, email)
	if task.Accepter != email {
		return errors.New("没有权限确认")
	}
	err := model.ConfirmPaid(id, email)
	if err != nil {
		return err
	}
	return nil
}
