package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/miaokobot/miaospeed/utils"
)

func InitConfig() *utils.GlobalConfig {
	gcfg := &utils.GCFG

	flag.StringVar(&gcfg.Token, "token", "", "specify the token used to sign request")
	flag.StringVar(&gcfg.Binder, "bind", "", "bind a socket, can be format like 0.0.0.0:8080 or /tmp/unix_socket")
	flag.Uint64Var(&gcfg.SpeedLimit, "speedlimit", 0, "speed ratelimit (in Bytes per Second), default with no limits")
	flag.UintVar(&gcfg.PauseSecond, "pausesecond", 0, "pause such period after each speed job (seconds)")
	flag.BoolVar(&gcfg.MiaoKoSignedTLS, "mtls", false, "enable miaoko certs for tls verification")
	flag.BoolVar(&gcfg.NoSpeedFlag, "nospeed", false, "decline all speedtest requests")

	verboseMode := flag.Bool("verbose", false, "whether to print out systems log")
	versionOnly := flag.Bool("version", false, "display version and exit")
	whiteList := flag.String("whitelist", "", "bot id whitelist, can be format like 1111,2222,3333")

	flag.Parse()

	if *verboseMode {
		utils.VerboseLevel = utils.LTLog
	}

	if *versionOnly {
		fmt.Println(utils.VERSION)
		os.Exit(0)
	}

	gcfg.WhiteList = make([]string, 0)
	if *whiteList != "" {
		gcfg.WhiteList = strings.Split(*whiteList, ",")
	}

	return gcfg
}
