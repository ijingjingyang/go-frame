package render

// 错误码定义：利用go工具链 generate 自动生成错误码映射map
//  -- 定义CE类型并通过注释给出对应文案
const (
	SuccessCode     = 0
	UnknownError CE = -1 // 未知错误
	// 基础服务错误
	DbError             CE = 100001 // 系统错误
	RedisError          CE = 100002 // 系统错误
	LostConnectionError CE = 100003 // 系统错误

	// 系统/服务错误
	InvalidParams    CE = 101001 // 请求参数错误
	InvalidRequest   CE = 101002 // 请求错误
	SystemBusy       CE = 101003 // 系统繁忙请稍后再试
	DbRecordNotExist CE = 101004 // 数据不存在
	ErrDefineWithMsg CE = 101005 // %s
)
