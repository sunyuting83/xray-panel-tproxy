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
		// tlsa   bool // 弃用旧的布尔值，直接映射到 StreamSecurity
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
			if len(x) < 2 {
				continue
			} // 安全检查
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

		// 核心对齐逻辑
		streamSec := ""
		if tls == "1" || tls == "tls" {
			streamSec = "tls"
		}

		j = &config.CodeList{
			Type:           "vmess", // 确保类型正确
			Title:          ps,
			Address:        host,
			Port:           port,
			ID:             uuid, // UUID 对应新结构的 ID
			AlterID:        alterid,
			Network:        obfs, // 原来叫 Net, 现在叫 Network
			Path:           path,
			Host:           obfsParam,
			StreamSecurity: streamSec, // 统一安全层
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

// V2rayToJsons 处理标准 VMess JSON 分享链接 (Base64 解码后的)
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

	// 你原来的兼容端口为字符串形式的 hack 逻辑保留
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

	// 处理 TLS 映射
	streamSec := ""
	if r.TLS { // config.Vary 里的 TLS 是布尔值
		streamSec = "tls"
	}

	result = &config.CodeList{
		Type:           "vmess", // 统一使用 Type
		Title:          r.Ps,
		Address:        r.Add,
		Port:           r.Port,
		ID:             r.ID, // UUID 对应 ID
		AlterID:        r.Aid,
		Network:        r.Net, // 对齐 Network
		Path:           r.Path,
		Host:           r.Host,
		StreamSecurity: streamSec, // 映射为字符串 "tls"
		Sni:            r.Host,    // VMess 开启 TLS 时，SNI 通常等于 Host
		Fingerprint:    "chrome",  // 预设一个指纹，现代节点通用
	}
	return
}
