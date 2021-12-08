package entry

/* 本文件内定义的struct全部均为公用的swagger文档包装 不实际使用到业务逻辑 */

// ErrorResWrapper 通用/错误响应
// swagger:response ErrorRes
type ErrorResWrapper struct {
	// in: body
	Body struct {
		BaseRes
		Data string `json:"data"`
	}
}

// NumResWrapper 数字响应返回
// swagger:response NumRes
type NumResWrapper struct {
	// in: body
	Body struct {
		BaseRes
		Data int64 `json:"data"`
	}
}

// BoolResWrapper Bool响应返回
// swagger:response BoolRes
type BoolResWrapper struct {
	// in: body
	Body struct {
		BaseRes
		Data bool `json:"data"`
	}
}

// StringResWrapper String响应返回
// swagger:response StringRes
type StringResWrapper struct {
	// in: body
	Body struct {
		BaseRes
		Data string `json:"data"`
	}
}
