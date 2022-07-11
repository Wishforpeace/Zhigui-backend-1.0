package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	Auth "zhigui/handler/auth"
	Forum "zhigui/handler/forum"
	Rank "zhigui/handler/rank"
	Task "zhigui/handler/task"
	TaskBar "zhigui/handler/taskbar"
	User "zhigui/handler/user"
	"zhigui/router/middleware"
)

func Router() *gin.Engine {
	r := gin.New()
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API router.")
	})
	// 注册,登录

	g1 := r.Group("/api/v1/auth")
	{

		g1.POST("/register", Auth.Register)
		g1.POST("/login", Auth.Login)
	}
	g2 := r.Group("/api/v1/user").Use(middleware.Auth())
	{
		g2.GET("/info", User.GetInfo)
		g2.GET("/others", User.GetOtherInfo)
		g2.POST("/info", User.UpdateInfo)
		g2.POST("/avatar", User.UploadProfile)
		g2.POST("/payment", User.UploadPayment)
	}
	g3 := r.Group("/api/v1/tasks").Use(middleware.Auth())
	{
		// 获取个人的任务，发布的和接受的
		g3.GET("/personal", Task.GetPersonalTasks)
		// 获取所有任务
		g3.GET("/prefecture", Task.GetTasks)
		// 查看任务细节
		g3.GET("/details", Task.GetTaskDetails)
		// 发布任务
		g3.POST("/publish", Task.PublishTasks)
		// 上传任务图片
		g3.POST("/publish/picture", Task.UploadPicture)
		// 删除任务图片
		g3.DELETE("/pictures", Task.DeletePictures)
		// 修改任务
		g3.POST("/details", Task.UpdateTask)
		// 删除任务
		g3.DELETE("/delete", Task.DeleteTask)
		// 接受任务
		g3.POST("/accept", Task.AcceptTask)
		// 选择接受者
		g3.POST("/select", Task.SelectAccepter)
		// 确认完成
		g3.POST("/confirm", Task.ConfirmTask)
		// 支付
		g3.POST("/confirm/payment", Task.PayBill)
		// 确认支付
		g3.POST("/confirm/paid", Task.ConfirmPaid)
	}
	g4 := r.Group("/api/v1/forum").Use(middleware.Auth())
	{
		// 获取全部帖子，带分页
		g4.GET("", Forum.GetPosts)
		// 查看详细帖子
		g4.GET("/post", Forum.GetDetails)
		// 发布帖子
		g4.POST("/publish", Forum.PostMessage)
		// 上传帖子图片
		g4.POST("/publish/pictures", Forum.UploadPicture)
		// 发布评论
		g4.POST("/comments", Forum.PostComment)
		// 获取评论
		g4.GET("/comments", Forum.GetComments)
		// 点赞获取消点赞
		g4.POST("/like", Forum.GiveLike)
		// 获取个人发布的帖子
		g4.GET("/personal/posts", Forum.GetMyPosts)
		// 获取个人发布的评论
		g4.GET("/personal/comments", Forum.GetMyComments)
		// 删除帖子
		g4.DELETE("/personal/posted", Forum.DeletePost)
		// 删除评论
		g4.DELETE("/personal/comments", Forum.DeleteComment)
	}
	g5 := r.Group("/api/v1/taskbar")
	{
		g5.GET("", TaskBar.GetList)
	}
	g6 := r.Group("/api/v1/ranking")
	{
		g6.GET("", Rank.GetRank)
	}
	return r
}
