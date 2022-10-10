package models

type User struct {
	Name  string `xorm:"varchar(255)" json:"name"`
	Phone string `xorm:"varchar(255)" json:"phone"`
	Sex   string `xorm:"varchar(255)" json:"sex"`
	ID    int32 `xorm:"int" json:"id"`
}
