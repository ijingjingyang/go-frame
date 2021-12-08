package render

import "fmt"

// see golang.org/x/tools/cmd/stringer
//go:generate $GOPATH/bin/stringer -type CE -linecomment  -output string.go

// CE 简单错误码错误
type CE int

// Error 错误信息
func (e CE) Error() string {
	return e.String()
}

// Wrap 包裹一个错误 返回 E
func (e CE) Wrap(err error, args ...interface{}) E {
	return E{err: err, CE: e, args: args}
}

// Code 错误码
func (e CE) Code() int {
	return int(e)
}

// E 组合错误
type E struct {
	CE
	args []interface{} // CE 格式化需要的参数
	err  error
}

// Format 格式化错误输出给客户端
func (e E) Format() string {
	ceMsg := e.CE.Error()
	if len(e.args) > 0 {
		ceMsg = fmt.Sprintf(ceMsg, e.args...)
	}
	return ceMsg
}

// Error 错误信息--一般用户日志记录 开发人员使用
func (e E) Error() string {
	ceMsg := e.CE.Error()
	if len(e.args) > 0 {
		ceMsg = fmt.Sprintf(ceMsg, e.args...)
	}
	if e.err == nil {
		return ceMsg
	}
	return fmt.Sprintf("%s (%s)", ceMsg, e.err.Error())
}

// Unwrap 解包
func (e E) Unwrap() error {
	return e.err
}
