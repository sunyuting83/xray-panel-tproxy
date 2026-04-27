package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	config "xpanel/Config"
)

// GenerateConfig 终极版：全量模板化，无硬编码拼接
func GenerateConfig(currentPath, action string) error {
	// 1. 获取当前节点
	idx := GetCurrentUID(currentPath)
	// fmt.Println(idx)
	node, err := GetNodeByUID(currentPath, idx)
	if err != nil {
		return err
	}

	// 2. 渲染单个协议零件 (如 vless.tmpl)
	proxyJson, err := renderOutbound(currentPath, node)
	if err != nil {
		return err
	}

	// 3. 渲染 outbounds.tmpl (将协议零件嵌入数组)
	outboundsArray, err := renderOutboundsArray(currentPath, proxyJson)
	if err != nil {
		return err
	}

	// 4. 准备全量模板数据
	data := map[string]interface{}{
		"DNS":               readTmplContent(currentPath, "dns_servers.tmpl", ""),
		"Outbounds":         outboundsArray,
		"ProxyDomainRules":  readTmplContent(currentPath, "rule_proxy_domain.tmpl", ""),
		"ProxyIPRules":      readTmplContent(currentPath, "rule_proxy_ip.tmpl", ""),
		"DirectDomainRules": readTmplContent(currentPath, "rule_direct_domain.tmpl", ""),
		"DirectIPRules":     readTmplContent(currentPath, "rule_direct_ip.tmpl", ""),
		"BlockDomainRules":  readTmplContent(currentPath, "rule_block_domain.tmpl", ""),
		"DirectAppRules":    readTmplContent(currentPath, "rule_direct_app.tmpl", ""),
	}

	// 5. 渲染主模板 main.tmpl
	mainTmplRaw, err := os.ReadFile(filepath.Join(currentPath, "template", "main.tmpl"))
	if err != nil {
		return err
	}

	t, err := template.New("main").Parse(string(mainTmplRaw))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	// 6. 处理 JSON 冗余 (如果零件为空，可能会产生 ,, )
	finalJson := strings.ReplaceAll(buf.String(), ",,", ",")
	finalJson = strings.ReplaceAll(finalJson, ",]", "]")

	// 7. 写入并重启
	configPath := filepath.Join(currentPath, "Core", "config.json")
	os.WriteFile(configPath, []byte(finalJson), 0644)

	RunXray(currentPath, action, node.Title)
	return nil
}

// renderOutboundsArray：直接将 proxyOutbound 填入 outbounds.tmpl
func renderOutboundsArray(currentPath string, proxyOutbound string) (string, error) {
	path := filepath.Join(currentPath, "template", "outbounds.tmpl")
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取 outbounds.tmpl 失败: %v", err)
	}

	data := struct {
		ProxyOutbound string
	}{
		ProxyOutbound: proxyOutbound,
	}

	// 解析并渲染
	tmpl, err := template.New("outbounds").Parse(string(raw))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// readTmplContent：确保返回的是合法的 JSON 片段
func readTmplContent(currentPath, fileName, defaultVal string) string {
	path := filepath.Join(currentPath, "template", fileName)
	content, err := os.ReadFile(path)

	// 如果文件不存在或内容为空
	if err != nil || len(strings.TrimSpace(string(content))) == 0 {
		if defaultVal != "" {
			return defaultVal
		}
		// 关键：由于 main.tmpl 里的 rules 数组是用逗号连接的
		// 如果零件为空，返回一个“不产生影响”的 JSON 对象
		return "{\"type\":\"field\",\"outboundTag\":\"outBound_DIRECT\",\"domain\":[\"none:none\"]}"
	}
	return string(content)
}

// --- 数据读取核心函数 ---

// 内部辅助函数：渲染具体的协议零件
func renderOutbound(currentPath string, node *config.CodeList) (string, error) {
	tmplPath := filepath.Join(currentPath, "template", "outbounds", node.Type+".tmpl")
	tmplRaw, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("找不到协议模板: %s", tmplPath)
	}

	tmpl, err := template.New("node").Parse(string(tmplRaw))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, node); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetCurrentUID 从 settings.json 获取当前选中的节点 UID
func GetCurrentUID(currentPath string) string {
	settingsPath := filepath.Join(currentPath, "data", "settings.json")

	// 如果文件不存在，默认返回空字符串
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return ""
	}

	var settings struct {
		CurrentUID string `json:"current_uid"`
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		return ""
	}

	return settings.CurrentUID
}

// GetNodeByUID 从 data.json 中根据 UID 查找节点
func GetNodeByUID(currentPath string, uid string) (*config.CodeList, error) {
	// 定义两个数据源路径
	dataSources := []string{
		filepath.Join(currentPath, "data", "data.json"),
		filepath.Join(currentPath, "data", "manual.json"),
	}

	for _, path := range dataSources {
		// 检查文件是否存在，不存在则跳过（防止初次运行没有 manual.json 报错）
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			// 这里记录日志但不中断，尝试下一个文件
			continue
		}

		var nodeList []*config.CodeList
		if err := json.Unmarshal(data, &nodeList); err != nil {
			continue
		}

		// 遍历查找匹配的 UID
		for _, node := range nodeList {
			if node.UID == uid {
				return node, nil
			}
		}
	}

	return nil, fmt.Errorf("在订阅和手动节点中均未找到 UID 为 %s 的节点", uid)
}

// --- 辅助写入函数 (给 SelectNode 使用) ---

// SaveCurrentUID 保存当前选中的 UID 到 settings.json
func SaveCurrentUID(currentPath string, uid string) error {
	settingsPath := filepath.Join(currentPath, "data", "settings.json")

	settings := struct {
		CurrentUID string `json:"current_uid"`
	}{
		CurrentUID: uid,
	}

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}

// --- 业务对接函数 ---

// SelectNode 供 router 调用的节点切换函数
func SelectNode(currentPath string, uid string) bool {
	// 1. 保存选择
	if err := SaveCurrentUID(currentPath, uid); err != nil {
		return false
	}

	// 2. 触发合成引擎 (昨天讨论的核心函数)
	if err := GenerateConfig(currentPath, "reload"); err != nil {
		fmt.Println("合成配置失败:", err)
		return false
	}

	return true
}

// GetNodeList 从 data.json 获取完整的节点列表
func GetNodeList(currentPath string) ([]config.CodeList, error) {
	// 1. 定义两个文件的路径
	subPath := filepath.Join(currentPath, "data", "data.json")
	manualPath := filepath.Join(currentPath, "data", "manual.json")

	var allNodes []config.CodeList

	// 3. 读取手动节点 (manualFile)
	manualData, err := os.ReadFile(manualPath)
	if err == nil && len(manualData) > 0 {
		var manualList []config.CodeList
		if err := json.Unmarshal(manualData, &manualList); err == nil {
			allNodes = append(allNodes, manualList...)
		}
	}

	// 2. 读取订阅节点 (dataFile)
	subData, err := os.ReadFile(subPath)
	if err == nil && len(subData) > 0 {
		var subList []config.CodeList
		if err := json.Unmarshal(subData, &subList); err == nil {
			allNodes = append(allNodes, subList...)
		}
	}

	// 4. 重新校准 Index
	// 因为两个文件合并后，原始记录里的 Index 可能会重复或断层
	// 统一按合并后的物理顺序重新赋值，方便前端表格展示
	for i := range allNodes {
		allNodes[i].Index = i
	}

	return allNodes, nil
}

// 定义所有受支持的规则零件名
var SupportedRuleTypes = []string{
	"block_domain",
	"direct_app",
	"direct_domain",
	"direct_ip",
	"proxy_domain",
}

// GetRules 一次性读取所有路由规则
func GetRules(currentPath string) (map[string]string, error) {
	results := make(map[string]string)

	fieldMap := map[string]string{
		"block_domain":  "domain",
		"direct_app":    "process",
		"direct_domain": "domain",
		"direct_ip":     "ip",
		"proxy_domain":  "domain",
	}

	for _, rt := range SupportedRuleTypes {
		fileName := fmt.Sprintf("rule_%s.tmpl", rt)
		path := filepath.Join(currentPath, "template", fileName)

		content, err := os.ReadFile(path)
		if err != nil {
			results[rt] = "" // 文件不存在返回空字符串
			continue
		}

		var ruleData map[string]any
		if err := json.Unmarshal(content, &ruleData); err != nil {
			results[rt] = ""
			continue
		}

		targetKey := fieldMap[rt]
		if val, ok := ruleData[targetKey]; ok {
			// 将接口类型断言为切片
			if list, ok := val.([]any); ok {
				var cleanList []string
				for _, item := range list {
					str := fmt.Sprintf("%v", item)

					// 针对域名类的规则进行特殊脱壳处理
					if rt == "block_domain" || rt == "direct_domain" || rt == "proxy_domain" {
						// 去掉 domain: 前缀
						str = strings.TrimPrefix(str, "domain:")
						// 去掉 geosite: 前缀 (如 geosite:google -> google)
						str = strings.TrimPrefix(str, "geosite:")
					}

					cleanList = append(cleanList, str)
				}
				// 用换行符连接，方便前端 textarea 直接显示
				results[rt] = strings.Join(cleanList, "\n")
			}
		} else {
			results[rt] = ""
		}
	}

	return results, nil
}

// UpdateRule 更新单条规则零件
func UpdateRule(currentPath string, ruleType string, rawContent string) error {
	// 1. 定义字段映射
	fieldMap := map[string]string{
		"block_domain":  "domain",
		"direct_app":    "process",
		"direct_domain": "domain",
		"direct_ip":     "ip",
		"proxy_domain":  "domain",
	}

	targetKey, ok := fieldMap[ruleType]
	if !ok {
		return fmt.Errorf("不支持的规则类型: %s", ruleType)
	}

	// 2. 读取原有的 tmpl 文件内容
	fileName := fmt.Sprintf("rule_%s.tmpl", ruleType)
	filePath := filepath.Join(currentPath, "template", fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取模板失败: %v", err)
	}

	// 3. 解析原有 JSON 结构
	var ruleData map[string]any
	if err := json.Unmarshal(content, &ruleData); err != nil {
		return fmt.Errorf("解析模板失败 (请检查 JSON 格式): %v", err)
	}

	// 4. 处理前端传来的新数据
	rawLines := strings.Split(strings.ReplaceAll(rawContent, "\r\n", "\n"), "\n")
	var cleanList []string

	for _, line := range rawLines {
		item := strings.TrimSpace(line)
		if item == "" {
			continue
		}

		// 域名类前缀自动包装逻辑
		if ruleType == "block_domain" || ruleType == "direct_domain" || ruleType == "proxy_domain" {
			item = strings.TrimPrefix(item, "domain:")
			item = strings.TrimPrefix(item, "geosite:")
			if strings.Contains(item, ".") {
				item = "domain:" + item
			} else {
				item = "geosite:" + item
			}
		}
		cleanList = append(cleanList, item)
	}

	// 5. 【关键修复】只覆盖目标字段，保留 outboundTag 和 type
	ruleData[targetKey] = cleanList

	// 6. 序列化回 JSON 并存回文件
	finalJSON, err := json.Marshal(ruleData)
	if err != nil {
		return fmt.Errorf("JSON 编码失败: %v", err)
	}

	if err := os.WriteFile(filePath, finalJSON, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	// 7. 联动合成配置
	return GenerateConfig(currentPath, "reload")
}

// SetNodeAndReload 核心切换函数：记录 Index -> 渲染配置 -> 重启服务
func SetNodeAndReload(currentPath string, uid string) bool {
	// 1. 校验：先看 data.json 里有没有这个节点，防止存入一个不存在的 Index
	_, err := GetNodeByUID(currentPath, uid)
	if err != nil {
		fmt.Printf("切换失败，节点不存在: %v\n", err)
		return false
	}

	// 2. 存储：将当前选中的 Index 写入 settings.json
	// 这是为了下次启动或刷新时，系统知道该用哪一个
	if err := SaveCurrentUID(currentPath, uid); err != nil {
		fmt.Printf("保存设置失败: %v\n", err)
		return false
	}

	// 3. 合成：调用我们之前写好的终极引擎
	// 它会读取 settings.json 里的 index，找到节点，渲染模板，生成 config.json，最后 RunXray
	if err := GenerateConfig(currentPath, "reload"); err != nil {
		fmt.Printf("配置合成失败: %v\n", err)
		return false
	}

	return true
}

func GenerateUID(node *config.CodeList) string {
	// 将核心参数拼接
	raw := fmt.Sprintf("%s:%d:%s:%s", node.Address, node.Port, node.ID, node.Password)
	// 使用 MD5 或简单的 Hash
	hasher := md5.New()
	hasher.Write([]byte(raw))
	return hex.EncodeToString(hasher.Sum(nil))[:12] // 取前12位即可
}
