// Package model defines the database schema for the model.
//
//	update 2024-06-22 09:33:43
package model

import (
	"time"
)

type (
	// Permission string 权限
	//	update 2024-09-21 01:34:29
	Permission string

	// PermissionLevel int8 权限等级
	//	update 2024-09-21 01:34:29
	PermissionLevel int8

	// Quota int8 配额
	//	update 2024-12-09 16:13:24
	Quota int8

	// Platform string 平台
	//	update 2024-09-21 01:34:12
	Platform string
)

const (

	// PlatformGithub github user
	//	update 2024-06-22 10:05:13
	PlatformGithub Platform = "github"

	// PlatformGoogle google user
	PlatformGoogle Platform = "google"

	// PermissionReader general permission
	//	update 2024-06-22 10:05:15
	PermissionReader Permission = "reader"

	// PermissionCreator creator permission
	//	update 2024-06-22 10:05:17
	PermissionCreator Permission = "creator"

	// PermissionAdmin admin permission
	//	update 2024-06-22 10:05:17
	PermissionAdmin Permission = "admin"
)

// PermissionLevelMapping 权限等级映射
//
//	update 2024-09-21 01:34:29
var (
	PermissionLevelMapping = map[Permission]int8{
		PermissionReader:  1,
		PermissionCreator: 2,
		PermissionAdmin:   3,
	}
)

// User 用户数据库模型
//
//	author centonhuang
//	update 2024-06-22 09:36:22
type User struct {
	BaseModel
	Name         string     `json:"name" gorm:"column:name;unique;not null;comment:用户名"`
	Email        string     `json:"email" gorm:"column:email;unique;not null;comment:邮箱"`
	Avatar       string     `json:"avatar" gorm:"column:avatar;not null;comment:头像"`
	Permission   Permission `json:"permission" gorm:"column:permission;not null;default:'reader';comment:权限"`
	LastLogin    time.Time  `json:"last_login" gorm:"column:last_login;comment:最后登录时间"`
	GithubBindID string     `json:"-" gorm:"unique;comment:Github绑定ID"`
	QQBindID     string     `json:"-" gorm:"unique;comment:QQ绑定ID"`
	GoogleBindID string     `json:"-" gorm:"unique;comment:Google绑定ID"`
}
