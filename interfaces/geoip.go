package interfaces

import (
	"sort"
	"strings"
)

type IPStacks struct {
	IPv4 []string
	IPv6 []string
}

func (ips *IPStacks) Init() *IPStacks {
	if ips == nil {
		ips = &IPStacks{}
	}
	ips.IPv4 = []string{}
	ips.IPv6 = []string{}
	return ips
}

func (ips *IPStacks) Count() int {
	if ips == nil {
		return 0
	}
	return len(ips.IPv4) + len(ips.IPv6)
}

type GeoInfo struct {
	Org           string  `json:"organization"`
	Lon           float32 `json:"longitude"`
	Lat           float32 `json:"latitude"`
	TimeZone      string  `json:"timezone"`
	ISP           string  `json:"isp"`
	ASN           int     `json:"asn"`
	ASNOrg        string  `json:"asn_organization"`
	Country       string  `json:"country"`
	IP            string  `json:"ip"`
	ContinentCode string  `json:"continent_code"`
	CountryCode   string  `json:"country_code"`

	StackType string `json:"stackType"`
}

func (gi *GeoInfo) IsV6() bool {
	return gi != nil && gi.IP != "" && strings.Contains(gi.IP, ":")
}

type MultiStacks struct {
	Domain    string   // 域组，作为 In 时为域名，Out 时则为线路本身
	MainStack *GeoInfo // deprecating
	IPv4Stack []*GeoInfo
	IPv6Stack []*GeoInfo
}

func (tms *MultiStacks) Repr() string {
	repr := []string{}
	if tms == nil || tms.Count() == 0 {
		return ""
	}
	for _, v4 := range tms.IPv4Stack {
		repr = append(repr, v4.IP)
	}
	for _, v6 := range tms.IPv6Stack {
		repr = append(repr, v6.IP)
	}

	sort.Strings(repr)
	return strings.Join(repr, ",")
}

func (tms *MultiStacks) FirstV2(tag string) *GeoInfo {
	if tms == nil || tms.Count() == 0 {
		return nil
	}

	// check tags
	if tag == "" {
		tag = "46"
	}
	if len(tag) > 2 {
		return nil
	} else if len(tag) == 2 && tag[0] == tag[1] {
		return nil
	} else {
		for _, r := range tag {
			if r != '4' && r != '6' {
				return nil
			}
		}
	}

	// get ordered by the sequence order of tags
	stacks := []*GeoInfo{}
	for _, r := range tag {
		if r == '4' && len(tms.IPv4Stack) > 0 {
			stacks = append(stacks, tms.IPv4Stack...)
		} else if r == '6' && len(tms.IPv6Stack) > 0 {
			stacks = append(stacks, tms.IPv6Stack...)
		}
	}

	for _, ip := range stacks {
		if ip.IP != "" {
			return ip
		}
	}
	return nil
}

func (tms *MultiStacks) First(tag string) *GeoInfo {
	if tms == nil || tms.Count() == 0 {
		return nil
	}

	if tag != "v6" {
		for _, v4 := range tms.IPv4Stack {
			if v4.IP != "" {
				return v4
			}
		}
	}

	if tag != "v4" {
		for _, v6 := range tms.IPv6Stack {
			if v6.IP != "" {
				return v6
			}
		}
	}

	return nil
}

func (tms *MultiStacks) ForEach(assignedMain *GeoInfo) map[int][]*GeoInfo {
	result := make(map[int][]*GeoInfo)
	if assignedMain != nil && (tms == nil || tms.Count() == 0) {
		result[assignedMain.ASN] = []*GeoInfo{assignedMain}
		return result
	}
	if tms == nil {
		return result
	}
	for _, v4 := range tms.IPv4Stack {
		result[v4.ASN] = append(result[v4.ASN], v4)
	}
	for _, v6 := range tms.IPv6Stack {
		result[v6.ASN] = append(result[v6.ASN], v6)
	}

	return result
}

func (tms *MultiStacks) Count() int {
	if tms == nil {
		return 0
	}
	a, b := tms.V46StackCount()
	return a + b
}

func (tms *MultiStacks) V46StackCount() (int, int) {
	if tms == nil {
		return 0, 0
	}
	return len(tms.IPv4Stack), len(tms.IPv6Stack)
}

func (tms *MultiStacks) V46StackInfo() string {
	v4, v6 := tms.V46StackCount()
	ret := "N/A"
	if v4 > 0 {
		ret = "4⃣"
	}
	if v6 > 0 {
		ret = "6⃣"
	}
	if v4 > 0 && v6 > 0 {
		ret = "4⃣6⃣"
	}
	return ret
}
