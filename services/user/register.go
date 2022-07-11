package user

import (
	"fmt"
	"zhigui/model"
)

func Register(Email, NickName, Password string) error {
	// 判断用户是否存在
	err := model.IfExist(Email, NickName)
	fmt.Println("err:", err)
	if err != nil {
		return err
	} else if err == nil {
		err1 := model.CreateAccount(Email, NickName, Password)
		return err1

	}
	return nil
}
