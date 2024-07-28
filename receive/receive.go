package receive

import (
	"EutopiaQQBot/send"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

type Sender struct {
	UserID   json.Number `json:"user_id"`
	Nickname string      `json:"nickname"`
	Card     string      `json:"card"`
	Role     string      `json:"role"`
}

type MessageData struct {
	Text string `json:"text"`
}

type Message struct {
	Data MessageData `json:"data"`
	Type string      `json:"type"`
}

type PostData struct {
	SelfID        json.Number `json:"self_id"`
	UserID        json.Number `json:"user_id"`
	Time          json.Number `json:"time"`
	MessageID     json.Number `json:"message_id"`
	MessageSeq    json.Number `json:"message_seq"`
	RealID        json.Number `json:"real_id"`
	MessageType   string      `json:"message_type"`
	Sender        Sender      `json:"sender"`
	RawMessage    string      `json:"raw_message"`
	Font          int         `json:"font"`
	SubType       string      `json:"sub_type"`
	Message       []Message   `json:"message"`
	MessageFormat string      `json:"message_format"`
	PostType      string      `json:"post_type"`
	GroupID       json.Number `json:"group_id"`
}

type JiTingRecord struct {
	num_  int
	man_  string
	date_ string
}

// 群号 -> 机厅名字 -> 记录
var QQqun map[string]map[string]JiTingRecord

func parseText(text string) (jiTing string, num int, hasNum bool, eYiShuru bool) {
	neg := 1
	flag := false
	num = 0
	for _, c := range text {
		if c == '$' {
			continue
		}
		if (c >= '0' && c <= '9') || c == '-' {
			flag = true
			if c == '-' {
				neg = -1
			}
		}
		if !flag {
			jiTing = jiTing + string(c)
		} else {
			if c >= '0' && c <= '9' {
				num = num*10 + int(c-'0')
			} else {
				if c != '-' {
					eYiShuru = true
					break
				}
			}
		}
	}
	return jiTing, num * neg, flag, eYiShuru
}

func getTime() string {
	now := time.Now()

	// 提取年、月、日、时和分
	year, month, day := now.Date()
	hour, min, second := now.Clock()

	// 打印结果
	return fmt.Sprintf("%d年%d月%d日\n%d时%d分%d秒", year, month, day, hour, min, second)
}

func checkQQGroupID(groupID string) bool {
	if groupID == "935956174" {
		return true
	} else if groupID == "821652948" {
		return true
	}

	return false
}

var data PostData
var groupID string

func GetTextFromMsg(c *gin.Context) {

	if err := c.BindJSON(&data); err != nil {
		// 如果请求体不是有效的JSON，返回一个400错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 返回一个空的JSON对象
	c.JSON(http.StatusOK, gin.H{})

	// 不接受私人消息和非文字消息
	if data.MessageType != "group" || data.Message[0].Type != "text" {
		return
	} else {
		groupID = string(data.GroupID)
	}
	// 支持的群ID
	if !checkQQGroupID(groupID) {
		return
	}
	text := data.Message[0].Data.Text
	// 复读特判
	{
		if text == "羡慕出勤" {
			send.SendText(groupID, "羡慕出勤")
			return
		}
	}
	// 判断是否是机器人指令
	if text[0] == '$' {
		// 机厅名字，人数， 是否有数字，是否恶意输入
		jiTing, num, flag, eYiShuru := parseText(text)
		// 判断一下恶意输入
		if eYiShuru {
			send.SendText(groupID, "别瞎搞，有你好果汁吃")
			return
		}
		if text == "$qa" {
			// 查询所有机厅
			QueryAllJiTing()
		} else if text == "$clearall" {
			// 清除所有机厅
			ClearAll()
		} else if text == "$clearalldata" {
			// 清除所有机厅人数
			ClearAllData()
		} else if text == "$help" {
			SendHelp()
		} else if len(text) >= 6 && text[:6] == "$clear" {
			ClearOneJiTing(text[6:])
		} else {
			// 上报人数
			ReportNum(jiTing, num, flag, eYiShuru)
		}
	}
}

func InitRoute() {
	QQqun = make(map[string]map[string]JiTingRecord)

	router := gin.Default()

	// 对 "/onebot" 路径设置GET请求的处理函数
	router.POST("/onebot", GetTextFromMsg)
	router.Run("127.0.0.1:11453")

}

/*
上报机厅人数                                    $ + 机厅名字 + 人数
上报机厅减少减少人数(负数)                  	 $ + 机厅名字 + 人数
查询全部机厅(QueryAll)                          $ + qa
清空某个机厅的数据                              $ + clear + 机厅名字
清空所有机厅数据                                $ + clearall
清空所有机厅人数(不清空机厅)                     $ + clearalldata
*/
