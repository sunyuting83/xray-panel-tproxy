package utils

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	config "xpanel/Config"
)

func AddManualNode(p string, newNode *config.CodeList) error {
	// 1. 定位手动节点文件
	manualFile := filepath.Join(p, "data", "manual.json")

	// 2. 读取现有数据
	var manualNodes []*config.CodeList
	data, err := os.ReadFile(manualFile)
	if err == nil && len(data) > 0 {
		// 处理可能存在的 null 截断
		if index := bytes.IndexByte(data, 0); index != -1 {
			data = data[:index]
		}
		json.Unmarshal(data, &manualNodes)
	}

	// 3. 核心清洗与增强
	// 生成唯一的 UID（复用你之前的 Hash 算法）
	newNode.UID = GenerateUID(newNode)

	// 针对 Shadowsocks 的“去壳”处理，防止 Xray 报错
	if newNode.Type == "shadowsocks" || newNode.Type == "ss" {
		newNode.Network = ""
		newNode.StreamSecurity = ""
		newNode.Path = ""
		newNode.Host = ""
	}

	// 4. 追加并写回
	manualNodes = append(manualNodes, newNode)
	saveConfig, err := json.MarshalIndent(manualNodes, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(manualFile, saveConfig, 0644)
}
