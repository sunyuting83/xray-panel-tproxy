package datafactory

import (
	"net/url"
	"strconv"
	"strings"
	config "xpanel/Config"
)

// VlessToJSON 解析 vless 链接并映射到新的 CodeList 结构体
func VlessToJSON(a string) (j *config.CodeList) {
	var (
		flow          string
		security      string // 链接中的 security (tls/reality/none)
		fp            string
		net           string
		path          string
		host          string
		pbk           string
		sid           string
		sni           string
		allowInsecure bool // 专门处理 insecure/allowInsecure
		a1            []string
		a2            []string
		a3            []string
		a4            []string
		a5            []string
	)

	// 标准 vless://uuid@host:port?params#title 拆分
	a1 = strings.Split(a, "@")
	if len(a1) < 2 {
		return nil
	} // 基础防呆

	a2 = strings.Split(a1[1], "?")
	if len(a2) < 2 {
		return nil
	}

	a3 = strings.Split(a2[0], ":")
	a4 = strings.Split(a2[1], "#")
	a5 = strings.Split(a4[0], "&")

	for _, item := range a5 {
		x := strings.Split(item, "=")
		if len(x) < 2 {
			continue
		}

		key := x[0]
		val := x[1]

		switch key {
		case "flow":
			flow = val
		case "security":
			security = val
		case "fp":
			fp = val
		case "type":
			net = val
		case "host":
			host = val
		case "path":
			path = val
		case "sni":
			sni = val
		case "pbk":
			pbk = val
		case "sid":
			sid = val
		case "insecure", "allowInsecure":
			// 只要出现 1，就视为允许不安全连接
			if val == "1" || strings.ToLower(val) == "true" {
				allowInsecure = true
			}
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

	// --- 核心逻辑转换 ---

	// 1. 处理 SNI 和 Host 的逻辑归正
	if sni == "" && host != "" {
		sni = host
	}

	// 2. 将链接里的 security 映射到 CodeList 的 StreamSecurity
	// 链接里 security=tls 对应 config 的 "tls"
	// 链接里 security=reality 对应 config 的 "reality"
	streamSecurity := ""
	if security == "tls" || security == "reality" {
		streamSecurity = security
	}

	j = &config.CodeList{
		Type:           "vless",
		Title:          title,
		Address:        a3[0],
		Port:           port,
		ID:             a1[0],  // VLESS 的 UUID 对应 ID 字段
		Security:       "none", // VLESS 内部加密固定为 none
		Network:        net,    // 原 net 现在对应 Network
		Path:           path,
		Host:           host,
		StreamSecurity: streamSecurity, // 对齐新字段
		Sni:            sni,
		Fingerprint:    fp,
		AllowInsecure:  allowInsecure, // 写入新处理的布尔值
		PublicKey:      pbk,           // REALITY 参数归位
		ShortId:        sid,           // REALITY 参数归位
		Flow:           flow,
	}

	return
}
