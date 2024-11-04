package receive

import (
	"fmt"
	"time"
)

func ParseNum(text string) (num int, BadMan bool) {
	num = 0
	BadMan = false
	for _, c := range text {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		} else {
			BadMan = true
			break
		}
	}
	return num, BadMan
}

func checkQQGroupID(groupID string) bool {
	if groupID == "935956174" {
		return true
	} else if groupID == "821652948" {
		return true
	}

	return false
}

func getTime() string {
	now := time.Now()

	// 提取年、月、日、时和分
	year, month, day := now.Date()
	hour, min, second := now.Clock()

	// 打印结果
	return fmt.Sprintf("%d年%d月%d日\n%d时%d分%d秒", year, month, day, hour, min, second)
}
