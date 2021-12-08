package conf

// redis redis config
type redis struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Password        string `json:"password"`
	Database        int    `json:"database"`
	ConnectTimeout  int    `json:"connect_timeout"`
	ReadTimeout     int    `json:"read_timeout"`
	WriteTimeout    int    `json:"write_timeout"`
	PoolMaxIdle     int    `json:"pool_max_idle"`
	PoolMaxOpen     int    `json:"pool_max_open"`
	PoolMaxTime     int    `json:"pool_max_time"`
	QueueConcurrent int64  `json:"queue_concurrent"`
	QueueWarnNum    int64  `json:"queue_warn_num"`
}
