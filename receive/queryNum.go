package receive

import (
	"EutopiaQQBot/send"
	"fmt"
	"strings"
)

func ConstructString(origin string, length int) (returnString string) {
	returnString += origin
	for len(returnString) <= length {
		returnString += " "
	}
	return returnString
}

func QueryAllJiTing() {
	if jiTings, ok := QQqun[groupID]; ok {
		var messageSend strings.Builder
		messageSend.WriteString("以下是已经上报的机厅及人数\n")
		for jiTingName, jiTingRecord := range jiTings {
			messageSend.WriteString(ConstructString(jiTingName, 13))
			messageSend.WriteString(fmt.Sprintf("%d人", jiTingRecord.num_))
			messageSend.WriteString("\n")
		}
		send.SendText(groupID, messageSend.String())
	} else {
		send.SendText(groupID, "目前还没有人上报机厅哦=。=")
	}

}
