package like

import (
	"zhigui/model"
)

func GiveLike(id string, email string) (int, error) {
	err := model.CheckLikeList(id, email)
	if err == nil {
		return 0, err
	} else if err != nil {
		e := model.CreateLike(id, email)
		return 1, e
	}
	return 2, nil
}
