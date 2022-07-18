package clash

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/constant"
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils"
	"gopkg.in/yaml.v2"
)

func parseProxy(proxyName, proxyPayload string) constant.Proxy {
	var payload map[string]any
	yaml.Unmarshal([]byte(proxyPayload), &payload)
	proxy, err := adapter.ParseProxy(payload)

	if err != nil {
		utils.DLogf("Vendor Parser | Parse clash profile error, error=%v", err.Error())
	}

	return proxy
}

func extractFirstProxy(proxyName, proxyPayload string) constant.Proxy {
	proxy := parseProxy(proxyName, proxyPayload)

	if proxy != nil && interfaces.Parse(proxy.Type().String()) != interfaces.ProxyInvalid {
		return proxy
	}

	return nil
}
