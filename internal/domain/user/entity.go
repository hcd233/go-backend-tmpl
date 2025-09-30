// Package user 用户领域模型
//
//	@update 2025-09-30 00:00:00
package user

import (
	"errors"
	"time"
)

// Permission 权限值对象
//
//	@update 2025-09-30 00:00:00
type Permission string

const (
	// PermissionReader 读者权限
	PermissionReader Permission = "reader"
	// PermissionCreator 创建者权限
	PermissionCreator Permission = "creator"
	// PermissionAdmin 管理员权限
	PermissionAdmin Permission = "admin"
)

// PermissionLevelMapping 权限等级映射
var PermissionLevelMapping = map[Permission]int8{
	PermissionReader:  1,
	PermissionCreator: 2,
	PermissionAdmin:   3,
}

// IsValid 验证权限是否有效
func (p Permission) IsValid() bool {
	_, ok := PermissionLevelMapping[p]
	return ok
}

// Level 获取权限等级
func (p Permission) Level() int8 {
	return PermissionLevelMapping[p]
}

// HasPermission 判断是否具有指定权限
func (p Permission) HasPermission(required Permission) bool {
	return p.Level() >= required.Level()
}

// Email 邮箱值对象
//
//	@update 2025-09-30 00:00:00
type Email struct {
	value string
}

// NewEmail 创建邮箱值对象
func NewEmail(email string) (*Email, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	// 这里可以添加更复杂的邮箱验证逻辑
	return &Email{value: email}, nil
}

// Value 获取邮箱值
func (e Email) Value() string {
	return e.value
}

// User 用户实体（聚合根）
//
//	@update 2025-09-30 00:00:00
type User struct {
	id           uint
	name         string
	email        Email
	avatar       string
	permission   Permission
	lastLogin    time.Time
	createdAt    time.Time
	updatedAt    time.Time
	githubBindID string
	qqBindID     string
	googleBindID string
}

// NewUser 创建新用户
func NewUser(name string, email Email, avatar string, permission Permission) (*User, error) {
	if name == "" {
		return nil, errors.New("user name cannot be empty")
	}

	if !permission.IsValid() {
		return nil, errors.New("invalid permission")
	}

	now := time.Now().UTC()
	return &User{
		name:       name,
		email:      email,
		avatar:     avatar,
		permission: permission,
		lastLogin:  now,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// ReconstructUser 从持久化数据重建用户实体
func ReconstructUser(
	id uint,
	name string,
	email Email,
	avatar string,
	permission Permission,
	lastLogin, createdAt, updatedAt time.Time,
	githubBindID, qqBindID, googleBindID string,
) *User {
	return &User{
		id:           id,
		name:         name,
		email:        email,
		avatar:       avatar,
		permission:   permission,
		lastLogin:    lastLogin,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		githubBindID: githubBindID,
		qqBindID:     qqBindID,
		googleBindID: googleBindID,
	}
}

// ID 获取用户ID
func (u *User) ID() uint {
	return u.id
}

// Name 获取用户名
func (u *User) Name() string {
	return u.name
}

// Email 获取邮箱
func (u *User) Email() Email {
	return u.email
}

// Avatar 获取头像
func (u *User) Avatar() string {
	return u.avatar
}

// Permission 获取权限
func (u *User) Permission() Permission {
	return u.permission
}

// LastLogin 获取最后登录时间
func (u *User) LastLogin() time.Time {
	return u.lastLogin
}

// CreatedAt 获取创建时间
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt 获取更新时间
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// GithubBindID 获取Github绑定ID
func (u *User) GithubBindID() string {
	return u.githubBindID
}

// QQBindID 获取QQ绑定ID
func (u *User) QQBindID() string {
	return u.qqBindID
}

// GoogleBindID 获取Google绑定ID
func (u *User) GoogleBindID() string {
	return u.googleBindID
}

// UpdateName 更新用户名
func (u *User) UpdateName(name string) error {
	if name == "" {
		return errors.New("user name cannot be empty")
	}
	u.name = name
	u.updatedAt = time.Now().UTC()
	return nil
}

// UpdateAvatar 更新头像
func (u *User) UpdateAvatar(avatar string) {
	u.avatar = avatar
	u.updatedAt = time.Now().UTC()
}

// UpdatePermission 更新权限
func (u *User) UpdatePermission(permission Permission) error {
	if !permission.IsValid() {
		return errors.New("invalid permission")
	}
	u.permission = permission
	u.updatedAt = time.Now().UTC()
	return nil
}

// RecordLogin 记录登录
func (u *User) RecordLogin() {
	u.lastLogin = time.Now().UTC()
	u.updatedAt = time.Now().UTC()
}

// BindGithub 绑定Github
func (u *User) BindGithub(githubID string) error {
	if githubID == "" {
		return errors.New("github ID cannot be empty")
	}
	u.githubBindID = githubID
	u.updatedAt = time.Now().UTC()
	return nil
}

// BindQQ 绑定QQ
func (u *User) BindQQ(qqID string) error {
	if qqID == "" {
		return errors.New("QQ ID cannot be empty")
	}
	u.qqBindID = qqID
	u.updatedAt = time.Now().UTC()
	return nil
}

// BindGoogle 绑定Google
func (u *User) BindGoogle(googleID string) error {
	if googleID == "" {
		return errors.New("google ID cannot be empty")
	}
	u.googleBindID = googleID
	u.updatedAt = time.Now().UTC()
	return nil
}

// HasPermission 判断用户是否具有指定权限
func (u *User) HasPermission(required Permission) bool {
	return u.permission.HasPermission(required)
}

// SetID 设置用户ID（仅用于持久化后回填）
func (u *User) SetID(id uint) {
	u.id = id
}

