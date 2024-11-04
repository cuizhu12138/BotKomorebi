package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	JitingDB *gorm.DB
	err      error
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
	Num_  int
	Man_  string
	Date_ string
}

// 群号 -> 机厅名字 -> 记录
var QQqun map[string]map[string]JiTingRecord

func InitDatabase() {
	if JitingDB, err = gorm.Open(sqlite.Open("Jiting.db"), &gorm.Config{}); err != nil {
		fmt.Println(err)
		os.Exit(1001)
	}
	QQqun = make(map[string]map[string]JiTingRecord)

	// 自动建表
	JitingDB.AutoMigrate(&Jiting{})
	fmt.Println("自动建表成功")
	// 启动时获取所有机厅数据
	var Jitings []Jiting
	JitingDB.Find(&Jitings)

	// 存入缓存

	for _, record := range Jitings {
		if _, ok := QQqun[record.QQqun]; !ok {
			QQqun[record.QQqun] = make(map[string]JiTingRecord)
		}
		num, _ := strconv.Atoi(record.PeopleNum)
		QQqun[record.QQqun][record.JitingName] = JiTingRecord{
			Num_:  num,
			Man_:  record.ReportPeople,
			Date_: record.ReportTime,
		}

	}

}
