package receive

import (
	"EutopiaQQBot/database"
	"EutopiaQQBot/send"
	_ "fmt"
	"net/http"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

var data database.PostData
var groupID string

// 解析一下有没有机厅名字，有就对应返回机厅长度
func CanGetJitingName(text string) (jitingLength int, okk bool) {
	var nowQQqun = database.QQqun[groupID]
	jitingLength = 0
	for i, _ := range text {
		if _, ok := nowQQqun[text[:i]]; ok {
			okk = true
			jitingLength = i
			break
		}
	}
	return jitingLength, okk
}

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

	//fmt.Println(text[:12])

	textLength := len(text)

	// 获取帮助
	if textLength >= 6 && text[:6] == "helpme" {
		SendHelp()
	} else if textLength >= 12 && text[:12] == "添加机厅" {
		AddJiting(text[12:])
	} else if textLength >= 4 && text[:4] == "allj"  || (text == "jtj"){
		ReportAllJiting()
	} else if (textLength >= 8 && text[:8] == "clearall"){
		ClearAllJiting()
	} else if textLength >= 5 && text[:5] == "clear" {
		Clear(text[5:])
	} else if jitingLength, ok := CanGetJitingName(text); ok { // 机厅j 和 机厅 +-x人
		if ok { // 解析出了机厅名字
			if jitingLength < len(text) {
				if text[jitingLength] == 'j' {
					QueryJiting(text[:jitingLength])
				} else if text[jitingLength] == '+' || text[jitingLength] == '-' {
					ReportJiTing(string(text[jitingLength]), text[jitingLength+1:], text[:jitingLength])
				} else if text[jitingLength] >= '0' && text[jitingLength] <= '9' {
					ReportJiTing("", text[jitingLength:], text[:jitingLength])
				}
			}
		}
	} else if text[len(text) - 1] == 'j' {
		send.SendText(groupID, "要先添加机厅才能查询机厅人数")		
	}

}

func InitRoute() {

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
