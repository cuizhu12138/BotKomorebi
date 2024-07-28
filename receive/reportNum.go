package receive

import (
	"EutopiaQQBot/send"
	"fmt"
)

func ReportNum(jiTing string, num int, flag bool, eYiShuru bool) {
	var jiTings map[string]JiTingRecord
	var ok bool
	// 没有这样的qq群需要创建一下map
	if jiTings, ok = QQqun[groupID]; !ok {
		QQqun[groupID] = make(map[string]JiTingRecord)
		jiTings = QQqun[groupID]
	}
	// 恶意输入
	if num > 1000000 {
		send.SendText(groupID, "别瞎搞，有你好果汁吃")
		return
	}
	// 恶意输入
	if jiTing == "" || (len(jiTing) == 1 && !(jiTing[0] >= 'a' && jiTing[0] <= 'z' || jiTing[0] <= 'A' && jiTing[0] >= 'Z')) {
		send.SendText(groupID, "你在干嘛。。。。？？")
		return
	}
	if flag {
		if num >= 0 {
			jiTings[jiTing] = JiTingRecord{
				num_:  num,
				man_:  data.Sender.Nickname,
				date_: getTime(),
			}
		} else {
			if jiTings[jiTing].num_ + num < 0 {
				var str = fmt.Sprintf("%s只有%d人哦", jiTing, jiTings[jiTing].num_)
				send.SendText(groupID, str)
				return
			} else {
				jiTings[jiTing] = JiTingRecord{
					num_:  jiTings[jiTing].num_ + num,
					man_:  data.Sender.Nickname,
					date_: getTime(),
				}
			}
		}
		send.SendText(groupID, fmt.Sprintf("已收到消息%s有%d人", jiTing, jiTings[jiTing].num_))
	} else {
		if jiTings[jiTing].num_ != 0 {
			var str = fmt.Sprintf("%s有%d人\n由%s在%s报告", jiTing, jiTings[jiTing].num_, jiTings[jiTing].man_, jiTings[jiTing].date_)
			send.SendText(groupID, str)
		} else {
			var str = fmt.Sprintf("%s有%d人", jiTing, jiTings[jiTing].num_)
			send.SendText(groupID, str)
		}
	}
}
