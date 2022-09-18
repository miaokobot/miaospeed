package geo

import (
	"time"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils"
	"github.com/miaokobot/miaospeed/utils/structs"
	"github.com/miaokobot/miaospeed/utils/structs/memutils"
	"github.com/miaokobot/miaospeed/utils/structs/obliviousmap"
)

var GeoCache *obliviousmap.ObliviousMap[*interfaces.GeoInfo]

func RunGeoCheck(p interfaces.Vendor, script string, ip string, retry int, network interfaces.RequestOptionsNetwork) *interfaces.GeoInfo {
	if mmdbRet := RunMMDBCheck(ip); mmdbRet != nil {
		return mmdbRet
	}

	if r, ok := GeoCache.Get(ip); ok && r != nil {
		return r
	}
	var ret *interfaces.GeoInfo = nil
	for i := 0; i < structs.WithIn(retry, 1, 3) && (ret == nil || ret.IP != ""); i++ {
		ret = ExecGeoCheck(p, script, ip, network)
	}
	if ret == nil {
		ret = &interfaces.GeoInfo{}
	}
	proxyName := "NoProxy"
	if p != nil {
		proxyName = p.ProxyInfo().Name
	}
	if ret != nil && ret.IP != "" {
		GeoCache.Set(ret.IP, ret)
		utils.DLogf("GetIP Resolver | Resolved IP=%s proxy=%v ASN=%d ASOrg=%s", ip, proxyName, ret.ASN, ret.ASNOrg)
	} else {
		utils.DWarnf("GeoIP Resolver | Fail to resolve IP=%s proxy=%v", ip, proxyName)
	}
	return ret
}

func init() {
	memGeoInfo := memutils.MemDriverMemory[*interfaces.GeoInfo]{}
	memGeoInfo.Init()
	GeoCache = obliviousmap.NewObliviousMap[*interfaces.GeoInfo]("GeoCache/", time.Hour*6, true, &memGeoInfo)
}
