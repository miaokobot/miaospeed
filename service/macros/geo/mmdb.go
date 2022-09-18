package geo

import (
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils"
)

func RunMMDBCheck(rawIp string) *interfaces.GeoInfo {
	if record := utils.QueryMaxMindDB(rawIp); record != nil {
		return &interfaces.GeoInfo{
			ASN:    record.ASN,
			ASNOrg: record.ASNOrg,
			Org:    record.ASNOrg,
			ISP:    record.ASNOrg, // inaccurate, just fallback
			IP:     rawIp,

			Country:       record.Country.Names.EN,
			CountryCode:   record.Country.ISOCode,
			ContinentCode: record.Continent.Code,
			TimeZone:      record.Location.TimeZone,
			Lat:           record.Location.Latitude,
			Lon:           record.Location.Longitude,
		}
	}

	return nil
}
