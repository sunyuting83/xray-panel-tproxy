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
		flow     string
		security string
		fp       string
		net      string
		path     string
		host     string
		pbk      string
		sid      string
		sni      string
		a1       []string
		a2       []string
		a3       []string
		a4       []string
		a5       []string
	)
	a1 = strings.Split(a, "@")
	a2 = strings.Split(a1[1], "?")
	a3 = strings.Split(a2[0], ":")
	a4 = strings.Split(a2[1], "#")
	a5 = strings.Split(a4[0], "&")
	for _, item := range a5 {
		x := strings.Split(item, "=")
		switch x[0] {
		case "flow":
			flow = x[1]
		case "security":
			security = x[1]
		case "fp":
			fp = x[1]
		case "type":
			net = x[1]
		case "host":
			host = x[1]
		case "sni":
			sni = x[1]
		case "fingerprint":
			path = x[1]
		case "pbk":
			pbk = x[1]
		case "sid":
			sid = x[1]
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
	if host == "" {
		host = sni
	}
	title, _ := url.QueryUnescape(t)
	j = &config.CodeList{
		Password:  a1[0],
		Address:   a3[0],
		Port:      port,
		Flow:      flow,
		Security:  security,
		Fp:        fp,
		Net:       net,
		Host:      host,
		Path:      path,
		Title:     title,
		Types:     "vless",
		Obfs:      pbk,
		ObfsParam: sid,
	}
	return
}
