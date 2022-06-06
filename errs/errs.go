package errs

import "errors"

var HasExistedKeyError = errors.New("hasExistedKey")
var NotFoundKey = errors.New("notFoundKey")
var InvalidDBFormat = errors.New("节点无法创建 找不到可以存储的磁盘块")
