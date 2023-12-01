package datafactory

import (
	config "xpanel/Config"
)

// makeDate make date
func MakeDate(a string, n string) (j *config.CodeList) {
	switch a {
	case "vless":
		j = VlessToJSON(n)
	case "ss":
		j = SsToJSON(n)
	case "trojan":
		j = TrojanToJSON(n)
	case "vmess":
		j = V2rayToJSON(n)
	case "ssr":
		j = SsrToJSON(n)
	}
	return
}
