package taskbar

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zhigui/handler"
	"zhigui/pkg/errno"
	"zhigui/services/taskbar"
)

// @Summary 查看任务栏
// @Description 一个list
// @Tags taskbar
// @Accept  json/application
// @Produce  json/application
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Success 200 {object}  []model.Task{} "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /taskbar [get]
func GetList(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var err error
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

	contents, num, err := taskbar.GetAll(limit*page, limit)
	if err != nil {
		handler.SendError(c, "查询失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "查询成功", map[string]interface{}{
		"NUM":   num,
		"TASKS": contents,
	})
}
