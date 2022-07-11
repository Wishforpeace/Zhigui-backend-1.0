package main

import (
	"flag"
	"fmt"
	"log"
	// "blackboard/config"
	"zhigui/model"
	"zhigui/router"
	"zhigui/services/flag_handle"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// "github.com/spf13/viper"
)

// @title Zhigui API
// @version 1.0.0
// @description 纸贵API
// @termsOfService http://swagger.io/terrms/
// @contact.name Wishiforpeace
// @contact.email 1903180340@qq.com
// @host 119.3.2.168:2333
// @BasePath /api/v1
// @Schemes http

var err error

func main() {
	// err := config.Init("./conf/config.yaml", "")
	// if err != nil {
	// 	panic(err)
	// }
	// dbMap := viper.GetStringMapString("db")
	// dsn := fmt.Sprintf("%s:%s@/%s?parseTime=True", dbMap["username"], dbMap["password"], dbMap["name"])
	dsn := "root:root&1234@tcp(127.0.0.1:3306)/Zhigui?charset=utf8mb4&parseTime=True&loc=Local"
	model.DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("数据库连接失败")
		panic(err)
	}
	//gin.SetMode(gin.ReleaseMode)
	model.Migrate(model.DB)
	link := "http://119.3.2.168:" + flag_handle.PORT
	log.Println("监听端口:", link, "请不要关闭终端")
	defer model.DB.Close()
	r := router.Router()
	err = r.Run(":" + flag_handle.PORT)
	if err != nil {
		panic(err)
	}
}

func init() {
	port := flag.String("port", "2333", "本地监听的端口")
	platform := flag.String("platform", "github", "平台名称，支持gitee/github")
	token := flag.String("token", "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", "Gitee/Github 的用户授权码")
	owner := flag.String("owner", "Wishforpeace", "仓库所属空间地址(企业、组织或个人的地址path)")
	repo := flag.String("repo", "ZhiguiImage", "仓库路径(path)")
	//path := flag.String("path", "", "文件的路径")
	branch := flag.String("branch", "master", "分支")
	flag.Parse()
	flag_handle.PORT = *port
	flag_handle.OWNER = *owner
	flag_handle.REPO = *repo
	//flag_handle.PATH = *path
	flag_handle.TOKEN = *token
	flag_handle.PLATFORM = *platform
	flag_handle.BRANCH = *branch

	if flag_handle.TOKEN == "" {
		panic("token 必须！")
	}
}
