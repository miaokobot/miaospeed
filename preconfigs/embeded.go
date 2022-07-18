package preconfigs

import (
	_ "embed"
)

//go:embed embeded/miaokoCA/miaoko.crt
var MIAOKO_TLS_CRT string

//go:embed embeded/miaokoCA/miaoko.key
var MIAOKO_TLS_KEY string

//go:embed embeded/ca-certificates.crt
var MIAOKO_ROOT_CA []byte
