// Package scripts sufficient接口swagger文档
//
// sufficient接口文档
//
// 测试BaseApi：https://sufficient-api.jing.dev
//
// 生产BaseApi：https://api.jing.dev
//
// 说明:
//
// 管理后台
//    1. 鉴权认证: header 头信息 `Authorization: Bearer {manage_token}`
//    2. ----其中 customer_token 登录接口（`/manage/login`）
//    3. 特定错误码 xxx => 账号封禁
//    4. 特定错误码 xxx => 鉴权失败，需要登录
//
//     Schemes: http, https
//     Host: api.jing.dev
//     BasePath: /
//     Version: 0.0.1
//     Contact: Jea杨<jjonline@jjonline.cn>
//
//     Consumes:
//     - application/json
//	   - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: Authorization
//          in: header
//          example: Bearer <token>
// swagger:meta
package scripts

/* Notice 该文件仅用于生成swagger文档说明，项目实际不引用 */
