package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ijingjingyang/go-frame/render"
	"net/http"
)

// 中间件 -- 本质上是一个路由操作方法等同 gin.HandlerFunc

// notRoute 找不到路由输出
func notRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, render.H(http.StatusNotFound, http.StatusText(http.StatusNotFound), ""))
	return
}
