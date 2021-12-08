# go-frame

日常开发过程中提炼封装的go语言开发骨架，**骨架非框架**！

## 一、项目描述

使用了 `gin`、`gorm`、`go-redis`等开源组件构造了一套自定义后端api开发骨架并尽量搭建好单元测试结构。

## 二、项目结构

项目核心目录树和介绍如下：

````
├── app                App应用目录
│   ├── console        终端控制台应用汇总
│   │   ├── command    终端命令应用
│   │   └── job        队列任务类job
│   ├── controller     控制器
│   ├── entry          实体类
│   ├── model          model模型
│   └── service        service服务类
├── client             全局句柄客户端
├── conf               全局配置
├── define             全局常量变量
├── extend             扩展功能集
├── migrations         数据库迁移文件
├── render             渲染器
├── route              路由
├── utils              utils帮助类
├── conf.toml.example  配置文件样例
└── main.go            main包入口
````
## 三、开发相关

### 1、环境准备

* 安装go最新稳定版本（>=1.16），设置好`GOROOT`、`GOPATH`环境变量，确保`go`命令可用
* 安装工具`stringer`、`swagger`、`golint`

````
# 1.进入`GOPATH`目录
cd `go env GOPATH`

# 2.安装工具

# go<1.16
go get -u -v golang.org/x/tools/cmd/stringer
go get -u -v github.com/go-swagger/go-swagger/cmd/swagger
go get -u -v golang.org/x/lint/golint

# go>=1.16
go install golang.org/x/tools/cmd/stringer@latest
go install github.com/go-swagger/go-swagger/cmd/swagger
go install golang.org/x/lint/golint@latest

# 3、确认安装的工具可用
# 看到有对应二进制文件即表示安装成功
ls `go env GOPATH`/bin/
````
### 2、变量及命名规范

* 变量：驼峰
* 常量：驼峰
* 导出的变量、struct、func 需要加上注释，见下方示例：

````
// DefaultAreaService 默认区域服务单例
var DefaultAreaService AreaService

// User 用户表模型
type User struct {
   // code
}

// GetUserByID 根据ID获取用户模型
func (s UserService) GetUserByID(id int) (int, error) {
   // code
}
````
* 使用`golint`工具检查是否规范：项目根目录下执行：`golint ./...`

### 3、swagger文档生成

* swagger文档JSON文件生成
````
swagger generate spec -m  -o ./runtime/sf.json
````
* 文档本地预览
````
swagger serve --host=0.0.0.0 --port=20218 ./runtime/sf.json
````
> 执行命令后将自动打开浏览器显示文档界面，如果不想自动打开加上 `--no-open`参数

> 在项目根目录执行`./scripts/swagger.sh` shell脚本可一键生成并打开浏览器显示swagger文档

* 文档Http服务镜像
````
docker create --name yy-doc -p 7071:7071 \
    --restart=always \
    -v /--YourDir--/runtime:/www/web/swagger \
    quay.io/goswagger/swagger serve --no-open --port=7071 /www/web/swagger/sf.json
````
> 此命令为docker容器命令，即`swagger`命令的docker版本，本地开发使用swagger命令即可


### 4、错误码管理方案

> 错误码管理使用go工具链里的 generate

* 在`render\codes.go`新增`CE`类型常量并带上注释，合理划分预留好错误码区间段
* 然后在项目根目录下执行 `go generate ./...` 即可自动生成或更新`render\string.go`文件
* 因`CE`是一个`error`包装类型，在需要抛出错误的位置返回自定义的错误码常量即可，具体输出时底层会转换为`code`和`msg`
* 建议错误码只在 service 和 controller 层的 导出方法 中使用
> 请勿手动编辑`render\string.go`文件，也不要把该文件提交到代码库，CI/CD流程里build前通过`go generate ./...`维护

### 5、配置文件相关提示

* 查找顺序如下
* 按命令行参数 `-config` 指定则优先使用
* 在当前工作目录查找 `conf.toml`
* 在可执行文件所在目录查找 `conf.toml`

### 6、日志

* 日志组件使用：`zap`
* 日志输出选项：`stdout` `stderr` | 日志目录 (配置文件`Log.Path`指定存储目录 或者 命令行参数`--log`指定)
* 日志级别选项：`panic` `fatal` `error` `warning` `info` `debug`（配置文件`Log.Level`指定日志级别 或者 命令行参数`--level`指定）
* 建议生产环境：`info`
* 建议开发测试：`debug`

### 7、路由 & 中间件

* 都放在 `route` 目录下
* `route_manage.go` 存放的是管理后台路由
* `middleware.go` 存放通用路由中间件实现方法
* `middleware_manage.go` 存放管理后台路由中间件实现方法
* 其他以此类推

### 8、队列 & 定时任务 & 命令行工具

* 队列在 `app/console/job` 目录新增job任务，并在`app/console/intstance.go`注册成普通队列Or延时队列
* 定时任务在 `app/console/crontab` 目录新增定时任务，并在`app/console/intstance.go`注册
* 命令行工具在 `app/console/commands` 目录按`cobra`规则新增
> 队列、定时任务、命令行工具均是一个go文件对应一个功能

### 9、本地内存缓存 & redis缓存

* `utils/memory`目录下封装有本地内存缓存的设置、获取等方法
* 其他缓存及分布式锁可使用redis
* 本地内存缓存重启即丢失需重建，但速度优于redis缓存，酌情使用
* 多实例部署时本地内存缓存同步清理通过redis的subscribe机制实现，参考`extend/subscrible`目录

> 注意：本地内存缓存的过期清理是10分钟1次，请勿精确依赖本地内存缓存的过期时间

### 10、数据库迁移

* 查看帮助：`go run main.go migrate`
* 创建迁移文件：`go run main.go migrate create your-filename`(文件名建议使用中划线隔开)
* 查看迁移文件执行情况：`go run main.go migrate status`
* 执行迁移文件：`go run main.go migrate up`
* 执行回滚操作（每次执行只回滚一个迁移文件）：`go run main.go migrate down [filename, 可选]`
* 迁移文件示例：`migrations/example.txt`
* 首次使用需配置好数据库信息后执行数据库迁移初始化基础数据表

## 四、业务相关
