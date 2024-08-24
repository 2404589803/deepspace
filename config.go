package main

import "gopkg.in/yaml.v3"

// 全局变量，用于保存配置
var DeepConfig Config

// Config 结构体用于表示配置文件的结构
type Config struct {
	Endpoint string       `yaml:"endpoint"` // 配置中的 "endpoint" 字段
	Start    *StartConfig `yaml:"start"`    // 配置中的 "start" 字段，它是 StartConfig 类型的指针
}

// init 函数在程序启动时自动调用
func init() {
	// 获取配置文件的文件句柄
	if file := getConfig(); file != nil {
		// 使用 YAML 解码器将文件内容解码到 DeepConfig 变量中
		if err := yaml.NewDecoder(file).Decode(&DeepConfig); err != nil {
			// 如果解码过程中发生错误，调用 logFatal 函数处理错误
			logFatal(err)
		}
	}
}
