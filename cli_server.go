package main

import (
	"flag"
	"strings"

	"github.com/miaokobot/miaospeed/service"
	"github.com/miaokobot/miaospeed/utils"
)

func InitConfigServer() *utils.GlobalConfig {
	gcfg := &utils.GCFG

	sflag := flag.NewFlagSet(cmdName+" server", flag.ExitOnError)
	sflag.StringVar(&gcfg.Token, "token", "", "specify the token used to sign request")
	sflag.StringVar(&gcfg.Binder, "bind", "", "bind a socket, can be format like 0.0.0.0:8080 or /tmp/unix_socket")
	sflag.Uint64Var(&gcfg.SpeedLimit, "speedlimit", 0, "speed ratelimit (in Bytes per Second), default with no limits")
	sflag.UintVar(&gcfg.PauseSecond, "pausesecond", 0, "pause such period after each speed job (seconds)")
	sflag.BoolVar(&gcfg.MiaoKoSignedTLS, "mtls", false, "enable miaoko certs for tls verification")
	sflag.BoolVar(&gcfg.NoSpeedFlag, "nospeed", false, "decline all speedtest requests")

	whiteList := sflag.String("whitelist", "", "bot id whitelist, can be format like 1111,2222,3333")
	parseFlag(sflag)

	gcfg.WhiteList = make([]string, 0)
	if *whiteList != "" {
		gcfg.WhiteList = strings.Split(*whiteList, ",")
	}

	return gcfg
}

func RunCliServer() {
	InitConfigServer()
	utils.DWarnf("MiaoSpeed speedtesting client %s", utils.VERSION)

	// start task server
	go service.StartTaskServer()

	// start api server
	service.CleanUpServer()
	go service.InitServer()

	<-utils.MakeSysChan()

	// clean up
	service.CleanUpServer()
	utils.DLog("shutting down.")
}
