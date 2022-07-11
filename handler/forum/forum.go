package forum

import (
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strconv"
	"time"
	"zhigui/handler"
	"zhigui/model"
	"zhigui/pkg/errno"
	"zhigui/services"
	"zhigui/services/connector"
	"zhigui/services/forum"
	"zhigui/services/like"
	"zhigui/services/random"
)

type ForumRequest struct {
	Title   string `json:"title"binding:"required"`
	Content string `json:"content"binding:"required"`
}
type CommentRequest struct {
	PostID  uint   `json:"post_id"`
	Content string `json:"content"`
}

// @Summary 查看帖子
// @Description 分页查看全部帖子
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query integer true "page--限制指定要检索的记录数 "
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum [get]
func GetPosts(c *gin.Context) {
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
	posts, num, err := model.GetPosts(limit*page, limit)
	if err != nil {
		handler.SendError(c, "查看失败", err.Error(), errno.ErrDatabase)
		return
	}

	handler.SendResponse(c, "查看成功", gin.H{
		"NUM":   num,
		"Posts": posts,
	})
}

// @Summary 查看帖子具体内容
// @Description 分页查看具体帖子
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id query integer true "帖子的id"
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query integer true "page--限制指定要检索的记录数 "
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/post [get]
func GetDetails(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	var err error
	var limit, page int

	ID, err := strconv.Atoi(c.DefaultQuery("id", "0"))

	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}

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
	post, num, err := model.GetPostDetails(limit*page, limit, ID)
	if err != nil {
		handler.SendError(c, nil, nil, errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "查询成功", map[string]interface{}{
		"Comments_Num": num,
		"Post":         post,
	})
}

// @Summary 发布帖子
// @Description 用户根据需要发布帖子
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param request body ForumRequest   true "id--任务的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/publish [post]
func PostMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	var Request ForumRequest
	if err := c.Bind(&Request); err != nil {
		handler.SendBadRequest(c, "请输入要发布的内容", err.Error(), errno.ErrBind)
		return
	}
	post, err := forum.PostMessage(email, Request.Content, Request.Title)
	if err != nil {
		handler.SendError(c, "发布失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "成功发布", post)
}

// @Summary 上传帖子图片
// @Description 上传新的图片
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData string   true "id--帖子的id"
// @Param file formData file true "文件"
// @Success 200 {object}  []model.Task{} "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/publish/pictures [post]
func UploadPicture(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	file, err := c.FormFile("file")
	ID := c.PostForm("id")
	PATH := "Forum"
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
	file.Filename = ID + " of Posts " + timeStr + " " + random.GetRandomString(16) + fileExt

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

	var post model.POST
	e := post.UploadPostPictures(ID, picUrl, picSha, picPath)

	if picUrl == "" || e != nil {
		handler.SendBadRequest(c, "上传失败,请检查token与其他配置参数是否正确", nil, errno.ErrPermissionDenied)
		return
	}

	handler.SendResponse(c, "上传成功", map[string]interface{}{
		"Post": post,
	})
}

// @Summary 获取用户发布的帖子
// @Description 用户发布过的帖子
// @Tags task
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query integer true "page--限制指定要检索的记录数 "
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/personal/posts [get]
func GetMyPosts(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}
	posts, num, err := forum.GetMyPosts(email, limit*page, limit)
	if err != nil {
		handler.SendError(c, "获取失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "获取成功", map[string]interface{}{
		"NUM":   num,
		"POSTS": posts,
	})
}

// @Summary 获取用户评论的帖子
// @Description 用户发布过的帖子
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query integer true "page--限制指定要检索的记录数 "
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/comments [get]
func GetMyComments(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}
	content, num, err := forum.GetMyComments(email, limit*page, limit)
	if err != nil {
		handler.SendError(c, "获取失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "获取成功", map[string]interface{}{
		"NUM":      num,
		"Comments": content,
	})
}

// @Summary 发布评论
// @Description 用户给已发布的帖子进行评论
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param request body CommentRequest    true "id--任务的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/comments [post]
func PostComment(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	var Request CommentRequest
	if err := c.Bind(&Request); err != nil {
		handler.SendBadRequest(c, "内容有误", err.Error(), errno.ErrBind)
		return
	}
	comment, err := forum.PostComment(Request.PostID, Request.Content, email)

	if err != nil {
		handler.SendError(c, "评论失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "评论成功", comment)
}

// @Summary 获取某帖子评论
// @Description 查看已发布帖子的评论内容
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id query integer true "id--帖子的id"
// @Param limit query integer true "limit--偏移量指定开始返回记录之前要跳过的记录数 "
// @Param page  query integer true "page--限制指定要检索的记录数 "
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/comments [get]
func GetComments(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	id := c.Query("id")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		handler.SendBadRequest(c, nil, nil, errno.ErrQuery)
		return
	}
	comments, num, err := forum.GetComments(id, limit*page, limit)
	if err != nil {
		handler.SendError(c, "查询评论失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "查询评论成功", map[string]interface{}{
		"NUM":      num,
		"Comments": comments,
	})
}

// @Summary 点赞
// @Description 给已发布的帖子点赞
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData integer true "帖子的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/like [post]
func GiveLike(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	id := c.PostForm("id")
	choice, err := like.GiveLike(id, email)
	if choice == 0 && err == nil {
		handler.SendResponse(c, "取消点赞", nil)
		return
	}
	if choice == 1 && err == nil {
		handler.SendResponse(c, "点赞成功", nil)
		return
	}
	if choice == 2 {
		handler.SendError(c, "操作失败", err.Error(), errno.ErrDatabase)
	}
}

// @Summary 删除某帖子
// @Description 删除用户发布的帖子
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData integer true "帖子的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/personal/posted [post]
func DeletePost(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	id := c.PostForm("id")
	err := forum.DeletePost(email, id)
	if err != nil {
		handler.SendError(c, "删除失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "删除成功", err)
}

// @Summary 删除某评论
// @Description 删除用户发布的帖子
// @Tags forum
// @Accept  json/application
// @Produce  json/application
// @Param Authorization header string true  "获取email"
// @Param id formData integer true "评论的id"
// @Success 200 {string}  json "{"code":0,"message":"OK","data":{}}"
// @Failure 400 {object} errno.Errno
// @Failure 404 {object} errno.Errno
// @Failure 500 {object} errno.Errno
// @Router /forum/personal/comments [post]
func DeleteComment(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	email := c.MustGet("email").(string)
	id := c.PostForm("id")
	err := forum.DeleteComment(email, id)
	if err != nil {
		handler.SendError(c, "删除失败", err.Error(), errno.ErrDatabase)
		return
	}
	handler.SendResponse(c, "删除成功", err)
}
