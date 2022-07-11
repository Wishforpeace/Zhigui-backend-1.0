package taskbar

import (
	"errors"
	"zhigui/model"
)

func GetAll(offset int, limit int) ([]*model.Task, int, error) {
	content, num, err := model.GetAllTasks(offset, limit)
	if err != nil {
		return nil, num, errors.New("没有任务")
	}
	return content, num, nil
}
