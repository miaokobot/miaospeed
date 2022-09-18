package utils

import (
	"net"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

type MMDBResult struct {
	ASN    int    `maxminddb:"autonomous_system_number"`
	ASNOrg string `maxminddb:"autonomous_system_organization"`

	City struct {
		GeoNameID int `maxminddb:"geoname_id"`
		Names     struct {
			EN string `maxminddb:"en"`
			JA string `maxminddb:"ja"`
			ZH string `maxminddb:"zh-CN"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`

	Continent struct {
		Code      string `maxminddb:"code"`
		GeoNameID int    `maxminddb:"geoname_id"`
		Names     struct {
			EN string `maxminddb:"en"`
			JA string `maxminddb:"ja"`
			ZH string `maxminddb:"zh-CN"`
		} `maxminddb:"names"`
	} `maxminddb:"continent"`

	Country struct {
		ISOCode   string `maxminddb:"iso_code"`
		GeoNameID int    `maxminddb:"geoname_id"`
		Names     struct {
			EN string `maxminddb:"en"`
			JA string `maxminddb:"ja"`
			ZH string `maxminddb:"zh-CN"`
		} `maxminddb:"names"`
	} `maxminddb:"country"`

	Location struct {
		Accuracy  int     `maxminddb:"accuracy_radius"`
		Latitude  float32 `maxminddb:"latitude"`
		Longitude float32 `maxminddb:"longitude"`
		TimeZone  string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

var MaxMindDBs []*maxminddb.Reader

func LoadMaxMindDB(pathList string) {
	if pathList == "" {
		return
	}

	MaxMindDBs = []*maxminddb.Reader{}
	for _, path := range strings.Split(pathList, ",") {
		DWarnf("Maxmind Database | Loading maxmind database, path=%v", path)
		mmdb, err := maxminddb.Open(path)
		if err != nil {
			DErrorf("Maxmind Database | Cannot load maxmind database, err=%v", err.Error())
			os.Exit(1)
		}
		MaxMindDBs = append(MaxMindDBs, mmdb)
	}
}

func QueryMaxMindDB(rawIp string) *MMDBResult {
	if len(MaxMindDBs) == 0 {
		return nil
	}

	result := MMDBResult{}

	ip := net.ParseIP(rawIp)
	if ip == nil {
		DErrorf("Maxmind Database | Cannot parse ip address, ip=%v", rawIp)
		return &result
	}

	for _, db := range MaxMindDBs {
		err := db.Lookup(ip, &result)
		if err != nil {
			DErrorf("Maxmind Database | Cannot query mmdb table, ip=%v err=%v", rawIp, err.Error())
		}
	}

	return &result
}
