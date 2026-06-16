package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorServiceBusy     = errors.New("DB繁忙")
	ErrorInvalidId       = errors.New("不合法ID")
	ErrorNoRows          = errors.New("数据不存在")
)
