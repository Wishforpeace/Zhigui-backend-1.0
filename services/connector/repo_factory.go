package connector

import (
	"zhigui/services"
	"zhigui/services/flag_handle"
	"zhigui/services/github"
)

//定义serve的映射关系
var serveMap = map[string]services.RepoInterface{
	"github": &github.GithubServe{},
}

func RepoCreate() services.RepoInterface {
	return serveMap[flag_handle.PLATFORM]
}
