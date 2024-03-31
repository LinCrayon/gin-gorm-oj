package models

import "gorm.io/gorm"

type TestCase struct {
	gorm.Model
	Identity        string `gorm:"column:identity;type:varchar(36);" json:"identity"`                 //用户表的唯一标识
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"` //用户名
	Input           string `gorm:"column:input;type:text;" json:"input"`                              //密码
	Output          string `gorm:"column:output;type:text;" json:"output"`
}

func (table *TestCase) TableName() string {
	return "test_case"
}
