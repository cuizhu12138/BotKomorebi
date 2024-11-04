package database

type Jiting struct {
	QQqun        string `gorm:"column:QQqun"`
	JitingName   string `gorm:"column:JitingName"`
	PeopleNum    string `gorm:"column:PeopleNum"`
	ReportPeople string `gorm:"column:ReportPeople"`
	ReportTime   string `gorm:"column:ReportTime"`
}

func (Jiting) TableName() string {
	return "Jiting"
}
