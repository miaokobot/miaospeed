package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/miaokobot/miaospeed/utils"
)

const MAXMIND_DB_DOWNLOAD_URL = "https://download.maxmind.com/app/geoip_download?edition_id=%s&license_key=%s&suffix=tar.gz"

var MAXMIND_EDITION_MATRIX = []string{
	"GeoLite2-ASN", "GeoLite2-City",
}

type MiscCliParams struct {
	MaxmindLicenseKey string
}

func InitConfigMisc() *MiscCliParams {
	stcp := &MiscCliParams{}

	sflag := flag.NewFlagSet(cmdName+" misc", flag.ExitOnError)
	sflag.StringVar(&stcp.MaxmindLicenseKey, "maxmind-update-license", "", "specify a maxmind license to update database.")

	parseFlag(sflag)

	return stcp
}

func RunCliMisc() {
	stcp := InitConfigMisc()

	if stcp.MaxmindLicenseKey != "" {
		// update maxmind database
		mmdbFilter := regexp.MustCompile(`\.mmdb$`)
		for _, edition := range MAXMIND_EDITION_MATRIX {
			url := fmt.Sprintf(MAXMIND_DB_DOWNLOAD_URL, edition, stcp.MaxmindLicenseKey)
			if downloadBytes, err := utils.DownloadBytes(url); err != nil {
				utils.DErrorf("Maxmind Updater | Cannot fetch content from server, edition=%s err=%s", edition, err.Error())
			} else if archiveEntries, err := utils.FindAndExtract(bytes.NewBuffer(downloadBytes), *mmdbFilter); err != nil || len(archiveEntries) == 0 {
				utils.DErrorf("Maxmind Updater | Cannot extract content from gzip file, edition=%s size=%d", edition, len(downloadBytes))
			} else {
				for file, fileBytes := range archiveEntries {
					if err := os.WriteFile(file, fileBytes, 0644); err != nil {
						utils.DErrorf("Maxmind Updater | Create local file, edition=%s size=%d file=%s err=%s", edition, len(fileBytes), file, err.Error())
					} else {
						utils.DWarnf("Maxmind Updater | File updated, edition=%s size=%d file=%s", edition, len(fileBytes), file)
					}
				}
			}
		}
		return
	}

	fmt.Printf("You have not specify any options, please call %s misc -help to see all available commands.\n", cmdName)
}
