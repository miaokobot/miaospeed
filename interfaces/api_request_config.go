package interfaces

import (
	"github.com/miaokobot/miaospeed/preconfigs"
	"github.com/miaokobot/miaospeed/utils/structs"
)

type SlaveRequestConfigs struct {
	STUNURL           string `yaml:"stunURL,omitempty" cf:"name=🫙 STUN 地址"`
	DownloadURL       string `yaml:"downloadURL,omitempty" cf:"name=📃 测速文件"`
	DownloadDuration  int64  `yaml:"downloadDuration,omitempty" cf:"name=⏱️ 测速时长 (单位: 秒)"`
	DownloadThreading uint   `yaml:"downloadThreading,omitempty" cf:"name=🧶 测速线程数"`

	PingAverageOver uint16 `yaml:"pingAverageOver,omitempty" cf:"name=🧮 多次 Ping 求均值,value"`
	PingAddress     string `yaml:"pingAddress,omitempty" cf:"name=🏫 URL Ping 地址"`

	TaskRetry  uint     `yaml:"taskRetry,omitempty" cf:"name=🐛 测试重试次数"`
	DNSServers []string `yaml:"dnsServers,omitempty" cf:"name=💾 自定义DNS服务器,childvalue"`

	TaskTimeout uint     `yaml:"-" fw:"readonly"`
	Scripts     []Script `yaml:"-" fw:"readonly"`
}

func (src *SlaveRequestConfigs) DescriptionText() string {
	hint := structs.X("案例:\ndownloadDuration: 取值范围 [1,30]\ndownloadThreading: 取值范围 [1,8]\ntaskThreading: 取值范围 [1,32]\ntaskRetry: 取值范围 [1,10]\n\n当前:\n")
	cont := "empty"
	if src != nil {
		cont = structs.X("downloadDuration: %d\ndownloadThreading: %d\ntaskRetry: %d\n", src.DownloadDuration, src.DownloadThreading, src.TaskRetry)
	}
	return hint + cont
}

func (src *SlaveRequestConfigs) Clone() *SlaveRequestConfigs {
	return &SlaveRequestConfigs{
		DownloadURL:       src.DownloadURL,
		DownloadDuration:  src.DownloadDuration,
		DownloadThreading: src.DownloadThreading,

		PingAverageOver: src.PingAverageOver,
		PingAddress:     src.PingAddress,

		TaskRetry:  src.TaskRetry,
		DNSServers: cloneSlice(src.DNSServers),

		TaskTimeout: src.TaskTimeout,
		Scripts:     src.Scripts,
	}
}

func (src *SlaveRequestConfigs) Merge(from *SlaveRequestConfigs) *SlaveRequestConfigs {
	ret := src.Clone()
	if from.DownloadURL != "" {
		ret.DownloadURL = from.DownloadURL
	}
	if from.DownloadDuration != 0 {
		ret.DownloadDuration = from.DownloadDuration
	}
	if from.DownloadThreading != 0 {
		ret.DownloadThreading = from.DownloadThreading
	}

	if from.PingAverageOver != 0 {
		ret.PingAverageOver = from.PingAverageOver
	}
	if from.PingAddress != "" {
		ret.PingAddress = from.PingAddress
	}

	if from.TaskRetry != 0 {
		ret.TaskRetry = from.TaskRetry
	}

	if from.DNSServers != nil {
		ret.DNSServers = from.DNSServers[:]
	}

	if from.TaskTimeout != 0 {
		ret.TaskTimeout = from.TaskTimeout
	}
	if from.Scripts != nil {
		ret.Scripts = from.Scripts
	}

	return ret
}

func (cfg *SlaveRequestConfigs) Check() *SlaveRequestConfigs {
	if cfg == nil {
		cfg = &SlaveRequestConfigs{}
	}

	if cfg.STUNURL == "" {
		cfg.STUNURL = preconfigs.PROXY_DEFAULT_STUN_SERVER
	}
	if cfg.DownloadURL == "" {
		cfg.DownloadURL = preconfigs.SPEED_DEFAULT_LARGE_FILE_DEFAULT
	}
	if cfg.DownloadDuration < 1 || cfg.DownloadDuration > 30 {
		cfg.DownloadDuration = preconfigs.SPEED_DEFAULT_DURATION
	}
	if cfg.DownloadThreading < 1 || cfg.DownloadThreading > 32 {
		cfg.DownloadThreading = preconfigs.SPEED_DEFAULT_THREADING
	}

	if cfg.TaskRetry < 1 || cfg.TaskRetry > 10 {
		cfg.TaskRetry = preconfigs.SLAVE_DEFAULT_RETRY
	}

	if cfg.PingAddress == "" {
		cfg.PingAddress = preconfigs.SLAVE_DEFAULT_PING
	}
	if cfg.PingAverageOver == 0 || cfg.PingAverageOver > 16 {
		cfg.PingAverageOver = 1
	}

	if cfg.DNSServers == nil {
		cfg.DNSServers = make([]string, 0)
	}

	if cfg.TaskTimeout < 10 || cfg.TaskTimeout > 10000 {
		cfg.TaskTimeout = preconfigs.SLAVE_DEFAULT_TIMEOUT
	}
	if cfg.Scripts == nil {
		cfg.Scripts = make([]Script, 0)
	}

	return cfg
}
