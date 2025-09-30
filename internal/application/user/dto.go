// Package user 用户应用层DTO
//
//	@update 2025-09-30 00:00:00
package user

import "time"

// GetUserInfoQuery 获取用户信息查询
type GetUserInfoQuery struct {
	UserID uint
}

// GetCurUserInfoQuery 获取当前用户信息查询
type GetCurUserInfoQuery struct {
	UserID uint
}

// UpdateUserInfoCommand 更新用户信息命令
type UpdateUserInfoCommand struct {
	UserID      uint
	UpdatedName string
}

// UserDTO 用户数据传输对象
type UserDTO struct {
	UserID    uint   `json:"userID"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
	LastLogin string `json:"lastLogin"`
}

// CurUserDTO 当前用户数据传输对象（包含权限信息）
type CurUserDTO struct {
	UserID     uint   `json:"userID"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	CreatedAt  string `json:"createdAt"`
	LastLogin  string `json:"lastLogin"`
	Permission string `json:"permission"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	User *UserDTO `json:"user"`
}

// CurUserInfoResponse 当前用户信息响应
type CurUserInfoResponse struct {
	User *CurUserDTO `json:"user"`
}

// UpdateUserInfoResponse 更新用户信息响应
type UpdateUserInfoResponse struct {
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

