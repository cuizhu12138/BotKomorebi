package receive

import (
	"EutopiaQQBot/send"
	"fmt"
)

func ClearAll() {
	if _, ok := QQqun[groupID]; ok {
		QQqun[groupID] = make(map[string]JiTingRecord)
		send.SendText(groupID, "成功清空所有数据")
	}
}

func ClearAllData() {
	if jiTings, ok := QQqun[groupID]; ok {
		for key := range jiTings {
			jiTings[key] = JiTingRecord{}
		}
		send.SendText(groupID, "成功清空所有人数")
	}
}

func ClearOneJiTing(jiTing string) {
	if jiTings, ok := QQqun[groupID]; ok {
		if _, ok2 := jiTings[jiTing]; ok2 {
			jiTings[jiTing] = JiTingRecord{}
			send.SendText(groupID, fmt.Sprintf("成功清空%s的人数", jiTing))
		} else {
			send.SendText(groupID, "没有这样的机厅哦")
		}
	}
}
