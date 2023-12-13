package utils

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	config "xpanel/Config"
	datafactory "xpanel/DataFactory"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type Auths struct {
	User     string `json:"user"`
	Password string `json:"pass"`
}

// CORSMiddleware cors middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// checkType check type
func checkType(a string) (c bool) {
	l := []string{"vless", "ss", "vmess", "ssr", "trojan"}
	for _, item := range l {
		if item == a {
			return true
		}
	}
	return false
}

// GetCurrentPath Get Current Path
func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
}

// GetConfig get config
func GetConfig() (j *config.Config) {
	path, _ := os.Executable()
	dir := filepath.Dir(path)
	jsonFile := strings.Join([]string{dir, "data/config.json"}, "/")
	configByte, _ := os.ReadFile(jsonFile)
	var (
		index int = len(configByte)
	)
	index = bytes.IndexByte(configByte, 0)
	if index != -1 {
		configByte = configByte[:index]
	}
	if err := json.Unmarshal(configByte, &j); err != nil {
		return
	}
	return
}

/* upDate function start */

// IgnoreTag Ignore Tag
func IgnoreTag(a []*config.CodeList, ignore string) []*config.CodeList {
	if len(ignore) > 0 {
		var ignoreList []string
		if strings.Contains(ignore, "|") {
			ignoreList = strings.Split(ignore, "|")
		} else {
			ignoreList = append(ignoreList, ignore)
		}
		var temp []*config.CodeList
		for _, item := range a {
			// val := reflect.ValueOf(item)
			// title := val.Elem().Field(val.Elem().NumField() - 1).Interface().(string)
			exist := false
			for _, ig := range ignoreList {
				if strings.Contains(item.Title, ig) {
					exist = true
				}
			}
			if !exist {
				temp = append(temp, item)
			}
		}
		return temp
	}
	return a
}

// RemoveRepeatedElement Remove Repeated Element
func RemoveRepeatedElement(personList []*config.CodeList) (result []*config.CodeList) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < len(personList); j++ {
			if personList[i].Types == personList[j].Types &&
				personList[i].Address == personList[j].Address &&
				personList[i].Port == personList[j].Port &&
				personList[i].Password == personList[j].Password &&
				personList[i].Security == personList[j].Security &&
				personList[i].Net == personList[j].Net &&
				personList[i].Path == personList[j].Path &&
				personList[i].TLS == personList[j].TLS &&
				personList[i].Aid == personList[j].Aid &&
				personList[i].Method == personList[j].Method &&
				personList[i].Protocol == personList[j].Protocol &&
				personList[i].ProtocolParam == personList[j].ProtocolParam &&
				personList[i].Obfs == personList[j].Obfs &&
				personList[i].ObfsParam == personList[j].ObfsParam &&
				personList[i].Host == personList[j].Host {
				repeat = true
				break
			}
		}
		if !repeat {
			if personList[i].Port != 0 {
				result = append(result, personList[i])
			}
		}
	}
	return
}

// makeDates make dates
func MakeDates(a []string) (b []*config.CodeList) {
	for _, v := range a {
		list := strings.Split(DeCodeBytes(v), "\n")
		for _, item := range list {
			if len(item) > 0 {
				if strings.Contains(item, "://") {
					a := strings.Split(item, "://")
					t := a[0]
					n := a[1]
					if checkType(t) {
						b = append(b, datafactory.MakeDate(t, n))
					}
				}
			}
		}
	}
	return
}

/* upDate function end
------------------------
nodeList function start
*/

// ListToJsons fun
func ListToJsons(s []byte) (result *[]config.CodeList) {
	var (
		index int = len(s)
	)
	index = bytes.IndexByte(s, 0)
	if index != -1 {
		s = s[:index]
	}
	if err := json.Unmarshal(s, &result); err != nil {
		return
	}
	return
}

/* nodeList function end
-------------------------
setNode function start
*/

// FormToJSON Form To JSON
func FormToJSON(s string) (result *config.CodeList, err error) {
	decode := DeCodeBytes(s)
	code := []byte(decode)
	var (
		index int = len(code)
	)
	index = bytes.IndexByte(code, 0)
	if index != -1 {
		code = code[:index]
	}
	if err = json.Unmarshal(code, &result); err != nil {
		return
	}
	return
}

// SetNodeToUnix Set Node To Unix
func SetNodeToUnix(i int, p string) (b bool) {
	j := GetNode(i, p)
	if j == nil {
		return false
	}
	/*
		switch j.Types {
		case "vless":
			b = SetVmess(j, p, c, r, cu)
		case "vmess":
			b = SetVmess(j, p, c, r, cu)
		case "ss":
			b = SetVmess(j, p, c, r, cu)
		case "ssr":
			b = SetSSR(j, p, c, r, cu)
		case "trojan":
			b = SetVmess(j, p, c, r, cu)
		}
	*/
	b = SetVmess(j, p)
	return
}

// ReSetNodeToUnix Restart Set Node To Unix
func GetCurrentNode(p string) int {
	Index := 0
	configs := GetConfig()
	jsonFile := strings.Join([]string{p, "data/dataFile"}, "/")
	data, _ := os.ReadFile(jsonFile)
	var result []config.CodeList
	if len(data) > 0 {
		var (
			index int = len(data)
		)
		index = bytes.IndexByte(data, 0)
		if index != -1 {
			data = data[:index]
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return Index
		}
		Arrlen := len(result)
		if Arrlen > 0 {
			for index, item := range result {
				if item.Title == configs.Current {
					Index = index
					break
				}
			}
		}
		return Index
	}
	return Index
}

// ReSetNodeToUnix Restart Set Node To Unix
func ReSetNodeToUnix(p string) {
	configs := GetConfig()
	jsonFile := strings.Join([]string{p, "data/dataFile"}, "/")
	data, _ := os.ReadFile(jsonFile)
	var result []config.CodeList
	if len(data) > 0 {
		var (
			index int = len(data)
		)
		index = bytes.IndexByte(data, 0)
		if index != -1 {
			data = data[:index]
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return
		}
		Arrlen := len(result)
		if Arrlen > 0 {
			for index, item := range result {
				if item.Title == configs.Current {
					j := GetNode(index, p)
					SetVmess(j, p)
					break
				}
			}
		}
	}
}

func DeleteNode(i int, p string) (b bool) {
	jsonFile := strings.Join([]string{p, "data/dataFile"}, "/")
	data, _ := os.ReadFile(jsonFile)
	var result []config.CodeList
	if len(data) > 0 {
		var (
			index int = len(data)
		)
		index = bytes.IndexByte(data, 0)
		if index != -1 {
			data = data[:index]
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return
		}
		Arrlen := len(result)
		if Arrlen > 0 {
			for index := range result {
				if index == i {
					if i == Arrlen-1 {
						result = result[0:i]
					} else {
						result = append(result[0:i], result[i+1:]...)
					}
					saveConfig, _ := json.Marshal(result)
					os.WriteFile(jsonFile, saveConfig, 0644)
					return true
				}
			}
		}
	}
	return false
}

func GetRules(p string) ([]interface{}, string, map[string]interface{}, error) {
	jsonFile := strings.Join([]string{p, "template/tempEnd"}, "/")
	d, _ := os.ReadFile(jsonFile)
	data := string(d)
	ignoreIndex := strings.Index(data, "255}}}],") + 8
	ignoreString := data[ignoreIndex:]
	content := strings.Join([]string{"{", ignoreString}, "")
	m := make(map[string]interface{})
	resolve := make([]interface{}, 0)
	err := json.Unmarshal([]byte(content), &m)
	if err != nil {
		return resolve, "", m, err
	}
	rules := m["routing"].(map[string]interface{})["rules"].([]interface{})
	return rules, data[0:ignoreIndex], m, nil
}

func GetDomains(p string) (map[string]interface{}, map[string]interface{}, string, []string, []string, error) {
	resolve := make(map[string]interface{}, 0)
	rules, startStr, jsons, err := GetRules(p)
	if err != nil {
		return resolve, jsons, "", make([]string, 0), make([]string, 0), err
	}
	proxyDomain := make([]string, 0)
	directDomain := make([]string, 0)
	for index, item := range rules {
		if index == 7 {
			proxyd := item.(map[string]interface{})["domain"].([]interface{})
			for _, domain := range proxyd {
				doString := fmt.Sprint(domain)
				if strings.Contains(doString, ":") {
					do := strings.Split((doString), ":")[1]
					proxyDomain = append(proxyDomain, do)
				}
			}
		}
		if index == 8 {
			direct := item.(map[string]interface{})["domain"].([]interface{})
			for _, domain := range direct {
				doString := fmt.Sprint(domain)
				if strings.Contains(doString, ":") {
					do := strings.Split((doString), ":")[1]
					directDomain = append(directDomain, do)
				}
			}
		}
	}
	resolve["proxyDomain"] = strings.Join(proxyDomain, "\n")
	resolve["directDomain"] = strings.Join(directDomain, "\n")
	return resolve, jsons, startStr, proxyDomain, directDomain, nil
}

func SetDomains(p string, domains config.Domains) bool {
	_, jsons, startStr, proxy, direct, err := GetDomains(p)
	if err != nil {
		return false
	}
	var (
		formProxy  []string
		formDirect []string
	)
	if strings.Contains(domains.Proxy, "\n") {
		formProxy = RemoveRepeatedSingle(strings.Split(domains.Proxy, "\n"))
	}
	if strings.Contains(domains.Direct, "\n") {
		formDirect = RemoveRepeatedSingle(strings.Split(domains.Direct, "\n"))
	}
	for index, item := range formProxy {
		if strings.Contains(item, ".") {
			formProxy[index] = strings.Join([]string{"domain", item}, ":")
		} else {
			formProxy[index] = strings.Join([]string{"geosite", item}, ":")
		}
	}
	for index, item := range formDirect {
		formDirect[index] = strings.Join([]string{"domain", item}, ":")
	}
	ignoreProxy := StringSliceToInterfaceSlice(IgnoreRepeated(formProxy, proxy))
	ignoreDirect := StringSliceToInterfaceSlice(IgnoreRepeated(formDirect, direct))
	rules := jsons["routing"].(map[string]interface{})["rules"].([]interface{})
	for index, item := range rules {
		if index == 7 {
			item.(map[string]interface{})["domain"] = ignoreProxy
		}
		if index == 8 {
			item.(map[string]interface{})["domain"] = ignoreDirect
		}
	}
	// fmt.Println(startStr)
	newRules, _ := json.Marshal(jsons)
	newData := strings.Join([]string{startStr, string(newRules)[1:]}, "")
	jsonFile := strings.Join([]string{p, "template/tempEnd"}, "/")
	os.WriteFile(jsonFile, []byte(newData), 0644)
	return true
}

func GetSubscribes(p string) (string, error) {
	jsonFile := strings.Join([]string{p, "data/subUrl"}, "/")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func GetIgnore(p string) (string, error) {
	jsonFile := strings.Join([]string{p, "data/ignore"}, "/")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetDns(p string) (interface{}, error) {
	var m interface{}
	jsonFile := strings.Join([]string{p, "template/tempStart"}, "/")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return m, err
	}
	dataStr := string(data)
	if strings.Contains(dataStr, `"dns":{"servers":`) {
		d1 := strings.Split(dataStr, `"dns":{"servers":`)[1]
		d2 := strings.Split(d1, `},"inbounds":`)[0]
		if len(d2) != 0 {
			err := json.Unmarshal([]byte(d2), &m)
			if err != nil {
				return m, err
			}
		}
		return m, nil
	}
	return m, err
}

func SetDns(p, data string) bool {
	var m interface{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return false
	}
	jsonFile := strings.Join([]string{p, "template/tempStart"}, "/")
	d0, err := os.ReadFile(jsonFile)
	if err != nil {
		return false
	}
	dataStr := string(d0)
	if strings.Contains(dataStr, `"dns":{"servers":`) {
		d1 := strings.Split(dataStr, `"dns":{"servers":`)[1]
		d2 := strings.Split(d1, `},"inbounds":`)[1]
		d3 := strings.Join([]string{`{"dns":{"servers":`, data, `},"inbounds":`, d2}, "")
		os.WriteFile(jsonFile, []byte(d3), 0644)
		return true
	}
	return false
}

func GetLocalSocks(p string) map[string]interface{} {
	var m map[string]interface{} = make(map[string]interface{})
	m["status"] = 0
	m["SocksStatus"] = false
	jsonFile := strings.Join([]string{p, "template/tempStart"}, "/")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		m["status"] = 1
		return m
	}
	dataStr := string(data)
	if strings.Contains(dataStr, `"protocol":"socks"`) {
		d1 := strings.Split(dataStr, `"auth":"`)[1]
		statuStr := strings.Split(d1, `","ip":"0.0.0.0","udp"`)[0]
		// fmt.Println(statuStr)
		accountListStr := strings.Split(d1, `"accounts":`)[1]
		accountListStr = strings.Split(accountListStr, `},"sniffing"`)[0]
		// fmt.Println(accountListStr)
		var accounts *[]Auths
		err := json.Unmarshal([]byte(accountListStr), &accounts)
		if err != nil {
			fmt.Println(err)
		}
		if len(statuStr) != 0 {
			if statuStr == "password" {
				m["SocksStatus"] = true
			}
			m["Auths"] = accounts
		}
		return m
	}
	m["status"] = 1
	return m
}

func SetLocalSocks(p, socks, auths string) bool {
	var accounts *[]Auths
	err := json.Unmarshal([]byte(auths), &accounts)
	if err != nil {
		return false
	}
	jsonFile := strings.Join([]string{p, "template/tempStart"}, "/")
	d0, err := os.ReadFile(jsonFile)
	if err != nil {
		return false
	}
	dataStr := string(d0)
	if strings.Contains(dataStr, `"protocol":"socks"`) {
		d1 := strings.Split(dataStr, `"auth":"`)
		d2 := strings.Split(d1[1], `","ip":"0.0.0.0","udp"`)
		accountLis := strings.Split(d2[1], `"accounts":`)
		// fmt.Println(statuStr)
		accountListStr := strings.Split(accountLis[1], `},"sniffing"`)
		d3 := strings.Join([]string{d1[0], `"auth":"`, socks, `","ip":"0.0.0.0","udp"`, accountLis[0], `"accounts":`, auths, `},"sniffing"`, accountListStr[1]}, "")
		os.WriteFile(jsonFile, []byte(d3), 0644)
		return true
	}
	return false
}

func SetSubscribes(p, data string) bool {
	var newSub []string
	if strings.Contains(data, "\n") {
		newSub = RemoveRepeatedSingle(strings.Split(data, "\n"))
		target := ""
		index := -1
		for i, num := range newSub {
			if num == target {
				index = i
				break
			}
		}
		if index != -1 {
			newSub = append(newSub[:index], newSub[index+1:]...)
		}
	} else {
		newSub = append(newSub, data)
	}
	if len(newSub) > 0 {
		jsonFile := strings.Join([]string{p, "data/subUrl"}, "/")
		subContent := strings.Join(newSub, "\n")
		os.WriteFile(jsonFile, []byte(subContent), 0644)
		return true
	}
	return false
}

func SetIgnore(p, data string) bool {
	jsonFile := strings.Join([]string{p, "data/ignore"}, "/")
	os.WriteFile(jsonFile, []byte(data), 0644)
	return true
}

func GetNode(i int, p string) (j *config.CodeList) {
	jsonFile := strings.Join([]string{p, "data/dataFile"}, "/")
	data, _ := os.ReadFile(jsonFile)
	var result []*config.CodeList
	if len(data) > 0 {
		var (
			index int = len(data)
		)
		index = bytes.IndexByte(data, 0)
		if index != -1 {
			data = data[:index]
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return
		}
		if len(result) >= i {
			for index, item := range result {
				if index == i {
					j = item
					break
				}
			}
		}
	}
	return
}

func StringSliceToInterfaceSlice(input []string) []interface{} {
	result := make([]interface{}, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}

func IgnoreRepeated(postList, dataList []string) []string {
	if len(dataList) != 0 {
		var temp []string
		for _, item := range postList {
			exist := false
			for _, ig := range dataList {
				if item == ig {
					exist = true
				}
			}
			if !exist {
				temp = append(temp, item)
			}
		}
		// fmt.Println(temp)
		return temp
	}
	return postList
}

// RemoveRepeatedSingle Remove Repeated Element
func RemoveRepeatedSingle(personList []string) (result []string) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if personList[i] == personList[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			result = append(result, personList[i])
		}
	}
	return
}

// ReadConfigFile read config file
func ReadConfigFile(p, types string) (code []byte) {
	T := strings.ToLower(types)
	tempStart := strings.Join([]string{p, "template", "tempStart"}, "/")
	s, _ := os.ReadFile(tempStart)
	jsonFile := strings.Join([]string{p, "template", T}, "/")
	c, _ := os.ReadFile(jsonFile)
	tempEnd := strings.Join([]string{p, "template", "tempEnd"}, "/")
	e, _ := os.ReadFile(tempEnd)
	result := append(append(s, c...), e...)
	return result
}

// SaveConfigFile save config file
func SaveConfigFile(pid string, r string) {
	getConfig := GetConfig()
	getConfig.Current = pid
	saveConfig, _ := json.Marshal(getConfig)
	os.WriteFile(r, saveConfig, 0644)
}

// SetVmess set vmess
func SetVmess(j *config.CodeList, p string) (b bool) {
	code := ReadConfigFile(p, j.Types)
	m := make(map[string]interface{})
	json.Unmarshal([]byte(code), &m)
	protocol := j.Types
	if j.Types == "ss" {
		protocol = "shadowsocks"
	}
	if j.Types == "trojan" {
		protocol = "trojan"
	}
	m["outbounds"].([]interface{})[0].(map[string]interface{})["protocol"] = protocol
	switch j.Types {
	case "vmess":
		vnext := m["outbounds"].([]interface{})[0].(map[string]interface{})["settings"].(map[string]interface{})["vnext"].([]interface{})[0]
		vnext.(map[string]interface{})["address"] = j.Address
		vnext.(map[string]interface{})["port"] = j.Port
		users := vnext.(map[string]interface{})["users"].([]interface{})[0]
		users.(map[string]interface{})["id"] = j.Password
		users.(map[string]interface{})["alterId"] = j.Aid
	case "vless":
		vnext := m["outbounds"].([]interface{})[0].(map[string]interface{})["settings"].(map[string]interface{})["vnext"].([]interface{})[0]
		vnext.(map[string]interface{})["address"] = j.Address
		vnext.(map[string]interface{})["port"] = j.Port
		users := vnext.(map[string]interface{})["users"].([]interface{})[0]
		users.(map[string]interface{})["id"] = j.Password
	case "ss":
		users := m["outbounds"].([]interface{})[0].(map[string]interface{})["settings"].(map[string]interface{})["servers"].([]interface{})[0]
		users.(map[string]interface{})["address"] = j.Address
		users.(map[string]interface{})["port"] = j.Port
		users.(map[string]interface{})["method"] = j.Method
		users.(map[string]interface{})["password"] = j.Password
	case "trojan":
		users := m["outbounds"].([]interface{})[0].(map[string]interface{})["settings"].(map[string]interface{})["servers"].([]interface{})[0]
		users.(map[string]interface{})["address"] = j.Address
		users.(map[string]interface{})["port"] = j.Port
		users.(map[string]interface{})["password"] = j.Password
		streamSettings := m["outbounds"].([]interface{})[0].(map[string]interface{})["streamSettings"].(map[string]interface{})["tlsSettings"]
		streamSettings.(map[string]interface{})["allowInsecure"] = true
		streamSettings.(map[string]interface{})["serverName"] = j.Host
	}
	if j.Types != "ss" {
		if j.Types != "trojan" {
			streamSettings := m["outbounds"].([]interface{})[0].(map[string]interface{})["streamSettings"]
			streamSettings.(map[string]interface{})["network"] = j.Net
			tls := "none"
			if j.TLS {
				tls = "tls"
			}
			streamSettings.(map[string]interface{})["security"] = tls
			streamSettings.(map[string]interface{})["wsSettings"].(map[string]interface{})["headers"].(map[string]interface{})["Host"] = j.Host
			streamSettings.(map[string]interface{})["wsSettings"].(map[string]interface{})["path"] = j.Path
		}
	}
	saveData, _ := json.Marshal(m)
	c := strings.Join([]string{p, "Core/config.json"}, "/")
	err := os.WriteFile(c, saveData, 0644)
	if err != nil {
		return false
	}
	RunXray(p, "reload", j.Title)
	return true
}

// SetSSR set ssr
func SetSSR(j *config.CodeList, p string, c string, r string, cu string) (b bool) {
	s := ReadConfigFile(p, j.Types)
	var ssr *config.SSR
	//3.json解析到结构体
	if err := json.Unmarshal(s, &ssr); err != nil {
		return false
	}
	ssr.Server = j.Address
	ssr.ServerPort = j.Port
	ssr.Method = j.Method
	ssr.Protocol = j.Protocol
	ssr.ProtocolParam = j.ProtocolParam
	ssr.Obfs = j.Obfs
	ssr.ObfsParam = j.ObfsParam
	ssr.Password = j.Password
	saveData, _ := json.Marshal(&ssr)
	err := os.WriteFile(c, saveData, 0644)
	if err != nil {
		return false
	}
	if cu == "ssr" {
		RunCommand("/etc/config/sh/ssr.sh restart")
		return true
	}
	SaveConfigFile("ssr", r)
	return true
}

// SetTrojan set trojan
func SetTrojan(j *config.CodeList, p string, c string, r string, cu string) (b bool) {
	s := ReadConfigFile(p, j.Types)
	var trojan *config.Trojan
	if err := json.Unmarshal(s, &trojan); err != nil {
		return false
	}
	trojan.RemoteAddr = j.Address
	trojan.RemotePort = j.Port
	trojan.Password[0] = j.Password
	trojan.Ssl.Sni = j.Host
	saveData, _ := json.Marshal(&trojan)
	err := os.WriteFile(c, saveData, 0644)
	if err != nil {
		return false
	}
	if cu == "trojan" {
		RunCommand("/etc/config/sh/trojan.sh restart")
		return true
	}
	go RunCommand("/etc/config/sh/trojan.sh start")
	SaveConfigFile("trojan", r)
	return true
}

// RunCommand run command
func RunCommand(command string) (pidstr string) {
	// fmt.Println(command)
	cmd := exec.Command("/bin/sh", "-c", command, " &")
	cmd.Start()
	pid := cmd.Process.Pid
	cmd.Wait()
	pidstr = strconv.Itoa(pid)
	return
}

// deCodeBytes
func DeCodeBytes(a string) (b string) {
	var str []byte = []byte(a)
	decodeBytes := make([]byte, base64.StdEncoding.DecodedLen(len(str))) // 计算解码后的长度
	base64.StdEncoding.Decode(decodeBytes, str)
	return string(decodeBytes)
}

// GetSubUrl
func GetSubUrl(urList []byte, current_path string) []string {
	var l []string
	if len(urList) > 0 {
		str := string(urList)
		list := strings.Split(str, "\n")
		return list
	}
	ConfigFilePath := strings.Join([]string{current_path, "data/config"}, "/")
	ConfigFile, _ := os.ReadFile(ConfigFilePath)
	var config *config.Configs
	//3.json解析到结构体
	if err := json.Unmarshal(ConfigFile, &config); err != nil {
		return l
	}
	l = append(l, config.DefaultSubUrl)
	return l
}

func DownloadFile(URL, filepath string) error {
	// 创建HTTP客户端，并设置15秒超时时间
	client := http.Client{
		Timeout: 600 * time.Second,
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// fmt.Println(resp.StatusCode)
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to download file, status code: %d", resp.StatusCode)
	}
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 下载文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
func Unzip(source, destination string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		filePath := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		fileDir := filepath.Dir(filePath)
		err := os.MkdirAll(fileDir, 0755)
		if err != nil {
			return err
		}

		writer, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer writer.Close()

		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetPlatform(platform string) string {
	if strings.Contains(platform, "amd") {
		return strings.ReplaceAll(platform, "amd", "")
	}
	if strings.Contains(platform, "arm64") {
		return "arm64-v8a"
	}
	if strings.Contains(platform, "arm32") {
		return "arm64-v7a"
	}
	return platform
}

func CheckVersion(CurrentPath string, proxy bool) map[string]interface{} {
	var data map[string]interface{} = make(map[string]interface{})
	data["status"] = 0
	data["CoreVersion"] = false
	data["GeoVersion"] = false
	config := GetConfig()
	OS := runtime.GOOS
	platform := runtime.GOARCH
	Platform := GetPlatform(platform)
	CoreFile := strings.Join([]string{"Xray", OS, Platform}, "-")
	CoreZip := strings.Join([]string{CoreFile, "zip"}, ".")

	CoreUri, _, _, version, err := GetVersionData(config.ProxyUrl, config.GetCoreVersionUrl, CoreZip, false, proxy)
	if err != nil {
		data["status"] = 1
		data["message"] = err.Error()
		data["CoreVersion"] = false
	}
	if len(version) != 0 && config.CoreVersion != version {
		data["CoreVersion"] = true
		data["CoreUri"] = CoreUri
		data["LocalCoreVersion"] = config.CoreVersion
		data["CurrentCoreVersion"] = version
	}
	_, GeoIP, GeoSite, GeoVersion, err := GetVersionData(config.ProxyUrl, config.GeoVersionUrl, "", true, proxy)
	if err != nil {
		data["status"] = 1
		data["message"] = err.Error()
		data["GeoVersion"] = false
	}
	if len(GeoVersion) != 0 && config.GeoVersion != GeoVersion {
		data["GeoVersion"] = true
		data["GeoIP"] = GeoIP
		data["GeoSite"] = GeoSite
		data["LocalGeoVersion"] = config.GeoVersion
		data["CurrentGeoVersion"] = GeoVersion
	}
	return data
}

func CheckCore(OS, platform, CurrentPath string) {
	config := GetConfig()
	Platform := GetPlatform(platform)
	CoreFile := strings.Join([]string{"Xray", OS, Platform}, "-")
	CorePath := strings.Join([]string{CurrentPath, "Core"}, "/")
	if !IsExist(CorePath) {
		os.MkdirAll(CorePath, 0755)
	}
	CoreZip := strings.Join([]string{CoreFile, "zip"}, ".")
	CoreZipFileName := strings.Join([]string{CorePath, CoreZip}, "/")
	CoreGeoIPFileName := strings.Join([]string{CorePath, "geoip.dat"}, "/")
	CoreGeoSiteFileName := strings.Join([]string{CorePath, "geosite.dat"}, "/")
	CoreFileName := strings.Join([]string{CorePath, "xray"}, "/")
	if !IsExist(CoreFileName) {
		fmt.Println("核心不存在,请等待下载,整个过程预计半小时,取决于网络环境")
		downUri, _, _, version, err := GetVersionData(config.ProxyUrl, config.GetCoreVersionUrl, CoreZip, false, false)
		if err != nil {
			fmt.Println("获取核心网址失败,请重新启动")
			os.Exit(0)
		}
		MD5Uri := strings.Join([]string{downUri, "dgst"}, ".")
		MD5, err := GetVersionMD5(config.ProxyUrl, MD5Uri)
		if err != nil {
			fmt.Println("获取核心网址失败,请重新启动")
			os.Exit(0)
		}
		for _, item := range config.ProxyUrl {
			uri := strings.Join([]string{item, downUri}, "")
			err := DownloadFile(uri, CoreZipFileName)
			fmt.Println("下载核心开始,请耐心等待")
			if err != nil {
				fmt.Println("忽略下面的错误,开始魔法下载.速度很慢,请耐心等待")
				fmt.Println(err)
			} else {
				fmt.Println("下载核心完成")
				break
			}
		}
		fileMD5 := Md5File(CoreZipFileName)
		if fileMD5 == MD5 {
			fmt.Println("开始下载规则库")
			Unzip(CoreZipFileName, CorePath)
			os.Remove(CoreGeoIPFileName)
			os.Remove(CoreGeoSiteFileName)
			_, GeoIPUri, GeoSiteUri, GeoVersion, err := GetVersionData(config.ProxyUrl, config.GeoVersionUrl, "", true, false)
			if err != nil {
				fmt.Println("获取IP规则库失败,请重新启动")
				os.Exit(0)
			}
			for _, item := range config.ProxyUrl {
				uri := strings.Join([]string{item, GeoIPUri}, "")
				err := DownloadFile(uri, CoreGeoIPFileName)
				fmt.Println("下载IP规则库")
				if err != nil {
					fmt.Println("忽略下面的错误,开始魔法下载.速度很慢,请耐心等待")
					fmt.Println(err)
				} else {
					fmt.Println("下载IP规则库成功")
					break
				}
			}
			for _, item := range config.ProxyUrl {
				uri := strings.Join([]string{item, GeoSiteUri}, "")
				err := DownloadFile(uri, CoreGeoSiteFileName)
				fmt.Println("下载域名规则库")
				if err != nil {
					fmt.Println("忽略下面的错误,开始魔法下载.速度很慢,请耐心等待")
					fmt.Println(err)
				} else {
					fmt.Println("下载域名规则库成功")
					break
				}
			}
			config.GeoVersion = GeoVersion
			config.CoreVersion = version
			saveConfig, _ := json.Marshal(config)
			path, _ := os.Executable()
			dir := filepath.Dir(path)
			jsonFile := strings.Join([]string{dir, "data/config.json"}, "/")
			os.WriteFile(jsonFile, saveConfig, 0644)
		}
		os.Remove(CoreZipFileName)
		os.Chmod(CoreFileName, 0777)
	}
}

func RunCommandWithRes(cmdExec string) (k string, err error) {
	cmd := exec.Command("/bin/sh", "-c", cmdExec, " &")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		return "", err
	}

	if len(bytesErr) != 0 {
		return "", errors.New("0")

	}

	bytes, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckXray() bool {
	comd := "ps | grep xray | grep -v grep"
	hasStatus := false
	str, err := RunCommandWithRes(comd)
	if err != nil {
		return false
	}
	if len(str) != 0 {
		if strings.Contains(str, "xray") {
			hasStatus = true
		}
	}
	return hasStatus
}

func RunXray(p, status, title string) {
	c := strings.Join([]string{p, "run.sh " + status}, "/")
	RunCommand(c)
	r := strings.Join([]string{p, "data/config.json"}, "/")
	SaveConfigFile(title, r)
}
func RunXrayWithoutConfig(status string) {
	p, _ := GetCurrentPath()
	c := strings.Join([]string{p, "run.sh " + status}, "/")
	RunCommand(c)
}

func IsExist(path string) bool {
	// 判断文件是否存在
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetVersionData(ProxyUri []string, VersionUrl, CoreZip string, geo, proxy bool) (uri, geoip, geosite, geoversion string, err error) {
	var versionData *config.JSONData
	var downUri string
	for _, item := range ProxyUri {
		uri := strings.Join([]string{item, VersionUrl}, "")
		data, err := datafactory.GetData(uri, proxy)
		if err == nil {
			var (
				index int = len(data)
			)
			index = bytes.IndexByte(data, 0)
			if index != -1 {
				data = data[:index]
			}
			if err = json.Unmarshal(data, &versionData); err != nil {
				return "", "", "", "", err
			}

			geoversion = versionData.TagName

			for _, v := range versionData.Assets {
				if geo {
					if v.Name == "geoip.dat" {
						geoip = v.BrowserDownloadURL
					}
					if v.Name == "geosite.dat" {
						geosite = v.BrowserDownloadURL
					}
				} else {
					if v.Name == CoreZip {
						downUri = v.BrowserDownloadURL
						break
					}
				}
			}
			return downUri, geoip, geosite, geoversion, nil
		}
	}
	return "", "", "", "", errors.New("has error")
}

func GetVersionMD5(ProxyUri []string, MD5uri string) (string, error) {
	var md5 string = ""
	for _, item := range ProxyUri {
		uri := strings.Join([]string{item, MD5uri}, "")
		data, err := datafactory.GetData(uri, false)
		if err == nil {
			dataSplit := strings.Split(string(data), "\n")
			for _, v := range dataSplit {
				if strings.Contains(v, "MD5= ") {
					md5 = strings.Split(v, "MD5= ")[1]
					break
				}
			}
			break
		}
	}
	return md5, nil
}

func Md5File(file string) string {
	f, _ := os.Open(file)
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

func Decimal(num float64) float64 {
	decimal := 2
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Floor(num*d)/d, 'f', -1, 64)
	floatNum, _ := strconv.ParseFloat(res, 64)
	return floatNum
}

func TestTCPing(m *melody.Melody, ID string) {
	current_path, _ := GetCurrentPath()
	jsonFile := strings.Join([]string{current_path, "data/dataFile"}, "/")
	data, _ := os.ReadFile(jsonFile)
	list := ListToJsons(data)
	if len(*list) > 0 {
		for i, item := range *list {
			port := strconv.Itoa(item.Port)
			index := strconv.Itoa(i)
			elapsedTime, err := TCPing(item.Address, port)
			speed := "0"
			if err == nil {
				speed = Float64ToStringWithPrecision(elapsedTime.Seconds()*1000, 2)
			}
			speeData := &config.Message{
				Type: "tcping",
				UUID: ID,
				Data: strings.Join([]string{index, speed}, "||||"),
			}
			sedData, _ := json.Marshal(speeData)
			m.Broadcast(sedData)
		}
	}
}

func Float64ToStringWithPrecision(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}

func TCPing(host string, port string) (time.Duration, error) {
	timeout := 5 * time.Second
	address := strings.Join([]string{host, port}, ":")

	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	elapsedTime := time.Since(startTime)

	return elapsedTime, nil
}
