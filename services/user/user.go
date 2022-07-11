package user

import "zhigui/model"

type Info struct {
	Name  string
	Email string
	IMG   model.UserImage
}

func GetInfo(email string) (Info, error) {
	user, err := model.GetInfo(email)
	info := Info{
		Name:  user.NickName,
		Email: email,
		IMG:   user.Image,
	}
	return info, err
}
