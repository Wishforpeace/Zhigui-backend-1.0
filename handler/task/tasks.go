package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	"strconv"
	"time"
	"zhigui/handler"
	"zhigui/model"
	"zhigui/pkg/errno"
	"zhigui/services"
	"zhigui/services/connector"
	"zhigui/services/random"
	"zhigui/services/tasks"
)

type TaskRequest struct {
	PrefectureID string `json:"prefecture_id" binding:"required"`
	TagID        string `json:"tag_id"binding:"required"`
	Method       string `json:"method"binding:"required"`
	Content      string `json:"content" binding:"required"`
	Award        int    `json:"award" binding:"required"`
}

type UpdateRequest struct {
	ID           string `json:"id"`
	TagID        string `json:"tag_id"`
	PrefectureID string `json:"prefecture_id"`
	Method       string `json:"method"`
	Content      string `json:"content"`
	Award        int    `json:"award"`
}

// @Summary 查看用户的发送接受的全部任务
// @Tags task
// @Description 查看发布与接受的任务
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token"
// @Param task query integer true "task:0---查看发布的任务；task:1---查看接受的任务"
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query  integer true "page--限制指定要检索的记录数"
// @Success 200 {object} []model.Task{} "{"msg":"查看成功"}"
// @Failure 500 {object} errno.Errno "{"msg":"Error occurred while getting url queries."}"
// @Router /tasks/personal [get]
func GetPersonalTasks(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	task := c.DefaultQuery("task", "0")
	var limit, page int
	var err error
	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}

	page, err = strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}
	variety, _ := strconv.Atoi(task)
	Tasks, num, err := model.GetPersonalTasks(variety, email, limit*page, limit)
	if err != nil {
		handler.SendError(c, "查询失败", task, errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "查询成功", map[string]interface{}{
		"NUM":   num,
		"Tasks": Tasks,
	})
}

// @Summary 获取专区的任务
// @Tags task
// @Description 查看分区的任务
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token"
// @Param prefecture_id query integer true "0:数理专区；1：英语专区；2：专业课专区；3：竞赛专区；4：体育专区；5：游戏专区；6：赏乐专区；7：吃喝专区"
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query integer true "page--限制指定要检索的记录数 "
// @Success 200 {object} []model.Task{} "{"msg":"查看成功"}"
// @Failure 500 {object} errno.Errno "{"msg":"Error occurred while getting url queries."}"
// @Router /tasks/prefecture [get]
func GetTasks(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var PrefectureID int
	var err error
	PrefectureID, err = strconv.Atoi(c.DefaultQuery("prefecture_id", "1"))
	if err != nil {
		handler.SendError(c, "内容缺失", nil, errno.ErrQuery)
		return
	}
	var limit, page int

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}

	page, err = strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}
	task, num, er := model.GetTasks(PrefectureID, limit*page, limit)
	if er != nil {
		handler.SendError(c, nil, er.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "获取成功", map[string]interface{}{
		"NUM":   num,
		"Tasks": task,
	})
}

// @Summary 查看某任务具体内容
// @Description 查看每一个任务的具体内容
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id query string true "id--任务的id"
// @Success 200 {object}  []model.Task{} "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/details [get]
func GetTaskDetails(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	ID := c.Query("id")
	id, _ := strconv.Atoi(ID)
	var task model.Task

	if will, err := task.GetTask(id, email); err != nil {
		handler.SendError(c, "数据库查找失败", nil, errno.ErrDatabase)
		return
	} else {
		handler.SendResponse(c, "获取成功", map[string]interface{}{
			"MyWilli": will,
			"Task":    task,
		})
		return
	}
}

// @Summary 发布任务
// @Description 发布具体任务内容
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param task body TaskRequest  true "id--任务的id"
// @Success 200 {object}  []model.Task{} "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/publish [post]
func PublishTasks(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	var Request TaskRequest
	if err := c.Bind(&Request); err != nil {
		handler.SendBadRequest(c, "输入内容有误", nil, errno.ErrBind)
		return
	}
	task, err := tasks.Publish(email, Request.PrefectureID, Request.TagID, Request.Method, Request.Content, Request.Award)
	if err != nil {
		handler.SendError(c, "创建失败", task, errno.ErrDatabase)
		return
	}

	handler.SendResponse(c, "创建", map[string]interface{}{
		"Task": task,
	})
}

// @Summary 修改任务
// @Description 修改具体任务内容
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param object body UpdateRequest  true "id--任务的id"
// @Success 200 {object}  []model.Task{} "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/details [post]
func UpdateTask(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var Request UpdateRequest
	if err := c.Bind(&Request); err != nil {
		handler.SendBadRequest(c, "内容有误", Request, errno.ErrBind)
		return
	}
	task, err := tasks.UpdateTask(Request.ID, Request.PrefectureID, Request.TagID, Request.Method, Request.Content, Request.Award)
	if err != nil {
		handler.SendError(c, "更新失败", task, errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "修改成功", task)

}

// @Summary 上传或修改任务图片
// @Description 上传新的图片
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData string   true "id--任务的id"
// @Param file formData file true "文件"
// @Success 200 {object}  []model.Task{} "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/publish/picture [post]
func UploadPicture(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	file, err := c.FormFile("file")
	ID := c.PostForm("id")
	PATH := "Tasks"
	//image := c.PostForm("image")
	if err != nil {
		handler.SendBadRequest(c, "上传失败", nil, errno.ErrBind)
		return
	}
	filepath := "./"
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	fileExt := path.Ext(filepath + file.Filename)
	timeStr := time.Now().Format("20060102150405")
	file.Filename = ID + " of Tasks " + timeStr + " " + random.GetRandomString(16) + fileExt

	filename := filepath + file.Filename

	if err := c.SaveUploadedFile(file, filename); err != nil {
		handler.SendBadRequest(c, "上传失败", nil, errno.ErrDatabase)
		return
	}

	// 上传新头像
	Base64 := services.ImagesToBase64(filename)
	//删除Base64 传入image
	picUrl, picPath, picSha := connector.RepoCreate().Push(PATH, file.Filename, Base64)
	//picUrl, picPath, picSha := connector.RepoCreate().Push(file.Filename, image)
	os.Remove(filename)

	var TASK model.Task
	e := TASK.UploadTask(ID, picUrl, picSha, picPath)
	log.Println("picURL", picUrl)
	log.Println("err:", e)
	if picUrl == "" || e != nil {
		handler.SendBadRequest(c, "上传失败,请检查token与其他配置参数是否正确", nil, errno.ErrPermissionDenied)
		return
	}

	handler.SendResponse(c, "上传成功", map[string]interface{}{
		"Task": TASK,
	})
}

// @Summary 删除任务
// @Description 删除用户发布的任务
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id query string   true "id--任务的id"
// @Success 200 {object} model.Task  "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/delete/:id [delete]
func DeleteTask(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.Query("id")
	fmt.Println("id", id)
	email := c.MustGet("email").(string)
	var task model.Task
	ID, _ := strconv.Atoi(id)
	task.GetTask(ID, email)
	fmt.Println("email", task.Publisher)
	if task.Publisher != email {
		handler.SendError(c, "无资格删除", nil, errno.ErrDatabase)
		return
	}
	err := model.DeleteTask(id)
	if err != nil {
		handler.SendError(c, "删除失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "删除成功", nil)
}

// @Summary 删除图片
// @Description 删除已经任务图片
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true "获取email"
// @Param id query int  true "id--图片的id"
// @Success 200 {object} model.Task  "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/pictures [delete]
func DeletePictures(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.Query("id")
	email := c.MustGet("email").(string)
	Id, err := model.GetPicTask(id)
	var ID int
	ID = int(Id)
	var task model.Task

	pic, err := model.GetTaskPicture(id)
	if err != nil {
		handler.SendError(c, "不存在该任务", nil, errno.ErrDatabase)
		return
	}
	task.GetTask(ID, email)
	if task.Publisher != email {
		handler.SendError(c, "无资格删除", nil, errno.ErrDatabase)
		return
	}
	connector.RepoCreate().Del(pic.Path, pic.Sha)
	err = model.DeletePicture(id)
	if err != nil {
		handler.SendError(c, "删除失败", nil, errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "修改成功", nil)
}

// @Summary 接受任务
// @Description 接受他人发布的
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param task_id formData string   true "id--任务的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/accept [post]
func AcceptTask(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	id := c.PostForm("task_id")

	err := tasks.AccpetTask(email, id)
	if err != nil {
		handler.SendBadRequest(c, "接受任务失败", nil, err.Error())
		return
	}
	handler.SendResponse(c, "接受任务成功", nil)
}

// @Summary 确认任务
// @Description 确认任务已完成
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData string   true "id--任务的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/confirm [post]
func ConfirmTask(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	id := c.PostForm("id")
	err := tasks.ConfirmFinish(email, id)
	if err != nil {
		handler.SendError(c, "确认失败", errno.ErrDatabase, err.Error())
		return
	}
	handler.SendResponse(c, "确认成功", nil)
}

// @Summary 支付订单
// @Description 订单完成后，用户向接受者支付
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData string   true "id--任务的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/payment [post]
func PayBill(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	id := c.PostForm("id")
	payment, err := tasks.Payment(id, email)
	if err != nil {
		handler.SendBadRequest(c, "无法获取对方收款码", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "成功获取收款码", payment)
}

// @Summary 选择任务接受者
// @Description 接受他人发布的
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param email formData string   true "有意愿者的email"
// @Param id formData string 	true "任务id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/select [post]
func SelectAccepter(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	accepter := c.PostForm("accepter")
	id := c.PostForm("id")
	if err := tasks.SelectAccepter(email, accepter, id); err != nil {
		handler.SendError(c, "无法确认", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "选择成功", nil)
}

// @Summary 接受者确认收到付款
// @Description 确认首款
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData string 	true "任务id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /tasks/select [post]
func ConfirmPaid(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.PostForm("id")
	email := c.MustGet("email").(string)
	err := tasks.ConfirmPaid(id, email)
	if err != nil {
		handler.SendError(c, "确认失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "成功确认收款", nil)
}
