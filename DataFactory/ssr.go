package datafactory

import (
	"net/url"
	"strconv"
	"strings"
	config "xpanel/Config"
)

// makeBase make base
func makeBase(a string) (b []string) {
	var (
		a3      []string
		address string
	)
	a3 = strings.Split(a, ":")
	address = a3[0]
	if strings.Contains(address, ".") {
		b = append(b, address, a3[1], a3[2], a3[3], a3[4], a3[5])
	} else {
		index := FindIndex(a3, "origin")
		ad := a3[0:index]
		address = strings.Join(ad, ":")
		b = append(b, address, a3[index+1], a3[index+2], a3[index+3], a3[index+4], a3[index+5])
	}
	return
}

// ssrToJSON ssr Json
func SsrToJSON(a string) (j *config.CodeList) {
	var (
		decode     string
		address    string
		port       int
		method     string
		password   string
		protocol   string = ""
		obfs       string = ""
		obfsparam  string = ""
		protoparam string = ""
		remarks    string
		a2         []string
	)
	decode = DecodeBytes(a)
	if strings.Contains(decode, ":") {
		if strings.Contains(decode, "/?") {
			a1 := strings.Split(decode, "/?")
			a2 = makeBase(a1[0])
			if len(a1[1]) > 0 {
				co, _ := url.QueryUnescape(a1[1])
				params := strings.Split(co, "&")
				for _, item := range params {
					x := strings.Split(item, "=")
					switch x[0] {
					case "obfsparam":
						if len(x) > 1 {
							obfsparam = DecodeBytes(x[1])
						}
					case "protoparam":
						if len(x) > 1 {
							protoparam = DecodeBytes(x[1])
						}
					case "remarks":
						if len(x) > 1 {
							remarks = DecodeBytes(x[1])
						}
					}
				}
			}
		} else {
			a2 = makeBase(decode)
		}
		address = a2[0]
		port, _ = strconv.Atoi(a2[1])
		protocol = a2[2]
		method = a2[3]
		obfs = a2[4]
		password = a2[5]
		if remarks == "" {
			remarks = address
		}
		j = &config.CodeList{
			Types:         "ssr",
			Title:         remarks,
			Address:       address,
			Port:          port,
			Password:      password,
			Method:        method,
			Protocol:      protocol,
			ProtocolParam: protoparam,
			Obfs:          obfs,
			ObfsParam:     obfsparam,
		}
	}
	return
}

// findIndex find index
func FindIndex(a []string, b string) (c int) {
	for i, item := range a {
		if item == b {
			c = i - 2
		}
	}
	return
}
