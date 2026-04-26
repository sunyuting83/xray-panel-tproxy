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
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	config "xpanel/Config"
	datafactory "xpanel/DataFactory"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/proxy"
)

type Auths struct {
	User     string `json:"user"`
	Password string `json:"pass"`
}

func GetData(u string, p bool) (s []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	if p {
		dialer, err := proxy.SOCKS5("tcp", "localhost:7891", nil, proxy.Direct)
		if err != nil {
			return []byte(""), err
		}
		client = &http.Client{
			Timeout:   time.Duration(15 * time.Second),
			Transport: &http.Transport{Dial: dialer.Dial},
		}
	}
	reqest, err := http.NewRequest("GET", u, nil)

	reqest.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	reqest.Header.Set("Content-Type", "application/json")
	reqest.Header.Set("X-Requested-With", "XMLHttpRequest")
	reqest.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	if err != nil {
		return []byte(""), err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return []byte(""), err
	}
	defer response.Body.Close()
	d, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}
	return d, nil
}

// SyncGetData get data
func SyncGetData(list []string, porxy bool) (d []string) {
	var wg sync.WaitGroup
	d = []string{}
	for i := 0; i < len(list); i++ {
		i0 := i
		wg.Add(1)
		go func() {
			res, errers := GetData(list[i0], porxy)
			if errers != nil {
				d = append(d, "")
			}
			d = append(d, string(res))
			wg.Done()
		}()
	}
	wg.Wait()
	return
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
			// title := val.Elem().Field(val.Elem().NumField() - 1).interface().(string)
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
	// key 为节点唯一性特征字符串，value 为占位布尔值
	seen := make(map[string]bool)

	for _, node := range personList {
		if node.Port == 0 {
			continue
		}

		// 生成节点的唯一指纹字符串
		// 我们将影响连接的核心参数拼接在一起
		fingerprint := fmt.Sprintf("%s|%s|%d|%s|%s|%s|%s|%s|%s|%s|%s",
			node.Type,
			node.Address,
			node.Port,
			node.ID,
			node.Password,
			node.Network,
			node.StreamSecurity,
			node.PublicKey,
			node.ShortId,
			node.Flow,
			node.Method,
		)

		if !seen[fingerprint] {
			seen[fingerprint] = true
			result = append(result, node)
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

func DeleteNode(uid string, p string) bool {
	// 1. 定义两个潜在的数据源
	files := []string{
		filepath.Join(p, "data", "data.json"),   // 订阅节点
		filepath.Join(p, "data", "manual.json"), // 手动节点
	}

	for _, jsonFile := range files {
		// 读取文件
		data, err := os.ReadFile(jsonFile)
		if err != nil || len(data) == 0 {
			continue
		}

		// 处理可能存在的 null 字符截断（沿用你原来的安全处理）
		if index := bytes.IndexByte(data, 0); index != -1 {
			data = data[:index]
		}

		var nodes []config.CodeList
		if err := json.Unmarshal(data, &nodes); err != nil {
			continue
		}

		// 2. 使用 UID 查找并删除
		found := false
		for i, node := range nodes {
			if node.UID == uid {
				// 执行删除：利用 Go 切片特性
				nodes = append(nodes[:i], nodes[i+1:]...)
				found = true
				break
			}
		}

		// 3. 如果在该文件中找到了并删除了，写回文件并返回成功
		if found {
			saveConfig, _ := json.Marshal(nodes) // 使用 Indent 方便调试查看
			err = os.WriteFile(jsonFile, saveConfig, 0644)
			return err == nil
		}
	}

	// 两个文件都没找到
	return false
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

func GetDns(p string) ([]any, error) {
	var m []any
	// 1. 拼接绝对路径，指向你现在的 DNS 模板文件
	dnsTmplPath := filepath.Join(p, "template", "dns_servers.tmpl")

	// 2. 读取文件内容
	data, err := os.ReadFile(dnsTmplPath)
	if err != nil {
		return m, fmt.Errorf("读取 DNS 模板失败: %v", err)
	}

	// 3. 定义一个临时结构体来匹配模板的 JSON 结构
	// 模板内容是 {"servers": [...]}
	var temp struct {
		Servers []any `json:"servers"`
	}

	// 4. 直接解析整个 JSON
	err = json.Unmarshal(data, &temp)
	if err != nil {
		return m, fmt.Errorf("解析 DNS JSON 失败: %v", err)
	}

	// 5. 返回数组部分
	return temp.Servers, nil
}

func SetDns(p, data string) bool {
	// 1. 验证前端传来的 data (即 servers 数组部分) 是否是合法的 JSON 数组
	var m []any
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		fmt.Printf("DNS 数组解析失败: %v\n", err)
		return false
	}

	// 2. 构造完整的 DNS 模板内容: {"servers": [...]}
	// 使用结构体序列化比字符串拼接更安全，能自动处理转义和格式
	dnsStruct := struct {
		Servers []any `json:"servers"`
	}{
		Servers: m,
	}

	finalJson, err := json.Marshal(dnsStruct)
	if err != nil {
		return false
	}

	// 3. 写入模板文件 (绝对路径)
	dnsTmplPath := filepath.Join(p, "template", "dns_servers.tmpl")
	err = os.WriteFile(dnsTmplPath, finalJson, 0644)
	if err != nil {
		fmt.Printf("写入 DNS 模板失败: %v\n", err)
		return false
	}

	// 4. 重点：触发核心引擎，重新合成 config.json 并重启 Xray
	// 这里直接调用我们之前写的 GenerateConfig
	if err := GenerateConfig(p, "reload"); err != nil {
		fmt.Printf("合成新配置失败: %v\n", err)
		return false
	}

	return true
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

// SaveConfigFile save config file
func SaveConfigFile(pid string, r string) {
	getConfig := GetConfig()
	getConfig.Current = pid
	saveConfig, _ := json.Marshal(getConfig)
	os.WriteFile(r, saveConfig, 0644)
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

func CheckVersion(CurrentPath string, proxy bool) map[string]any {
	var data map[string]any = make(map[string]any)
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
	// CoreGeoIPFileName := strings.Join([]string{CorePath, "geoip.dat"}, "/")
	// CoreGeoSiteFileName := strings.Join([]string{CorePath, "geosite.dat"}, "/")
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
			// fmt.Println("开始下载规则库")
			Unzip(CoreZipFileName, CorePath)
			// os.Remove(CoreGeoIPFileName)
			// os.Remove(CoreGeoSiteFileName)
			// _, GeoIPUri, GeoSiteUri, GeoVersion, err := GetVersionData(config.ProxyUrl, config.GeoVersionUrl, "", true, false)
			// if err != nil {
			// 	fmt.Println("获取IP规则库失败,请重新启动")
			// 	os.Exit(0)
			// }
			// for _, item := range config.ProxyUrl {
			// 	uri := strings.Join([]string{item, GeoIPUri}, "")
			// 	err := DownloadFile(uri, CoreGeoIPFileName)
			// 	fmt.Println("下载IP规则库")
			// 	if err != nil {
			// 		fmt.Println("忽略下面的错误,开始魔法下载.速度很慢,请耐心等待")
			// 		fmt.Println(err)
			// 	} else {
			// 		fmt.Println("下载IP规则库成功")
			// 		break
			// 	}
			// }
			// for _, item := range config.ProxyUrl {
			// 	uri := strings.Join([]string{item, GeoSiteUri}, "")
			// 	err := DownloadFile(uri, CoreGeoSiteFileName)
			// 	fmt.Println("下载域名规则库")
			// 	if err != nil {
			// 		fmt.Println("忽略下面的错误,开始魔法下载.速度很慢,请耐心等待")
			// 		fmt.Println(err)
			// 	} else {
			// 		fmt.Println("下载域名规则库成功")
			// 		break
			// 	}
			// }
			// config.GeoVersion = GeoVersion
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
	// 默认参数
	psParam := "-ef"

	// 如果是 Alpine 系统，去掉 -ef 参数
	if _, err := os.Stat("/etc/alpine-release"); err == nil {
		psParam = ""
	}

	// 动态拼接命令
	comd := fmt.Sprintf("ps %s | grep xray | grep -v grep", psParam)

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
		if item == "" {
			uri = VersionUrl
		}
		data, err := GetData(uri, proxy)
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
		data, err := GetData(uri, false)
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
func Float64ToStringWithPrecision(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}

// checkType check type
func checkType(a string) (c bool) {
	l := []string{"vless", "ss", "vmess", "ssr", "trojan"}
	return slices.Contains(l, a)
}

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
