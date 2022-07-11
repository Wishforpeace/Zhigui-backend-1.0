package user

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	response "zhigui/handler"
	"zhigui/model"
	"zhigui/pkg/errno"
	"zhigui/services"
	"zhigui/services/connector"
	"zhigui/services/random"
	"zhigui/services/user"
)

type UpdateInfoRequest struct {
	PhoneNum    string `json:"phone_num"`
	NickName    string `json:"nick_name"`
	Degree      string `json:"degree"`
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}
type InfoResponse struct {
	Email    string `json:"email" `
	PhoneNum string `json:"phone_num"`
	NickName string `json:"nick_name"`
	Gender   string `json:"gender" `
	Degree   string `json:"degree" `
	Avatar   string
	Earning  string
	Doing    int
	Done     int
}

// @Summary 得到用户信息
// @Description 得到用户所有的个人信息
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {object} InfoResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user/info [get]
func GetInfo(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	if details, err := model.GetInfo(email); err != nil {
		response.SendBadRequest(c, nil, nil, errno.ErrDatabase)
		return
	} else {
		var info InfoResponse
		info = InfoResponse{
			Email:    email,
			PhoneNum: details.PhoneNum,
			NickName: details.NickName,
			Gender:   details.Gender,
			Degree:   details.Degree,
			Avatar:   details.Image.Avatar,
			Earning:  details.Earning,
			Doing:    details.Doing,
			Done:     details.Done,
		}
		response.SendResponse(c, "获取成功", info)
		return
	}

}

// @Summary 查看其他用户信息
// @Description 得到其他用户的个人信息
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "token"
// @Param email query string true "用户邮箱"
// @Success 200 {object} InfoResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user/others [get]
func GetOtherInfo(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.Query("email")
	Info, err := user.GetInfo(email)
	if err != nil {
		response.SendError(c, "获取失败", err.Error(), errno.ErrDatabase)
		return
	}
	response.SendResponse(c, "获取成功", Info)
}

// @Summary 修改用户信息
// @Description 修改用户所有的个人信息
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "token"
// @Param req body UpdateInfoRequest true "需要修改对内容"
// @Success 200 {string} string  "{"NickName":"新名字","Degree":"新学历"}"
// @Router /user/info [post]
func UpdateInfo(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	var req UpdateInfoRequest
	if err := c.ShouldBind(&req); err != nil {
		response.SendBadRequest(c, "请确认修改信息", nil, errno.ErrBind)
		return
	}
	if req.Password != "" {
		user, _ := model.GetInfo(email)
		password, _ := base64.StdEncoding.DecodeString(user.Password)
		if string(password) != req.OldPassword {
			response.SendBadRequest(c, "旧密码错误", nil, errno.ErrBind)
		}
	}
	if err := model.UpdateUserInfo(email, req.PhoneNum, req.NickName, req.Degree, req.Password); err != nil {
		response.SendError(c, "修改失败", nil, errno.ErrDatabase)
		return
	}
	details, _ := model.GetInfo(email)

	response.SendResponse(c, "修改成功", map[string]interface{}{
		"PhoneNumber": details.PhoneNum,
		"NickName":    details.NickName,
		"Degree":      details.Degree,
	})
}

// @Summary	修改头像
// @Tags user
// @Description 修改用户头像
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token"
// @Param file formData file true "文件"
// @Success 200 {string} string "{"message":"上传成功","data":map[string]interface{"url","path","sha"}}"
// @Failure 400 {object} errno.Errno "上传失败"
// @Failure 400 {object} errno.Errno "上传失败,请检查token与其他配置参数是否正确"
// @Router /user/avatar [post]
func UploadProfile(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	file, err := c.FormFile("file")
	PATH := "Users"
	//image := c.PostForm("image")
	if err != nil {
		response.SendBadRequest(c, "上传失败", nil, errno.ErrBind)
		return
	}
	filepath := "./"
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	fileExt := path.Ext(filepath + file.Filename)

	file.Filename = email + random.GetRandomString(16) + fileExt

	filename := filepath + file.Filename

	if err := c.SaveUploadedFile(file, filename); err != nil {
		response.SendBadRequest(c, "上传失败", nil, errno.ErrDatabase)
		return
	}

	// 删除原头像
	image, _ := model.GetImage(email)
	fmt.Println("image", image.Path, image.Sha)
	if image.Avatar != "https://cdn.jsdelivr.net/gh/Wishforpeace/ZhiguiImage@master/Users/user.png" {
		if image.Path != "" && image.Sha != "" {
			connector.RepoCreate().Del(image.Path, image.Sha)
		}
	}

	// 上传新头像
	Base64 := services.ImagesToBase64(filename)
	//删除Base64 传入image
	picUrl, picPath, picSha := connector.RepoCreate().Push(PATH, file.Filename, Base64)
	//picUrl, picPath, picSha := connector.RepoCreate().Push(file.Filename, image)
	os.Remove(filename)

	_, e := model.UploadAvatar(email, picUrl, picSha, picPath)
	log.Println("picURL", picUrl)
	log.Println("err:", e)
	if picUrl == "" || e != nil {
		response.SendBadRequest(c, "上传失败,请检查token与其他配置参数是否正确", nil, errno.ErrPermissionDenied)
		return
	}

	response.SendResponse(c, "上传成功", map[string]interface{}{
		"url":  picUrl,
		"sha":  picSha,
		"path": picPath,
	})

}

// @Summary	上传收款码
// @Tags user
// @Description 修改用户收款码
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token"
// @Param file formData file true "文件"
// @Success 200 {string} string "{"message":"上传成功","data":map[string]interface{"url","path","sha"}}"
// @Failure 400 {object} errno.Errno "上传失败"
// @Failure 400 {object} errno.Errno "上传失败,请检查token与其他配置参数是否正确"
// @Router /user/payment [post]
func UploadPayment(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	file, err := c.FormFile("file")
	PATH := "Payment"
	//image := c.PostForm("image")
	if err != nil {
		response.SendBadRequest(c, "上传失败", nil, errno.ErrBind)
		return
	}
	filepath := "./"
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	fileExt := path.Ext(filepath + file.Filename)

	file.Filename = email + "(for Payment)" + fileExt

	filename := filepath + file.Filename

	if err := c.SaveUploadedFile(file, filename); err != nil {
		response.SendBadRequest(c, "上传失败", nil, errno.ErrDatabase)
		return
	}

	// 删除原头像
	image, _ := model.GetPayment(email)
	if image.Path != "" && image.Sha != "" {
		connector.RepoCreate().Del(image.Path, image.Sha)
	}

	// 上传新头像
	Base64 := services.ImagesToBase64(filename)
	//删除Base64 传入image
	picUrl, picPath, picSha := connector.RepoCreate().Push(PATH, file.Filename, Base64)
	//picUrl, picPath, picSha := connector.RepoCreate().Push(file.Filename, image)
	os.Remove(filename)

	_, e := model.UploadPayment(email, picUrl, picSha, picPath)
	log.Println("picURL", picUrl)
	log.Println("err:", e)
	if picUrl == "" || e != nil {
		response.SendBadRequest(c, "上传失败,请检查token与其他配置参数是否正确", nil, errno.ErrPermissionDenied)
		return
	}

	response.SendResponse(c, "上传成功", map[string]interface{}{
		"url":  picUrl,
		"sha":  picSha,
		"path": picPath,
	})

}
