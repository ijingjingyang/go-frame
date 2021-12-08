#!/bin/bash

# base method begin
outputFail() {
   echo "$(date "+%Y/%m/%d %H:%M:%S") $1 [fail]"
}

outputOk() {
    echo "$(date "+%Y/%m/%d %H:%M:%S") $1 [ok]"
}

# check go command install
isSwaggerExist() {
    which swagger &> /dev/null
    if [[ "$?" == 1 ]]; then
        outputFail "未找到swagger工具，请先在GOPATH下安装swagger并将GOBIN目录加入环境变量"
        exit 1
    fi
}
# base method end

# step1、检查swagger命令是否安装并且加入了PATH
isSwaggerExist

# step2、生成swagger到runtime目录
swagger_file="./runtime/ff-$(date "+%Y%m%d%H%M%S").json"

# step3、生成swagger文件
outputOk "文件名：${swagger_file}"
outputOk "执行swagger命令生成json文件中，请稍后..."

swagger generate spec -m  -o ${swagger_file}

# step4、开启swagger预览web服务并自动打开默认浏览器
swagger serve --host=0.0.0.0 --port=20218 ${swagger_file}
