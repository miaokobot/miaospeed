package engine

import _ "embed"

//go:embed embeded/predefined.js
var PREDEFINED_SCRIPT string

//go:embed embeded/default_geoip.js
var DEFAULT_GEOIP_SCRIPT string

//go:embed embeded/default_ip.js
var DEFAULT_IP_SCRIPT string
