package model

import (
	"context"
	"fmt"
	"github.com/ijingjingyang/go-frame/client"
	"github.com/ijingjingyang/go-frame/conf"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// region 定义model`基类`和基础字段结构

// 定义Where条件常量便于IDE追踪
const (
	In         = "IN"          // 在指定枚举值范围内
	NotIn      = "NOT IN"      // 不在指定枚举值范围内
	Between    = "BETWEEN"     // 在指定范围内
	NotBetween = "NOT BETWEEN" // 不在指定范围内
	Like       = "LIKE"        // 模糊查询like
	FindINSet  = "FIND_IN_SET" // find in set
	Raw        = "RAW"         // 原样输出
	EQ         = "="           // equal 等于<=>
	NEQ        = "<>"          // not equal 不等于(!=、<>)
	GT         = ">"           // greater than 大于
	GTE        = ">="          // greater than or equal 大于等于
	LT         = "<"           // less than 小于
	LTE        = "<"           // less than or equal 小于等于
)

// BaseField 通用基础字段
type BaseField struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt uint `gorm:"autoCreateTime"`
	UpdatedAt uint `gorm:"autoCreateTime,autoUpdateTime"`
}

// dbPrefix 统一表前缀
func dbPrefix() string {
	return conf.Config.Database.Prefix
}

// endregion

// region 定义where查询条件结构

// Where 定义where查询条件
//  切片Where多个条件最终SQL中将使用 AND 符拼接
//  - Field 查询的字段名称字符串
//  - Op    查询的条件符号，建议使用包内常量，全部大写[IN、NOT IN、BETWEEN、NOT BETWEEN、LIKE、FIND_IN_SET、=、<>、>=、<、<=、RAW等]
//    当使用RAW时将忽略Op值，直接使用gorm提供的原生Where方法构建<注意转义特殊字符特别留意SQL注入风险>
//  - Value 查询的条件值，当Op为LIKE时无需添加前后的百分号<%>，方法体自动添加前后百分号
type Where struct {
	Field string
	Op    string
	Value interface{}
}

// toWhere 将where传给db查询构造器
func (w Where) toWhere(query *gorm.DB) *gorm.DB {
	switch w.Op {
	case In:
		query = query.Where(fmt.Sprintf("`%s` IN (?)", w.Field), w.Value)
	case NotIn:
		query = query.Where(fmt.Sprintf("`%s` NOT IN (?)", w.Field), w.Value)
	case Between:
		v := w.Value.([]interface{})
		query = query.Where(fmt.Sprintf("`%s` BETWEEN ? AND ?", w.Field), v[0], v[1])
	case NotBetween:
		v := w.Value.([]interface{})
		query = query.Where(fmt.Sprintf("`%s` NOT BETWEEN ? AND ?", w.Field), v[0], v[1])
	case Like:
		query = query.Where(fmt.Sprintf("`%s` LIKE ?", w.Field), "%"+w.Value.(string)+"%")
	case FindINSet:
		query = query.Where(fmt.Sprintf("FIND_IN_SET(?, %s)", w.Field), w.Value)
	case Raw:
		query = query.Where(w.Field, w.Value.([]interface{})...)
	default:
		// "<>", "!=", ">", ">=", "=", "<", "<="
		query = query.Where(fmt.Sprintf("`%s` %s ?", w.Field, w.Op), w.Value)
	}
	return query
}

// parseWhere 解析where条件
func parseWhere(query *gorm.DB, wheres []Where) *gorm.DB {
	if wheres == nil {
		return query
	}
	for _, w := range wheres {
		query = w.toWhere(query)
	}
	return query
}

// endregion

// region 定义通用模型方法

// model 模型通用方法封装
//  - 将本结构体嵌入模型方法即可
type model struct {
	self schema.Tabler // 模型实例
}

// construct 模型初始化设置实例
//  - 模型定义 init 方法中设置通用方法依赖的模型结构实例
//  - @param instance 传参模型引用，即指针
func (m *model) construct(instance schema.Tabler) {
	m.self = instance
}

// FindByPrimary 单主键查询1条记录<仅支持单字段主键不支持符合主键>
//  - 查询不到记录返回 gorm.ErrRecordNotFound 的error
//     使用errors.Is(result.Error, gorm.ErrRecordNotFound)进行判断
//  - pVal：Primary Value 需要查询的主键字段的值<仅支持数值类型主键，字符串类型不支持>
// 	   自定义模型中主键字段通过tag标签 gorm:"primaryKey" 指定的即为主键
// 	   若未定义主键，则取定义模型中的第一个字段作为主键字段使用
//	   绝大多数情况下自增主键可使用；复合主键或字符串主键请使用 FindByWhere
//  - target：查询结果集模型引用，形参为 schema.Tabler 接口类型<注意：指针也是接口类型>，需传指针
//  - fields：查询的字段，可选不传表示默认查询出所有字段
// 	   可变参数查询多个字段，使用字符串切片可变参数展开模式可使用字符串切片数组
//	   []string{"name", "sex"}...
func (m *model) FindByPrimary(ctx context.Context, pVal interface{}, target schema.Tabler, fields ...string) (err error) {
	return client.DB.WithContext(ctx).Table(m.self.TableName()).Select(fields).Take(target, pVal).Error
}

// FindByWhere 按where条件查询1条记录<若有多条记录符合要求取按主键升序的第一条记录>
//  - 查询不到记录返回 gorm.ErrRecordNotFound 的error
// 	   使用errors.Is(result.Error, gorm.ErrRecordNotFound)进行判断
//  - where：查询条件 Where 结构体切片
//  - target：查询结果集模型引用，形参为 schema.Tabler 接口类型<注意：指针也是接口类型>，需传指针
//  - fields：查询的字段，可选不传表示默认查询出所有字段
// 	   可变参数查询多个字段，使用字符串切片可变参数展开模式可使用字符串切片数组
//	   []string{"name", "sex"}...
func (m *model) FindByWhere(ctx context.Context, where []Where, target schema.Tabler, fields ...string) (err error) {
	return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Select(fields).First(target).Error
}

// CountByWhere 按where条件统计记录总数
//  - 查询不到记录返回返回0
//  - where：查询条件 Where 结构体切片
//  - targetTotal：查询结果集总条数指针引用
func (m *model) CountByWhere(ctx context.Context, where []Where, targetTotal *int64) (err error) {
	return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Count(targetTotal).Error
}

// ListByWhere 按where条件查询多条确认数量有限的列表记录
//	特别注意：该方法用于查询有限数量的多条记录，不设置limit条件
//  - 查询不到记录返回 gorm.ErrRecordNotFound 的error
// 	   使用errors.Is(result.Error, gorm.ErrRecordNotFound)进行判断
//  - where：查询条件 Where 结构体切片
//  - target：查询结果集模型的切片引用，形参为 interface，需要传具体model实现<schema.Tabler>的切片引用：需传指针
//     例如：var a []Ad 传参 &a
//  - orderBy：排序条件字符串
//     例子1：`name` <表示按name字段升序>
//     例子2：`name` ASC <表示按name字段升序>
//     例子3：`name` ASC, `ID` DESC <表示按name字段升序后按ID降序>
//     例子4：给空字符串表示不设置排序条件
//  - fields：查询的字段，可选不传表示默认查询出所有字段
// 	   可变参数查询多个字段，使用字符串切片可变参数展开模式可使用字符串切片数组
//	   []string{"name", "sex"}...
func (m *model) ListByWhere(ctx context.Context, where []Where, target interface{}, orderBy string, fields ...string) (err error) {
	if orderBy == "" {
		return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Select(fields).Find(target).Error
	}
	return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Select(fields).Order(orderBy).Find(target).Error
}

// ListByOriginWhere 按gorm原始where查询多条确认数量有限的列表记录
//  ①、可用于构造查询条件较为复杂的SQL，包括查询条件包含or、and以及使用括号连起来的分组条件
//     例如：select * from tb where (id=1 OR name="tom") OR ((sex=0 OR title="leader") AND dept_id=10)
//	②、特别注意：该方法用于查询有限数量的多条记录，不设置limit条件
//  - 查询不到记录返回 gorm.ErrRecordNotFound 的error
// 	   使用errors.Is(result.Error, gorm.ErrRecordNotFound)进行判断
//  - originWhere：gorm原始查询条件，使用gorm全局句柄构建，参考下方例子伪代码构造
//  - target：查询结果集模型的切片引用，形参为 interface，需要传具体model实现<schema.Tabler>的切片引用：需传指针
//     例如：var a []Ad 传参 &a
//  - orderBy：排序条件字符串
//     例子1：`name` <表示按name字段升序>
//     例子2：`name` ASC <表示按name字段升序>
//     例子3：`name` ASC, `ID` DESC <表示按name字段升序后按ID降序>
//     例子4：给空字符串表示不设置排序条件
//  - fields：查询的字段，可选不传表示默认查询出所有字段
// 	   可变参数查询多个字段，使用字符串切片可变参数展开模式可使用字符串切片数组
//	   []string{"name", "sex"}...
//  例如上述列出SQL的伪代码(注意gg即为originWhere，由g1、g2组合构成)：
//     g1 := client.DB.Where("id=?", 1).Or("name=?", "tom")
//     g2 := client.DB.Where(client.DB.Where("sex=?", 0).Or("title=?", "leader")).Where("dept_id=?", 10)
//     gg := client.DB.Where(g1).Or(g2)
//     err := tbModel.ListByOriginWhere(ctx, gg, &target, "id")
func (m *model) ListByOriginWhere(ctx context.Context, originWhere *gorm.DB, target interface{}, orderBy string, fields ...string) (err error) {
	if orderBy == "" {
		return client.DB.WithContext(ctx).Table(m.self.TableName()).Where(originWhere).Select(fields).Find(target).Error
	}
	return client.DB.WithContext(ctx).Table(m.self.TableName()).Where(originWhere).Select(fields).Order(orderBy).Find(target).Error
}

// Columns 获取列数据
//  - where：查询条件 Where 结构体切片
//  - target：查询结果集切片引用，形参为 interface，例如: var a []uint32 \ var b []string 则传参 &a \ &b
//  - fields：查询的单个字段列名称
//  - orderBy：排序条件字符串
//     例子1：`name` <表示按name字段升序>
//     例子2：`name` ASC <表示按name字段升序>
//     例子3：`name` ASC, `ID` DESC <表示按name字段升序后按ID降序>
//     例子4：给空字符串表示不设置排序条件
func (m *model) Columns(ctx context.Context, where []Where, target interface{}, field string, orderBy string) error {
	if orderBy == "" {
		return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Pluck(field, target).Error
	}
	return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Order(orderBy).Pluck(field, target).Error
}

// Paginate 分页查询，按where条件分页查询获取分页列表数和总记录数
//  - 查询不到记录返回返回0和空切片
//  - where：查询条件 Where 结构体切片
//  - target：查询结果集模型的切片引用，形参为 interface，需要传具体model实现<schema.Tabler>的切片引用：需传指针
//     例如：var a []Ad 传参 &a
//  - targetTotal：查询结果集总条数指针引用
//  - page：查询分页的当前页码数，从1开始
//  - limit：查询分页的1页多少条限制
//  - orderBy：排序条件字符串
//     例子1：`name` <表示按name字段升序>
//     例子2：`name` ASC <表示按name字段升序>
//     例子3：`name` ASC, `ID` DESC <表示按name字段升序后按ID降序>
//     例子4：给空字符串表示不设置排序条件
//  - fields：查询的字段，可选不传表示默认查询出所有字段
// 	   可变参数查询多个字段，使用字符串切片可变参数展开模式可使用字符串切片数组
//	   []string{"name", "sex"}...
func (m *model) Paginate(ctx context.Context, where []Where, target interface{}, targetTotal *int64, page, limit int, orderBy string, fields ...string) (err error) {
	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	// calc total count
	err = parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Count(targetTotal).Error
	if err != nil {
		return
	}
	if orderBy == "" {
		return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Offset(offset).Limit(limit).Select(fields).Find(target).Error
	}
	return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Offset(offset).Limit(limit).Select(fields).Order(orderBy).Find(target).Error
}

// SimplePaginate 简单分页查询，按where条件分页简要查询获取分页列表数据<不返回总条数减少1次查询>
//  - 查询不到记录返回返回空切片
//  - where：查询条件 Where 结构体切片
//  - target：查询结果集模型的切片引用，形参为 interface，需要传具体model实现<schema.Tabler>的切片引用：需传指针
//     例如：var a []Ad 传参 &a
//  - page：查询分页的当前页码数，从1开始
//  - limit：查询分页的1页多少条限制
//  - orderBy：排序条件字符串
//     例子1：`name` <表示按name字段升序>
//     例子2：`name` ASC <表示按name字段升序>
//     例子3：`name` ASC, `ID` DESC <表示按name字段升序后按ID降序>
//     例子4：给空字符串表示不设置排序条件
//  - fields：查询的字段，可选不传表示默认查询出所有字段
// 	   可变参数查询多个字段，使用字符串切片可变参数展开模式可使用字符串切片数组
//	   []string{"name", "sex"}...
func (m *model) SimplePaginate(ctx context.Context, where []Where, target interface{}, page, limit int, orderBy string, fields ...string) (err error) {
	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	if orderBy == "" {
		return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Offset(offset).Limit(limit).Select(fields).Find(target).Error
	}
	return parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Offset(offset).Limit(limit).Select(fields).Order(orderBy).Find(target).Error
}

// InsertOne 通过赋值模型创建1条记录
//  - data：model实例对象，请注意不要给主键字段赋值，创建成功后主键字段将填充创建记录的主键值
// 	   例如：var a Ad; a.XX="yy"; 传参 &a
//  - fields：限定创建语句写入的字段名<指定的字段的零值也会写入>，不传留空则取结构体非零字段创建1条新记录
//	   []string{"name", "sex"}...
func (m *model) InsertOne(ctx context.Context, data schema.Tabler, fields ...string) error {
	return client.DB.WithContext(ctx).Table(m.self.TableName()).Select(fields).Create(data).Error
}

// InsertOneUseMap 通过map键值对创建1条记录
//  - 注意：创建成功后不会回填主键
//  - data：新增1条记录的map数据，注意map的键为数据库实际字段名而不是结构体FieldName，支持零值
func (m *model) InsertOneUseMap(ctx context.Context, data map[string]interface{}) error {
	return client.DB.WithContext(ctx).Table(m.self.TableName()).Create(data).Error
}

// MultiInsert 通过赋值模型批量创建
//  - data：model实例对象切片，请注意不要给主键字段赋值，创建成功后主键字段将填充创建记录的主键值
// 	   例如：var a []Ad; 传参 &a
//  - fields：限定创建语句写入的字段名<指定的字段的零值也会写入>，不传留空则取结构体非零字段创建1条新记录
//	   []string{"name", "sex"}...
func (m *model) MultiInsert(ctx context.Context, data interface{}, fields ...string) error {
	return client.DB.WithContext(ctx).Table(m.self.TableName()).Select(fields).Create(data).Error
}

// MultiInsertUseMap 通过map键值对切片批量创建
//  - 注意：创建成功后不会回填主键
//  - data：批量新增多条记录的map数据，注意map的键为数据库实际字段名而不是结构体FieldName，支持零值
func (m *model) MultiInsertUseMap(ctx context.Context, data []map[string]interface{}) error {
	return client.DB.WithContext(ctx).Table(m.self.TableName()).Create(data).Error
}

// UpdateOne 通过model的主键字段更新指定字段
//  - 注意：model对象指定需要更新的单条记录的值，主键字段必须指定值
//     然后通过第二个可选参数指定需要更新的字段<指定的字段的零值也会被更新为对应的零值，不指定则data结构体中零值字段不会更新>
//  - data：model实例对象，字段填充好需要更新的值
// 	   例如：var a Ad; a.ID=1;a.Name="Tom" 传参 &a
//  - fields：限定创建语句写入的字段名<指定的字段的零值也会写入>，不传留空则取结构体非零字段去更新
//	   []string{"name", "sex"}...
func (m *model) UpdateOne(ctx context.Context, data schema.Tabler, fields ...string) (int64, error) {
	result := client.DB.WithContext(ctx).Table(m.self.TableName()).Select(fields).Updates(data)
	return result.RowsAffected, result.Error
}

// UpdateByWhere 通过where条件更新记录
//  - 注意：model对象指定需要更新的单条记录的值，主键字段必须指定值
//     然后通过第二个参数指定需要更新的字段<指定的字段的零值也会被更新为对应的零值>
//  - where：查询条件 Where 结构体切片
//  - data：更新的字段map数据，注意map的键为数据库实际字段名而不是结构体FieldName，支持零值更新
func (m *model) UpdateByWhere(ctx context.Context, where []Where, data map[string]interface{}) (int64, error) {
	result := parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Updates(data)
	return result.RowsAffected, result.Error
}

// DeleteByPrimary 通过主键字段值硬删除1条或多条记录
//  - 注意：model对象指定需要删除的记录的主键值
//     gorm软删除特性需在model字段定义时使用 gorm.DeletedAt 类型的字段特性实现，本model不支持软删除
//     gorm软删除特性引入后查询条件将自动附加过滤已软删除记录的条件
//  - primaryKey：数值类型的主键值，1个或多个
// 	   例子1：1 单个主键
//     例子2："1" 单个主键字符串类型字面量，本质是个数值
//     例子3：1,3,4 多个数值类型<多个需要相同类型>
//     例子4："1","3","4" 多个数值字符串<多个需要相同类型>
func (m *model) DeleteByPrimary(ctx context.Context, primaryKey ...interface{}) (int64, error) {
	result := client.DB.WithContext(ctx).Table(m.self.TableName()).Delete(m.self, primaryKey)
	return result.RowsAffected, result.Error
}

// DeleteByWhere 通过model的主键字段硬删除1条记录
//  - 注意：软删除功能是一项功能特性，要么全部使用gorm的软删除特性，要么业务删除时调用update方法
//     gorm软删除特性需在model字段定义时使用 gorm.DeletedAt 类型的字段特性实现
//     gorm软删除特性引入后查询条件将自动附加过滤已软删除记录的条件，无需手动指定
//  - where：条件 Where 结构体切片
func (m *model) DeleteByWhere(ctx context.Context, where []Where) (int64, error) {
	result := parseWhere(client.DB.WithContext(ctx).Table(m.self.TableName()), where).Delete(m.self)
	return result.RowsAffected, result.Error
}

// endregion
