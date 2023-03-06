package preconfigs

const PROXY_DEFAULT_STUN_SERVER = "udp://stun.voipstunt.com:3478"

const NETCAT_HTTP_PAYLOAD = `GET %s HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Host: %s
User-Agent: HTTPie/3.0.2 MiaoSpeed/%s

`

const (
	SPEED_DEFAULT_DURATION  int64 = 3
	SPEED_DEFAULT_THREADING uint  = 1

	SPEED_DEFAULT_LARGE_FILE_STATIC_APPLE    string = "https://updates.cdn-apple.com/2019FallFCS/fullrestores/061-22552/374D62DE-E18B-11E9-A68D-B46496A9EC6E/iPhone12,1_13.1.2_17A860_Restore.ipsw"
	SPEED_DEFAULT_LARGE_FILE_STATIC_MSFT     string = "https://download.microsoft.com/download/2/0/E/20E90413-712F-438C-988E-FDAA79A8AC3D/dotnetfx35.exe"
	SPEED_DEFAULT_LARGE_FILE_STATIC_GOOGLE   string = "https://dl.google.com/android/studio/maven-google-com/stable/offline-gmaven-stable.zip"
	SPEED_DEFAULT_LARGE_FILE_STATIC_CACHEFLY string = "http://cachefly.cachefly.net/200mb.test"

	SPEED_DEFAULT_LARGE_FILE_DYN_INTL string = "DYNAMIC:INTL"
	SPEED_DEFAULT_LARGE_FILE_DYN_FAST string = "DYNAMIC:FAST"

	SPEED_DEFAULT_LARGE_FILE_DEFAULT = SPEED_DEFAULT_LARGE_FILE_DYN_INTL

	SLAVE_DEFAULT_PING         = "http://gstatic.com/generate_204"
	SLAVE_DEFAULT_RETRY   uint = 3
	SLAVE_DEFAULT_TIMEOUT uint = 5000
)
