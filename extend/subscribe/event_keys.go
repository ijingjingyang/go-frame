package subscribe

// 定义redis pub-sub订阅消息事件常量<本质就是redis订阅推送过来的消息本身>
const (
	EventCleanAppMemCache = "clean:app:memory:cache" // 清理各pod本地App相关内存缓存
)
