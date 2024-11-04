package receive

import (
	"EutopiaQQBot/database"
	"EutopiaQQBot/send"
	"fmt"
	"strconv"
	"strings"
)

func ConstructString(origin string, length int) (returnString string) {
	returnString += origin
	for len(returnString) <= length {
		returnString += " "
	}
	return returnString
}

// 新增机厅
func AddJiting(Jiting string) {
	var jiTings map[string]database.JiTingRecord
	var ok bool
	// 没有这样的qq群需要创建一下map
	if jiTings, ok = database.QQqun[groupID]; !ok {
		database.QQqun[groupID] = make(map[string]database.JiTingRecord)
		jiTings = database.QQqun[groupID]
	}

	if Jiting == "" || (len(Jiting) == 1 && !(Jiting[0] >= 'a' && Jiting[0] <= 'z' || Jiting[0] <= 'A' && Jiting[0] >= 'Z')) {
		send.SendText(groupID, "你在干嘛。。。。？？")
		return
	}

	newJiting := database.Jiting{}
	// 检查一下有没有这个机厅
	if _, ok := jiTings[Jiting]; ok {
		send.SendText(groupID, "已经有这个机厅啦")
		return
	}

	database.JitingDB.Where("QQqun = ? and JitingName = ?", groupID, Jiting).Take(&newJiting)

	if newJiting.QQqun != "" {
		send.SendText(groupID, "已经有这个机厅啦")
		return
	}

	// 先在缓存中加入
	jiTings[Jiting] = database.JiTingRecord{}

	// 塞入数据库中
	newJiting = database.Jiting{
		QQqun:      groupID,
		JitingName: Jiting,
	}

	// 加入数据库
	if err := database.JitingDB.Create(&newJiting).Error; err != nil {
		send.SendText(groupID, "新增机厅错误，请联系管理员查看日志")
		fmt.Printf("创建失败，错误信息：%v\n", err)
	} else {
		send.SendText(groupID, "新增机厅成功！出勤出勤！")
	}
}

// 通知所有已经上报的机厅人数
func ReportAllJiting() {
	if jiTings, ok := database.QQqun[groupID]; ok {
		var messageSend strings.Builder
		messageSend.WriteString("以下是已经上报的机厅及人数\n")
		for jiTingName, jiTingRecord := range jiTings {
			messageSend.WriteString(ConstructString(jiTingName, 13))
			messageSend.WriteString(fmt.Sprintf("%d人", jiTingRecord.Num_))
			messageSend.WriteString("\n")
		}
		send.SendText(groupID, messageSend.String())
	} else {
		send.SendText(groupID, "目前还没有人上报机厅哦=。=")
	}
}

// 清除所有机厅
func ClearAllJiting() {
	if _, ok := database.QQqun[groupID]; ok {
		database.QQqun[groupID] = make(map[string]database.JiTingRecord)
		send.SendText(groupID, "成功清空所有数据")
	}
	database.JitingDB.Where("1 = 1").Delete(database.Jiting{})
}

// 清除某个机厅的人数
func Clear(Jiting string) {
	if jiTings, ok := database.QQqun[groupID]; ok {
		if _, ok2 := jiTings[Jiting]; ok2 {
			jiTings[Jiting] = database.JiTingRecord{}
			send.SendText(groupID, fmt.Sprintf("成功清空%s的人数", Jiting))
		} else {
			send.SendText(groupID, "没有这样的机厅哦")
		}
	}
	database.JitingDB.Where("QQqun = ? and JitingName = ?", groupID, Jiting).Delete(&database.Jiting{})
}

// 查询某个机厅人数
func QueryJiting(jiTing string) {

	var jiTings map[string]database.JiTingRecord
	var ok bool
	// 没有这样的qq群需要创建一下map
	if jiTings, ok = database.QQqun[groupID]; !ok {
		database.QQqun[groupID] = make(map[string]database.JiTingRecord)
		jiTings = database.QQqun[groupID]
	}

	if jiTings[jiTing].Num_ != 0 {
		var str = fmt.Sprintf("%s有%d人\n由%s在%s报告", jiTing, jiTings[jiTing].Num_, jiTings[jiTing].Man_, jiTings[jiTing].Date_)
		send.SendText(groupID, str)
	} else {
		var str = fmt.Sprintf("%s有%d人", jiTing, jiTings[jiTing].Num_)
		send.SendText(groupID, str)
	}
}

// 上报某个机厅人数
func ReportJiTing(Signal string, NumText string, jiTing string) {
	num, BadMan := ParseNum(NumText)

	var jiTings map[string]database.JiTingRecord
	var ok bool
	// 没有这样的qq群需要创建一下map
	if jiTings, ok = database.QQqun[groupID]; !ok {
		database.QQqun[groupID] = make(map[string]database.JiTingRecord)
		jiTings = database.QQqun[groupID]
	}
	// 恶意输入
	if num > 1000000 {
		send.SendText(groupID, "别瞎搞，有你好果汁吃")
		return
	}

	// 数字不对
	if BadMan {
		send.SendText(groupID, "这是人数？")
		return
	}

	// 恶意输入
	if jiTing == "" || (len(jiTing) == 1 && !(jiTing[0] >= 'a' && jiTing[0] <= 'z' || jiTing[0] <= 'A' && jiTing[0] >= 'Z')) {
		send.SendText(groupID, "你在干嘛。。。。？？")
		return
	}

	if Signal == string('+') {
		jiTings[jiTing] = database.JiTingRecord{
			Num_:  jiTings[jiTing].Num_ + num,
			Man_:  data.Sender.Nickname,
			Date_: getTime(),
		}
	} else if Signal == string('-') {
		fmt.Println("负数判断：", jiTings[jiTing].Num_)
		fmt.Println("加后为：", jiTings[jiTing].Num_+num)
		if jiTings[jiTing].Num_-num < 0 {
			var str = fmt.Sprintf("%s只有%d人哦", jiTing, jiTings[jiTing].Num_)
			send.SendText(groupID, str)
			return
		} else {
			jiTings[jiTing] = database.JiTingRecord{
				Num_:  jiTings[jiTing].Num_ - num,
				Man_:  data.Sender.Nickname,
				Date_: getTime(),
			}
		}
	} else {
		jiTings[jiTing] = database.JiTingRecord{
			Num_:  num,
			Man_:  data.Sender.Nickname,
			Date_: getTime(),
		}
	}
	send.SendText(groupID, fmt.Sprintf("已收到消息%s有%d人", jiTing, jiTings[jiTing].Num_))
	fmt.Println(jiTings[jiTing].Num_)
	database.JitingDB.Model(&database.Jiting{}).Where("QQqun = ? and JitingName = ?", groupID, jiTing).Updates(&database.Jiting{
		QQqun:        groupID,
		JitingName:   jiTing,
		PeopleNum:    strconv.Itoa(jiTings[jiTing].Num_),
		ReportPeople: data.Sender.Nickname,
		ReportTime:   getTime(),
	})
}

func SendHelp() {
	var messageSend strings.Builder
	messageSend.WriteString("以下是bot使用方法\n\n")
	messageSend.WriteString("最新的使用流程需要先添加机厅才能查询机厅人数！！\n\n")
	messageSend.WriteString("获取帮助：\nhelpme\n\n")
	messageSend.WriteString("添加机厅：\n添加机厅 + 机厅名字\n\n")
	messageSend.WriteString("查询机厅：\n机厅名字 + j\n（eg.mmj可以查询mm的人数）\n\n")
	messageSend.WriteString("查询全部机厅：\nallj\n\n")
	messageSend.WriteString("上报人数：\n机厅名字 + 人数\n（ps.支持类似\"mm+1\",\"mm-1\",\"mm4\"之类的用法）\n\n")
	messageSend.WriteString("清空所有机厅数据：\nclearall\n（慎用，需要重新添加机厅）\n\n")
	messageSend.WriteString("清空某个机厅人数：\nclear + 机厅名字\n\n")
	send.SendText(groupID, messageSend.String())
}
