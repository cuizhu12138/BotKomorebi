package receive

import (
	"EutopiaQQBot/send"
)

func SendHelp() {
	str := "获取帮助\n$ + help\n\n上报机厅人数\n$ + 机厅名字 + 人数\n\n上报机厅减少人数(负数)\n$ + 机厅名字 + 人数\n\n查询全部机厅(QueryAll)\n$ + qa\n\n清空某个机厅的数据\n$ + clear + 机厅名字\n\n清空所有机厅数据\n$ + clearall\n\n清空所有机厅人数(不清空机厅)\n$ + clearalldata"
	send.SendText(groupID, str)
}
