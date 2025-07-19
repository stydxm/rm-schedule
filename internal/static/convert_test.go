package static

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
)

func TestCompleteForm(t *testing.T) {
	var err error
	err = convertAndSaveToJSON("./complete_form.tsv", "./complete_form.json", []string{"team"})
	if err != nil {
		t.Fatal(err)
	}
	err = convertAndSaveToJSON("./complete_form_rank.tsv", "./complete_form_rank.json", []string{"team"})
	if err != nil {
		t.Fatal(err)
	}
	err = convertAndSaveToJSON("./rank_score.tsv", "./rank_score.json", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = convertAndSaveToJSON("./bilibili_official.tsv", "./bilibili_official.json", nil)
	if err != nil {
		t.Fatal(err)
	}
}

// convertAndSaveToJSON 将 TSV 文件转换为 JSON 格式并保存
// stringColumns 是需要以字符串形式存储的列名列表
func convertAndSaveToJSON(inputFile, outputFile string, stringColumns []string) error {
	// 打开 TSV 文件
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// 可以选择记录日志，这里简单打印错误
			fmt.Printf("关闭文件 %s 失败: %v\n", inputFile, closeErr)
		}
	}()

	// 创建 CSV 读取器，设置分隔符为制表符
	reader := csv.NewReader(file)
	reader.Comma = '\t'

	// 读取表头
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("读取表头失败: %v", err)
	}

	// 读取所有行
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("读取数据失败: %v", err)
	}

	// 转换 TSV 数据为 JSON
	var data []map[string]interface{}
	for _, row := range rows {
		item := make(map[string]interface{})
		for i, header := range headers {
			// 检查是否需要保持字符串格式
			isStringColumn := false
			for _, col := range stringColumns {
				if col == header {
					isStringColumn = true
					break
				}
			}
			if isStringColumn {
				item[header] = row[i]
			} else {
				// 尝试转换为数字
				item[header] = convertToNumber(row[i])
			}
		}
		data = append(data, item)
	}

	// 将数据转换为 JSON 格式
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("转换为 JSON 失败: %v", err)
	}

	// 保存到文件
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("保存 JSON 文件失败: %v", err)
	}

	return nil
}

// convertToNumber 尝试将字符串转换为数字
func convertToNumber(s string) interface{} {
	var num float64
	_, err := fmt.Sscanf(s, "%f", &num)
	if err != nil {
		return s
	}
	if math.IsNaN(num) {
		return s
	}
	return num
}
