package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/klovercloud-ci-cd/terminal/config"
	"github.com/klovercloud-ci-cd/terminal/routers"
	"github.com/klovercloud-ci-cd/terminal/service"
)

// @title Klovercloud-ci-cd-terminal API
// @description Klovercloud-ci-cd-terminal API
func main() {
	config.InitEnvironmentVariables()
	config.KubeConfig = config.GetKubeConfig()
	routers.Init()
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowOrigins:     []string{"console.klovercloud.com", "console.klovercloud.io", "localhost:6342"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	service.SetPublicKey()
	beego.Run(":" + config.ServerPort)
}
