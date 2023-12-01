package datafactory

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	config "xpanel/Config"
)

// V2rayToJSON v2ray to json
func V2rayToJSON(item string) (j *config.CodeList) {
	var (
		newstr    string
		ps        string
		obfsParam string
		path      string
		obfs      string
		alterid   int
		tls       string
		tlsa      bool
	)

	if strings.Contains(item, "?remarks=") {
		strsss := strings.Split(item, "?")
		newstr = DecodeBytes(strsss[0])

		blen := len(newstr)
		a := strings.Index(newstr, ":")
		b := strings.Index(newstr, "@")
		c := strings.LastIndex(newstr, ":")
		uuid := newstr[a+1 : b]
		host := newstr[b+1 : c]
		port, _ := strconv.Atoi(newstr[c+1 : blen])

		params, _ := url.QueryUnescape(strsss[1])
		and := "&"
		if strings.Contains(params, "&amp;") {
			and = "&amp;"
		}
		l := strings.Split(params, and)
		for _, it := range l {
			x := strings.Split(it, "=")
			switch x[0] {
			case "remarks":
				ps = x[1]
			case "obfsParam":
				obfsParam = x[1]
			case "path":
				path = x[1]
			case "obfs":
				obfs = x[1]
			case "alterId":
				alterid, _ = strconv.Atoi(x[1])
			case "tls":
				tls = x[1]
			}
		}
		if ps == "" {
			ps = "未知名称"
		}
		if tls == "1" {
			tlsa = true
		} else {
			tlsa = false
		}
		j = &config.CodeList{
			Types:    "vmess",
			Title:    ps,
			Host:     obfsParam,
			Path:     path,
			TLS:      tlsa,
			Address:  host,
			Port:     port,
			Password: uuid,
			Aid:      alterid,
			Net:      obfs,
		}
	} else {
		newstr = DeBase(item)
		if !strings.Contains(newstr, "}") {
			newstr = strings.Join([]string{newstr, "}"}, "")
		}
		j = V2rayToJsons(newstr)
	}
	return
}

func UnmarshalJSON(t *config.Vary, data []byte) error {
	type VaryAlias config.Vary
	v2ray := &VaryAlias{
		Host:  "",
		Path:  "",
		TLS:   false,
		Ps:    "noname",
		Add:   "0.0.0.0",
		Port:  0,
		ID:    "",
		Aid:   0,
		Net:   "tcp",
		Type:  "none",
		Types: "vmess",
		Title: "noname",
	}

	_ = json.Unmarshal(data, v2ray)

	*t = config.Vary(*v2ray)
	return nil
}

// V2rayToJsons fun
func V2rayToJsons(s string) (result *config.CodeList) {
	var (
		a     []byte = []byte(s)
		index int    = len(a)
	)
	s = strings.Replace(s, " ", "", -1)
	if strings.Contains(s, "\n") {
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "\t", "", -1)
		s = strings.Replace(s, "\r", "", -1)
	}
	if strings.Contains(s, `"port":"`) {
		portString := `"port":"`
		overLen := len(s)
		plin := strings.Index(s, portString)
		firstStr := s[0 : plin+7]
		endStr := strings.Split(s, portString)[1]
		portLin := strings.Index(endStr, `"`)
		portNum := s[plin+8 : plin+8+portLin]
		overStr := s[plin+9+portLin : overLen]
		newStr := strings.Join([]string{firstStr, portNum, overStr}, "")
		a = []byte(newStr)
	}
	index = bytes.IndexByte(a, 0)
	if index != -1 {
		a = a[:index]
	}
	r := &config.Vary{}
	_ = json.Unmarshal(a, &r)
	result = &config.CodeList{
		Types:    "vmess",
		Title:    r.Ps,
		Host:     r.Host,
		Path:     r.Path,
		TLS:      r.TLS,
		Address:  r.Add,
		Port:     r.Port,
		Password: r.ID,
		Aid:      r.Aid,
		Net:      r.Net,
	}
	result.Types = "vmess"
	return
}
