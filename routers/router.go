package routers

import (
	"github.com/astaxie/beego"
	"github.com/klovercloud-ci-cd/terminal/service"
)

func Init() {
	beego.Handler("/terminal/ws", &service.TerminalSockjs{}, true)
}
