package rank

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zhigui/handler"
	"zhigui/pkg/errno"
	"zhigui/services/rank"
)

// @Summary 查看排行榜
// @Description 按照收入排名
// @Tags taskbar
// @Accept  json/application
// @Produce  json/application
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query integer true "page--限制指定要检索的记录数 "
// @Success 200 {object}  []model.User{} "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /ranking [get]
func GetRank(c *gin.Context) {
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
	item, num, err := rank.GetRank(limit*page, limit)
	if err != nil {
		handler.SendError(c, "获取失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "获取成功", map[string]interface{}{
		"NUM":   num,
		"Users": item,
	})
}
