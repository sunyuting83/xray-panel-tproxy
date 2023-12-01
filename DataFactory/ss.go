package datafactory

import (
	"net/url"
	"strconv"
	"strings"
	config "xpanel/Config"
)

// ssToJSON ss Json
func SsToJSON(a string) (j *config.CodeList) {
	var (
		address  string
		port     int
		method   string
		password string
		title    string
		a1       []string
		a2       []string
	)
	title = "默认节点"
	if strings.Contains(a, "#") {
		a1 = strings.Split(a, "#")
		title, _ = url.QueryUnescape(a1[1])
	} else {
		a1 = []string{a}
	}

	if strings.Contains(a1[0], "@") {
		a2 = strings.Split(a1[0], "@")
		a3 := strings.Split(a2[1], ":")
		address = a3[0]
		port, _ = strconv.Atoi(a3[1])
		if strings.Contains(a3[1], "?") {
			splitStr := "?"
			if strings.Contains(a3[1], "/") {
				splitStr = "/?"
			}
			portStr := strings.Split(a3[1], splitStr)[0]
			port, _ = strconv.Atoi(portStr)
		}
		a4 := make([]string, 1)
		if strings.Contains(a2[0], ":") {
			a4 = strings.Split(a2[0], ":")
		} else {
			decode := DecodeBytes(a2[0])
			a4 = strings.Split(decode, ":")
		}
		method = a4[0]
		password = a4[1]
	} else {
		a2 = []string{DecodeBytes(a)}
		a7 := make([]string, 1)

		a5 := strings.Split(a2[0], "@")
		a6 := strings.Split(a5[1], ":")
		if strings.Contains(a5[1], ":") {
			address = a6[0]
			port, _ = strconv.Atoi(a6[1])
			if strings.Contains(a6[1], "?") {
				splitStr := "?"
				if strings.Contains(a6[1], "/") {
					splitStr = "/?"
				}
				portStr := strings.Split(a6[1], splitStr)[0]
				port, _ = strconv.Atoi(portStr)
			}
		}
		if strings.Contains(a5[0], ":") {
			a7 = strings.Split(a5[0], ":")
		} else {
			decode := DecodeBytes(a5[0])
			a7 = strings.Split(decode, ":")
		}
		method = a7[0]
		password = a7[1]
	}

	j = &config.CodeList{
		Types:    "ss",
		Title:    title,
		Address:  address,
		Port:     port,
		Method:   method,
		Password: password,
	}
	return
}
