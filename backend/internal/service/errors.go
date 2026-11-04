package service

import "errors"

// 业务层统一错误，handler 据此映射 HTTP 状态码。
var (
	ErrNotFound = errors.New("资源不存在")
	ErrInUse    = errors.New("资源正在被使用")
	ErrConflict = errors.New("资源冲突")
	ErrBadInput = errors.New("参数不合法")
)
