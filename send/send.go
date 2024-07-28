package send

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func SendText(groupId string, text string) {
	// 创建要发送的数据
	data := map[string]interface{}{
		"group_id": groupId,
		"message": []Message{
			{
				Type: "text",
				Data: map[string]interface{}{
					"text": text,
				},
			},
		},
	}

	// 将数据编码为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON encoding error: %s", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", "http://127.0.0.1:11452/send_group_msg", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}

	// 设置请求头，指定内容类型为 JSON
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
	}

	// 关闭响应体，防止资源泄露
	defer resp.Body.Close()

	// 打印响应状态
	log.Println("Response status:", resp.Status)
}
