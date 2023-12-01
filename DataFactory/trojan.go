package datafactory

import (
	"net/url"
	"strconv"
	"strings"
	config "xpanel/Config"
)

// TrojanToJSON ss Json
func TrojanToJSON(a string) (j *config.CodeList) {
	var (
		address  string
		port     int
		password string
		title    string
		host     string
		a1       []string
		a2       []string
		a3       []string
	)
	title = "默认节点"
	if strings.Contains(a, "#") {
		a1 = strings.Split(a, "#")
		title, _ = url.QueryUnescape(a1[1])
	} else {
		a1 = make([]string, 1)
		a1 = append(a1, a)
	}
	if strings.Contains(a1[0], "?") {
		a4 := strings.Split(a1[0], "?")
		a2 = strings.Split(a4[0], "@")
		a3 = strings.Split(a2[1], ":")
		address = a3[0]
		port, _ = strconv.Atoi(a3[1])
		password = a2[0]
		if strings.Contains(a4[1], "sni=") {
			a5 := strings.Split(a4[1], "&")
			for _, item := range a5 {
				x := strings.Split(item, "=")
				if x[0] == "sni" {
					host = x[1]
					break
				}
			}
		} else {
			host = address
		}
	} else {
		a2 = strings.Split(a1[0], "@")
		a3 = strings.Split(a2[1], ":")
		address = a3[0]
		port, _ = strconv.Atoi(a3[1])
		password = a2[0]
		host = address
	}

	j = &config.CodeList{
		Types:    "trojan",
		Title:    title,
		Address:  address,
		Port:     port,
		Password: password,
		Host:     host,
	}
	return
}
