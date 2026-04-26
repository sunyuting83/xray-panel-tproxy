package router

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"xpanel/utils"

	"github.com/gin-gonic/gin"
)

// upData updata
func UpData(c *gin.Context) {
	// 1. 获取基础路径与参数
	currentPathRaw, _ := c.Get("current_path")
	currentPath := currentPathRaw.(string)
	proxyFlag := c.DefaultQuery("proxy", "0") != "0"

	// 2. 拼接文件路径 (统一到新架构路径)
	ignoreFile := filepath.Join(currentPath, "data", "ignore")
	subUrlFile := filepath.Join(currentPath, "data", "subUrl")
	dataFile := filepath.Join(currentPath, "data", "data.json") // 这里的路径要和你 utils 读的地方一致

	// 3. 读取配置
	ignore, _ := os.ReadFile(ignoreFile)
	subUrl, _ := os.ReadFile(subUrlFile)

	// 4. 获取订阅数据
	urlList := utils.GetSubUrl(subUrl, currentPath)
	rawNodes := utils.SyncGetData(urlList, proxyFlag)

	if len(rawNodes) > 0 {
		// A. 转换节点为 CodeList 结构体 (这里内部要确保 types -> type 的转换)
		base := utils.MakeDates(rawNodes)
		// fmt.Println(base)
		// B. 过滤与去重
		datas := utils.IgnoreTag(base, string(ignore))
		datas = utils.RemoveRepeatedElement(datas)

		if len(datas) > 0 {
			// C. 关键点：重新分配 Index，确保与 data.json 物理位置对齐
			for i := range datas {
				datas[i].UID = utils.GenerateUID(datas[i])
				datas[i].Index = i
			}
			// fmt.Println(datas)
			// D. 序列化并保存 (Marshal 会自动根据你 CodeList 的 tag 转换字段名)
			saveData, err := json.MarshalIndent(datas, "", "  ")
			if err == nil {
				os.WriteFile(dataFile, saveData, 0644)
			}
		}

		c.JSON(200, gin.H{
			"status":  0,
			"date":    datas,
			"message": fmt.Sprintf("成功同步 %d 个节点", len(datas)),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  1,
		"date":    []string{},
		"message": "未能获取到有效节点数据",
	})
}
