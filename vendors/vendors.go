package vendors

import (
	"github.com/miaokobot/miaospeed/interfaces"

	"github.com/miaokobot/miaospeed/vendors/clash"
	"github.com/miaokobot/miaospeed/vendors/invalid"
	"github.com/miaokobot/miaospeed/vendors/local"
)

var registeredList = map[interfaces.VendorType]func() interfaces.Vendor{
	interfaces.VendorLocal: func() interfaces.Vendor {
		return &local.Local{}
	},
	interfaces.VendorClash: func() interfaces.Vendor {
		return &clash.Clash{}
	},
}

func Find(vendorType interfaces.VendorType) interfaces.Vendor {
	if vendor, ok := registeredList[vendorType]; ok {
		return vendor()
	}

	return &invalid.Invalid{}
}
