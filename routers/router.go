package routers

import (
	"github.com/astaxie/beego"
	"github.com/klovercloud-ci-cd/terminal/service"
)

// Init router
func Init() {
	beego.Handler("/terminal/ws", &service.TerminalSockjs{}, true)
}
