package NacosUtils

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

type NacosConfig struct {
	ConfigParam NacosConfigParam
}

type NacosConfigParam struct {
	Host         string
	Port         uint64
	DataId       string
	Group        string
	ClientConfig constant.ClientConfig `json:"ClientConfig"`
}

// GetNacosConfig 从配置文件中获取nacos配置
func GetNacosConfig(nacosConfigFile string) (NacosConfig, error) {
	var nacosConfig NacosConfig
	newViper := viper.New()
	newViper.SetConfigFile(nacosConfigFile)
	//输出配置文件
	err := newViper.ReadInConfig()
	if err != nil {
		return nacosConfig, err
	}
	//解析配置文件
	err = newViper.Unmarshal(&nacosConfig)
	if err != nil {
		return nacosConfig, err
	}
	return nacosConfig, nil

}

// GetServerConfig 获取nacos服务配置
func GetServerConfig(nacosConfig NacosConfig) []constant.ServerConfig {

	return []constant.ServerConfig{
		*constant.NewServerConfig(nacosConfig.ConfigParam.Host, nacosConfig.ConfigParam.Port),
	}

}

// GetNacosClientConfig 获取nacos客户端配置
func GetNacosClientConfig(nacosConfig NacosConfig) *constant.ClientConfig {
	return &nacosConfig.ConfigParam.ClientConfig
}

// GetNacosClient 获取nacos配置客户端
func GetNacosClient(nacosConfig NacosConfig) (config_client.IConfigClient, error) {
	sc := GetServerConfig(nacosConfig)
	cc := GetNacosClientConfig(nacosConfig)
	return clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
}

// GetConfigContentByConfigFile 获取nacos配置文件内容
func GetConfigContentByConfigFile(nacosConfigFile string) (string, error) {
	//获取nacos配置
	nacosConfig, _ := GetNacosConfig(nacosConfigFile)
	//获取nacos配置客户端
	configClient, _ := GetNacosClient(nacosConfig)
	fmt.Println(nacosConfig.ConfigParam.DataId)
	//返回配置文件内容，如果配置文件不存在则返回错误
	return configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.ConfigParam.DataId,
		Group:  nacosConfig.ConfigParam.Group,
	})
}

// GetConfigContentByNacosConfig 根据nacos配置获取配置文件内容
func GetConfigContentByNacosConfig(nacosConfig NacosConfig) (string, error) {
	//获取nacos配置客户端
	configClient, _ := GetNacosClient(nacosConfig)
	fmt.Println(nacosConfig.ConfigParam.DataId)
	//返回配置文件内容，如果配置文件不存在则返回错误
	return configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.ConfigParam.DataId,
		Group:  nacosConfig.ConfigParam.Group,
	})

}

// ListenConfigContentByConfigFile TODO:  想想之后怎么实现 1.配置变化后，热重启或者重新加载配置
// ListenConfigContentByConfigFile 监听nacos配置文件
func ListenConfigContentByConfigFile(nacosConfigFile *string) {
	var nacosConfig, _ = GetNacosConfig(*nacosConfigFile)
	fmt.Println(nacosConfig)
	//监听配置文件
	configClient, err := GetNacosClient(nacosConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: nacosConfig.ConfigParam.DataId,
		Group:  nacosConfig.ConfigParam.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生了变化...")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}
