package geekhub

import (
	"strings"
)

var Version = "beta0.2"

var VersionPublishOrder = []string{
	"beta0.1",
	"beta0.2",
}

const VersionFile = "https://raw.githubusercontent.com/rrylee/geekterm/master/version.txt"

type NewVersion struct {
	t string
	s string
}

func CheckNewVersion() (hasNewVersion bool, newVersion *NewVersion, err error) {
	response, err := httpClient.R().Get(VersionFile)
	if err != nil {
		return false, nil, err
	}
	body := strings.TrimSpace(string(response.Body()))
	splitedBody := strings.Split(body, "\n")

	latestVersion := splitedBody[0]
	if strings.HasPrefix(latestVersion, "version:") {
		content := strings.Join(splitedBody[1:len(splitedBody)], "\n")
		latestVersion = strings.TrimLeft(latestVersion, "version:")

		if VersionPublishOrder[len(VersionPublishOrder)-1] == latestVersion {
			return false, nil, nil
		} else {
			return true, &NewVersion{
				t: latestVersion,
				s: content,
			}, nil
		}
	}
	return false, nil, nil
}
