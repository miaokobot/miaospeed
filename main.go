package main

import (
	"github.com/miaokobot/miaospeed/service"
	"github.com/miaokobot/miaospeed/utils"
)

var COMPILATIONTIME string
var BUILDCOUNT string
var COMMIT string
var BRAND string

func main() {
	utils.COMPILATIONTIME = COMPILATIONTIME
	utils.BUILDCOUNT = BUILDCOUNT
	utils.COMMIT = COMMIT
	utils.BRAND = BRAND

	InitConfig()
	utils.DInfof("MiaoSpeed speedtesting client %s", utils.VERSION)

	// start task server
	go service.StartTaskServer()

	// start api server
	service.CleanUpServer()
	go service.InitServer()

	<-utils.MakeSysChan()

	// clean up
	service.CleanUpServer()
	utils.DInfo("shutting down.")
}
