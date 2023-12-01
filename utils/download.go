package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	config "xpanel/Config"

	"github.com/olahol/melody"
)

func SendMessageToWs(m *melody.Melody, ID, types, message string) {
	speeData := &config.Message{
		Type: types,
		UUID: ID,
		Data: message,
	}
	sedData, _ := json.Marshal(speeData)
	m.Broadcast(sedData)
}

func DownloadFileWithHeaders(url, filePath, ID, FileName, CoreFile, dataPath, Version string, m *melody.Melody) {
	client := &http.Client{
		Timeout: 600 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Println("创建请求时发生错误:", err)
		message := strings.Join([]string{FileName, "创建请求时发生错误"}, "----")
		SendMessageToWs(m, ID, "error", message)
		return
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("下载文件时发生错误:", err)
		message := strings.Join([]string{FileName, "下载文件时发生错误"}, "----")
		SendMessageToWs(m, ID, "error", message)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// fmt.Printf("下载文件时返回非200状态码: %d\n", resp.StatusCode)
		message := strings.Join([]string{FileName, "下载文件时返回非200状态码"}, "----")
		SendMessageToWs(m, ID, "error", message)
		return
	}

	out, err := os.Create(filePath)
	if err != nil {
		// fmt.Println("创建文件时发生错误:", err)
		message := strings.Join([]string{FileName, "创建文件时发生错误"}, "----")
		SendMessageToWs(m, ID, "error", message)
		return
	}
	defer out.Close()

	contentLength := resp.Header.Get("Content-Length")
	var totalSize int64
	fmt.Sscan(contentLength, &totalSize)

	buffer := make([]byte, 1024)
	var downloadedSize int64

	percen := 0
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			// fmt.Println("下载文件时发生错误:", err)
			message := strings.Join([]string{FileName, "下载文件时发生错误"}, "----")
			SendMessageToWs(m, ID, "error", message)
			return
		}

		if n == 0 {
			break
		}

		downloadedSize += int64(n)
		percentage := int(float64(downloadedSize) / float64(totalSize) * 100)
		// fmt.Printf("下载进度: %d%%\n", percentage)
		if percentage > percen {
			percenStr := strconv.Itoa(percentage)
			message := strings.Join([]string{FileName, percenStr}, "----")
			SendMessageToWs(m, ID, "download", message)
		}
		percen = percentage

		_, err = out.Write(buffer[:n])
		if err != nil {
			// fmt.Println("写入文件时发生错误:", err)
			message := strings.Join([]string{FileName, "写入文件时发生错误"}, "----")
			SendMessageToWs(m, ID, "error", message)
			return
		}
	}

	message := strings.Join([]string{FileName, "100"}, "----")
	SendMessageToWs(m, ID, "download", message)

	config := GetConfig()
	if strings.Contains(CoreFile, "Xray") {
		config.CoreVersion = Version
		RunXrayWithoutConfig("stop")

		CurrentPath, _ := GetCurrentPath()
		tmpPath := strings.Join([]string{CurrentPath, "tmp"}, "/")
		CoreTmpPath := strings.Join([]string{tmpPath, "Xcore"}, "/")
		CorePath := strings.Join([]string{CurrentPath, "Core"}, "/")
		NewCoreFile := strings.Join([]string{CorePath, "xray"}, "/")
		tmpFileXray := strings.Join([]string{CoreTmpPath, "xray"}, "/")
		// fmt.Println(NewCoreFile)
		// fmt.Println(CoreTmpPath)
		// fmt.Println(tmpFileXray)
		Unzip(filePath, CoreTmpPath)
		os.Remove(filePath)
		os.Remove(NewCoreFile)
		os.Rename(tmpFileXray, NewCoreFile)
		deleteFiles(CoreTmpPath)
		os.Remove(CoreTmpPath)
		os.Chmod(NewCoreFile, 0777)
		RunXrayWithoutConfig("start")
		message = strings.Join([]string{FileName, "已更新到最新版本"}, "----")
		SendMessageToWs(m, ID, "error", message)
	} else {
		RunXrayWithoutConfig("stop")
		os.Remove(CoreFile)
		os.Rename(filePath, CoreFile)
		RunXrayWithoutConfig("start")
		message = strings.Join([]string{FileName, "已更新到最新版本"}, "----")
		SendMessageToWs(m, ID, "error", message)
		config.GeoVersion = Version
	}
	saveConfig, _ := json.Marshal(config)
	jsonFile := strings.Join([]string{dataPath, "config.json"}, "/")
	os.WriteFile(jsonFile, saveConfig, 0644)
}

func createDirectoryIfNotExists(path string) error {
	// 判断路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 如果路径不存在，尝试创建目录
		err := os.MkdirAll(path, 0755) // 使用0755权限创建目录
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func deleteFiles(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(dir + "/" + name)
		if err != nil {
			return err
		}
	}
	return nil
}

func MakeWsData(m *melody.Melody, ID, Data string) {
	if !strings.Contains(Data, "||||") {
		data := strings.Join([]string{Data, "数据格式错误"}, "----")
		SendMessageToWs(m, ID, "error", data)
	} else {
		dataSplit := strings.Split(Data, "||||")
		name := dataSplit[0]
		uri := dataSplit[1]
		Version := dataSplit[2]
		CurrentPath, _ := GetCurrentPath()
		dataPath := strings.Join([]string{CurrentPath, "data"}, "/")
		tmpPath := strings.Join([]string{CurrentPath, "tmp"}, "/")
		CorePath := strings.Join([]string{CurrentPath, "Core"}, "/")
		CoreFile := strings.Join([]string{CorePath, name}, "/")
		// fmt.Println(CoreFile)
		createDirectoryIfNotExists(tmpPath)
		filePath := strings.Join([]string{tmpPath, name}, "/")
		DownloadFileWithHeaders(uri, filePath, ID, name, CoreFile, dataPath, Version, m)
	}
}
