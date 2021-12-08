package conf

// server 环境、端口、版本、http读写超时等server级别的配置
type server struct {
	Env          string `json:"env"`           // Env环境：dev-开发环境，stg-测试环境，uat-商业验证环境，prod-生产环境
	Version      string `json:"version"`       // 当前服务端版本
	Port         int    `json:"port"`          // http服务绑定的端口号
	ReadTimeout  int    `json:"read_timeout"`  // http服务读取请求最大时长
	WriteTimeout int    `json:"write_timeout"` // http服务响应请求最大时长
	Cors         bool   `json:"cors"`          // 是否全局添加CORS跨域header头配置，开发时无nginx代理可添加，部署时由nginx完成
}
