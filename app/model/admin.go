package model

import (
	"gorm.io/gorm"
)

// AdminRole for admin
type AdminRole string

const (
	AdminRoleSystem AdminRole = "system"
	AdminRoleNormal AdminRole = "normal"
)

func (r AdminRole) String() string {
	return string(r)
}

func (r AdminRole) Name() string {
	switch r {
	case AdminRoleSystem:
		return "システム管理者"
	case AdminRoleNormal:
		return "通常管理者"
	}
	return ""
}

func AdminRoleToArray() []KeyName {
	return []KeyName{
		{
			Key:  string(AdminRoleSystem),
			Name: AdminRoleSystem.Name(),
		},
		{
			Key:  string(AdminRoleNormal),
			Name: AdminRoleNormal.Name(),
		},
	}
}

// Admin is auth user
type Admin struct {
	gorm.Model
	Name  string `gorm:"type:varchar(191);index;not null"`
	Email string `gorm:"type:varchar(191);uniqueIndex;not null"`
	Pass  []byte
	Role  AdminRole `gorm:"type:varchar(191);index;not null;default:admin"`
}
