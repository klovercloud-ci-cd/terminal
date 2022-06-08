package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/klovercloud-ci-cd/terminal/config"
	"github.com/klovercloud-ci-cd/terminal/routers"
	"github.com/klovercloud-ci-cd/terminal/service"
)

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

//go get k8s.io/api@kubernetes-1.12.9
//go get k8s.io/apimachinery@kubernetes-1.12.9
//go get k8s.io/client-go@kubernetes-1.12.9
// go get github.com/googleapis/gnostic@v0.4.0
