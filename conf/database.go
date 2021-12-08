package conf

// database DB配置
type database struct {
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Database    string `json:"database"`
	Prefix      string `json:"prefix"`
	User        string `json:"user"`
	Password    string `json:"password"`
	LogMode     bool   `json:"log_mode"`
	PoolMaxIdle int    `json:"pool_max_idle"`
	PoolMaxOpen int    `json:"pool_max_open"`
	PoolMaxTime int    `json:"pool_max_time"`
}
