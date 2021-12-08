package cfg

import (
	"bytes"
	"github.com/spf13/viper"
	"os"
)

// Viper 依赖viper实现的文件模式配置解析
// viper支持的格式见包内常量
type Viper struct{}

// Parse 解析映射配置至结构体变量
//  - @param resource 文件路径（string） 或 配置资源（[]byte）
//  - @param cType    配置文件类型：请使用包内Type开头的常量，例如： TypeToml
//  - @param target   解析结果集引用（struct pointer）
func (v Viper) Parse(resource interface{}, cType string, target interface{}) error {
	if !IsCfgTypeSupport(cType) {
		return ErrConfigTypeNotSupport
	}

	var stream []byte
	var err error

	switch resource.(type) {
	case string:
		// file dir for check file exist
		filePath := resource.(string)
		if !IsFileExist(filePath) {
			return ErrConfigFileNotExist
		}

		stream, err = os.ReadFile(filePath)
		if err != nil {
			return ErrConfigFileParseFailed
		}
	case []byte:
		stream = resource.([]byte)
	default:
		return ErrConfigFileNotExist
	}

	// use viper parse
	vip := viper.New()
	vip.SetConfigType(cType)

	if err = vip.ReadConfig(bytes.NewReader(stream)); err != nil {
		return ErrConfigFileParseFailed
	}

	// set config value
	return vip.Unmarshal(target)
}
