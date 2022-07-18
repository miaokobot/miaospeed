package utils

import (
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils/structs"
)

type GlobalConfig struct {
	Token            string
	Binder           string
	WhiteList        []string
	SpeedLimit       uint64
	PauseSecond      uint
	ConnTaskTreading uint
	MiaoKoSignedTLS  bool
	NoSpeedFlag      bool
}

func (gc *GlobalConfig) InWhiteList(invoker string) bool {
	if len(gc.WhiteList) == 0 {
		return true
	}

	return structs.Contains(gc.WhiteList, invoker)
}

func (gc *GlobalConfig) VerifyRequest(req *interfaces.SlaveRequest) bool {
	return req.Challenge == gc.SignRequest(req)
}

func (gc *GlobalConfig) SignRequest(req *interfaces.SlaveRequest) string {
	return SignRequest(gc.Token, req)
}

var GCFG GlobalConfig
