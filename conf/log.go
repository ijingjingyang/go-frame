package conf

// log 日志相关配置
type log struct {
	Level string `json:"level"` // 日志级别配置：panic|fatal|error|warning|info|debug|trace
	Path  string `json:"path"`  // 日志存储路径配置：stderr|stdout|file_path_only_dir
}
