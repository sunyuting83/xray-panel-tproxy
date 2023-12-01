package datafactory

import (
	"net/url"
	"strconv"
	"strings"
	config "xpanel/Config"
)

// vlessToJSON vless Json
func VlessToJSON(a string) (j *config.CodeList) {
	var (
		encryption string
		security   string
		headerType string
		net        string
		path       string
		host       string
		a1         []string
		a2         []string
		a3         []string
		a4         []string
		a5         []string
	)
	a1 = strings.Split(a, "@")
	a2 = strings.Split(a1[1], "?")
	a3 = strings.Split(a2[0], ":")
	a4 = strings.Split(a2[1], "#")
	a5 = strings.Split(a4[0], "&")
	for _, item := range a5 {
		x := strings.Split(item, "=")
		switch x[0] {
		case "encryption":
			encryption = x[1]
		case "security":
			security = x[1]
		case "headerType":
			headerType = x[1]
		case "type":
			net = x[1]
		case "host":
			host = x[1]
		case "path":
			path = x[1]
		}
	}
	port, err := strconv.Atoi(a3[1])
	if err != nil {
		port = 0
	}
	t := "测试节点"
	if len(a4) > 1 {
		t = a4[1]
	}
	title, _ := url.QueryUnescape(t)
	j = &config.CodeList{
		Password:   a1[0],
		Address:    a3[0],
		Port:       port,
		Encryption: encryption,
		Security:   security,
		HeaderType: headerType,
		Net:        net,
		Host:       host,
		Path:       path,
		Title:      title,
		Types:      "vless",
	}
	return
}
